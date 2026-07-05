// Package domain holds the Identity bounded context's entities, value objects,
// domain events and ports (repository/gateway interfaces). It has no dependency
// on any other layer — infrastructure implements its ports; application
// orchestrates it.
package domain

import "time"

// Account status values.
const (
	StatusActive      = "active"
	StatusSuspended   = "suspended"
	StatusDeactivated = "deactivated"
)

// RBAC role names. Mirrors the seeded `roles` table.
const (
	RoleJobSeeker = "job_seeker"
	RoleReferrer  = "referrer"
	RoleMentor    = "mentor"
	RoleRecruiter = "recruiter"
	RoleAdmin     = "admin"
)

// SelfAssignableRoles are the roles a user may grant or remove for themselves
// via PUT /users/me/roles. Admin and recruiter are intentionally excluded — privilege
// elevation remains an admin-only action.
var SelfAssignableRoles = map[string]bool{
	RoleJobSeeker: true,
	RoleReferrer:  true,
	RoleMentor:    true,
}

// OAuth provider names.
const (
	ProviderGoogle   = "google"
	ProviderLinkedIn = "linkedin"
)

// User is the aggregate root of the identity context.
type User struct {
	ID            string
	Email         string
	PasswordHash  string // empty for OAuth-only accounts
	FullName      string
	EmailVerified bool
	Status        string
	MFAEnabled    bool
	Roles         []string
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
}

// HasPassword reports whether the account can authenticate with a password.
func (u *User) HasPassword() bool { return u.PasswordHash != "" }

// IsActive reports whether the account may authenticate.
func (u *User) IsActive() bool { return u.Status == StatusActive }

// HasRole reports whether the user holds the given role.
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// RefreshToken is a single token in a rotation family. Reuse of a rotated
// (replaced) token signals theft and revokes the whole family.
type RefreshToken struct {
	ID         string
	UserID     string
	TokenHash  string
	FamilyID   string
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *string
}

// IsUsable reports whether the token can still be exchanged.
func (t *RefreshToken) IsUsable(now time.Time) bool {
	return t.RevokedAt == nil && t.ReplacedBy == nil && now.Before(t.ExpiresAt)
}

// DirectoryEntry is a public, searchable read projection of a user (joined with
// their profile headline/photo) — used by the people directory.
type DirectoryEntry struct {
	ID       string
	FullName string
	Email    string
	Headline string
	PhotoURL string
}

// MFACredential is a user's TOTP secret (stored encrypted at rest).
type MFACredential struct {
	UserID    string
	SecretEnc string
	Confirmed bool
}
