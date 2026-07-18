package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/profile/application"
	"workspace-app/internal/profile/domain"
)

type VisibilityReader interface {
	ProfileVisibility(ctx context.Context, userID string) (string, error)
}

type Handler struct {
	svc        *application.Service
	visibility VisibilityReader
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

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
	draft := r.URL.Query().Get("draft") != "false"
	var p *domain.Profile
	var err error
	if draft {
		p, err = h.svc.Get(r.Context(), uid)
	} else {
		p, err = h.svc.GetPublished(r.Context(), uid)
	}
	h.write(w, p, err)
}

// GetByID handles GET /profiles/{id} (public view of another user).
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	p, err := h.svc.GetPublished(r.Context(), targetID)
	h.write(w, p, err)
}

// e164Pattern validates an international phone number: leading '+', a country
// digit 1-9, then 7–14 more digits (ITU E.164, min national-number length).
var e164Pattern = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

// UpdateMe handles PUT /profiles/me.
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req updateProfileRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}

	if req.Phone != nil && *req.Phone != "" && !e164Pattern.MatchString(*req.Phone) {
		common.WriteValidationError(w, "phone must be in international E.164 format, e.g. +919876543210")
		return
	}
	if req.Email != nil && *req.Email != "" && !strings.Contains(*req.Email, "@") {
		common.WriteValidationError(w, "invalid email format")
		return
	}

	upd := req.toDomain()
	p, err := h.svc.UpdateProfile(r.Context(), uid, 0, upd)
	h.write(w, p, err)
}

// PublishMe handles POST /profiles/me/publish.
func (h *Handler) PublishMe(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	ip := r.RemoteAddr
	ua := r.Header.Get("User-Agent")
	p, err := h.svc.Publish(r.Context(), uid, uid, ip, ua)
	h.write(w, p, err)
}

// RollbackMe handles POST /profiles/me/rollback.
func (h *Handler) RollbackMe(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	ip := r.RemoteAddr
	ua := r.Header.Get("User-Agent")
	verStr := r.URL.Query().Get("version")
	ver, err := strconv.Atoi(verStr)
	if err != nil || ver <= 0 {
		common.WriteValidationError(w, "valid version query parameter is required")
		return
	}

	p, err := h.svc.Rollback(r.Context(), uid, ver, uid, ip, ua)
	h.write(w, p, err)
}

// GetVersions handles GET /profiles/me/versions.
func (h *Handler) GetVersions(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	vers, err := h.svc.ListVersions(r.Context(), uid)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, vers)
}

// GetAnalytics handles GET /profiles/me/analytics.
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	anal, err := h.svc.GetAnalytics(r.Context(), uid)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, anal)
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

// --- collections ---
func (h *Handler) SetSkills(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req skillsRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	skills := make([]domain.SkillItem, 0, len(req.Skills))
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
	langs := make([]domain.LanguageItem, 0, len(req.Languages))
	for _, l := range req.Languages {
		langs = append(langs, l.toDomain())
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
	links := make([]domain.ProjectItem, 0, len(req.Portfolio))
	for _, l := range req.Portfolio {
		links = append(links, l.toDomain())
	}
	p, err := h.svc.SetPortfolio(r.Context(), uid, links)
	h.write(w, p, err)
}

// --- consent & networking ---
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

	// Simple mock mapping to endorsementSummary
	sum := domain.EndorsementSummary{
		FromUserName: e.FromUserID,
		Relationship: e.Relationship,
		Text:         e.Text,
		CreatedAt:    time.Now(),
	}

	p, err := h.svc.AddEndorsement(r.Context(), req.ToUserID, &sum)
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
