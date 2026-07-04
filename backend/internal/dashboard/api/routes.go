package api

import "net/http"

// RegisterRoutes mounts the dashboard endpoints. All routes require an
// authenticated caller; the summary is always scoped to that caller.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/me/dashboard", auth(http.HandlerFunc(h.Summary)))
}
