package api

import "net/http"

// RegisterRoutes mounts referral endpoints behind the auth middleware.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}

	reg("POST /api/v1/referrals", h.Request)
	reg("GET /api/v1/referrals/incoming", h.Incoming)
	reg("GET /api/v1/referrals/outgoing", h.Outgoing)
	reg("POST /api/v1/referrals/{id}/accept", h.Accept)
	reg("POST /api/v1/referrals/{id}/decline", h.Decline)
	reg("PATCH /api/v1/referrals/{id}/outcome", h.UpdateOutcome)
}
