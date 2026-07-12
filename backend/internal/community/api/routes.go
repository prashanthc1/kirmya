package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}
	reg("GET /api/v1/communities", h.List)
	reg("POST /api/v1/communities", h.Create)
	reg("GET /api/v1/communities/{slug}", h.Get)
	reg("POST /api/v1/communities/{slug}/join", h.ToggleJoin)
	reg("GET /api/v1/communities/{slug}/posts", h.ListPosts)
	reg("POST /api/v1/communities/{slug}/posts", h.CreatePost)
	reg("GET /api/v1/communities/{slug}/tags", h.ListTags)
	reg("GET /api/v1/communities/{slug}/reports", h.ListReports)
	reg("DELETE /api/v1/communities/{slug}/posts/{id}", h.HidePost)
	reg("GET /api/v1/posts/{id}/comments", h.ListComments)
	reg("POST /api/v1/posts/{id}/comments", h.AddComment)
	reg("POST /api/v1/posts/{id}/reactions", h.ToggleReaction)
	reg("POST /api/v1/posts/{id}/polls", h.CreatePoll)
	reg("POST /api/v1/posts/{id}/report", h.ReportPost)
	reg("GET /api/v1/polls/{id}", h.GetPoll)
	reg("POST /api/v1/polls/{id}/vote", h.Vote)
}
