package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/identity/infrastructure/crypto"
	"workspace-app/internal/identity/infrastructure/jwtauth"
)

type fakeTOTP struct {
	secret string
}

func (f *fakeTOTP) Generate(email string) (string, string, error) {
	return f.secret, "otpauth://totp/Kirmya:" + email + "?secret=" + f.secret, nil
}

func (f *fakeTOTP) Validate(secretEnc, code string) bool {
	return secretEnc == f.secret && code == "123456"
}

func TestMFAReplayPrevention(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-secret")
	t.Setenv("EMAIL_VERIFICATION_REQUIRED", "false")

	users := newFakeUsers()
	cache := newFakeCache()

	s := NewService(Deps{
		Users:   users,
		Refresh: newFakeRefresh(),
		Verif:   newFakeVerif(),
		MFA:     newFakeMFA(),
		Audit:   noopAudit{},
		Hasher:  crypto.NewArgon2Hasher(),
		Tokens:  jwtauth.NewFactory(),
		TOTP:    &fakeTOTP{secret: "my-test-secret"},
		Cache:   cache,
		Mailer:  newFakeMailer(),
		Events:  noopEvents{},
	})

	ctx := context.Background()

	// 1. Create a user
	reg, err := s.Register(ctx, RegisterInput{Email: "mfa-test@example.com", Password: "supersecret"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	userID := reg.User.ID

	// 2. Setup MFA
	_, err = s.SetupMFA(ctx, userID)
	if err != nil {
		t.Fatalf("setup MFA: %v", err)
	}

	// 3. Confirm MFA (First use of code "123456")
	err = s.ConfirmMFA(ctx, userID, "123456")
	if err != nil {
		t.Fatalf("confirm MFA should succeed, got: %v", err)
	}

	// 4. Replay the same code to confirm MFA again — should be blocked!
	err = s.ConfirmMFA(ctx, userID, "123456")
	if !errors.Is(err, ErrInvalidMFACode) {
		t.Fatalf("expected ErrInvalidMFACode on replay, got: %v", err)
	}

	// 5. Try login. Since the code was spent in ConfirmMFA, it should be blocked.
	_, err = s.Login(ctx, LoginInput{Email: "mfa-test@example.com", Password: "supersecret", Code: "123456"})
	if !errors.Is(err, ErrInvalidMFACode) {
		t.Fatalf("expected login to be blocked due to spent code, got: %v", err)
	}

	// Let's clear the cache to simulate a new time step (new time window)
	cache.items = map[string][]byte{}

	// Login should now succeed
	loginRes, err := s.Login(ctx, LoginInput{Email: "mfa-test@example.com", Password: "supersecret", Code: "123456"})
	if err != nil {
		t.Fatalf("login should succeed after cache cleared, got: %v", err)
	}
	if loginRes.AccessToken == "" {
		t.Fatal("expected access token")
	}

	// Replay login with same code — should be blocked!
	_, err = s.Login(ctx, LoginInput{Email: "mfa-test@example.com", Password: "supersecret", Code: "123456"})
	if !errors.Is(err, ErrInvalidMFACode) {
		t.Fatalf("expected login replay to be blocked, got: %v", err)
	}

	// 6. Disable MFA. Replay code is blocked.
	err = s.DisableMFA(ctx, userID, "123456")
	if !errors.Is(err, ErrInvalidMFACode) {
		t.Fatalf("expected disable replay to be blocked, got: %v", err)
	}

	// Clear cache for a new time step
	cache.items = map[string][]byte{}

	// Disable MFA should succeed
	err = s.DisableMFA(ctx, userID, "123456")
	if err != nil {
		t.Fatalf("disable MFA should succeed, got: %v", err)
	}
}
