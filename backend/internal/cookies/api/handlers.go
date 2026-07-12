package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"workspace-app/internal/common"
	"workspace-app/internal/cookies/application"
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

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
}

func getCountry(r *http.Request) string {
	if cty := r.Header.Get("CF-IPCountry"); cty != "" {
		return cty
	}
	return "Unknown"
}

// GET /api/cookies/preferences
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())

	anonID := r.URL.Query().Get("anonymous_id")
	if anonID == "" {
		anonID = r.Header.Get("X-Anonymous-ID")
	}

	prefs, err := h.svc.GetPreferences(r.Context(), userID, anonID)
	if err != nil {
		common.WriteInternalError(w, err.Error())
		return
	}

	common.WriteSuccess(w, http.StatusOK, toDTO(prefs))
}

// POST/PUT /api/cookies/preferences
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	var req saveRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}

	userIDStr := common.UserIDFromContext(r.Context())
	var userID *string
	if userIDStr != "" {
		userID = &userIDStr
	}

	in := application.SaveInput{
		UserID:          userID,
		AnonymousID:     req.AnonymousID,
		Functional:      req.Functional,
		Analytics:       req.Analytics,
		Marketing:       req.Marketing,
		Performance:     req.Performance,
		Personalization: req.Personalization,
		AIPreferences:   req.AIPreferences,
		ConsentVersion:  req.ConsentVersion,
		IPAddress:       getIP(r),
		Country:         getCountry(r),
		UserAgent:       r.UserAgent(),
	}

	prefs, err := h.svc.SavePreferences(r.Context(), in)
	if err != nil {
		common.WriteInternalError(w, err.Error())
		return
	}

	common.WriteSuccess(w, http.StatusOK, toDTO(prefs))
}

// DELETE /api/cookies/preferences
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())

	anonID := r.URL.Query().Get("anonymous_id")
	if anonID == "" {
		anonID = r.Header.Get("X-Anonymous-ID")
	}

	if userID == "" && anonID == "" {
		common.WriteValidationError(w, "must provide user_id or anonymous_id to delete preferences")
		return
	}

	err := h.svc.DeletePreferences(r.Context(), userID, anonID)
	if err != nil {
		common.WriteInternalError(w, err.Error())
		return
	}

	common.WriteSuccess(w, http.StatusOK, map[string]bool{"success": true})
}
