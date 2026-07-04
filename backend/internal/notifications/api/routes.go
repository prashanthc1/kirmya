package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(http.HandlerFunc(fn)))
	}
	reg("GET /api/v1/notifications", h.List)
	reg("GET /api/v1/notifications/stream", h.Stream)
	reg("POST /api/v1/notifications/read-all", h.MarkAllRead)
	reg("POST /api/v1/notifications/{id}/read", h.MarkRead)
}
