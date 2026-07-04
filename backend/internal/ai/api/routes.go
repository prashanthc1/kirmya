package api

import "net/http"

// RegisterRoutes mounts AI endpoints behind the auth middleware.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}

	reg("POST /api/v1/ai/resume-review", h.ResumeReview)
	reg("POST /api/v1/ai/skill-gap", h.SkillGap)
	reg("POST /api/v1/ai/coach", h.Coach)
	reg("GET /api/v1/ai/coach/threads", h.ListThreads)
	reg("GET /api/v1/ai/coach/threads/{id}", h.GetThread)
}
