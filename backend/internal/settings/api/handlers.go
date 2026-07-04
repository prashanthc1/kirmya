package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/settings/application"
	"workspace-app/internal/settings/domain"
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
	case errors.Is(err, domain.ErrOptimisticLock):
		common.WriteError(w, common.NewConflictError("settings were modified concurrently; reload and retry"))
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "settings not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

// GET /api/v1/settings
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Get(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(s))
}

// PATCH /api/v1/settings/general
func (h *Handler) UpdateGeneral(w http.ResponseWriter, r *http.Request) {
	var req generalDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	s, err := h.svc.UpdateGeneral(r.Context(), common.UserIDFromContext(r.Context()), application.GeneralInput{
		Language:    req.Language,
		Timezone:    req.Timezone,
		Theme:       req.Theme,
		EmailDigest: req.EmailDigest,
	})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(s))
}

// PATCH /api/v1/settings/privacy
func (h *Handler) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	var req privacyDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	s, err := h.svc.UpdatePrivacy(r.Context(), common.UserIDFromContext(r.Context()), application.PrivacyInput{
		ProfileVisibility: req.ProfileVisibility,
		ShowEmail:         req.ShowEmail,
		Discoverable:      req.Discoverable,
		AllowMessages:     req.AllowMessages,
	})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(s))
}

// PATCH /api/v1/settings/notifications
func (h *Handler) UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	var req notificationsDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	s, err := h.svc.UpdateNotifications(r.Context(), common.UserIDFromContext(r.Context()), req.toDomain())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(s))
}

// PATCH /api/v1/settings/security
func (h *Handler) UpdateSecurity(w http.ResponseWriter, r *http.Request) {
	var req securityDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	s, err := h.svc.UpdateSecurity(r.Context(), common.UserIDFromContext(r.Context()), application.SecurityInput{
		LoginAlerts: req.LoginAlerts,
	})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(s))
}
