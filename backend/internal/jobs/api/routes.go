package api

import "net/http"

// RegisterRoutes mounts job endpoints. Seeker actions (search/get/apply/save and
// the caller's own lists) require only authentication; posting and managing jobs
// and applicants is gated to recruiters via recruiterOnly. More-specific literal
// paths ("saved", "applications") take precedence over "{id}".
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth, recruiterOnly func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(http.HandlerFunc(fn)))
	}
	recruiter := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, recruiterOnly(http.HandlerFunc(fn)))
	}

	// Seeker / any authenticated user.
	reg("GET /api/v1/jobs", h.Search)
	reg("GET /api/v1/jobs/saved", h.SavedJobs)
	reg("GET /api/v1/jobs/applications", h.MyApplications)
	reg("GET /api/v1/jobs/matches", h.Matches)
	reg("GET /api/v1/jobs/{id}", h.GetJob)
	reg("POST /api/v1/jobs/{id}/apply", h.Apply)
	reg("POST /api/v1/jobs/{id}/save", h.ToggleSave)

	// Recruiter only.
	recruiter("POST /api/v1/jobs", h.PostJob)
	recruiter("PUT /api/v1/jobs/{id}", h.UpdateJob)
	recruiter("DELETE /api/v1/jobs/{id}", h.DeleteJob)
	recruiter("GET /api/v1/jobs/{id}/applicants", h.JobApplicants)
	recruiter("PATCH /api/v1/applications/{id}", h.UpdateApplicationStatus)
}
