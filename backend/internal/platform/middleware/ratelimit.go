package middleware

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"workspace-app/internal/common"
)

// RateLimiter is a per-client-IP token-bucket limiter. The real client IP is
// taken from X-Forwarded-For (set by the frontend proxy / ingress) and falls
// back to RemoteAddr. Health and metrics endpoints are never limited.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	rate    rate.Limit
	burst   int
	ttl     time.Duration
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter reads RATE_LIMIT_RPS (default 50) and RATE_LIMIT_BURST
// (default 100) and starts a background sweeper for idle clients.
func NewRateLimiter() *RateLimiter {
	rps := envFloat("RATE_LIMIT_RPS", 50)
	burst := int(envFloat("RATE_LIMIT_BURST", 100))
	rl := &RateLimiter{
		clients: map[string]*client{},
		rate:    rate.Limit(rps),
		burst:   burst,
		ttl:     10 * time.Minute,
	}
	go rl.sweep()
	return rl
}

// Middleware enforces the limit, returning 429 with Retry-After when exceeded.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" || r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}
		if !rl.limiter(clientIP(r)).Allow() {
			w.Header().Set("Retry-After", "1")
			common.WriteError(w, &common.AppError{
				Code:    "rate_limited",
				Message: "too many requests, please slow down",
				Status:  http.StatusTooManyRequests,
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) limiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	c, ok := rl.clients[ip]
	if !ok {
		c = &client{limiter: rate.NewLimiter(rl.rate, rl.burst)}
		rl.clients[ip] = c
	}
	c.lastSeen = time.Now()
	return c.limiter
}

func (rl *RateLimiter) sweep() {
	ticker := time.NewTicker(rl.ttl)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > rl.ttl {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP prefers the first X-Forwarded-For hop, then RemoteAddr.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			return f
		}
	}
	return def
}
