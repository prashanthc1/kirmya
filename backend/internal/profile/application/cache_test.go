package application

import (
	"context"
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
	if _, err := svc.UpdateScalars(ctx, "u1", domain.Scalars{Headline: "Ops Lead"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// First Get populates the cache (write-through on the write above already did,
	// but a fresh read must still serve correctly).
	p, err := svc.Get(ctx, "u1")
	if err != nil || p.Headline != "Ops Lead" {
		t.Fatalf("first get: %v %+v", err, p)
	}
	if _, ok := cache.store[profileKey("u1")]; !ok {
		t.Fatal("expected cache to be populated for u1")
	}

	// Mutate the repo out-of-band (bypassing the service so the cache is stale).
	if err := repo.UpdateScalars(ctx, "u1", domain.Scalars{Headline: "Changed Underneath"}); err != nil {
		t.Fatalf("oob update: %v", err)
	}

	// Get must still return the cached value, proving the cache is consulted.
	p, err = svc.Get(ctx, "u1")
	if err != nil {
		t.Fatalf("cached get: %v", err)
	}
	if p.Headline != "Ops Lead" {
		t.Fatalf("expected cached headline, got %q", p.Headline)
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
	if _, err := svc.SetSkills(ctx, "u1", []string{"Leadership", "Budgeting"}); err != nil {
		t.Fatalf("set skills: %v", err)
	}

	// A subsequent read sees the fresh skills (served from the refreshed cache).
	p, err := svc.Get(ctx, "u1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(p.Skills) != 2 || p.Skills[0] != "Leadership" {
		t.Fatalf("expected refreshed skills in cache, got %v", p.Skills)
	}
}
