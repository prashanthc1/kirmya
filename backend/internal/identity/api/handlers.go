package api

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/identity/application"
	"workspace-app/internal/identity/domain"
)

const refreshCookieName = "refresh_token"
const refreshCookiePath = "/api/v1/auth"
const csrfCookieName = "csrf_token"

// oauthStateCookieName binds the OAuth `state` issued at OAuthStart to the
// browser that began the flow. The callback must echo the same value (CSRF /
// account-fixation defense). httpOnly so script cannot read or forge it.
const oauthStateCookieName = "oauth_state"

// Handler holds the identity HTTP handlers.
type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

// clientIP resolves the caller's IP. X-Forwarded-For is attacker-controlled
// unless the app sits behind a trusted proxy, so it is only honoured when
// TRUST_PROXY=true; in that case the left-most entry (the original client) is
// used. Otherwise we fall back to the transport peer address (RemoteAddr),
// stripped of its port. This prevents audit-log / rate-limit spoofing (M2).
func clientIP(r *http.Request) string {
	if os.Getenv("TRUST_PROXY") == "true" {
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			if first := strings.TrimSpace(strings.Split(fwd, ",")[0]); first != "" {
				return first
			}
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	if err := req.validate(); err != nil {
		common.WriteValidationError(w, err.Error())
		return
	}
	res, err := h.svc.Register(r.Context(), application.RegisterInput{
		Email: req.Email, Password: req.Password, FullName: req.FullName, Role: req.Role, IP: clientIP(r),
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	// When the verification gate is on, the account is created but no session is
	// issued: the client must route the user to verify before logging in.
	if res.AccessToken == "" {
		common.WriteSuccess(w, http.StatusCreated, map[string]any{
			"user":                  res.User,
			"verification_required": true,
		})
		return
	}
	// Auto-login: set the refresh cookie and return the session (201 Created).
	h.setRefreshCookie(w, r, res.RefreshToken)
	common.WriteSuccess(w, http.StatusCreated, map[string]any{
		"access_token": res.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   res.ExpiresIn,
		"user":         res.User,
	})
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	res, err := h.svc.Login(r.Context(), application.LoginInput{
		Email: req.Email, Password: req.Password, Code: req.Code, IP: clientIP(r), UserAgent: r.UserAgent(),
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if res.MFARequired {
		common.WriteSuccess(w, http.StatusOK, map[string]any{"mfa_required": true})
		return
	}
	h.writeSession(w, r, res)
}

// Refresh handles POST /auth/refresh using the refresh cookie.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if !verifyDoubleSubmitCSRF(r) {
		common.WriteForbiddenError(w, "missing or invalid CSRF token")
		return
	}
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil {
		common.WriteUnauthorizedError(w, "missing refresh token")
		return
	}
	res, err := h.svc.Refresh(r.Context(), cookie.Value)
	if err != nil {
		h.clearRefreshCookie(w, r)
		h.writeError(w, err)
		return
	}
	h.writeSession(w, r, res)
}

// Logout handles POST /auth/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !verifyDoubleSubmitCSRF(r) {
		common.WriteForbiddenError(w, "missing or invalid CSRF token")
		return
	}
	if cookie, err := r.Cookie(refreshCookieName); err == nil {
		_ = h.svc.Logout(r.Context(), cookie.Value)
	}
	h.clearRefreshCookie(w, r)
	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

// VerifyEmail handles POST /auth/verify-email.
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if !decode(r, &req) || req.Token == "" {
		common.WriteValidationError(w, "token is required")
		return
	}
	if err := h.svc.VerifyEmail(r.Context(), req.Token); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"verified": true})
}

// ResendVerification handles POST /auth/resend-verification.
func (h *Handler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req emailRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	if err := h.svc.ResendVerification(r.Context(), req.Email); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"sent": true})
}

// ForgotPassword handles POST /auth/forgot-password.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req emailRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	if err := h.svc.ForgotPassword(r.Context(), req.Email); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"sent": true})
}

// ResetPassword handles POST /auth/reset-password.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	if len(req.Password) < 8 {
		common.WriteValidationError(w, "password must be at least 8 characters")
		return
	}
	if err := h.svc.ResetPassword(r.Context(), req.Token, req.Password); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"reset": true})
}

// OAuthStart handles GET /auth/oauth/{provider}. It always mints a fresh,
// unpredictable state value, stores it in an httpOnly cookie bound to this
// browser, and returns it. The callback must present the same value.
func (h *Handler) OAuthStart(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	state := randomState()
	url, err := h.svc.OAuthAuthURL(provider, state)
	if err != nil {
		h.writeError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   secureCookies(r),
		SameSite: http.SameSiteLaxMode, // Lax: the provider redirect is a top-level GET.
		MaxAge:   600,                  // 10 minutes to complete the flow.
	})
	common.WriteSuccess(w, http.StatusOK, map[string]string{"url": url, "state": state})
}

// OAuthCallback handles POST /auth/oauth/{provider}/callback with the code. The
// request must echo the `state` issued by OAuthStart; it is compared in constant
// time against the httpOnly cookie to defeat login-CSRF / account fixation (H1).
func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	var req oauthCallbackRequest
	if !decode(r, &req) || req.Code == "" {
		common.WriteValidationError(w, "authorization code is required")
		return
	}
	cookie, err := r.Cookie(oauthStateCookieName)
	if err != nil || cookie.Value == "" || req.State == "" ||
		subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(req.State)) != 1 {
		h.clearOAuthStateCookie(w, r)
		common.WriteValidationError(w, "invalid or missing oauth state")
		return
	}
	h.clearOAuthStateCookie(w, r) // single-use
	res, err := h.svc.OAuthCallback(r.Context(), provider, req.Code, clientIP(r))
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeSession(w, r, res)
}

// MFASetup handles POST /auth/mfa/setup (authenticated).
func (h *Handler) MFASetup(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	url, err := h.svc.SetupMFA(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"otpauth_url": url})
}

// MFAVerify handles POST /auth/mfa/verify (authenticated) to confirm enrollment.
func (h *Handler) MFAVerify(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	var req mfaVerifyRequest
	if !decode(r, &req) || req.Code == "" {
		common.WriteValidationError(w, "code is required")
		return
	}
	if err := h.svc.ConfirmMFA(r.Context(), userID, req.Code); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"mfa_enabled": true})
}

// MFADisable handles POST /auth/mfa/disable (authenticated): turns off TOTP
// after validating a current code.
func (h *Handler) MFADisable(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	var req mfaVerifyRequest
	if !decode(r, &req) || req.Code == "" {
		common.WriteValidationError(w, "code is required")
		return
	}
	if err := h.svc.DisableMFA(r.Context(), userID, req.Code); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"mfa_enabled": false})
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword handles POST /auth/change-password (authenticated). It verifies
// the current password, enforces new-password complexity, and revokes all
// sessions so other devices must re-authenticate.
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	var req changePasswordRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	if req.CurrentPassword == "" {
		common.WriteValidationError(w, "current_password is required")
		return
	}
	if err := common.ValidatePassword(req.NewPassword); err != nil {
		if ae, ok := err.(*common.AppError); ok {
			common.WriteError(w, ae)
			return
		}
		common.WriteValidationError(w, "new password is invalid")
		return
	}
	if req.NewPassword == req.CurrentPassword {
		common.WriteValidationError(w, "new password must be different from the current password")
		return
	}
	if err := h.svc.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"changed": true})
}

// LogoutAll handles POST /auth/logout-all (authenticated): revoke every session.
// The current refresh cookie is left in place but its token is now revoked, so
// the next refresh fails and the client is effectively signed out everywhere.
func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.LogoutAll(r.Context(), common.UserIDFromContext(r.Context())); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"signed_out": true})
}

// DeactivateAccount handles DELETE /users/me (authenticated): soft-close the
// caller's account and revoke all sessions.
func (h *Handler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DeactivateAccount(r.Context(), common.UserIDFromContext(r.Context())); err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"deactivated": true})
}

// directoryEntryDTO is the public projection of a user returned by the
// directory/search endpoints. Email is intentionally omitted: exposing every
// member's email to any authenticated caller enables enumeration/harvesting
// (PII over-exposure, M4). The owner's own email is available via /users/me.
type directoryEntryDTO struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Headline string `json:"headline,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
}

// SearchUsers handles GET /users/search?q= (authenticated).
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.SearchUsers(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	out := make([]directoryEntryDTO, 0, len(entries))
	for _, e := range entries {
		out = append(out, directoryEntryDTO{ID: e.ID, FullName: e.FullName, Headline: e.Headline, PhotoURL: e.PhotoURL})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"users": out})
}

// GetUser handles GET /users/{id} (authenticated) — public directory entry.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	e, err := h.svc.Directory(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, directoryEntryDTO{
		ID: e.ID, FullName: e.FullName, Headline: e.Headline, PhotoURL: e.PhotoURL,
	})
}

// Me handles GET /users/me (authenticated).
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.svc.Me(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, user)
}

type updateRolesRequest struct {
	Roles []string `json:"roles"`
}

// UpdateMyRoles handles PUT /users/me/roles (authenticated): reconciles the
// caller's self-assignable roles (job_seeker, referrer, mentor, recruiter) to
// the provided set. Admin is not self-assignable and an existing admin role is
// preserved.
func (h *Handler) UpdateMyRoles(w http.ResponseWriter, r *http.Request) {
	var req updateRolesRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	user, err := h.svc.SetMyRoles(r.Context(), common.UserIDFromContext(r.Context()), req.Roles)
	if err != nil {
		h.writeError(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, user)
}

// CSRF handles GET /auth/csrf: issues a double-submit CSRF token.
func (h *Handler) CSRF(w http.ResponseWriter, r *http.Request) {
	token := randomState()
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   secureCookies(r),
		SameSite: http.SameSiteStrictMode,
	})
	common.WriteSuccess(w, http.StatusOK, map[string]string{"csrf_token": token})
}

// csrfDoubleSubmitEnabled reports whether the double-submit CSRF check is
// active. It is OPT-IN (CSRF_DOUBLE_SUBMIT=true), mirroring the conservative
// default of the Origin check (VerifyOrigin): turning it on requires every
// cookie-auth client to first fetch GET /auth/csrf and echo the token, so it is
// enabled explicitly per environment (it is on in docker-compose.prod.yml)
// rather than globally by default, to avoid rejecting clients that predate it.
func csrfDoubleSubmitEnabled() bool {
	return os.Getenv("CSRF_DOUBLE_SUBMIT") == "true"
}

// verifyDoubleSubmitCSRF enforces the double-submit cookie pattern on
// cookie-authenticated, state-changing endpoints (refresh, logout). The
// X-CSRF-Token header must be present and equal the non-httpOnly csrf_token
// cookie issued by GET /auth/csrf. A cross-site attacker can ride along the
// victim's cookies but cannot read csrf_token (to copy into the header), so a
// match proves the request came from our own first-party JavaScript. The
// compare is constant-time. Returns true when the request may proceed.
func verifyDoubleSubmitCSRF(r *http.Request) bool {
	if !csrfDoubleSubmitEnabled() {
		return true
	}
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	header := r.Header.Get("X-CSRF-Token")
	if header == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(header)) == 1
}

// --- helpers ---

func (h *Handler) writeSession(w http.ResponseWriter, r *http.Request, res application.AuthResult) {
	h.setRefreshCookie(w, r, res.RefreshToken)
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"access_token": res.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   res.ExpiresIn,
		"user":         res.User,
	})
}

func (h *Handler) setRefreshCookie(w http.ResponseWriter, r *http.Request, raw string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    raw,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   secureCookies(r),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		MaxAge:   30 * 24 * 3600, // 30 days in seconds
	})
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter, r *http.Request) {
	// MaxAge=-1 instructs the browser to delete the cookie immediately (RFC 6265).
	// Expires is also set to the Unix epoch so that older User-Agents that honour
	// Expires over MaxAge also delete it.  Both attributes must match the original
	// cookie's Path/Domain/Secure/HttpOnly/SameSite so the browser treats this
	// Set-Cookie as targeting the same cookie jar entry.
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   secureCookies(r),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func (h *Handler) clearOAuthStateCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   secureCookies(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// secureCookies reports whether the Secure attribute should be set. It is true
// in production and also whenever the request itself arrived over HTTPS
// (directly or via a TLS-terminating proxy that sets X-Forwarded-Proto), so
// cookies are not silently downgraded to plaintext on HTTPS deployments (M1).
func secureCookies(r *http.Request) bool {
	if os.Getenv("APP_ENV") == "production" {
		return true
	}
	if r != nil {
		if r.TLS != nil {
			return true
		}
		if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
			return true
		}
	}
	return false
}

func randomState() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// writeError maps application/domain errors to HTTP responses.
func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrInvalidCredentials),
		errors.Is(err, application.ErrInvalidMFACode):
		common.WriteUnauthorizedError(w, err.Error())
	case errors.Is(err, application.ErrEmailNotVerified):
		// Distinct code so the client can offer a "resend verification" action.
		common.WriteError(w, &common.AppError{Code: "email_not_verified", Message: err.Error(), Status: http.StatusForbidden})
	case errors.Is(err, application.ErrAccountInactive):
		common.WriteForbiddenError(w, err.Error())
	case errors.Is(err, application.ErrUnknownProvider),
		errors.Is(err, application.ErrInvalidRole),
		errors.Is(err, application.ErrNoRoles):
		common.WriteValidationError(w, err.Error())
	case errors.Is(err, domain.ErrEmailTaken):
		common.WriteError(w, common.NewConflictError("email already registered"))
	case errors.Is(err, domain.ErrUserNotFound):
		common.WriteNotFoundError(w, "user not found")
	case errors.Is(err, domain.ErrOptimisticLock):
		common.WriteError(w, common.NewConflictError("the record was modified; please retry"))
	default:
		log.Printf("identity: unhandled error: %v", err)
		common.WriteInternalError(w, "something went wrong")
	}
}
