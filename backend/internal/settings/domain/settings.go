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
		Version:     1,
	}
}
