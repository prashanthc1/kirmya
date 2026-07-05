package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
}

func TestSecurityHeaders(t *testing.T) {
	h := SecurityHeaders(okHandler())

	// API path → locked-down CSP.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil))
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("nosniff missing, got %q", got)
	}
	if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("frame-options, got %q", got)
	}
	if got := rec.Header().Get("Content-Security-Policy"); got != apiCSP {
		t.Fatalf("api CSP, got %q", got)
	}

	// Swagger path → relaxed CSP allowing the CDN.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger-ui/", nil))
	if got := rec.Header().Get("Content-Security-Policy"); got != swaggerCSP {
		t.Fatalf("swagger CSP, got %q", got)
	}
}

func TestRateLimiter(t *testing.T) {
	os.Setenv("RATE_LIMIT_RPS", "1")
	os.Setenv("RATE_LIMIT_BURST", "1")
	defer func() { os.Unsetenv("RATE_LIMIT_RPS"); os.Unsetenv("RATE_LIMIT_BURST") }()

	h := NewRateLimiter().Middleware(okHandler())
	req := func() int {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
		r.RemoteAddr = "203.0.113.7:1234"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}
	if code := req(); code != http.StatusOK {
		t.Fatalf("first request should pass, got %d", code)
	}
	if code := req(); code != http.StatusTooManyRequests {
		t.Fatalf("second request should be limited, got %d", code)
	}

	// Health is never limited.
	r := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	r.RemoteAddr = "203.0.113.7:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	if rec.Code != http.StatusOK {
		t.Fatalf("health must not be rate limited, got %d", rec.Code)
	}
}

func TestVerifyOrigin(t *testing.T) {
	os.Setenv("APP_URL", "https://app.example.com")
	os.Setenv("CSRF_VERIFY_ORIGIN", "true")
	defer func() { os.Unsetenv("APP_URL"); os.Unsetenv("CSRF_VERIFY_ORIGIN") }()
	h := VerifyOrigin(okHandler())

	send := func(method, origin string) int {
		r := httptest.NewRequest(method, "/api/v1/auth/refresh", nil)
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	if code := send(http.MethodPost, "https://evil.example.com"); code != http.StatusForbidden {
		t.Fatalf("cross-origin POST should be blocked, got %d", code)
	}
	if code := send(http.MethodPost, "https://app.example.com"); code != http.StatusOK {
		t.Fatalf("allowed origin should pass, got %d", code)
	}
	if code := send(http.MethodPost, ""); code != http.StatusOK {
		t.Fatalf("no-Origin (non-browser) should pass, got %d", code)
	}
	if code := send(http.MethodGet, "https://evil.example.com"); code != http.StatusOK {
		t.Fatalf("safe method should not be checked, got %d", code)
	}
}

func TestVerifyOriginExplicitlyDisabled(t *testing.T) {
	os.Setenv("CSRF_VERIFY_ORIGIN", "false")
	defer os.Unsetenv("CSRF_VERIFY_ORIGIN")
	h := VerifyOrigin(okHandler())
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	r.Header.Set("Origin", "https://evil.example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	if rec.Code != http.StatusOK {
		t.Fatalf("disabled check must pass everything, got %d", rec.Code)
	}
}

func TestVerifyOriginDefaultBehavior(t *testing.T) {
	os.Unsetenv("CSRF_VERIFY_ORIGIN")
	os.Unsetenv("APP_URL")
	h := VerifyOrigin(okHandler())

	send := func(origin string) int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		r.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	if code := send("https://evil.example.com"); code != http.StatusForbidden {
		t.Fatalf("expected forbidden for evil origin by default, got %d", code)
	}
	if code := send("http://localhost:3000"); code != http.StatusOK {
		t.Fatalf("expected OK for default developer origin, got %d", code)
	}
}

func TestCORS(t *testing.T) {
	os.Setenv("APP_URL", "https://prod.kirmya.app")
	defer os.Unsetenv("APP_URL")

	h := CORS(okHandler())

	// Preflight request
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/jobs", nil)
	req.Header.Set("Origin", "https://prod.kirmya.app")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 NoContent for preflight OPTIONS request, got %d", rec.Code)
	}
	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "https://prod.kirmya.app" {
		t.Fatalf("expected Access-Control-Allow-Origin header to match trusted origin, got %s", origin)
	}

	// Normal request with trusted origin
	req = httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Fatalf("expected Access-Control-Allow-Origin header to match dev server origin, got %s", origin)
	}

	// Normal request with untrusted origin
	req = httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	req.Header.Set("Origin", "https://evil.attacker.com")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK (CORS does not block requests by itself, only headers), got %d", rec.Code)
	}
	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "" {
		t.Fatalf("untrusted origin should not receive Access-Control-Allow-Origin, got %s", origin)
	}
}
