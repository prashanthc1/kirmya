package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/mentorship/application"
	"workspace-app/internal/mentorship/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, domain.ErrNotMentor), errors.Is(err, domain.ErrNotMentee), errors.Is(err, domain.ErrSlotNotOwned):
		common.WriteForbiddenError(w, "you cannot act on this resource")
	case errors.Is(err, domain.ErrNotComplete):
		common.WriteError(w, common.NewConflictError("session must be completed before review"))
	case errors.Is(err, domain.ErrSlotUnavailable):
		common.WriteError(w, common.NewConflictError("availability slot is already booked"))
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrSlotNotFound):
		common.WriteNotFoundError(w, "not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

type mentorDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Headline  string `json:"headline"`
	Bio       string `json:"bio,omitempty"`
	Expertise string `json:"expertise,omitempty"`
}

func toMentor(p *domain.MentorProfile) mentorDTO {
	return mentorDTO{ID: p.ID, UserID: p.UserID, Headline: p.Headline, Bio: p.Bio, Expertise: p.Expertise}
}

type sessionDTO struct {
	ID          string `json:"id"`
	MentorID    string `json:"mentor_id"`
	MenteeID    string `json:"mentee_id"`
	Topic       string `json:"topic,omitempty"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduled_at"`
}

func toSession(s *domain.Session) sessionDTO {
	return sessionDTO{
		ID: s.ID, MentorID: s.MentorID, MenteeID: s.MenteeID, Topic: s.Topic, Status: s.Status,
		ScheduledAt: s.ScheduledAt.UTC().Format(time.RFC3339),
	}
}

// --- mentors ---

type becomeMentorRequest struct {
	Headline  string `json:"headline"`
	Bio       string `json:"bio"`
	Expertise string `json:"expertise"`
}

func (h *Handler) BecomeMentor(w http.ResponseWriter, r *http.Request) {
	var req becomeMentorRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	p, err := h.svc.BecomeMentor(r.Context(), common.UserIDFromContext(r.Context()),
		application.MentorInput(req))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toMentor(p))
}

func (h *Handler) ListMentors(w http.ResponseWriter, r *http.Request) {
	mentors, err := h.svc.ListMentors(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]mentorDTO, 0, len(mentors))
	for i := range mentors {
		out = append(out, toMentor(&mentors[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"mentors": out})
}

func (h *Handler) GetMentor(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetMentor(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toMentor(p))
}

// --- sessions ---

type bookRequest struct {
	MentorID    string `json:"mentor_id"`
	Topic       string `json:"topic"`
	ScheduledAt string `json:"scheduled_at"`
	SlotID      string `json:"slot_id"`
}

func (h *Handler) Book(w http.ResponseWriter, r *http.Request) {
	var req bookRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	when, _ := time.Parse(time.RFC3339, req.ScheduledAt)
	sess, err := h.svc.Book(r.Context(), common.UserIDFromContext(r.Context()), application.BookInput{
		MentorID:    req.MentorID,
		Topic:       req.Topic,
		ScheduledAt: when,
		SlotID:      req.SlotID,
	})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toSession(sess))
}

// --- availability ---

type slotDTO struct {
	ID       string `json:"id"`
	MentorID string `json:"mentor_id"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	IsBooked bool   `json:"is_booked"`
}

func toSlot(s *domain.AvailabilitySlot) slotDTO {
	return slotDTO{
		ID: s.ID, MentorID: s.MentorID,
		StartsAt: s.StartsAt.UTC().Format(time.RFC3339),
		EndsAt:   s.EndsAt.UTC().Format(time.RFC3339),
		IsBooked: s.IsBooked,
	}
}

type addAvailabilityRequest struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

func (h *Handler) AddAvailability(w http.ResponseWriter, r *http.Request) {
	var req addAvailabilityRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	starts, _ := time.Parse(time.RFC3339, req.StartsAt)
	ends, _ := time.Parse(time.RFC3339, req.EndsAt)
	slot, err := h.svc.AddAvailability(r.Context(), common.UserIDFromContext(r.Context()), starts, ends)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toSlot(slot))
}

func (h *Handler) MentorAvailability(w http.ResponseWriter, r *http.Request) {
	slots, err := h.svc.MentorAvailability(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]slotDTO, 0, len(slots))
	for i := range slots {
		out = append(out, toSlot(&slots[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"slots": out})
}

func (h *Handler) MySessions(w http.ResponseWriter, r *http.Request) {
	asMentee, asMentor, err := h.svc.Sessions(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"as_mentee": mapSessions(asMentee),
		"as_mentor": mapSessions(asMentor),
	})
}

func mapSessions(list []domain.Session) []sessionDTO {
	out := make([]sessionDTO, 0, len(list))
	for i := range list {
		out = append(out, toSession(&list[i]))
	}
	return out
}

type statusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req statusRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	sess, err := h.svc.UpdateStatus(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Status)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toSession(sess))
}

type reviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	var req reviewRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	rv, err := h.svc.Review(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Rating, req.Comment)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, map[string]any{"id": rv.ID, "rating": rv.Rating, "comment": rv.Comment})
}

type reviewDTO struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) MentorReviews(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.MentorReviews(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]reviewDTO, 0, len(stats.Reviews))
	for _, rv := range stats.Reviews {
		out = append(out, reviewDTO{
			ID: rv.ID, SessionID: rv.SessionID, Rating: rv.Rating, Comment: rv.Comment,
			CreatedAt: rv.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"reviews":        out,
		"average_rating": stats.AverageRating,
		"count":          stats.Count,
	})
}
