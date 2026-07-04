package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(http.HandlerFunc(fn)))
	}
	reg("GET /api/v1/conversations", h.List)
	reg("POST /api/v1/conversations", h.Start)
	reg("GET /api/v1/conversations/stream", h.Stream)
	reg("GET /api/v1/conversations/{id}/messages", h.Messages)
	reg("POST /api/v1/conversations/{id}/messages", h.Send)
	reg("POST /api/v1/conversations/{id}/read", h.MarkRead)
	reg("POST /api/v1/conversations/{id}/typing", h.Typing)
}
