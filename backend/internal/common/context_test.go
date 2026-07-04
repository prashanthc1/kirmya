package common

import (
	"context"
	"testing"
)

func TestContextWithUserID(t *testing.T) {
	ctx := context.Background()
	withID := ContextWithUserID(ctx, "user-123")

	if got := UserIDFromContext(withID); got != "user-123" {
		t.Fatalf("expected user ID %q, got %q", "user-123", got)
	}

	if got := UserIDFromContext(ctx); got != "" {
		t.Fatalf("expected empty user ID from background context, got %q", got)
	}
}

func TestContextWithAuthUser(t *testing.T) {
	ctx := context.Background()
	withUser := ContextWithAuthUser(ctx, AuthUser{
		ID:    "user-123",
		Email: "test@example.com",
		Role:  "admin",
	})

	if got := UserIDFromContext(withUser); got != "user-123" {
		t.Fatalf("expected user ID %q, got %q", "user-123", got)
	}
	if got := UserEmailFromContext(withUser); got != "test@example.com" {
		t.Fatalf("expected email %q, got %q", "test@example.com", got)
	}
	if got := UserRoleFromContext(withUser); got != "admin" {
		t.Fatalf("expected role %q, got %q", "admin", got)
	}
}
