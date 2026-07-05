package api

import "net/http"

// RegisterRoutes mounts profile endpoints, all behind the auth middleware.
func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}

	reg("GET /api/v1/profiles/me", h.GetMe)
	reg("PUT /api/v1/profiles/me", h.UpdateMe)

	reg("PUT /api/v1/profiles/me/skills", h.SetSkills)
	reg("PUT /api/v1/profiles/me/languages", h.SetLanguages)
	reg("PUT /api/v1/profiles/me/portfolio", h.SetPortfolio)

	reg("POST /api/v1/profiles/me/experiences", h.AddExperience)
	reg("PUT /api/v1/profiles/me/experiences/{id}", h.UpdateExperience)
	reg("DELETE /api/v1/profiles/me/experiences/{id}", h.DeleteExperience)

	reg("POST /api/v1/profiles/me/educations", h.AddEducation)
	reg("PUT /api/v1/profiles/me/educations/{id}", h.UpdateEducation)
	reg("DELETE /api/v1/profiles/me/educations/{id}", h.DeleteEducation)

	reg("POST /api/v1/profiles/me/certifications", h.AddCertification)
	reg("PUT /api/v1/profiles/me/certifications/{id}", h.UpdateCertification)
	reg("DELETE /api/v1/profiles/me/certifications/{id}", h.DeleteCertification)

	// New endpoints
	reg("POST /api/v1/profiles/me/consent", h.AddConsentLog)
	reg("POST /api/v1/profiles/me/endorsements", h.AddEndorsement)
	reg("POST /api/v1/profiles/me/references", h.AddReference)
	reg("PUT /api/v1/profiles/me/references/{id}", h.UpdateReference)
	reg("DELETE /api/v1/profiles/me/references/{id}", h.DeleteReference)

	// Public view of another user's profile (most-specific "me" routes win).
	reg("GET /api/v1/profiles/{id}", h.GetByID)
}
