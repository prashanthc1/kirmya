package application

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"workspace-app/internal/profile/domain"
)

// memCache is an in-memory Cache for tests (ignores TTL).
type memCache struct {
	mu      sync.Mutex
	store   map[string][]byte
	hits    int
	sets    int
	deletes int
}

func newMemCache() *memCache { return &memCache{store: map[string][]byte{}} }

func (c *memCache) Get(_ context.Context, key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	b, ok := c.store[key]
	if ok {
		c.hits++
	}
	return b, ok
}

func (c *memCache) Set(_ context.Context, key string, value []byte, _ time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
	c.sets++
}

func (c *memCache) Delete(_ context.Context, keys ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.store, k)
	}
	c.deletes++
}

func TestGetIsCacheAside(t *testing.T) {
	repo := newFakeRepo()
	cache := newMemCache()
	svc := NewService(repo, nil, cache)
	ctx := context.Background()

	// Seed the repo with a known headline.
	isDraft := true
	if _, err := svc.UpdateProfile(ctx, "u1", 0, domain.AggregateUpdate{
		Identity: &domain.IdentitySection{Headline: "Ops Lead"},
		IsDraft:  &isDraft,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// First Get populates the cache
	p, err := svc.Get(ctx, "u1")
	if err != nil || p.Identity.Headline != "Ops Lead" {
		t.Fatalf("first get: %v %+v", err, p)
	}
	if _, ok := cache.store[profileKey("u1")]; !ok {
		t.Fatal("expected cache to be populated for u1")
	}

	// Mutate the repo out-of-band (bypassing the service so the cache is stale).
	err = repo.UpdateAggregate(ctx, "u1", 0, domain.AggregateUpdate{
		Identity: &domain.IdentitySection{Headline: "Changed Underneath"},
	})
	if err != nil {
		t.Fatalf("oob update: %v", err)
	}

	// Get must still return the cached value, proving the cache is consulted.
	p, err = svc.Get(ctx, "u1")
	if err != nil {
		t.Fatalf("cached get: %v", err)
	}
	if p.Identity.Headline != "Ops Lead" {
		t.Fatalf("expected cached headline, got %q", p.Identity.Headline)
	}
	if cache.hits == 0 {
		t.Fatal("expected at least one cache hit")
	}
}

func TestWriteRefreshesCache(t *testing.T) {
	repo := newFakeRepo()
	cache := newMemCache()
	svc := NewService(repo, nil, cache)
	ctx := context.Background()

	// Prime the cache.
	if _, err := svc.Get(ctx, "u1"); err != nil {
		t.Fatalf("prime: %v", err)
	}
	// A write through the service refreshes the cache write-through.
	if _, err := svc.SetSkills(ctx, "u1", []domain.SkillItem{{Name: "Leadership"}, {Name: "Budgeting"}}); err != nil {
		t.Fatalf("set skills: %v", err)
	}

	// Verify new values are immediately visible via the cache.
	b, ok := cache.Get(ctx, profileKey("u1"))
	if !ok {
		t.Fatal("expected cached profile after write")
	}
	var cached domain.Profile
	if err := json.Unmarshal(b, &cached); err != nil {
		t.Fatalf("unmarshal cached: %v", err)
	}
	if len(cached.Skills) != 2 || cached.Skills[0].Name != "Leadership" {
		t.Fatalf("expected refreshed cached skills, got %+v", cached.Skills)
	}
}
