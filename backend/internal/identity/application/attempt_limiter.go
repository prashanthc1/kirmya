package application

import (
	"sync"
	"time"
)

// attemptLimiter is a small in-memory fixed-window counter that caps the number
// of failed attempts per key within a window. It is used to throttle TOTP/MFA
// code submissions during login (L3), defeating brute force of the 6-digit code.
//
// Scope/limitation: state is per-process. In a multi-instance deployment a
// distributed/shared store (e.g. Redis) should back this; the interface stays
// the same. It is intentionally dependency-free so the application layer remains
// pure and unit-testable.
type attemptLimiter struct {
	max    int
	window time.Duration

	mu      sync.Mutex
	buckets map[string]*attemptBucket
}

type attemptBucket struct {
	count    int
	resetsAt time.Time
}

func newAttemptLimiter(max int, window time.Duration) *attemptLimiter {
	return &attemptLimiter{max: max, window: window, buckets: make(map[string]*attemptBucket)}
}

// allow records an attempt for key and reports whether it is within the limit.
// Returns false once the cap for the current window is exceeded.
func (l *attemptLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.resetsAt) {
		l.buckets[key] = &attemptBucket{count: 1, resetsAt: now.Add(l.window)}
		return true
	}
	if b.count >= l.max {
		return false
	}
	b.count++
	return true
}

// reset clears the counter for key (e.g. after a successful authentication).
func (l *attemptLimiter) reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}
