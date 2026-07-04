package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/identity/domain"
	"workspace-app/internal/identity/infrastructure/crypto"
	"workspace-app/internal/identity/infrastructure/jwtauth"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	t.Setenv("JWT_SECRET", "unit-test-secret")
	// Tests that exercise login right after register opt out of the email-
	// verification gate (the documented dev/seed default); the gate itself is
	// covered by TestLoginBlockedUntilVerified.
	t.Setenv("EMAIL_VERIFICATION_REQUIRED", "false")
	return NewService(Deps{
		Users:   newFakeUsers(),
		Refresh: newFakeRefresh(),
		Verif:   newFakeVerif(),
		Audit:   noopAudit{},
		Hasher:  crypto.NewArgon2Hasher(),
		Tokens:  jwtauth.NewFactory(),
		Mailer:  newFakeMailer(),
		Events:  noopEvents{},
	})
}

func TestRegisterThenLogin(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()

	if _, err := s.Register(ctx, RegisterInput{Email: "Alex@Example.com", Password: "supersecret", FullName: "Alex"}); err != nil {
		t.Fatalf("register: %v", err)
	}

	res, err := s.Login(ctx, LoginInput{Email: "alex@example.com", Password: "supersecret"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatal("expected tokens on login")
	}
	if res.User.Roles[0] != domain.RoleJobSeeker {
		t.Errorf("default role = %v", res.User.Roles)
	}
}

// Email delivery is best-effort: when the mailer fails (e.g. the production
// log-mailer that refuses to send), registration must still create the account
// rather than returning an error to the caller.
func TestRegisterSucceedsWhenMailerFails(t *testing.T) {
	s := newTestService(t)
	s.mailer.(*fakeMailer).verifyErr = errors.New("no production mailer configured")
	ctx := context.Background()

	reg, err := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret", FullName: "A"})
	if err != nil {
		t.Fatalf("register should not fail on mailer error: %v", err)
	}
	if reg.User.ID == "" {
		t.Fatal("expected a created user despite the mailer failure")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, _ = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"})

	_, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "nope"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRefreshRotationAndReuseDetection(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, _ = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"})
	login, _ := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "supersecret"})

	// First refresh rotates the token.
	rot, err := s.Refresh(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if rot.RefreshToken == login.RefreshToken {
		t.Fatal("refresh token should rotate")
	}

	// New token works.
	if _, err := s.Refresh(ctx, rot.RefreshToken); err != nil {
		t.Fatalf("second refresh: %v", err)
	}

	// Reusing the original (already-rotated) token must fail AND revoke family,
	// so the latest token is now invalid too.
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected reuse rejection, got %v", err)
	}
	if _, err := s.Refresh(ctx, rot.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatal("expected family revocation after reuse")
	}
}

func TestLoginBlockedUntilVerified(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-secret")
	t.Setenv("EMAIL_VERIFICATION_REQUIRED", "true")
	s := NewService(Deps{
		Users:   newFakeUsers(),
		Refresh: newFakeRefresh(),
		Verif:   newFakeVerif(),
		Audit:   noopAudit{},
		Hasher:  crypto.NewArgon2Hasher(),
		Tokens:  jwtauth.NewFactory(),
		Mailer:  newFakeMailer(),
		Events:  noopEvents{},
	})
	ctx := context.Background()

	if _, err := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Unverified login is rejected.
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "supersecret"}); !errors.Is(err, ErrEmailNotVerified) {
		t.Fatalf("expected ErrEmailNotVerified, got %v", err)
	}
	// After verifying, login succeeds.
	raw := s.mailer.(*fakeMailer).lastVerifyRaw
	if err := s.VerifyEmail(ctx, raw); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "supersecret"}); err != nil {
		t.Fatalf("login after verify: %v", err)
	}
}

// --- Logout tests ---

// TestLogoutRevokesToken verifies that after a successful logout the refresh
// token can no longer be used to obtain a new session.
func TestLogoutRevokesToken(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, _ = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"})
	login, _ := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "supersecret"})

	if err := s.Logout(ctx, login.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}

	// The revoked token must be rejected on the next Refresh attempt.
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials after logout, got %v", err)
	}
}

// TestLogoutEmptyTokenIsNoOp verifies that passing an empty string to Logout
// is a no-op and returns nil (idempotency / no leak on missing cookie).
func TestLogoutEmptyTokenIsNoOp(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()

	if err := s.Logout(ctx, ""); err != nil {
		t.Fatalf("expected no error on empty token, got %v", err)
	}
}

// TestLogoutUnknownTokenIsNoOp verifies that presenting an unrecognised token
// (never issued or already expired/deleted from storage) is silently ignored
// and returns nil — callers must not infer token existence from the result.
func TestLogoutUnknownTokenIsNoOp(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()

	if err := s.Logout(ctx, "this-token-was-never-issued"); err != nil {
		t.Fatalf("expected no error on unknown token, got %v", err)
	}
}

// TestLogoutThenRefreshRejected is an end-to-end sequence that confirms a
// logged-out token cannot be silently rotated back via Refresh: the Refresh
// path must return ErrInvalidCredentials, NOT silently issue new tokens.
func TestLogoutThenRefreshRejected(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, _ = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"})
	login, _ := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "supersecret"})

	// Rotate once so we have a second-generation token.
	rotated, err := s.Refresh(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("first refresh: %v", err)
	}

	// Log out using the second-generation token.
	if err := s.Logout(ctx, rotated.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}

	// Neither the original nor the rotated token should produce a new session.
	if _, err := s.Refresh(ctx, rotated.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected rejection of logged-out token, got %v", err)
	}
	// The original (already-replaced) token was also revoked by MarkReplaced
	// during rotation, so it must also be rejected.
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected rejection of replaced token, got %v", err)
	}
}

func TestVerifyEmailFlow(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, _ = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "supersecret"})

	raw := s.mailer.(*fakeMailer).lastVerifyRaw
	if raw == "" {
		t.Fatal("expected a verification token to be emailed")
	}
	if err := s.VerifyEmail(ctx, raw); err != nil {
		t.Fatalf("verify: %v", err)
	}
	u, _ := s.users.GetByEmail(ctx, "a@b.com")
	if !u.EmailVerified {
		t.Fatal("email should be verified")
	}
}
