package api

import "net/http"

// RegisterRoutes mounts the career endpoints. The ladder is reference data, but
// the route still requires authentication to match the rest of the app surface.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/career/paths", auth(http.HandlerFunc(h.Paths)))
}
