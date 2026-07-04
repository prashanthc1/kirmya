package api

import "net/http"

// RegisterRoutes mounts admin endpoints. The /admin/* routes are gated by
// adminOnly (RBAC admin role); the public report-filing endpoint is available to
// any authenticated user via auth.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth, adminOnly func(http.Handler) http.Handler) {
	admin := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, adminOnly(fn))
	}

	// Analytics.
	admin("GET /api/v1/admin/stats", h.Stats)

	// User management.
	admin("GET /api/v1/admin/users", h.ListUsers)
	admin("GET /api/v1/admin/users/{id}", h.GetUser)
	admin("PATCH /api/v1/admin/users/{id}/status", h.SetUserStatus)
	admin("POST /api/v1/admin/users/{id}/roles", h.AssignRole)
	admin("DELETE /api/v1/admin/users/{id}/roles/{role}", h.RevokeRole)

	// Content moderation.
	admin("DELETE /api/v1/admin/posts/{id}", h.RemovePost)
	admin("DELETE /api/v1/admin/comments/{id}", h.RemoveComment)

	// Report queue (admin triage).
	admin("GET /api/v1/admin/reports", h.ListReports)
	admin("PATCH /api/v1/admin/reports/{id}", h.ResolveReport)

	// File a report — open to any authenticated user.
	mux.Handle("POST /api/v1/reports", auth(http.HandlerFunc(h.FileReport)))
}
