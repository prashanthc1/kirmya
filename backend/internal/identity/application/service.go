// Package application implements the identity use cases (commands/queries). It
// depends only on domain ports; infrastructure adapters are injected in
// module.go. Other modules consume the exported Service interface.
package application

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"workspace-app/internal/identity/domain"
)

// Errors surfaced to the api layer.
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountInactive    = errors.New("account is not active")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrMFARequired        = errors.New("mfa code required")
	ErrInvalidMFACode     = errors.New("invalid mfa code")
	ErrUnknownProvider    = errors.New("unknown oauth provider")
	// ErrInvalidRole is returned when a self-serve role update names a role the
	// caller may not assign to themselves (unknown, or admin).
	ErrInvalidRole = errors.New("invalid or non-self-assignable role")
	// ErrNoRoles is returned when a self-serve role update would leave the
	// account with no self-assignable roles.
	ErrNoRoles = errors.New("at least one role is required")
)

// OAuthProvider is an application port for an external identity provider.
type OAuthProvider interface {
	Name() string
	AuthURL(state string) string
	Exchange(ctx context.Context, code string) (providerUID, email, fullName string, err error)
}

// PublicUser is the safe projection returned to callers (no secrets).
type PublicUser struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	FullName      string   `json:"full_name"`
	EmailVerified bool     `json:"email_verified"`
	MFAEnabled    bool     `json:"mfa_enabled"`
	Roles         []string `json:"roles"`
}

// AuthResult carries tokens after a successful authentication. RefreshToken is
// the raw token the api layer sets as an httpOnly cookie.
type AuthResult struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
	User         PublicUser
	MFARequired  bool
}

// TxManager coordinates database transactions.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// Service is the identity use-case service.
type Service struct {
	users      domain.UserRepository
	refresh    domain.RefreshTokenRepository
	verif      domain.VerificationRepository
	oauth      domain.OAuthRepository
	mfa        domain.MFARepository
	audit      domain.AuditRepository
	hasher     domain.PasswordHasher
	tokens     domain.TokenFactory
	totp       domain.TOTP
	mailer     domain.Mailer
	events     domain.EventPublisher
	providers  map[string]OAuthProvider
	tx         TxManager
	refreshTTL time.Duration

	// requireVerify, when true, blocks login for password accounts whose email
	// is not yet verified (H3). Controlled by EMAIL_VERIFICATION_REQUIRED
	// (default false). When off, accounts can log in immediately and verify
	// later from profile settings. OAuth accounts are always created verified.
	requireVerify bool

	// mfaThrottle limits TOTP attempts per account to defeat brute force of the
	// 6-digit code during the login MFA step (L3).
	mfaThrottle *attemptLimiter
}

// Deps bundles the injected dependencies.
type Deps struct {
	Users     domain.UserRepository
	Refresh   domain.RefreshTokenRepository
	Verif     domain.VerificationRepository
	OAuth     domain.OAuthRepository
	MFA       domain.MFARepository
	Audit     domain.AuditRepository
	Hasher    domain.PasswordHasher
	Tokens    domain.TokenFactory
	TOTP      domain.TOTP
	Mailer    domain.Mailer
	Events    domain.EventPublisher
	Providers map[string]OAuthProvider
	Tx        TxManager
}

func NewService(d Deps) *Service {
	return &Service{
		users: d.Users, refresh: d.Refresh, verif: d.Verif, oauth: d.OAuth,
		mfa: d.MFA, audit: d.Audit, hasher: d.Hasher, tokens: d.Tokens,
		totp: d.TOTP, mailer: d.Mailer, events: d.Events, providers: d.Providers,
		tx:            d.Tx,
		refreshTTL:    refreshTTL(),
		requireVerify: emailVerificationRequired(),
		mfaThrottle:   newAttemptLimiter(5, 15*time.Minute),
	}
}

// emailVerificationRequired returns whether unverified password accounts are
// blocked at login. Defaults to false so login stays frictionless and users
// verify later from profile settings; set EMAIL_VERIFICATION_REQUIRED=true to
// re-enable the hard gate (e.g. for stricter deployments).
func emailVerificationRequired() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("EMAIL_VERIFICATION_REQUIRED")), "true")
}

func refreshTTL() time.Duration {
	if v := os.Getenv("JWT_REFRESH_TTL"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil {
			return time.Duration(secs) * time.Second
		}
	}
	return 30 * 24 * time.Hour
}

const verificationTTL = 24 * time.Hour
const resetTTL = time.Hour

// RegisterInput is the register command.
type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Role     string
	IP       string
}

// Register creates a password account, assigns a role, sends a verification
// email, and (unless the verification gate is on) logs the new user in
// immediately by returning a fresh session. When EMAIL_VERIFICATION_REQUIRED is
// true the AuthResult carries only the User (no tokens) so the caller routes
// them to verify first.
func (s *Service) Register(ctx context.Context, in RegisterInput) (AuthResult, error) {
	var res AuthResult
	err := s.runInTx(ctx, func(ctx context.Context) error {
		email := strings.TrimSpace(strings.ToLower(in.Email))
		hash, err := s.hasher.Hash(in.Password)
		if err != nil {
			return err
		}
		u := &domain.User{Email: email, PasswordHash: hash, FullName: strings.TrimSpace(in.FullName), Status: domain.StatusActive}
		if err := s.users.Create(ctx, u); err != nil {
			return err
		}
		role := in.Role
		if role == "" {
			role = domain.RoleJobSeeker
		}
		if err := s.users.AssignRole(ctx, u.ID, role); err != nil {
			return err
		}
		u.Roles = []string{role}

		// Email delivery is best-effort: a missing/failed mailer must not block
		// account creation. The user can re-request verification later
		// (ResendVerification). Login enforcement is governed separately by
		// EMAIL_VERIFICATION_REQUIRED.
		if err := s.sendVerification(ctx, u); err != nil {
			log.Printf("register: verification email not sent for user %s: %v", u.ID, err)
		}
		s.record(ctx, u.ID, "user.register", "user", u.ID, in.IP)
		_ = s.events.Publish(ctx, domain.EventUserRegistered, u.ID, map[string]any{"email": u.Email})

		// Auto-login: issue a session unless the verification gate is active, in
		// which case the account exists but cannot log in until verified.
		if s.requireVerify && !u.EmailVerified {
			res = AuthResult{User: toPublic(u)}
			return nil
		}
		var errSess error
		res, errSess = s.issueSession(ctx, u)
		return errSess
	})
	return res, err
}

func (s *Service) runInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if s.tx == nil {
		return fn(ctx)
	}
	return s.tx.RunInTx(ctx, fn)
}

// LoginInput is the login command. A non-empty Code completes an MFA challenge.
type LoginInput struct {
	Email     string
	Password  string
	Code      string
	IP        string
	UserAgent string
}

// Login authenticates with email + password (and TOTP if MFA is enabled).
func (s *Service) Login(ctx context.Context, in LoginInput) (AuthResult, error) {
	u, err := s.users.GetByEmail(ctx, in.Email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}
	if !u.IsActive() {
		return AuthResult{}, ErrAccountInactive
	}
	if !u.HasPassword() {
		return AuthResult{}, ErrInvalidCredentials
	}
	ok, err := s.hasher.Verify(in.Password, u.PasswordHash)
	if err != nil || !ok {
		return AuthResult{}, ErrInvalidCredentials
	}
	// Block unverified password accounts (H3). The email claim must not be
	// trusted by downstream features unless ownership was proven.
	if s.requireVerify && !u.EmailVerified {
		return AuthResult{}, ErrEmailNotVerified
	}
	if u.MFAEnabled {
		if in.Code == "" {
			return AuthResult{MFARequired: true}, nil
		}
		// Rate-limit code submissions per account (L3) before validating.
		if s.mfaThrottle != nil && !s.mfaThrottle.allow(u.ID) {
			return AuthResult{}, ErrInvalidMFACode
		}
		cred, err := s.mfa.Get(ctx, u.ID)
		if err != nil || !s.totp.Validate(cred.SecretEnc, in.Code) {
			return AuthResult{}, ErrInvalidMFACode
		}
		if s.mfaThrottle != nil {
			s.mfaThrottle.reset(u.ID) // successful code clears the counter
		}
	}
	s.record(ctx, u.ID, "user.login", "user", u.ID, in.IP)
	_ = s.events.Publish(ctx, domain.EventUserLoggedIn, u.ID, nil)
	return s.issueSession(ctx, u)
}

// issueSession mints an access token + a new refresh-token family.
func (s *Service) issueSession(ctx context.Context, u *domain.User) (AuthResult, error) {
	access, expiresIn, err := s.tokens.IssueAccessToken(u.ID, u.Email, u.Roles)
	if err != nil {
		return AuthResult{}, err
	}
	raw, hash, err := s.tokens.GenerateOpaqueToken()
	if err != nil {
		return AuthResult{}, err
	}
	rt := &domain.RefreshToken{UserID: u.ID, TokenHash: hash, FamilyID: newFamilyID(), ExpiresAt: time.Now().Add(s.refreshTTL)}
	if err := s.refresh.Store(ctx, rt); err != nil {
		return AuthResult{}, err
	}
	_ = s.users.UpdateLastLogin(ctx, u.ID)
	return AuthResult{AccessToken: access, ExpiresIn: expiresIn, RefreshToken: raw, User: toPublic(u)}, nil
}

// Refresh rotates a refresh token. Reuse of an already-rotated token revokes
// the entire family (theft detection).
func (s *Service) Refresh(ctx context.Context, rawRefresh string) (AuthResult, error) {
	hash := s.tokens.HashOpaqueToken(rawRefresh)
	rt, err := s.refresh.FindByHash(ctx, hash)
	if errors.Is(err, domain.ErrTokenNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}
	if !rt.IsUsable(time.Now()) {
		// A replaced token being reused means the family may be compromised.
		if rt.ReplacedBy != nil {
			_ = s.refresh.RevokeFamily(ctx, rt.FamilyID)
		}
		return AuthResult{}, ErrInvalidCredentials
	}
	u, err := s.users.GetByID(ctx, rt.UserID)
	if err != nil {
		return AuthResult{}, err
	}
	access, expiresIn, err := s.tokens.IssueAccessToken(u.ID, u.Email, u.Roles)
	if err != nil {
		return AuthResult{}, err
	}
	raw, newHash, err := s.tokens.GenerateOpaqueToken()
	if err != nil {
		return AuthResult{}, err
	}
	next := &domain.RefreshToken{UserID: u.ID, TokenHash: newHash, FamilyID: rt.FamilyID, ExpiresAt: time.Now().Add(s.refreshTTL)}
	if err := s.refresh.Store(ctx, next); err != nil {
		return AuthResult{}, err
	}
	if err := s.refresh.MarkReplaced(ctx, rt.ID, next.ID); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{AccessToken: access, ExpiresIn: expiresIn, RefreshToken: raw, User: toPublic(u)}, nil
}

// Logout revokes the presented refresh token.
func (s *Service) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	rt, err := s.refresh.FindByHash(ctx, s.tokens.HashOpaqueToken(rawRefresh))
	if errors.Is(err, domain.ErrTokenNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.refresh.Revoke(ctx, rt.ID)
}

// Me returns the current user's public projection.
func (s *Service) Me(ctx context.Context, userID string) (PublicUser, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return PublicUser{}, err
	}
	return toPublic(u), nil
}

// SetMyRoles reconciles the caller's self-assignable roles to the desired set.
// Roles outside domain.SelfAssignableRoles are rejected; any non-self-assignable
// role the user already holds (notably admin) is preserved and never removed
// through this path. Returns the updated user.
func (s *Service) SetMyRoles(ctx context.Context, userID string, desired []string) (PublicUser, error) {
	want := map[string]bool{}
	for _, role := range desired {
		if !domain.SelfAssignableRoles[role] {
			return PublicUser{}, ErrInvalidRole
		}
		want[role] = true
	}
	if len(want) == 0 {
		return PublicUser{}, ErrNoRoles
	}

	current, err := s.users.GetRoles(ctx, userID)
	if err != nil {
		return PublicUser{}, err
	}
	have := map[string]bool{}
	for _, role := range current {
		have[role] = true
	}

	for role := range want {
		if !have[role] {
			if err := s.users.AssignRole(ctx, userID, role); err != nil {
				return PublicUser{}, err
			}
		}
	}
	for role := range have {
		if domain.SelfAssignableRoles[role] && !want[role] {
			if err := s.users.RemoveRole(ctx, userID, role); err != nil {
				return PublicUser{}, err
			}
		}
	}

	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return PublicUser{}, err
	}
	if s.events != nil {
		_ = s.events.Publish(ctx, domain.EventUserRolesUpdated, userID, map[string]any{"roles": u.Roles})
	}
	return toPublic(u), nil
}

// GetUser satisfies the cross-module Service interface.
func (s *Service) GetUser(ctx context.Context, userID string) (PublicUser, error) {
	return s.Me(ctx, userID)
}

// SearchUsers returns directory entries matching the query (name or email).
func (s *Service) SearchUsers(ctx context.Context, query string) ([]domain.DirectoryEntry, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []domain.DirectoryEntry{}, nil
	}
	return s.users.Search(ctx, query, 20)
}

// Directory returns a single user's public directory entry.
func (s *Service) Directory(ctx context.Context, id string) (*domain.DirectoryEntry, error) {
	return s.users.GetDirectory(ctx, id)
}

func (s *Service) record(ctx context.Context, actorID, action, targetType, targetID, ip string) {
	_ = s.audit.Record(ctx, actorID, action, targetType, targetID, nil, ip)
}

// LogoutAll revokes every refresh token for the user — "sign out on all devices".
func (s *Service) LogoutAll(ctx context.Context, userID string) error {
	if err := s.refresh.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}
	s.record(ctx, userID, "user.logout_all", "user", userID, "")
	return nil
}

// DeactivateAccount soft-closes the caller's account (status=deactivated) and
// revokes all sessions. Deactivated accounts can no longer authenticate (see
// User.IsActive); reactivation is an admin/support action.
func (s *Service) DeactivateAccount(ctx context.Context, userID string) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.Status == domain.StatusDeactivated {
		return nil
	}
	u.Status = domain.StatusDeactivated
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	if err := s.refresh.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}
	s.record(ctx, userID, "user.deactivate", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventAccountDeactivated, userID, nil)
	return nil
}

func toPublic(u *domain.User) PublicUser {
	return PublicUser{ID: u.ID, Email: u.Email, FullName: u.FullName, EmailVerified: u.EmailVerified, MFAEnabled: u.MFAEnabled, Roles: u.Roles}
}
