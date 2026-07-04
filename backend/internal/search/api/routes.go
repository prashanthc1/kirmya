package api

import "net/http"

// RegisterRoutes mounts search endpoints behind the auth middleware.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/search", auth(http.HandlerFunc(h.Search)))
	mux.Handle("GET /api/v1/search/autocomplete", auth(http.HandlerFunc(h.Autocomplete)))
}
