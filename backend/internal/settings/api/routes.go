package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	reg := func(pattern string, fn http.HandlerFunc) {
		mux.Handle(pattern, auth(fn))
	}
	// General / Settings
	reg("GET /api/v1/settings", h.Get)
	reg("PATCH /api/v1/settings", h.UpdateSettingsUnified)

	// Profile Settings Controls
	reg("GET /api/v1/profile/settings", h.GetProfileSettings)
	reg("PATCH /api/v1/profile/settings", h.UpdateProfileSettings)

	// Security Settings & Active Sessions
	reg("GET /api/v1/security/activity", h.GetSecurityActivity)
	reg("POST /api/v1/security/password/change", h.ChangePassword)
	reg("POST /api/v1/security/logout-device", h.LogoutDevice)

	// Privacy
	reg("GET /api/v1/privacy/settings", h.GetPrivacySettings)
	reg("PATCH /api/v1/privacy/settings", h.UpdatePrivacySettings)

	// Notifications
	reg("GET /api/v1/notifications/preferences", h.GetNotificationPreferences)
	reg("PATCH /api/v1/notifications/preferences", h.UpdateNotificationPreferences)

	// Cookie Consents
	reg("GET /api/v1/privacy/cookies", h.GetCookies)
	reg("POST /api/v1/privacy/cookies", h.UpdateCookies)
	reg("PATCH /api/v1/privacy/cookies", h.UpdateCookies)

	// Legacy endpoints (for safety and backward compatibility)
	reg("PATCH /api/v1/settings/general", h.UpdateGeneral)
	reg("PATCH /api/v1/settings/privacy", h.UpdatePrivacy)
	reg("PATCH /api/v1/settings/notifications", h.UpdateNotifications)
	reg("PATCH /api/v1/settings/security", h.UpdateSecurity)
}
