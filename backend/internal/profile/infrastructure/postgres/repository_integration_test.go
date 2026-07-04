//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/profile/domain"
	"workspace-app/internal/testsupport"
)

// TestProfileLifecycle exercises the profile repository end to end against a
// real PostgreSQL: lazy row creation on Get, scalar updates with the optimistic
// version bump, and round-tripping a child collection (experiences) plus skills.
func TestProfileLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	userID := testsupport.InsertUser(t, db, "profile@cb.test", "Pat Profile")

	// Get lazily creates the row and returns empty scalars at the initial
	// version 1 — the codebase-wide optimistic-lock default (see migration
	// 005_create_profile_tables.sql and the users/referrals tables).
	p, err := repo.Get(ctx, userID)
	if err != nil {
		t.Fatalf("get (create): %v", err)
	}
	if p.UserID != userID {
		t.Fatalf("expected profile for %s, got %s", userID, p.UserID)
	}
	if p.Version != 1 {
		t.Fatalf("expected initial version 1, got %d", p.Version)
	}
	if len(p.Experiences) != 0 || len(p.Skills) != 0 {
		t.Fatalf("expected empty child collections, got %d experiences / %d skills",
			len(p.Experiences), len(p.Skills))
	}

	// UpdateScalars persists fields and bumps the version.
	if err := repo.UpdateScalars(ctx, userID, domain.Scalars{
		Headline: "Recovering PM",
		About:    "Back after a career gap.",
		Location: "Remote",
	}); err != nil {
		t.Fatalf("update scalars: %v", err)
	}

	p, err = repo.Get(ctx, userID)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if p.Headline != "Recovering PM" || p.Location != "Remote" {
		t.Fatalf("scalars not persisted: %+v", p)
	}
	if p.Version != 2 {
		t.Fatalf("expected version 2 after one update, got %d", p.Version)
	}

	// AddExperience writes a child row and returns its generated id.
	exp := &domain.WorkExperience{
		Title:     "Product Manager",
		Company:   "Acme",
		StartDate: "2020-01-01",
		EndDate:   "2023-06-01",
		IsCurrent: false,
	}
	if err := repo.AddExperience(ctx, userID, exp); err != nil {
		t.Fatalf("add experience: %v", err)
	}
	if exp.ID == "" {
		t.Fatal("expected experience id to be populated")
	}

	// SetSkills replaces the skill set wholesale.
	if err := repo.SetSkills(ctx, userID, []string{"Strategy", "SQL", "Roadmapping"}); err != nil {
		t.Fatalf("set skills: %v", err)
	}

	p, err = repo.Get(ctx, userID)
	if err != nil {
		t.Fatalf("get after children: %v", err)
	}
	if len(p.Experiences) != 1 || p.Experiences[0].Company != "Acme" {
		t.Fatalf("expected one Acme experience, got %+v", p.Experiences)
	}
	if len(p.Skills) != 3 {
		t.Fatalf("expected 3 skills, got %d (%v)", len(p.Skills), p.Skills)
	}

	// Deleting the experience leaves the collection empty.
	if err := repo.DeleteExperience(ctx, userID, exp.ID); err != nil {
		t.Fatalf("delete experience: %v", err)
	}
	p, err = repo.Get(ctx, userID)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if len(p.Experiences) != 0 {
		t.Fatalf("expected no experiences after delete, got %d", len(p.Experiences))
	}
}
