package application

import (
	"context"
	"testing"

	"workspace-app/internal/platform/search"
)

// fakeEngine records indexed docs and returns canned hits. ready toggles whether
// queries use the engine or fall back to the source.
type fakeEngine struct {
	ready   bool
	indexed []search.Doc
	deleted []string
	hits    []search.Hit
}

func (f *fakeEngine) Ready() bool                           { return f.ready }
func (f *fakeEngine) Index(_ context.Context, d search.Doc) { f.indexed = append(f.indexed, d) }
func (f *fakeEngine) Delete(_ context.Context, typ, id string) {
	f.deleted = append(f.deleted, typ+":"+id)
}
func (f *fakeEngine) Search(context.Context, string, []string, int) ([]search.Hit, error) {
	return f.hits, nil
}
func (f *fakeEngine) Suggest(context.Context, string, int) ([]search.Hit, error) {
	return f.hits, nil
}

// fakeSource returns canned rows and records fallback calls.
type fakeSource struct {
	users          []UserRow
	jobs           []JobRow
	comms          []CommunityRow
	skills         []string
	fallbackHits   []search.Hit
	fallbackCalled bool
}

func (s *fakeSource) User(_ context.Context, id string) (*UserRow, error) {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i], nil
		}
	}
	return nil, nil
}
func (s *fakeSource) Job(_ context.Context, id string) (*JobRow, error) {
	for i := range s.jobs {
		if s.jobs[i].ID == id {
			return &s.jobs[i], nil
		}
	}
	return nil, nil
}
func (s *fakeSource) AllUsers(context.Context) ([]UserRow, error)            { return s.users, nil }
func (s *fakeSource) AllJobs(context.Context) ([]JobRow, error)              { return s.jobs, nil }
func (s *fakeSource) AllCommunities(context.Context) ([]CommunityRow, error) { return s.comms, nil }
func (s *fakeSource) AllSkills(context.Context) ([]string, error)            { return s.skills, nil }
func (s *fakeSource) FallbackSearch(_ context.Context, _ string, _ []string, _ int) ([]search.Hit, error) {
	s.fallbackCalled = true
	return s.fallbackHits, nil
}

func TestQueryUsesEngineWhenReady(t *testing.T) {
	eng := &fakeEngine{ready: true, hits: []search.Hit{{Type: "job", RefID: "j1", Title: "Ops Manager"}}}
	src := &fakeSource{fallbackHits: []search.Hit{{Type: "user", RefID: "u1"}}}
	svc := NewService(eng, src)

	hits, err := svc.Query(context.Background(), "ops", nil, 10)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(hits) != 1 || hits[0].Title != "Ops Manager" {
		t.Fatalf("expected engine hit, got %+v", hits)
	}
	if src.fallbackCalled {
		t.Fatal("fallback should not be used when engine is ready")
	}
}

func TestQueryFallsBackToDBWhenNotReady(t *testing.T) {
	eng := &fakeEngine{ready: false}
	src := &fakeSource{fallbackHits: []search.Hit{{Type: "user", RefID: "u1", Title: "Asha"}}}
	svc := NewService(eng, src)

	hits, err := svc.Query(context.Background(), "asha", []string{"user"}, 10)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !src.fallbackCalled {
		t.Fatal("expected DB fallback when engine not ready")
	}
	if len(hits) != 1 || hits[0].Title != "Asha" {
		t.Fatalf("expected fallback hit, got %+v", hits)
	}
}

func TestEmptyQueryShortCircuits(t *testing.T) {
	src := &fakeSource{}
	svc := NewService(&fakeEngine{ready: true}, src)
	hits, err := svc.Query(context.Background(), "   ", nil, 10)
	if err != nil || len(hits) != 0 {
		t.Fatalf("expected empty result, got %v %v", hits, err)
	}
	if src.fallbackCalled {
		t.Fatal("empty query must not hit the source")
	}
}

func TestIndexUserBuildsDoc(t *testing.T) {
	eng := &fakeEngine{ready: true}
	src := &fakeSource{users: []UserRow{{ID: "u1", FullName: "Asha Rao", Email: "asha@cb.io", Headline: "Ops Lead"}}}
	svc := NewService(eng, src)

	svc.IndexUser(context.Background(), "u1")
	if len(eng.indexed) != 1 {
		t.Fatalf("expected one indexed doc, got %d", len(eng.indexed))
	}
	d := eng.indexed[0]
	if d.Type != "user" || d.RefID != "u1" || d.Title != "Asha Rao" || d.Subtitle != "Ops Lead" {
		t.Fatalf("unexpected user doc %+v", d)
	}
}

func TestBackfillIndexesEverything(t *testing.T) {
	eng := &fakeEngine{ready: true}
	src := &fakeSource{
		users:  []UserRow{{ID: "u1", FullName: "A"}},
		jobs:   []JobRow{{ID: "j1", Title: "Job"}},
		comms:  []CommunityRow{{ID: "c1", Slug: "tech", Name: "Technology"}},
		skills: []string{"Go", "PostgreSQL"},
	}
	svc := NewService(eng, src)

	svc.Backfill(context.Background())
	// 1 user + 1 job + 1 community + 2 skills = 5
	if len(eng.indexed) != 5 {
		t.Fatalf("expected 5 indexed docs, got %d", len(eng.indexed))
	}
}

func TestBackfillSkippedWhenEngineNotReady(t *testing.T) {
	eng := &fakeEngine{ready: false}
	svc := NewService(eng, &fakeSource{users: []UserRow{{ID: "u1"}}})
	svc.Backfill(context.Background())
	if len(eng.indexed) != 0 {
		t.Fatal("backfill should be skipped when engine not ready")
	}
}
