package domain

import (
	"testing"
	"time"
)

func TestUserHasPassword(t *testing.T) {
	if (&User{}).HasPassword() {
		t.Error("empty hash should report no password (OAuth-only account)")
	}
	if !(&User{PasswordHash: "argon2id$..."}).HasPassword() {
		t.Error("a set hash should report a password")
	}
}

func TestUserIsActive(t *testing.T) {
	cases := map[string]bool{
		StatusActive:      true,
		StatusSuspended:   false,
		StatusDeactivated: false,
		"":                false,
	}
	for status, want := range cases {
		if got := (&User{Status: status}).IsActive(); got != want {
			t.Errorf("IsActive(%q) = %v, want %v", status, got, want)
		}
	}
}

func TestUserHasRole(t *testing.T) {
	u := &User{Roles: []string{RoleJobSeeker, RoleMentor}}
	if !u.HasRole(RoleMentor) {
		t.Error("expected user to hold the mentor role")
	}
	if u.HasRole(RoleAdmin) {
		t.Error("did not expect the user to hold the admin role")
	}
	if (&User{}).HasRole(RoleJobSeeker) {
		t.Error("a user with no roles should hold no role")
	}
}

func TestRefreshTokenIsUsable(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)
	revoked := now.Add(-time.Minute)
	replacement := "next-token-id"

	tests := []struct {
		name string
		tok  RefreshToken
		want bool
	}{
		{"fresh", RefreshToken{ExpiresAt: future}, true},
		{"expired", RefreshToken{ExpiresAt: past}, false},
		{"revoked", RefreshToken{ExpiresAt: future, RevokedAt: &revoked}, false},
		{"rotated/replaced", RefreshToken{ExpiresAt: future, ReplacedBy: &replacement}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.tok.IsUsable(now); got != tc.want {
				t.Errorf("IsUsable() = %v, want %v", got, tc.want)
			}
		})
	}
}
