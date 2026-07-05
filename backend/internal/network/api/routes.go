package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}

	reg("POST /api/v1/network/requests", h.SendRequest)
	reg("PUT /api/v1/network/requests/{id}/accept", h.AcceptRequest)
	reg("PUT /api/v1/network/requests/{id}/reject", h.RejectRequest)
	reg("GET /api/v1/network/connections", h.ListConnections)
	reg("GET /api/v1/network/requests/incoming", h.ListIncomingRequests)
	reg("GET /api/v1/network/status/{userID}", h.GetStatus)
}
