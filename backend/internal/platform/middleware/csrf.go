package middleware

import (
	"net/http"
	"os"

	"workspace-app/internal/common"
)

// safeMethods never change state, so they are exempt from the Origin check.
var safeMethods = map[string]bool{
	http.MethodGet: true, http.MethodHead: true, http.MethodOptions: true,
}

// VerifyOrigin is an opt-in CSRF defense-in-depth check. The primary CSRF
// defense is already in place: the refresh-token cookie is SameSite=Strict and
// every other endpoint authenticates with a Bearer token (immune to CSRF). This
// adds an Origin allowlist on top, but it is OFF unless CSRF_VERIFY_ORIGIN=true,
// because it requires APP_URL to exactly match the browser origin (it will
// reject legitimate traffic when the app is reached via multiple hosts or a
// proxy that rewrites the origin).
//
// When enabled: state-changing requests that carry an Origin header must match
// APP_URL or the request's own host. Requests without an Origin (curl,
// server-to-server, native apps) are always allowed.
func VerifyOrigin(next http.Handler) http.Handler {
	// Verify Origin is disabled by default, unless explicitly enabled via "true".
	enabled := os.Getenv("CSRF_VERIFY_ORIGIN") == "true"
	allowed := map[string]bool{}
	if app := os.Getenv("APP_URL"); app != "" {
		allowed[app] = true
	} else {
		// Provide secure developer defaults when APP_URL is empty
		allowed["http://localhost:3000"] = true
		allowed["http://127.0.0.1:3000"] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enabled || safeMethods[r.Method] {
			next.ServeHTTP(w, r)
			return
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r) // non-browser client
			return
		}
		if allowed[origin] || origin == "http://"+r.Host || origin == "https://"+r.Host {
			next.ServeHTTP(w, r)
			return
		}
		common.WriteForbiddenError(w, "cross-origin request blocked")
	})
}
