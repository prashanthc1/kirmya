package application

import (
	"context"
	"errors"
	"testing"
)

func TestChangePasswordRevokesSessionsAndUpdatesCredential(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()

	reg, err := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "Supersecret1", FullName: "A"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	login, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	// Wrong current password is rejected.
	if err := s.ChangePassword(ctx, reg.User.ID, "wrongcurrent", "Newpassword1"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	// Correct change succeeds.
	if err := s.ChangePassword(ctx, reg.User.ID, "Supersecret1", "Newpassword1"); err != nil {
		t.Fatalf("change password: %v", err)
	}

	// Existing sessions are revoked: the old refresh token no longer works.
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected old session revoked, got %v", err)
	}

	// Old password no longer logs in; new password does.
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("old password should fail, got %v", err)
	}
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Newpassword1"}); err != nil {
		t.Fatalf("login with new password: %v", err)
	}
}
