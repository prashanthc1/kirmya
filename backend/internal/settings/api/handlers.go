package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/settings/application"
	"workspace-app/internal/settings/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
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
		// Map DB errors to nice user messages if possible
		if strings.Contains(err.Error(), "uq_users_username") {
			common.WriteValidationError(w, "this username is already taken")
			return
		}
		if strings.Contains(err.Error(), "uq_profiles_custom_url") {
			common.WriteValidationError(w, "this custom profile URL is already taken")
			return
		}
		common.WriteInternalError(w, err.Error())
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

// PATCH /api/v1/settings
func (h *Handler) UpdateSettingsUnified(w http.ResponseWriter, r *http.Request) {
	var req unifiedSettingsPatchRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	s, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	if req.General != nil {
		s.Language = req.General.Language
		s.Timezone = req.General.Timezone
		s.Theme = req.General.Theme
		s.EmailDigest = req.General.EmailDigest
	}
	if req.Privacy != nil {
		s.ProfileVisibility = req.Privacy.ProfileVisibility
		s.ShowEmail = req.Privacy.ShowEmail
		s.Discoverable = req.Privacy.Discoverable
		s.AllowMessages = req.Privacy.AllowMessages
	}
	if req.Notifications != nil {
		s.Notifications = req.Notifications.toDomain()
	}
	if req.Security != nil {
		s.LoginAlerts = req.Security.LoginAlerts
	}
	if req.Accessibility != nil {
		s.FontSize = req.Accessibility.FontSize
		s.HighContrast = req.Accessibility.HighContrast
		s.ReducedMotion = req.Accessibility.ReducedMotion
		s.CompactMode = req.Accessibility.CompactMode
		s.DefaultLandingPage = req.Accessibility.DefaultLandingPage
		s.AccessibilityKeyboardNavigation = req.Accessibility.AccessibilityKeyboardNavigation
		s.AccessibilityScreenReader = req.Accessibility.AccessibilityScreenReader
		s.AccessibilityFocusIndicators = req.Accessibility.AccessibilityFocusIndicators
	}
	if req.AI != nil {
		s.EnableAIAssistant = req.AI.EnableAIAssistant
		s.AIJobRecommendations = req.AI.AIJobRecommendations
		s.AIResumeSuggestions = req.AI.AIResumeSuggestions
		s.AIRoadmapSuggestions = req.AI.AIRoadmapSuggestions
		s.AISkillGapAnalysis = req.AI.AISkillGapAnalysis
		s.AIInterviewPrep = req.AI.AIInterviewPrep
		s.AILearningRecommendations = req.AI.AILearningRecommendations
	}
	if req.Learning != nil {
		s.LearningGoals = req.Learning.LearningGoals
		s.TechnologiesOfInterest = req.Learning.TechnologiesOfInterest
		s.CertificationGoals = req.Learning.CertificationGoals
		s.LearningReminders = req.Learning.LearningReminders
	}

	updated, err := h.svc.UpdateUserSettings(r.Context(), s)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(updated))
}

// GET /api/v1/profile/settings
func (h *Handler) GetProfileSettings(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	username, customURL, profileVisibility, fieldVisibility, openToWork, referralEligible, willingToMentor, err := h.svc.GetProfileSettings(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, profileSettingsDTO{
		Username:          username,
		CustomURL:         customURL,
		ProfileVisibility: profileVisibility,
		FieldVisibility:   fieldVisibility,
		OpenToWork:        openToWork,
		ReferralEligible:  referralEligible,
		WillingToMentor:   willingToMentor,
	})
}

// PATCH /api/v1/profile/settings
func (h *Handler) UpdateProfileSettings(w http.ResponseWriter, r *http.Request) {
	var req profileSettingsDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	err := h.svc.UpdateProfileSettings(
		r.Context(), userID, req.Username, req.CustomURL, req.ProfileVisibility,
		req.FieldVisibility, req.OpenToWork, req.ReferralEligible, req.WillingToMentor,
	)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"success": true})
}

// GET /api/v1/security/activity
func (h *Handler) GetSecurityActivity(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	sessions, err := h.svc.ListActiveSessions(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	history, err := h.svc.ListSecurityHistory(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	var sessionDTOs []activeSessionDTO
	for _, s := range sessions {
		sessionDTOs = append(sessionDTOs, activeSessionDTO{
			ID:        s.ID,
			UserAgent: s.UserAgent,
			IPAddress: s.IPAddress,
			CreatedAt: s.CreatedAt.UTC().Format(time.RFC3339),
			ExpiresAt: s.ExpiresAt.UTC().Format(time.RFC3339),
		})
	}

	var historyDTOs []securityHistoryDTO
	for _, hist := range history {
		historyDTOs = append(historyDTOs, securityHistoryDTO{
			ID:        hist.ID,
			Action:    hist.Action,
			IPAddress: hist.IPAddress,
			CreatedAt: hist.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	// Fetch connected accounts
	connAccounts, err := h.svc.ListConnectedAccounts(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	var connDTOs []connectedAccountDTO
	for _, c := range connAccounts {
		connDTOs = append(connDTOs, connectedAccountDTO{
			ID:          c.ID,
			Provider:    c.Provider,
			ProviderUID: c.ProviderUID,
			CreatedAt:   c.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"sessions":           sessionDTOs,
		"history":            historyDTOs,
		"connected_accounts": connDTOs,
	})
}

// POST /api/v1/security/password/change
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req changePasswordRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	err := h.svc.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword, getIP(r))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"success": true})
}

// POST /api/v1/security/logout-device
func (h *Handler) LogoutDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	err := h.svc.RevokeSession(r.Context(), userID, req.SessionID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"success": true})
}

// GET /api/v1/privacy/settings
func (h *Handler) GetPrivacySettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Get(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"profile_visibility": s.ProfileVisibility,
		"show_email":         s.ShowEmail,
		"discoverable":       s.Discoverable,
		"allow_messages":     s.AllowMessages,
		"enable_ai_assistant": s.EnableAIAssistant,
	})
}

// PATCH /api/v1/privacy/settings
func (h *Handler) UpdatePrivacySettings(w http.ResponseWriter, r *http.Request) {
	var req privacyDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	s, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	s.ProfileVisibility = req.ProfileVisibility
	s.ShowEmail = req.ShowEmail
	s.Discoverable = req.Discoverable
	s.AllowMessages = req.AllowMessages

	updated, err := h.svc.UpdateUserSettings(r.Context(), s)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toDTO(updated))
}

// GET /api/v1/notifications/preferences
func (h *Handler) GetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Get(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toNotificationsDTO(s.Notifications))
}

// PATCH /api/v1/notifications/preferences
func (h *Handler) UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	var req notificationsDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	s, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	s.Notifications = req.toDomain()
	updated, err := h.svc.UpdateUserSettings(r.Context(), s)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toNotificationsDTO(updated.Notifications))
}

// GET /api/v1/privacy/cookies
func (h *Handler) GetCookies(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromContext(r.Context())
	cc, err := h.svc.GetCookieConsent(r.Context(), userID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, cookieConsentDTO{
		Functional:        cc.Functional,
		Analytics:         cc.Analytics,
		AIPersonalization: cc.AIPersonalization,
	})
}

// POST/PATCH /api/v1/privacy/cookies
func (h *Handler) UpdateCookies(w http.ResponseWriter, r *http.Request) {
	var req cookieConsentDTO
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	userID := common.UserIDFromContext(r.Context())
	cc, err := h.svc.SaveCookieConsent(r.Context(), userID, req.Functional, req.Analytics, req.AIPersonalization)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, cookieConsentDTO{
		Functional:        cc.Functional,
		Analytics:         cc.Analytics,
		AIPersonalization: cc.AIPersonalization,
	})
}

// --- Legacy segment updates ---

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
