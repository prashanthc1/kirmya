package application

import (
	"context"
	"errors"
	"testing"
)

func TestResetPasswordRevokesSessions(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()

	if _, err := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "Supersecret1", FullName: "A"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	login, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if err := s.ForgotPassword(ctx, "a@b.com"); err != nil {
		t.Fatalf("forgot: %v", err)
	}
	raw := s.mailer.(*fakeMailer).lastResetRaw
	if raw == "" {
		t.Fatal("expected a reset token to be emailed")
	}
	if err := s.ResetPassword(ctx, raw, "Newpassword1"); err != nil {
		t.Fatalf("reset: %v", err)
	}

	// Existing sessions are revoked by the reset.
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected old session revoked, got %v", err)
	}
	// The new password works; the old one does not.
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("old password should fail, got %v", err)
	}
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Newpassword1"}); err != nil {
		t.Fatalf("login with new password: %v", err)
	}
}
