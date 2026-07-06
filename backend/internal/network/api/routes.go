package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler, limit func(string) func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}
	regLimit := func(pattern string, action string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(limit(action)(fn)))
	}

	regLimit("POST /api/v1/network/requests", "send_connection_request", h.SendRequest)
	reg("PUT /api/v1/network/requests/{id}/accept", h.AcceptRequest)
	reg("PUT /api/v1/network/requests/{id}/reject", h.RejectRequest)
	reg("POST /api/v1/network/block", h.Block)
	reg("DELETE /api/v1/network/connections/{userID}", h.Unconnect)
	reg("GET /api/v1/network/connections", h.ListConnections)
	reg("GET /api/v1/network/requests/incoming", h.ListIncomingRequests)
	reg("GET /api/v1/network/status/{userID}", h.GetStatus)
}
