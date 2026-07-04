package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(http.HandlerFunc(fn)))
	}
	reg("POST /api/v1/mentorship/mentors", h.BecomeMentor)
	reg("GET /api/v1/mentorship/mentors", h.ListMentors)
	reg("GET /api/v1/mentorship/mentors/{id}", h.GetMentor)
	reg("GET /api/v1/mentorship/mentors/{id}/reviews", h.MentorReviews)
	reg("GET /api/v1/mentorship/mentors/{id}/availability", h.MentorAvailability)
	reg("POST /api/v1/mentorship/availability", h.AddAvailability)
	reg("POST /api/v1/mentorship/sessions", h.Book)
	reg("GET /api/v1/mentorship/sessions", h.MySessions)
	reg("PATCH /api/v1/mentorship/sessions/{id}", h.UpdateStatus)
	reg("POST /api/v1/mentorship/sessions/{id}/review", h.Review)
}
