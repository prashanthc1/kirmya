package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/identity/domain"
)

func rolePresent(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}

func TestSetMyRoles(t *testing.T) {
	f := newFakeUsers()
	svc := NewService(Deps{Users: f, Events: noopEvents{}})
	ctx := context.Background()

	u := &domain.User{Email: "marcus@example.com", FullName: "Marcus", Status: domain.StatusActive}
	if err := f.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := f.AssignRole(ctx, u.ID, domain.RoleJobSeeker); err != nil {
		t.Fatalf("assign: %v", err)
	}

	// Add mentor alongside the existing job_seeker role.
	pu, err := svc.SetMyRoles(ctx, u.ID, []string{domain.RoleJobSeeker, domain.RoleMentor})
	if err != nil {
		t.Fatalf("set roles: %v", err)
	}
	if !rolePresent(pu.Roles, domain.RoleJobSeeker) || !rolePresent(pu.Roles, domain.RoleMentor) {
		t.Fatalf("expected job_seeker+mentor, got %v", pu.Roles)
	}

	// Drop job_seeker, keep mentor — a self-assignable role can be removed.
	pu, err = svc.SetMyRoles(ctx, u.ID, []string{domain.RoleMentor})
	if err != nil {
		t.Fatalf("set roles: %v", err)
	}
	if rolePresent(pu.Roles, domain.RoleJobSeeker) || !rolePresent(pu.Roles, domain.RoleMentor) {
		t.Fatalf("expected only mentor, got %v", pu.Roles)
	}

	// An existing admin role must be preserved and never stripped via self-serve.
	if err := f.AssignRole(ctx, u.ID, domain.RoleAdmin); err != nil {
		t.Fatalf("assign admin: %v", err)
	}
	pu, err = svc.SetMyRoles(ctx, u.ID, []string{domain.RoleRecruiter})
	if err != nil {
		t.Fatalf("set roles: %v", err)
	}
	if !rolePresent(pu.Roles, domain.RoleAdmin) {
		t.Fatalf("admin role must be preserved, got %v", pu.Roles)
	}
	if !rolePresent(pu.Roles, domain.RoleRecruiter) || rolePresent(pu.Roles, domain.RoleMentor) {
		t.Fatalf("expected recruiter with mentor removed (plus admin), got %v", pu.Roles)
	}
}

func TestSetMyRolesValidation(t *testing.T) {
	f := newFakeUsers()
	svc := NewService(Deps{Users: f, Events: noopEvents{}})
	ctx := context.Background()

	u := &domain.User{Email: "ada@example.com", FullName: "Ada", Status: domain.StatusActive}
	if err := f.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := svc.SetMyRoles(ctx, u.ID, []string{domain.RoleAdmin}); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("expected ErrInvalidRole for admin, got %v", err)
	}
	if _, err := svc.SetMyRoles(ctx, u.ID, []string{"wizard"}); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("expected ErrInvalidRole for unknown role, got %v", err)
	}
	if _, err := svc.SetMyRoles(ctx, u.ID, nil); !errors.Is(err, ErrNoRoles) {
		t.Fatalf("expected ErrNoRoles for empty set, got %v", err)
	}
}
