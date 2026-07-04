package application

import (
	"context"
	"testing"

	"workspace-app/internal/platform/search"
)

// Covers Ready, Suggest, IndexJob, the limit clamping and the userDoc email
// fallback, reusing the in-package fakeEngine/fakeSource.

func TestReadyReflectsEngine(t *testing.T) {
	if !NewService(&fakeEngine{ready: true}, &fakeSource{}).Ready() {
		t.Fatal("expected Ready=true when engine is ready")
	}
	if NewService(&fakeEngine{ready: false}, &fakeSource{}).Ready() {
		t.Fatal("expected Ready=false when engine is not ready")
	}
}

func TestSuggestEngineAndFallback(t *testing.T) {
	ctx := context.Background()

	// Engine ready: suggestions come from the engine, not the source.
	eng := &fakeEngine{ready: true, hits: []search.Hit{{Type: "skill", RefID: "go", Title: "Go"}}}
	src := &fakeSource{}
	svc := NewService(eng, src)
	hits, err := svc.Suggest(ctx, "g", 5)
	if err != nil {
		t.Fatalf("suggest: %v", err)
	}
	if len(hits) != 1 || hits[0].Title != "Go" {
		t.Fatalf("expected engine suggestion, got %+v", hits)
	}
	if src.fallbackCalled {
		t.Fatal("fallback must not be used when engine is ready")
	}

	// Engine not ready: falls back to the source.
	src2 := &fakeSource{fallbackHits: []search.Hit{{Type: "user", RefID: "u1", Title: "Asha"}}}
	svc2 := NewService(&fakeEngine{ready: false}, src2)
	if _, err := svc2.Suggest(ctx, "as", 5); err != nil {
		t.Fatalf("suggest fallback: %v", err)
	}
	if !src2.fallbackCalled {
		t.Fatal("expected DB fallback when engine not ready")
	}

	// Blank query short-circuits to an empty slice without touching the source.
	src3 := &fakeSource{}
	if hits, err := NewService(&fakeEngine{ready: true}, src3).Suggest(ctx, "   ", 5); err != nil || len(hits) != 0 {
		t.Fatalf("expected empty suggestion, got %v %v", hits, err)
	}
	if src3.fallbackCalled {
		t.Fatal("blank query must not hit the source")
	}
}

func TestIndexJobBuildsDoc(t *testing.T) {
	eng := &fakeEngine{ready: true}
	src := &fakeSource{jobs: []JobRow{{ID: "j1", Title: "Ops Manager", Company: "Acme", Location: "Remote", Description: "Lead ops"}}}
	svc := NewService(eng, src)

	svc.IndexJob(context.Background(), "j1")
	if len(eng.indexed) != 1 {
		t.Fatalf("expected one indexed doc, got %d", len(eng.indexed))
	}
	d := eng.indexed[0]
	if d.Type != "job" || d.RefID != "j1" || d.Title != "Ops Manager" || d.Subtitle != "Acme · Remote" {
		t.Fatalf("unexpected job doc %+v", d)
	}

	// Unknown id indexes nothing.
	svc.IndexJob(context.Background(), "missing")
	if len(eng.indexed) != 1 {
		t.Fatalf("expected no extra index for unknown job, got %d", len(eng.indexed))
	}
}

func TestQueryClampsLimitBounds(t *testing.T) {
	ctx := context.Background()
	src := &fakeSource{fallbackHits: []search.Hit{{Type: "user", RefID: "u1"}}}
	svc := NewService(&fakeEngine{ready: false}, src)

	// limit <= 0 and limit > 50 both go through clampLimit without error.
	if _, err := svc.Query(ctx, "x", nil, 0); err != nil {
		t.Fatalf("query limit 0: %v", err)
	}
	if _, err := svc.Query(ctx, "x", nil, 100); err != nil {
		t.Fatalf("query limit 100: %v", err)
	}
}

func TestUserDocEmailFallbackTitle(t *testing.T) {
	eng := &fakeEngine{ready: true}
	// No FullName: the indexed doc title falls back to the email.
	src := &fakeSource{users: []UserRow{{ID: "u9", Email: "noname@cb.io"}}}
	svc := NewService(eng, src)

	svc.IndexUser(context.Background(), "u9")
	if len(eng.indexed) != 1 || eng.indexed[0].Title != "noname@cb.io" {
		t.Fatalf("expected title to fall back to email, got %+v", eng.indexed)
	}
}
