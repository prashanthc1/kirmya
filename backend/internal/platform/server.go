package platform

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"workspace-app/internal/platform/middleware"
	"workspace-app/internal/platform/observability"
)

// NewServer creates the HTTP server for the platform. The handler chain is:
// logging → security headers → rate limit → CSRF/Origin → OTel spans →
// Prometheus metrics → router.
func NewServer(port string, db *sql.DB) *http.Server {
	rateLimiter := middleware.NewRateLimiter()
	handler := loggingMiddleware(
		middleware.SecurityHeaders(
			rateLimiter.Middleware(
				middleware.VerifyOrigin(
					observability.WrapHandler(
						observability.MetricsMiddleware(NewRouter(db)),
					),
				),
			),
		),
	)
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
