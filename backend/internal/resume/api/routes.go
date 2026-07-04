package api

import "net/http"

// RegisterRoutes mounts resume endpoints behind the auth middleware.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}

	reg("POST /api/v1/resumes", h.Upload)
	reg("GET /api/v1/resumes", h.List)
	reg("GET /api/v1/resumes/{id}", h.Get)
	reg("DELETE /api/v1/resumes/{id}", h.Delete)
	reg("POST /api/v1/resumes/{id}/versions", h.AddVersion)
	reg("GET /api/v1/resumes/{id}/versions", h.Versions)
	reg("GET /api/v1/resumes/{id}/score", h.Score)
	reg("POST /api/v1/resumes/{id}/review", h.Review)
}
