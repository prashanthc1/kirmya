// Package search provides the platform's full-text search engine port and an
// OpenSearch-backed implementation (thin REST client — no heavy SDK), plus a
// Noop used when OpenSearch is not configured or unreachable. Indexing is
// best-effort (errors logged, never returned); querying returns errors so
// callers can fall back to a database search. Modules depend on a structurally
// identical Engine interface declared in their own application package.
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Doc is a single searchable document. ID is unique per (Type, RefID).
type Doc struct {
	Type     string // user | job | community | skill
	RefID    string // entity id in its owning context
	Title    string
	Subtitle string
	Body     string
	URL      string
}

func (d Doc) docID() string { return d.Type + ":" + d.RefID }

// Hit is a single search result.
type Hit struct {
	Type     string  `json:"type"`
	RefID    string  `json:"ref_id"`
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	URL      string  `json:"url"`
	Score    float64 `json:"score"`
}

// Engine is the platform search port.
type Engine interface {
	// Ready reports whether a real search backend is available. When false,
	// callers should fall back to a database query.
	Ready() bool
	// Index upserts a document (best-effort).
	Index(ctx context.Context, doc Doc)
	// Delete removes a document by type+id (best-effort).
	Delete(ctx context.Context, typ, refID string)
	// Search runs a fuzzy multi-field query, optionally filtered by types.
	Search(ctx context.Context, query string, types []string, limit int) ([]Hit, error)
	// Suggest runs a prefix query for autocomplete.
	Suggest(ctx context.Context, query string, limit int) ([]Hit, error)
}

const indexName = "kirmya"

// New returns an OpenSearch-backed engine when OPENSEARCH_URL is set and the
// cluster is reachable; otherwise a Noop engine (Ready()==false). It never
// returns an error so the platform degrades gracefully without OpenSearch.
func New() Engine {
	base := strings.TrimRight(os.Getenv("OPENSEARCH_URL"), "/")
	if base == "" {
		log.Printf("[search] OPENSEARCH_URL not set; full-text search disabled (DB fallback)")
		return Noop{}
	}
	eng := &OpenSearch{
		base:   base,
		user:   os.Getenv("OPENSEARCH_USER"),
		pass:   os.Getenv("OPENSEARCH_PASSWORD"),
		client: &http.Client{Timeout: 5 * time.Second},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := eng.ensureIndex(ctx); err != nil {
		log.Printf("[search] OpenSearch unavailable (%v); full-text search disabled (DB fallback)", err)
		return Noop{}
	}
	log.Printf("[search] OpenSearch connected; index %q ready", indexName)
	return eng
}

// OpenSearch is a thin REST client implementing Engine.
type OpenSearch struct {
	base   string
	user   string
	pass   string
	client *http.Client
}

func (o *OpenSearch) Ready() bool { return true }

func (o *OpenSearch) Index(ctx context.Context, doc Doc) {
	body := map[string]any{
		"type": doc.Type, "ref_id": doc.RefID, "title": doc.Title,
		"subtitle": doc.Subtitle, "body": doc.Body, "url": doc.URL,
	}
	path := fmt.Sprintf("/%s/_doc/%s?refresh=false", indexName, doc.docID())
	if _, err := o.do(ctx, http.MethodPut, path, body); err != nil {
		log.Printf("[search] index %s: %v", doc.docID(), err)
	}
}

func (o *OpenSearch) Delete(ctx context.Context, typ, refID string) {
	path := fmt.Sprintf("/%s/_doc/%s:%s?refresh=false", indexName, typ, refID)
	if _, err := o.do(ctx, http.MethodDelete, path, nil); err != nil {
		log.Printf("[search] delete %s:%s: %v", typ, refID, err)
	}
}

func (o *OpenSearch) Search(ctx context.Context, query string, types []string, limit int) ([]Hit, error) {
	q := map[string]any{
		"multi_match": map[string]any{
			"query":     query,
			"fields":    []string{"title^3", "subtitle^2", "body"},
			"fuzziness": "AUTO",
			"type":      "best_fields",
		},
	}
	return o.runQuery(ctx, q, types, limit)
}

func (o *OpenSearch) Suggest(ctx context.Context, query string, limit int) ([]Hit, error) {
	q := map[string]any{
		"match_phrase_prefix": map[string]any{"title": map[string]any{"query": query}},
	}
	return o.runQuery(ctx, q, nil, limit)
}

func (o *OpenSearch) runQuery(ctx context.Context, inner map[string]any, types []string, limit int) ([]Hit, error) {
	query := inner
	if len(types) > 0 {
		query = map[string]any{
			"bool": map[string]any{
				"must":   inner,
				"filter": map[string]any{"terms": map[string]any{"type": types}},
			},
		}
	}
	body := map[string]any{"size": limit, "query": query}
	raw, err := o.do(ctx, http.MethodPost, "/"+indexName+"/_search", body)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Hits struct {
			Hits []struct {
				Score  float64 `json:"_score"`
				Source struct {
					Type     string `json:"type"`
					RefID    string `json:"ref_id"`
					Title    string `json:"title"`
					Subtitle string `json:"subtitle"`
					URL      string `json:"url"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	hits := make([]Hit, 0, len(parsed.Hits.Hits))
	for _, h := range parsed.Hits.Hits {
		hits = append(hits, Hit{
			Type: h.Source.Type, RefID: h.Source.RefID, Title: h.Source.Title,
			Subtitle: h.Source.Subtitle, URL: h.Source.URL, Score: h.Score,
		})
	}
	return hits, nil
}

// ensureIndex creates the index with an autocomplete-friendly mapping if absent.
func (o *OpenSearch) ensureIndex(ctx context.Context) error {
	status, err := o.status(ctx, http.MethodHead, "/"+indexName)
	if err != nil {
		return err
	}
	if status == http.StatusOK {
		return nil
	}
	mapping := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"type":     map[string]any{"type": "keyword"},
				"ref_id":   map[string]any{"type": "keyword"},
				"title":    map[string]any{"type": "text"},
				"subtitle": map[string]any{"type": "text"},
				"body":     map[string]any{"type": "text"},
				"url":      map[string]any{"type": "keyword", "index": false},
			},
		},
	}
	_, err = o.do(ctx, http.MethodPut, "/"+indexName, mapping)
	return err
}

func (o *OpenSearch) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, o.base+path, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if o.user != "" {
		req.SetBasicAuth(o.user, o.pass)
	}
	res, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	data, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("opensearch %s %s: %d %s", method, path, res.StatusCode, truncate(data))
	}
	return data, nil
}

func (o *OpenSearch) status(ctx context.Context, method, path string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, o.base+path, nil)
	if err != nil {
		return 0, err
	}
	if o.user != "" {
		req.SetBasicAuth(o.user, o.pass)
	}
	res, err := o.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = res.Body.Close() }()
	_, _ = io.Copy(io.Discard, res.Body)
	return res.StatusCode, nil
}

func truncate(b []byte) string {
	s := string(b)
	if len(s) > 200 {
		return s[:200]
	}
	return s
}

// Noop is the engine used when OpenSearch is not configured. Ready() is false so
// callers fall back to a database search; index/delete are no-ops.
type Noop struct{}

func (Noop) Ready() bool                            { return false }
func (Noop) Index(context.Context, Doc)             {}
func (Noop) Delete(context.Context, string, string) {}
func (Noop) Search(context.Context, string, []string, int) ([]Hit, error) {
	return nil, nil
}
func (Noop) Suggest(context.Context, string, int) ([]Hit, error) { return nil, nil }
