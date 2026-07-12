// Package domain holds the Settings bounded context: a single per-user
// preferences aggregate spanning general, privacy, notification and
// security-preference settings.
package domain

import (
	"errors"
	"time"
)

// Sentinel errors returned by ports and the service, mapped to HTTP statuses in api/.
var (
	ErrNotFound       = errors.New("settings not found")
	ErrOptimisticLock = errors.New("stale update (version mismatch)")
)

// Enum value sets. Persisted columns are CHECK-constrained to match these.
const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"

	DigestOff    = "off"
	DigestDaily  = "daily"
	DigestWeekly = "weekly"

	VisibilityPublic  = "public"
	VisibilityNetwork = "network"
	VisibilityPrivate = "private"

	MessagesEveryone = "everyone"
	MessagesNetwork  = "network"
	MessagesNone     = "none"
)

var (
	ValidThemes        = map[string]bool{ThemeLight: true, ThemeDark: true, ThemeSystem: true}
	ValidDigests       = map[string]bool{DigestOff: true, DigestDaily: true, DigestWeekly: true}
	ValidVisibilities  = map[string]bool{VisibilityPublic: true, VisibilityNetwork: true, VisibilityPrivate: true}
	ValidMessagePolicy = map[string]bool{MessagesEveryone: true, MessagesNetwork: true, MessagesNone: true}
)

// NotificationPrefs groups the per-channel x per-category toggles.
type NotificationPrefs struct {
	EmailJobs       bool
	EmailMentorship bool
	EmailMessages   bool
	EmailReferrals  bool
	InAppJobs       bool
	InAppMentorship bool
	InAppMessages   bool
	InAppReferrals  bool
}

// UserSettings is the preferences aggregate for a single user. Version powers
// optimistic locking on Update.
type UserSettings struct {
	UserID string

	// General
	Language    string
	Timezone    string
	Theme       string
	EmailDigest string

	// Privacy
	ProfileVisibility string
	ShowEmail         bool
	Discoverable      bool
	AllowMessages     string

	// Notifications
	Notifications NotificationPrefs

	// Security preferences
	LoginAlerts bool

	// Accessibility
	FontSize                        string
	HighContrast                    bool
	ReducedMotion                   bool
	CompactMode                     bool
	DefaultLandingPage              string
	AccessibilityKeyboardNavigation bool
	AccessibilityScreenReader       bool
	AccessibilityFocusIndicators    bool

	// AI Preferences
	EnableAIAssistant         bool
	AIJobRecommendations      bool
	AIResumeSuggestions       bool
	AIRoadmapSuggestions      bool
	AISkillGapAnalysis        bool
	AIInterviewPrep           bool
	AILearningRecommendations bool

	// Learning Preferences
	LearningGoals          []string
	TechnologiesOfInterest []string
	CertificationGoals     []string
	LearningReminders      bool

	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Defaults returns the baseline settings used when a user has no row yet. It
// mirrors the column defaults in migration 017 so reads are consistent whether
// the row was materialised by the DB default or by this helper.
func Defaults(userID string) UserSettings {
	return UserSettings{
		UserID:            userID,
		Language:          "en",
		Timezone:          "UTC",
		Theme:             ThemeSystem,
		EmailDigest:       DigestWeekly,
		ProfileVisibility: VisibilityPublic,
		ShowEmail:         false,
		Discoverable:      true,
		AllowMessages:     MessagesEveryone,
		Notifications: NotificationPrefs{
			EmailJobs: true, EmailMentorship: true, EmailMessages: true, EmailReferrals: true,
			InAppJobs: true, InAppMentorship: true, InAppMessages: true, InAppReferrals: true,
		},
		LoginAlerts: true,

		// Accessibility defaults
		FontSize:                        "medium",
		HighContrast:                    false,
		ReducedMotion:                   false,
		CompactMode:                     false,
		DefaultLandingPage:              "dashboard",
		AccessibilityKeyboardNavigation: false,
		AccessibilityScreenReader:       false,
		AccessibilityFocusIndicators:    false,

		// AI Preferences defaults
		EnableAIAssistant:         true,
		AIJobRecommendations:      true,
		AIResumeSuggestions:       true,
		AIRoadmapSuggestions:      true,
		AISkillGapAnalysis:        true,
		AIInterviewPrep:           true,
		AILearningRecommendations: true,

		// Learning Preferences defaults
		LearningGoals:          []string{},
		TechnologiesOfInterest: []string{},
		CertificationGoals:     []string{},
		LearningReminders:      true,

		Version: 1,
	}
}

// ConnectedAccount represents a linked external provider identity (OAuth).
type ConnectedAccount struct {
	ID          string    `json:"id"`
	Provider    string    `json:"provider"`
	ProviderUID string    `json:"provider_uid"`
	CreatedAt   time.Time `json:"created_at"`
}

// CookieConsent represents the user's cookie preferences.
type CookieConsent struct {
	UserID            string    `json:"user_id"`
	Essential         bool      `json:"essential"`
	Functional        bool      `json:"functional"`
	Analytics         bool      `json:"analytics"`
	AIPersonalization bool      `json:"ai_personalization"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ActiveSession represents a logged-in session mapped to an unexpired refresh token.
type ActiveSession struct {
	ID        string    `json:"id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SecurityHistoryEntry represents a logged security event or audit log.
type SecurityHistoryEntry struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
