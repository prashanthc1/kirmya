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

	// Get lazily creates the row and returns empty scalars at the initial version 1.
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
		Headline:         "Recovering PM",
		About:            "Back after a career gap.",
		Location:         "Remote",
		CareerStatus:     "career_break",
		TransitionReason: "layoff",
		SalaryMin:        8000,
		SalaryMax:        12000,
		SalaryCurrency:   "AED",
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
	if p.TransitionReason != "layoff" || p.SalaryMin != 8000 || p.SalaryCurrency != "AED" {
		t.Fatalf("decrypted fields mismatch: %+v", p)
	}
	if p.Version != 2 {
		t.Fatalf("expected version 2 after one update, got %d", p.Version)
	}

	// Assert Encryption at rest in DB
	var transEnc, salMinEnc string
	err = db.QueryRow(`SELECT transition_reason_enc, salary_min_enc FROM profiles WHERE user_id = $1`, userID).Scan(&transEnc, &salMinEnc)
	if err != nil {
		t.Fatalf("direct DB query failed: %v", err)
	}
	if transEnc == "layoff" || transEnc == "" {
		t.Errorf("expected transition_reason to be encrypted in DB, got %q", transEnc)
	}
	if salMinEnc == "8000" || salMinEnc == "" {
		t.Errorf("expected salary_min to be encrypted in DB, got %q", salMinEnc)
	}

	// AddExperience writes a child row and returns its generated id.
	exp := &domain.WorkExperience{
		Title:        "Product Manager",
		Company:      "Acme",
		StartDate:    "2020-01-01",
		EndDate:      "2023-06-01",
		IsCurrent:    false,
		Achievements: []string{"Delivered project A", "Reduced cost by 15%"},
	}
	if err := repo.AddExperience(ctx, userID, exp); err != nil {
		t.Fatalf("add experience: %v", err)
	}
	if exp.ID == "" {
		t.Fatal("expected experience id to be populated")
	}

	// SetSkills replaces the skill set wholesale.
	if err := repo.SetSkills(ctx, userID, []domain.ProfileSkill{
		{Name: "Strategy", ProficiencyLevel: "expert"},
		{Name: "SQL", ProficiencyLevel: "intermediate"},
		{Name: "Roadmapping", ProficiencyLevel: "expert"},
	}); err != nil {
		t.Fatalf("set skills: %v", err)
	}

	p, err = repo.Get(ctx, userID)
	if err != nil {
		t.Fatalf("get after children: %v", err)
	}
	if len(p.Experiences) != 1 || p.Experiences[0].Company != "Acme" {
		t.Fatalf("expected one Acme experience, got %+v", p.Experiences)
	}
	if len(p.Experiences[0].Achievements) != 2 || p.Experiences[0].Achievements[0] != "Delivered project A" {
		t.Fatalf("expected experience achievements, got %+v", p.Experiences[0].Achievements)
	}
	if len(p.Skills) != 3 || p.Skills[0].ProficiencyLevel != "expert" {
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
