package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/profile/application"
	"workspace-app/internal/profile/domain"
)

// VisibilityReader reports a user's profile visibility so the public-view handler
// can hide private profiles from other users. A nil reader disables the check.
type VisibilityReader interface {
	ProfileVisibility(ctx context.Context, userID string) (string, error)
}

type Handler struct {
	svc        *application.Service
	visibility VisibilityReader
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// SetVisibilityReader injects the profile-visibility reader.
func (h *Handler) SetVisibilityReader(v VisibilityReader) { h.visibility = v }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) write(w http.ResponseWriter, p *domain.Profile, err error) {
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toResponse(p))
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "not found")
	case errors.Is(err, domain.ErrOptimisticLock):
		common.WriteError(w, common.NewConflictError("profile was modified by another request; reload and retry"))
	default:
		common.WriteValidationError(w, err.Error())
	}
}

// GetMe handles GET /profiles/me.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	p, err := h.svc.Get(r.Context(), uid)
	h.write(w, p, err)
}

// GetByID handles GET /profiles/{id} (public view of another user).
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	viewerID := common.UserIDFromContext(r.Context())
	viewerRole := common.UserRoleFromContext(r.Context())

	p, err := h.svc.Get(r.Context(), targetID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	isOwner := targetID == viewerID

	// Profile level visibility check
	if !isOwner && !hasAccess(p.VisibilityProfile, viewerRole, false) {
		common.WriteNotFoundError(w, "profile not found")
		return
	}

	// Apply privacy settings/visibility checks for public views
	if !isOwner {
		if !hasAccess(p.VisibilitySalary, viewerRole, false) {
			p.SalaryMin = 0
			p.SalaryMax = 0
			p.SalaryCurrency = ""
		}
		if !hasAccess(p.VisibilityTransitionReason, viewerRole, false) {
			p.TransitionReason = ""
		}
		if !hasAccess(p.VisibilityExperience, viewerRole, false) {
			p.Experiences = nil
		}
		if !hasAccess(p.VisibilityEducation, viewerRole, false) {
			p.Educations = nil
		}
		if !hasAccess(p.VisibilityCertifications, viewerRole, false) {
			p.Certifications = nil
		}
		if !hasAccess(p.VisibilitySkills, viewerRole, false) {
			p.Skills = nil
		}
		if !hasAccess(p.VisibilityPortfolio, viewerRole, false) {
			p.Portfolio = nil
		}
		if !hasAccess(p.VisibilityReferences, viewerRole, false) {
			p.References = nil
		}
	}

	h.write(w, p, nil)
}

func hasAccess(vis string, viewerRole string, isOwner bool) bool {
	if isOwner {
		return true
	}
	switch vis {
	case "public", "":
		return true
	case "recruiters_only":
		return viewerRole == "recruiter" || viewerRole == "admin"
	case "mentors_only":
		return viewerRole == "mentor" || viewerRole == "admin"
	case "private":
		return false
	default:
		return false
	}
}

// UpdateMe handles PUT /profiles/me. It applies the scalar fields plus any
// provided child collections atomically (single transaction) with an optimistic
// version check. A stale req.Version yields HTTP 409.
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req updateScalarsRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}

	upd := domain.AggregateUpdate{Scalars: toScalars(req)}

	if req.Experiences != nil {
		items := make([]domain.WorkExperience, 0, len(*req.Experiences))
		for _, e := range *req.Experiences {
			items = append(items, e.toDomain())
		}
		upd.Experiences = &items
	}
	if req.Educations != nil {
		items := make([]domain.Education, 0, len(*req.Educations))
		for _, e := range *req.Educations {
			items = append(items, e.toDomain())
		}
		upd.Educations = &items
	}
	if req.Certifications != nil {
		items := make([]domain.Certification, 0, len(*req.Certifications))
		for _, c := range *req.Certifications {
			items = append(items, c.toDomain())
		}
		upd.Certifications = &items
	}
	if req.Skills != nil {
		items := make([]domain.ProfileSkill, 0, len(*req.Skills))
		for _, sk := range *req.Skills {
			items = append(items, sk.toDomain())
		}
		upd.Skills = &items
	}
	if req.Languages != nil {
		items := make([]domain.Language, 0, len(*req.Languages))
		for _, l := range *req.Languages {
			items = append(items, domain.Language(l))
		}
		upd.Languages = &items
	}
	if req.Portfolio != nil {
		items := make([]domain.PortfolioLink, 0, len(*req.Portfolio))
		for _, l := range *req.Portfolio {
			items = append(items, domain.PortfolioLink(l))
		}
		upd.Portfolio = &items
	}
	if req.References != nil {
		items := make([]domain.Reference, 0, len(*req.References))
		for _, ref := range *req.References {
			items = append(items, ref.toDomain())
		}
		upd.References = &items
	}

	p, err := h.svc.UpdateProfile(r.Context(), uid, req.Version, upd)
	h.write(w, p, err)
}

// --- experiences ---

func (h *Handler) AddExperience(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req experienceDTO
	if !decode(r, &req) || req.Title == "" || req.Company == "" {
		common.WriteValidationError(w, "title and company are required")
		return
	}
	e := req.toDomain()
	p, err := h.svc.AddExperience(r.Context(), uid, &e)
	h.write(w, p, err)
}

func (h *Handler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req experienceDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	req.ID = r.PathValue("id")
	p, err := h.svc.UpdateExperience(r.Context(), uid, req.toDomain())
	h.write(w, p, err)
}

func (h *Handler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	p, err := h.svc.DeleteExperience(r.Context(), uid, r.PathValue("id"))
	h.write(w, p, err)
}

// --- educations ---

func (h *Handler) AddEducation(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req educationDTO
	if !decode(r, &req) || req.School == "" {
		common.WriteValidationError(w, "school is required")
		return
	}
	e := req.toDomain()
	p, err := h.svc.AddEducation(r.Context(), uid, &e)
	h.write(w, p, err)
}

func (h *Handler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req educationDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	req.ID = r.PathValue("id")
	p, err := h.svc.UpdateEducation(r.Context(), uid, req.toDomain())
	h.write(w, p, err)
}

func (h *Handler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	p, err := h.svc.DeleteEducation(r.Context(), uid, r.PathValue("id"))
	h.write(w, p, err)
}

// --- certifications ---

func (h *Handler) AddCertification(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req certificationDTO
	if !decode(r, &req) || req.Name == "" {
		common.WriteValidationError(w, "name is required")
		return
	}
	c := req.toDomain()
	p, err := h.svc.AddCertification(r.Context(), uid, &c)
	h.write(w, p, err)
}

func (h *Handler) UpdateCertification(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req certificationDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	req.ID = r.PathValue("id")
	p, err := h.svc.UpdateCertification(r.Context(), uid, req.toDomain())
	h.write(w, p, err)
}

func (h *Handler) DeleteCertification(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	p, err := h.svc.DeleteCertification(r.Context(), uid, r.PathValue("id"))
	h.write(w, p, err)
}

// --- sets ---

func (h *Handler) SetSkills(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req skillsRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	skills := make([]domain.ProfileSkill, 0, len(req.Skills))
	for _, sk := range req.Skills {
		skills = append(skills, sk.toDomain())
	}
	p, err := h.svc.SetSkills(r.Context(), uid, skills)
	h.write(w, p, err)
}

func (h *Handler) SetLanguages(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req languagesRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	langs := make([]domain.Language, 0, len(req.Languages))
	for _, l := range req.Languages {
		langs = append(langs, domain.Language(l))
	}
	p, err := h.svc.SetLanguages(r.Context(), uid, langs)
	h.write(w, p, err)
}

func (h *Handler) SetPortfolio(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req portfolioRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	links := make([]domain.PortfolioLink, 0, len(req.Portfolio))
	for _, l := range req.Portfolio {
		links = append(links, domain.PortfolioLink(l))
	}
	p, err := h.svc.SetPortfolio(r.Context(), uid, links)
	h.write(w, p, err)
}

// --- new endpoints handlers ---

func (h *Handler) AddConsentLog(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req consentDTO
	if !decode(r, &req) || req.ConsentType == "" {
		common.WriteValidationError(w, "consent_type is required")
		return
	}
	cl := req.toDomain()
	cl.UserID = uid
	if cl.IPAddress == "" {
		cl.IPAddress = r.RemoteAddr
	}
	if cl.UserAgent == "" {
		cl.UserAgent = r.Header.Get("User-Agent")
	}

	err := h.svc.AddConsentLog(r.Context(), &cl)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, map[string]string{"status": "consent registered"})
}

func (h *Handler) AddEndorsement(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req endorsementDTO
	if !decode(r, &req) || req.ToUserID == "" || req.Text == "" {
		common.WriteValidationError(w, "to_user_id and text are required")
		return
	}
	e := req.toDomain()
	e.FromUserID = uid

	p, err := h.svc.AddEndorsement(r.Context(), req.ToUserID, &e)
	h.write(w, p, err)
}

func (h *Handler) AddReference(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req referenceDTO
	if !decode(r, &req) || req.Name == "" || req.Relationship == "" || req.ContactInfo == "" {
		common.WriteValidationError(w, "name, relationship, and contact_info are required")
		return
	}
	rf := req.toDomain()
	p, err := h.svc.AddReference(r.Context(), uid, &rf)
	h.write(w, p, err)
}

func (h *Handler) UpdateReference(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req referenceDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	req.ID = r.PathValue("id")
	p, err := h.svc.UpdateReference(r.Context(), uid, req.toDomain())
	h.write(w, p, err)
}

func (h *Handler) DeleteReference(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	p, err := h.svc.DeleteReference(r.Context(), uid, r.PathValue("id"))
	h.write(w, p, err)
}
