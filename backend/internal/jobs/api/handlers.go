package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"workspace-app/internal/common"
	"workspace-app/internal/jobs/application"
	"workspace-app/internal/jobs/domain"
)

type Handler struct {
	svc      *application.Service
	matchSvc *application.MatchService
}

func NewHandler(svc *application.Service, matchSvc *application.MatchService) *Handler {
	return &Handler{svc: svc, matchSvc: matchSvc}
}

// Matches handles GET /jobs/matches — AI-ranked job recommendations for the
// current seeker (heuristic fallback when AI is unavailable).
func (h *Handler) Matches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.matchSvc.Matches(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"matches": toMatches(matches)})
}

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, application.ErrForbidden):
		common.WriteForbiddenError(w, "you do not own this resource")
	case errors.Is(err, domain.ErrJobNotFound):
		common.WriteNotFoundError(w, "job not found")
	case errors.Is(err, domain.ErrApplicationNotFound):
		common.WriteNotFoundError(w, "application not found")
	case errors.Is(err, domain.ErrAlreadyApplied):
		common.WriteError(w, common.NewConflictError("you have already applied to this job"))
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

// PostJob handles POST /jobs.
func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	var req postJobRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	job, err := h.svc.PostJob(r.Context(), common.UserIDFromContext(r.Context()), application.PostJobInput(req))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toJobResponse(job))
}

// Search handles GET /jobs.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 0
	if v := q.Get("limit"); v != "" {
		limit, _ = strconv.Atoi(v)
	}
	// Recruiter "my postings": ?mine=true or ?posted_by=me resolves to the
	// authenticated caller; an explicit ?posted_by=<id> filters by that poster.
	postedBy := q.Get("posted_by")
	if postedBy == "me" || q.Get("mine") == "true" {
		postedBy = common.UserIDFromContext(r.Context())
	}
	jobs, err := h.svc.SearchJobs(r.Context(), domain.Filter{
		Keyword: q.Get("q"), Location: q.Get("location"), JobType: q.Get("type"), PostedBy: postedBy, Limit: limit,
	})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"jobs": toJobResponses(jobs)})
}

// GetJob handles GET /jobs/{id}.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	job, err := h.svc.GetJob(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toJobResponse(job))
}

// UpdateJob handles PUT /jobs/{id}.
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	var req postJobRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	job, err := h.svc.UpdateJob(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), application.PostJobInput(req))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toJobResponse(job))
}

// DeleteJob handles DELETE /jobs/{id}.
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DeleteJob(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Apply handles POST /jobs/{id}/apply.
func (h *Handler) Apply(w http.ResponseWriter, r *http.Request) {
	var req applyRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	app, err := h.svc.Apply(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.CoverLetter)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toAppResponse(app))
}

// ToggleSave handles POST /jobs/{id}/save.
func (h *Handler) ToggleSave(w http.ResponseWriter, r *http.Request) {
	saved, err := h.svc.ToggleSave(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"saved": saved})
}

// SavedJobs handles GET /jobs/saved.
func (h *Handler) SavedJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.svc.SavedJobs(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"jobs": toJobResponses(jobs)})
}

// MyApplications handles GET /jobs/applications.
func (h *Handler) MyApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.svc.MyApplications(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"applications": toAppResponses(apps)})
}

// JobApplicants handles GET /jobs/{id}/applicants.
func (h *Handler) JobApplicants(w http.ResponseWriter, r *http.Request) {
	apps, err := h.svc.JobApplicants(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"applications": toAppResponses(apps)})
}

// UpdateApplicationStatus handles PATCH /applications/{id}.
func (h *Handler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	var req statusRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	app, err := h.svc.UpdateApplicationStatus(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Status)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toAppResponse(app))
}
