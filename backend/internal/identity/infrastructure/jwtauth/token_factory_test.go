package jwtauth

import (
	"testing"
)

func TestIssueAndParseAccessToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	f := NewFactory()

	token, expiresIn, err := f.IssueAccessToken("user-123", "a@example.com", []string{"job_seeker", "mentor"})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if expiresIn <= 0 {
		t.Fatalf("expected positive expiry, got %d", expiresIn)
	}

	claims, err := f.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Errorf("subject = %q, want user-123", claims.Subject)
	}
	if claims.Email != "a@example.com" {
		t.Errorf("email = %q", claims.Email)
	}
	if len(claims.Roles) != 2 || claims.Roles[0] != "job_seeker" {
		t.Errorf("roles = %v", claims.Roles)
	}
}

func TestParseRejectsTamperedToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	f := NewFactory()
	token, _, _ := f.IssueAccessToken("u", "e@e.com", nil)

	if _, err := f.Parse(token + "x"); err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestOpaqueTokenHashIsDeterministic(t *testing.T) {
	f := NewFactory()
	raw, hash, err := f.GenerateOpaqueToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if raw == hash {
		t.Fatal("raw and hash should differ")
	}
	if f.HashOpaqueToken(raw) != hash {
		t.Fatal("hash should be reproducible from raw token")
	}
}
