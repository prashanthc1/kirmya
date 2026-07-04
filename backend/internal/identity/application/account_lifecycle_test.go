package application

import (
	"context"
	"errors"
	"testing"
)

func TestLogoutAllRevokesEverySession(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	reg, _ := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "Supersecret1"})
	login, _ := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"})

	if err := s.LogoutAll(ctx, reg.User.ID); err != nil {
		t.Fatalf("logout-all: %v", err)
	}
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected session revoked, got %v", err)
	}
}

func TestDeactivateAccountRevokesAndBlocksLogin(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	reg, _ := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "Supersecret1"})
	login, _ := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"})

	if err := s.DeactivateAccount(ctx, reg.User.ID); err != nil {
		t.Fatalf("deactivate: %v", err)
	}
	if _, err := s.Refresh(ctx, login.RefreshToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected sessions revoked, got %v", err)
	}
	if _, err := s.Login(ctx, LoginInput{Email: "a@b.com", Password: "Supersecret1"}); !errors.Is(err, ErrAccountInactive) {
		t.Fatalf("expected ErrAccountInactive after deactivation, got %v", err)
	}
}
