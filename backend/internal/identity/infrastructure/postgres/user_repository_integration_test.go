//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/identity/domain"
	"workspace-app/internal/testsupport"
)

func TestUserRepository_CRUDAndRoles(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	u := &domain.User{Email: "asha@cb.test", PasswordHash: "hash", FullName: "Asha Rao", EmailVerified: true, Status: domain.StatusActive}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if u.ID == "" || u.Version != 1 {
		t.Fatalf("expected id + version=1, got %+v", u)
	}

	// Duplicate email is rejected.
	dup := &domain.User{Email: "ASHA@cb.test", FullName: "dupe"}
	if err := repo.Create(ctx, dup); err != domain.ErrEmailTaken {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}

	// Get by email (case-insensitive).
	got, err := repo.GetByEmail(ctx, "Asha@CB.test")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if got.ID != u.ID || got.FullName != "Asha Rao" {
		t.Fatalf("unexpected user %+v", got)
	}

	// Roles round-trip.
	if err := repo.AssignRole(ctx, u.ID, domain.RoleRecruiter); err != nil {
		t.Fatalf("assign role: %v", err)
	}
	if err := repo.AssignRole(ctx, u.ID, domain.RoleRecruiter); err != nil {
		t.Fatalf("assign role (idempotent): %v", err)
	}
	roles, err := repo.GetRoles(ctx, u.ID)
	if err != nil {
		t.Fatalf("get roles: %v", err)
	}
	if len(roles) != 1 || roles[0] != domain.RoleRecruiter {
		t.Fatalf("expected [recruiter], got %v", roles)
	}

	got, err = repo.GetByID(ctx, u.ID)
	if err != nil || !got.HasRole(domain.RoleRecruiter) {
		t.Fatalf("expected recruiter role on reload, got %+v err=%v", got, err)
	}
}

func TestUserRepository_Search(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	for _, u := range []*domain.User{
		{Email: "ben@cb.test", FullName: "Ben Carter", Status: domain.StatusActive},
		{Email: "carla@cb.test", FullName: "Carla Mendes", Status: domain.StatusActive},
	} {
		if err := repo.Create(ctx, u); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	hits, err := repo.Search(ctx, "carl", 10)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 1 || hits[0].FullName != "Carla Mendes" {
		t.Fatalf("expected Carla, got %+v", hits)
	}
}
