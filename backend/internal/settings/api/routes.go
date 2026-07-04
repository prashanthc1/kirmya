package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}
	reg("GET /api/v1/settings", h.Get)
	reg("PATCH /api/v1/settings/general", h.UpdateGeneral)
	reg("PATCH /api/v1/settings/privacy", h.UpdatePrivacy)
	reg("PATCH /api/v1/settings/notifications", h.UpdateNotifications)
	reg("PATCH /api/v1/settings/security", h.UpdateSecurity)
}
