package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"workspace-app/internal/common"
)

type RedisRateLimiter struct {
	rdb *redis.Client
	db  *sql.DB
}

func NewRedisRateLimiter(db *sql.DB) *RedisRateLimiter {
	var rdb *redis.Client
	url := os.Getenv("REDIS_URL")
	if url == "" {
		if addr := os.Getenv("REDIS_ADDR"); addr != "" {
			url = "redis://" + addr
		}
	}
	if url != "" {
		opts, err := redis.ParseURL(url)
		if err == nil {
			rdb = redis.NewClient(opts)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err != nil {
				rdb = nil
			}
		}
	}
	return &RedisRateLimiter{rdb: rdb, db: db}
}

// Limit checks the token bucket in Redis and rate limits the request.
func (rl *RedisRateLimiter) Limit(action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl.rdb == nil {
				next.ServeHTTP(w, r)
				return
			}

			uid := common.UserIDFromContext(r.Context())
			if uid == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Default limits (Low trust / unverified email)
			rate := 0.05 // 1 token per 20 seconds (3 per minute)
			capacity := 3.0

			// Query user trust level
			var emailVerified bool
			var role string
			err := rl.db.QueryRowContext(r.Context(), `
				SELECT email_verified, COALESCE((
					SELECT r.role_name 
					FROM user_roles ur 
					JOIN roles r ON ur.role_id = r.id 
					WHERE ur.user_id = users.id 
					LIMIT 1
				), 'job_seeker')
				FROM users WHERE id = $1
			`, uid).Scan(&emailVerified, &role)

			if err == nil {
				if emailVerified {
					if action == "send_message" {
						rate = 0.5 // 30 messages per minute
						capacity = 30.0
					} else { // connection requests
						rate = 0.16 // 10 requests per minute
						capacity = 10.0
					}
				} else {
					if action == "send_message" {
						rate = 0.05 // 3 messages per minute
						capacity = 3.0
					} else {
						rate = 0.03 // 2 requests per minute
						capacity = 2.0
					}
				}

				if role == "admin" || role == "recruiter" {
					rate = 5.0
					capacity = 50.0
				}
			}

			key := fmt.Sprintf("kirmya:ratelimit:%s:%s", uid, action)
			now := time.Now().UnixNano()

			// Lua Script to atomic check and update token bucket
			script := `
				local key = KEYS[1]
				local capacity = tonumber(ARGV[1])
				local rate = tonumber(ARGV[2])
				local now = tonumber(ARGV[3])
				local cost = tonumber(ARGV[4])

				local state = redis.call('HMGET', key, 'tokens', 'last_update')
				local tokens = tonumber(state[1])
				local last_update = tonumber(state[2])

				if not tokens then
					tokens = capacity
					last_update = now
				else
					local elapsed = (now - last_update) / 1e9
					tokens = math.min(capacity, tokens + elapsed * rate)
					last_update = now
				end

				if tokens >= cost then
					tokens = tokens - cost
					redis.call('HMSET', key, 'tokens', tokens, 'last_update', last_update)
					redis.call('EXPIRE', key, 86400)
					return 1
				else
					redis.call('HMSET', key, 'tokens', tokens, 'last_update', last_update)
					return 0
				end
			`

			res, err := rl.rdb.Eval(r.Context(), script, []string{key}, capacity, rate, now, 1).Int()
			if err != nil {
				log.Printf("[rate-limiter] Redis Eval error: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			if res == 0 {
				w.Header().Set("Retry-After", "10")
				common.WriteError(w, &common.AppError{
					Code:    "rate_limited",
					Message: fmt.Sprintf("too many requests for %s, please slow down", strings.ReplaceAll(action, "_", " ")),
					Status:  http.StatusTooManyRequests,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
