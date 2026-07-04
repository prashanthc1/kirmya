package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// requestWith builds a POST request carrying the given csrf cookie / header
// values (empty string = omit).
func requestWith(t *testing.T, cookie, header string) *http.Request {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: csrfCookieName, Value: cookie})
	}
	if header != "" {
		r.Header.Set("X-CSRF-Token", header)
	}
	return r
}

func TestVerifyDoubleSubmitCSRF(t *testing.T) {
	t.Setenv("CSRF_DOUBLE_SUBMIT", "true")

	cases := []struct {
		name           string
		cookie, header string
		want           bool
	}{
		{"matching token passes", "abc123", "abc123", true},
		{"missing header fails", "abc123", "", false},
		{"missing cookie fails", "", "abc123", false},
		{"mismatched token fails", "abc123", "different", false},
		{"both empty fails", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := verifyDoubleSubmitCSRF(requestWith(t, tc.cookie, tc.header))
			if got != tc.want {
				t.Fatalf("verifyDoubleSubmitCSRF(cookie=%q header=%q) = %v, want %v",
					tc.cookie, tc.header, got, tc.want)
			}
		})
	}
}

func TestVerifyDoubleSubmitCSRF_DisabledAllowsAll(t *testing.T) {
	t.Setenv("CSRF_DOUBLE_SUBMIT", "false")
	if !verifyDoubleSubmitCSRF(requestWith(t, "", "")) {
		t.Fatal("expected check to be bypassed when CSRF_DOUBLE_SUBMIT=false")
	}
}

// Refresh must reject a cookie-auth request with no CSRF token before touching
// the service, so a nil-svc Handler is safe here.
func TestRefresh_RejectsMissingCSRF(t *testing.T) {
	t.Setenv("CSRF_DOUBLE_SUBMIT", "true")
	h := &Handler{}
	w := httptest.NewRecorder()
	h.Refresh(w, requestWith(t, "", ""))
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for missing CSRF token, got %d", w.Code)
	}
}

func TestLogout_RejectsMismatchedCSRF(t *testing.T) {
	t.Setenv("CSRF_DOUBLE_SUBMIT", "true")
	h := &Handler{}
	w := httptest.NewRecorder()
	h.Logout(w, requestWith(t, "cookie-val", "header-val"))
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for mismatched CSRF token, got %d", w.Code)
	}
}
