// Package cache provides the platform's cache-aside layer. It exposes a small,
// error-free Cache port (failures are swallowed and logged so caching can never
// break a request) with two implementations: a Redis-backed cache and a no-op
// used when Redis is not configured or unreachable. Modules depend on a
// structurally-identical interface declared in their own application package, so
// they stay decoupled from this package.
package cache

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"workspace-app/internal/platform/observability"
)

// Cache is the platform cache port. Get reports whether the key was present.
// Set and Delete are best-effort: errors are logged, not returned.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	Delete(ctx context.Context, keys ...string)
}

// New returns a Redis-backed cache when REDIS_URL (or REDIS_ADDR) is configured
// and reachable; otherwise it returns a thread-safe in-memory cache so the
// platform has working cache-aside even without Redis. It never returns an error.
func New() Cache {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		if addr := os.Getenv("REDIS_ADDR"); addr != "" {
			url = "redis://" + addr
		}
	}
	if url == "" {
		log.Printf("[cache] REDIS_URL not set; falling back to in-memory cache")
		return NewMemory()
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Printf("[cache] invalid REDIS_URL (%v); falling back to in-memory cache", err)
		return NewMemory()
	}
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[cache] Redis unreachable (%v); falling back to in-memory cache", err)
		_ = client.Close()
		return NewMemory()
	}

	log.Printf("[cache] Redis connected; cache-aside enabled")
	return &Redis{client: client}
}

// Redis is the Redis-backed Cache implementation.
type Redis struct{ client *redis.Client }

func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool) {
	b, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		observability.RecordCacheMiss()
		return nil, false
	}
	if err != nil {
		log.Printf("[cache] get %s: %v", key, err)
		observability.RecordCacheMiss()
		return nil, false
	}
	observability.RecordCacheHit()
	return b, true
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		log.Printf("[cache] set %s: %v", key, err)
	}
}

func (r *Redis) Delete(ctx context.Context, keys ...string) {
	if len(keys) == 0 {
		return
	}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		log.Printf("[cache] del %v: %v", keys, err)
	}
}

// Close releases the underlying client (Noop/Memory has nothing to close).
func (r *Redis) Close() error { return r.client.Close() }

// Memory is a thread-safe, size-bounded in-memory Cache implementation.
type Memory struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

type memoryItem struct {
	value     []byte
	expiresAt time.Time
}

// NewMemory returns a new initialized Memory cache.
func NewMemory() *Memory {
	return &Memory{
		items: make(map[string]memoryItem),
	}
}

func (m *Memory) Get(ctx context.Context, key string) ([]byte, bool) {
	m.mu.RLock()
	item, ok := m.items[key]
	m.mu.RUnlock()

	if !ok {
		observability.RecordCacheMiss()
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		m.mu.Lock()
		delete(m.items, key)
		m.mu.Unlock()
		observability.RecordCacheMiss()
		return nil, false
	}
	observability.RecordCacheHit()
	return item.value, true
}

func (m *Memory) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Evict expired items if size grows past threshold to prevent leaks
	if len(m.items) > 1000 {
		now := time.Now()
		for k, item := range m.items {
			if now.After(item.expiresAt) {
				delete(m.items, k)
			}
		}
	}

	m.items[key] = memoryItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (m *Memory) Delete(ctx context.Context, keys ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, key := range keys {
		delete(m.items, key)
	}
}

// Noop is kept for backward compatibility (e.g. if code checks for type Noop).
type Noop struct{}

func (Noop) Get(context.Context, string) ([]byte, bool)         { return nil, false }
func (Noop) Set(context.Context, string, []byte, time.Duration) {}
func (Noop) Delete(context.Context, ...string)                  {}
