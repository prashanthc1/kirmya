package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"workspace-app/internal/ai/application"
	"workspace-app/internal/ai/domain"
	"workspace-app/internal/common"
)

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, application.ErrForbidden):
		common.WriteForbiddenError(w, "you do not own this conversation")
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "not found")
	case errors.Is(err, domain.ErrLLMNotReady):
		common.WriteError(w, &common.AppError{
			Code: "ai_unavailable", Status: http.StatusServiceUnavailable,
			Message: "AI features are not configured on this server",
		})
	default:
		// Upstream LLM/transport failure.
		common.WriteError(w, &common.AppError{
			Code: "ai_error", Status: http.StatusBadGateway,
			Message: "the AI provider could not complete the request",
		})
	}
}

// ResumeReview handles POST /ai/resume-review.
func (h *Handler) ResumeReview(w http.ResponseWriter, r *http.Request) {
	var req resumeReviewRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	review, err := h.svc.ReviewResume(r.Context(), common.UserIDFromContext(r.Context()), req.ResumeText)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, review)
}

// SkillGap handles POST /ai/skill-gap.
func (h *Handler) SkillGap(w http.ResponseWriter, r *http.Request) {
	var req skillGapRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	gap, err := h.svc.SkillGap(r.Context(), common.UserIDFromContext(r.Context()), req.CurrentRole, req.TargetRole, req.CurrentSkills)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, gap)
}

// Coach handles POST /ai/coach.
func (h *Handler) Coach(w http.ResponseWriter, r *http.Request) {
	var req coachRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	reply, threadID, err := h.svc.Coach(r.Context(), common.UserIDFromContext(r.Context()), req.ThreadID, req.Message)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"thread_id": threadID, "reply": reply})
}

// ListThreads handles GET /ai/coach/threads.
func (h *Handler) ListThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := h.svc.ListThreads(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]threadSummary, 0, len(threads))
	for i := range threads {
		out = append(out, toThreadSummary(&threads[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"threads": out})
}

// GetThread handles GET /ai/coach/threads/{id}.
func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	t, err := h.svc.GetThread(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toThreadDetail(t))
}
