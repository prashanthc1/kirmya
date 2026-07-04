package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/referrals/application"
	"workspace-app/internal/referrals/domain"
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
		common.WriteForbiddenError(w, "you may not act on this referral")
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "referral not found")
	case errors.Is(err, domain.ErrAlreadyDecided):
		common.WriteError(w, common.NewConflictError("this referral has already been decided"))
	case errors.Is(err, domain.ErrNotAccepted):
		common.WriteError(w, common.NewConflictError("referral must be accepted before tracking an outcome"))
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

// Request handles POST /referrals.
func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
	var req requestRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	ref, err := h.svc.Request(r.Context(), common.UserIDFromContext(r.Context()), application.RequestInput(req))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toResponse(ref))
}

// Incoming handles GET /referrals/incoming.
func (h *Handler) Incoming(w http.ResponseWriter, r *http.Request) {
	refs, err := h.svc.Incoming(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"referrals": toResponses(refs)})
}

// Outgoing handles GET /referrals/outgoing.
func (h *Handler) Outgoing(w http.ResponseWriter, r *http.Request) {
	refs, err := h.svc.Outgoing(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"referrals": toResponses(refs)})
}

// Accept handles POST /referrals/{id}/accept.
func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	ref, err := h.svc.Accept(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toResponse(ref))
}

// Decline handles POST /referrals/{id}/decline.
func (h *Handler) Decline(w http.ResponseWriter, r *http.Request) {
	ref, err := h.svc.Decline(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toResponse(ref))
}

// UpdateOutcome handles PATCH /referrals/{id}/outcome.
func (h *Handler) UpdateOutcome(w http.ResponseWriter, r *http.Request) {
	var req outcomeRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	ref, err := h.svc.UpdateOutcome(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Outcome)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toResponse(ref))
}
