package domain

import "context"

// Repository is the persistence port for the Settings context. The concrete
// adapter lives in infrastructure/postgres.
type Repository interface {
	// Get returns the user's settings, or ErrNotFound if no row exists yet.
	Get(ctx context.Context, userID string) (*UserSettings, error)

	// EnsureDefaults inserts the default row for a user if one does not already
	// exist and returns the current settings. It is idempotent.
	EnsureDefaults(ctx context.Context, userID string) (*UserSettings, error)

	// Update persists the aggregate. It is version-checked and returns
	// ErrOptimisticLock on a stale write.
	Update(ctx context.Context, s *UserSettings) error

	// Connected Accounts
	ListConnectedAccounts(ctx context.Context, userID string) ([]ConnectedAccount, error)
	DisconnectAccount(ctx context.Context, userID string, provider string) error

	// Cookie Consent
	GetCookieConsent(ctx context.Context, userID string) (*CookieConsent, error)
	SaveCookieConsent(ctx context.Context, cc *CookieConsent) error

	// Active Sessions
	ListActiveSessions(ctx context.Context, userID string) ([]ActiveSession, error)
	RevokeSession(ctx context.Context, userID string, tokenID string) error

	// Security history / audit logs
	ListSecurityHistory(ctx context.Context, userID string) ([]SecurityHistoryEntry, error)
	WriteSecurityLog(ctx context.Context, userID string, action string, ip string) error

	// Profile Settings Bridge
	GetProfileSettings(ctx context.Context, userID string) (username, customURL, profileVisibility string, fieldVisibility map[string]string, openToWork, referralEligible, willingToMentor bool, err error)
	UpdateProfileSettings(ctx context.Context, userID string, username, customURL, profileVisibility string, fieldVisibility map[string]string, openToWork, referralEligible, willingToMentor bool) error
}
