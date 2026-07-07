package domain

import (
	"context"
	"errors"
	"time"
)

// Sentinel errors returned by ports and mapped to HTTP statuses in the api layer.
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmailTaken     = errors.New("email already registered")
	ErrTokenNotFound  = errors.New("token not found")
	ErrOptimisticLock = errors.New("stale update (version mismatch)")
)

// UserRepository persists the User aggregate and its role assignments.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error // optimistic-locked on Version
	AssignRole(ctx context.Context, userID, roleName string) error
	RemoveRole(ctx context.Context, userID, roleName string) error
	GetRoles(ctx context.Context, userID string) ([]string, error)
	SetEmailVerified(ctx context.Context, userID string) error
	SetPasswordHash(ctx context.Context, userID, hash string) error
	SetMFAEnabled(ctx context.Context, userID string, enabled bool) error
	UpdateLastLogin(ctx context.Context, userID string) error

	// Directory (people search) read projections.
	Search(ctx context.Context, query string, limit int) ([]DirectoryEntry, error)
	GetDirectory(ctx context.Context, id string) (*DirectoryEntry, error)
}

// RefreshTokenRepository stores hashed refresh tokens and supports rotation
// and family revocation (reuse detection).
type RefreshTokenRepository interface {
	Store(ctx context.Context, t *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	MarkReplaced(ctx context.Context, id, replacedBy string) error
	Revoke(ctx context.Context, id string) error
	RevokeFamily(ctx context.Context, familyID string) error
	// RevokeAllForUser revokes every active refresh token for a user. Used after
	// a password change to force re-authentication on all devices.
	RevokeAllForUser(ctx context.Context, userID string) error
}

// VerificationRepository stores single-use email-verification and
// password-reset tokens (hashed).
type VerificationRepository interface {
	StoreEmailToken(ctx context.Context, userID, hash string, expiresAt time.Time) error
	ConsumeEmailToken(ctx context.Context, hash string) (userID string, err error)
	StorePasswordToken(ctx context.Context, userID, hash string, expiresAt time.Time) error
	ConsumePasswordToken(ctx context.Context, hash string) (userID string, err error)
}

// OAuthRepository links external identities to local users.
type OAuthRepository interface {
	FindUserIDByProvider(ctx context.Context, provider, providerUID string) (userID string, found bool, err error)
	Link(ctx context.Context, userID, provider, providerUID string) error
}

// MFARepository stores TOTP credentials.
type MFARepository interface {
	Upsert(ctx context.Context, c *MFACredential) error
	Get(ctx context.Context, userID string) (*MFACredential, error)
	Confirm(ctx context.Context, userID string) error
}

// AuditRepository appends to the audit log.
type AuditRepository interface {
	Record(ctx context.Context, actorID, action, targetType, targetID string, metadata map[string]any, ip string) error
}

// Mailer sends transactional auth emails (verification, reset).
type Mailer interface {
	SendVerificationEmail(ctx context.Context, email, rawToken string) error
	SendPasswordResetEmail(ctx context.Context, email, rawToken string) error
}

// EventPublisher publishes domain events onto the (in-process, NATS-ready) bus.
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

// PasswordHasher hashes and verifies passwords (Argon2id).
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, encoded string) (bool, error)
}

// TokenFactory issues access tokens and generates/Hashes opaque tokens.
type TokenFactory interface {
	IssueAccessToken(userID, email string, roles []string) (token string, expiresIn int, err error)
	GenerateOpaqueToken() (raw string, hash string, err error)
	HashOpaqueToken(raw string) string
}

// TOTP provisions and validates time-based one-time-password MFA secrets.
// The returned secretEnc is encrypted-at-rest; otpauthURL embeds the plaintext
// secret for QR enrollment and must never be persisted.
type TOTP interface {
	Generate(accountEmail string) (secretEnc, otpauthURL string, err error)
	Validate(secretEnc, code string) bool
}

// Cache is a decoupled caching port for identity context needs (e.g. spent code tracking).
type Cache interface {
	Get(ctx context.Context, key string) (value []byte, ok bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
}
