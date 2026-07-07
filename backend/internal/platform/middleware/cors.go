package middleware

import (
	"net/http"
	"os"
)

// CORS is a secure-by-default Cross-Origin Resource Sharing middleware.
// It matches incoming Origin headers against configured or safe defaults (such
// as APP_URL, localhost:3000, 127.0.0.1:3000, and the host origin itself)
// instead of using a wildcard "*", enabling secure credential handling (cookies).
func CORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{}
	if app := os.Getenv("APP_URL"); app != "" {
		allowedOrigins[app] = true
	}
	// Fallback/standard developer origins
	allowedOrigins["http://localhost:3000"] = true
	allowedOrigins["http://127.0.0.1:3000"] = true
	allowedOrigins["http://[::1]:3000"] = true
	allowedOrigins["https://localhost:3000"] = true
	allowedOrigins["https://127.0.0.1:3000"] = true
	allowedOrigins["https://[::1]:3000"] = true

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Origin validation: must match allowlist, request Host, or local scheme
			isAllowed := allowedOrigins[origin] ||
				origin == "http://"+r.Host ||
				origin == "https://"+r.Host

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours cache limit
			}
		}

		// Preflight OPTIONS requests: respond immediately with status 204
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
