package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/network/application"
	"workspace-app/internal/network/domain"
)

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler {
	return &Handler{svc: svc}
}

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		common.WriteNotFoundError(w, "not found")
		return
	}
	if errors.Is(err, domain.ErrDuplicateRequest) || errors.Is(err, domain.ErrSelfConnection) || errors.Is(err, domain.ErrInvalidTransition) {
		common.WriteValidationError(w, err.Error())
		return
	}
	common.WriteInternalError(w, err.Error())
}

func (h *Handler) SendRequest(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	var req sendRequest
	if !decode(r, &req) || req.ReceiverID == "" {
		common.WriteValidationError(w, "receiver_id is required")
		return
	}

	c, err := h.svc.SendRequest(r.Context(), uid, req.ReceiverID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	common.WriteSuccess(w, http.StatusCreated, toResponse(*c))
}

func (h *Handler) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	connectionID := r.PathValue("id")

	err := h.svc.AcceptRequest(r.Context(), uid, connectionID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "accepted"})
}

func (h *Handler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	connectionID := r.PathValue("id")

	err := h.svc.RejectRequest(r.Context(), uid, connectionID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (h *Handler) ListConnections(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())

	list, err := h.svc.GetConnections(r.Context(), uid)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	res := make([]connectionResponse, 0, len(list))
	for _, c := range list {
		res = append(res, toResponse(c))
	}

	common.WriteSuccess(w, http.StatusOK, res)
}

func (h *Handler) ListIncomingRequests(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())

	list, err := h.svc.GetIncomingRequests(r.Context(), uid)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	res := make([]connectionResponse, 0, len(list))
	for _, c := range list {
		res = append(res, toResponse(c))
	}

	common.WriteSuccess(w, http.StatusOK, res)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	targetID := r.PathValue("userID")

	status, reqID, err := h.svc.GetConnectionStatus(r.Context(), uid, targetID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	common.WriteSuccess(w, http.StatusOK, statusResponse{
		Status:      string(status),
		RequesterID: reqID,
	})
}
