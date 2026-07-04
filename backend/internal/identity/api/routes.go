package api

import "net/http"

// RegisterRoutes mounts the identity endpoints on the platform mux.
func RegisterRoutes(mux *http.ServeMux, h *Handler, mw *Middleware) {
	// Public auth endpoints.
	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", h.Logout)
	mux.HandleFunc("POST /api/v1/auth/verify-email", h.VerifyEmail)
	mux.HandleFunc("POST /api/v1/auth/resend-verification", h.ResendVerification)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", h.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", h.ResetPassword)
	mux.HandleFunc("GET /api/v1/auth/csrf", h.CSRF)

	// OAuth.
	mux.HandleFunc("GET /api/v1/auth/oauth/{provider}", h.OAuthStart)
	mux.HandleFunc("POST /api/v1/auth/oauth/{provider}/callback", h.OAuthCallback)

	// Authenticated endpoints.
	mux.Handle("POST /api/v1/auth/mfa/setup", mw.Authenticate(http.HandlerFunc(h.MFASetup)))
	mux.Handle("POST /api/v1/auth/mfa/verify", mw.Authenticate(http.HandlerFunc(h.MFAVerify)))
	mux.Handle("POST /api/v1/auth/mfa/disable", mw.Authenticate(http.HandlerFunc(h.MFADisable)))
	mux.Handle("POST /api/v1/auth/change-password", mw.Authenticate(http.HandlerFunc(h.ChangePassword)))
	mux.Handle("POST /api/v1/auth/logout-all", mw.Authenticate(http.HandlerFunc(h.LogoutAll)))
	mux.Handle("DELETE /api/v1/users/me", mw.Authenticate(http.HandlerFunc(h.DeactivateAccount)))
	mux.Handle("GET /api/v1/users/me", mw.Authenticate(http.HandlerFunc(h.Me)))
	mux.Handle("PUT /api/v1/users/me/roles", mw.Authenticate(http.HandlerFunc(h.UpdateMyRoles)))
	mux.Handle("GET /api/v1/users/search", mw.Authenticate(http.HandlerFunc(h.SearchUsers)))
	mux.Handle("GET /api/v1/users/{id}", mw.Authenticate(http.HandlerFunc(h.GetUser)))
}
