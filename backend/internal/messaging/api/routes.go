package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler, limit func(string) func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}
	regLimit := func(pattern string, action string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(limit(action)(fn)))
	}

	reg("GET /api/v1/conversations", h.List)
	reg("POST /api/v1/conversations", h.Start)
	reg("GET /api/v1/conversations/stream", h.Stream)
	reg("GET /api/v1/conversations/{id}/messages", h.Messages)
	regLimit("POST /api/v1/conversations/{id}/messages", "send_message", h.Send)
	reg("DELETE /api/v1/conversations/{id}/messages/{messageID}", h.DeleteMessage)
	reg("POST /api/v1/conversations/{id}/read", h.MarkRead)
	reg("POST /api/v1/conversations/{id}/typing", h.Typing)
	reg("POST /api/v1/conversations/{id}/archive", h.Archive)
	reg("POST /api/v1/conversations/{id}/pin", h.Pin)

	// WebSocket endpoint - handles auth internally via token query param
	mux.Handle("GET /api/v1/ws", limit("websocket_connect")(http.HandlerFunc(h.WebSocket)))
}
