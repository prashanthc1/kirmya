// Package middleware holds the platform's cross-cutting HTTP middleware:
// security headers, rate limiting, and CSRF/Origin verification. These wrap the
// router to implement the OWASP baseline controls.
package middleware

import (
	"net/http"
	"strings"
)

// apiCSP is a locked-down policy for JSON/metrics responses (no resources are
// ever loaded from them). swaggerCSP relaxes just enough for the bundled
// Swagger UI, which pulls its assets from the jsDelivr CDN.
const (
	apiCSP = "default-src 'none'; frame-ancestors 'none'"

	swaggerCSP = "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
		"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
		"img-src 'self' data: https://cdn.jsdelivr.net; " +
		"connect-src 'self'; frame-ancestors 'none'"
)

// SecurityHeaders applies the OWASP secure-response-header baseline. The
// Content-Security-Policy is path-aware so the API stays locked down while the
// Swagger UI page keeps working.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// Honoured only over HTTPS; harmless over plain HTTP.
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/metrics" {
			h.Set("Content-Security-Policy", apiCSP)
		} else {
			h.Set("Content-Security-Policy", swaggerCSP)
		}

		next.ServeHTTP(w, r)
	})
}
