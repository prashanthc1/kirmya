//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"workspace-app/internal/profile/domain"
	"workspace-app/internal/testsupport"
)

func TestProfileLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	userID := testsupport.InsertUser(t, db, "profile@cb.test", "Pat Profile")

	// Get lazily creates the row and returns empty sections at the initial version 1.
	p, err := repo.Get(ctx, userID, true)
	if err != nil {
		t.Fatalf("get (create): %v", err)
	}
	if p.UserID != userID {
		t.Fatalf("expected profile for %s, got %s", userID, p.UserID)
	}
	if p.Version != 1 {
		t.Fatalf("expected initial version 1, got %d", p.Version)
	}

	// UpdateAggregate persists fields and bumps version.
	err = repo.UpdateAggregate(ctx, userID, 0, domain.AggregateUpdate{
		Identity: &domain.IdentitySection{
			Headline: "Recovering PM",
			Location: "Remote",
		},
		Summary: &domain.SummarySection{
			ExecutiveSummary: "Back after a career gap.",
		},
		Preferences: &domain.CareerPreferences{
			SalaryMin:      8000,
			SalaryMax:      12000,
			SalaryCurrency: "AED",
		},
	})
	if err != nil {
		t.Fatalf("update aggregate: %v", err)
	}

	p, err = repo.Get(ctx, userID, true)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if p.Identity.Headline != "Recovering PM" || p.Identity.Location != "Remote" {
		t.Fatalf("scalars not persisted: %+v", p)
	}
	if p.Summary.ExecutiveSummary != "Back after a career gap." || p.Preferences.SalaryMin != 8000 || p.Preferences.SalaryCurrency != "AED" {
		t.Fatalf("fields mismatch: %+v", p)
	}

	// Add experience
	exp := domain.WorkExperience{
		Position:     "Product Manager",
		Company:      "Acme",
		StartDate:    time.Now().AddDate(-1, 0, 0),
		EndDate:      time.Now(),
		IsCurrent:    false,
		Achievements: []string{"Delivered project A", "Reduced cost by 15%"},
	}
	err = repo.UpdateAggregate(ctx, userID, 0, domain.AggregateUpdate{
		Experiences: &[]domain.WorkExperience{exp},
	})
	if err != nil {
		t.Fatalf("update experience: %v", err)
	}

	// Set skills
	err = repo.UpdateAggregate(ctx, userID, 0, domain.AggregateUpdate{
		Skills: &[]domain.SkillItem{
			{Name: "Strategy", Level: "expert"},
			{Name: "SQL", Level: "intermediate"},
		},
	})
	if err != nil {
		t.Fatalf("set skills: %v", err)
	}

	p, err = repo.Get(ctx, userID, true)
	if err != nil {
		t.Fatalf("get after children: %v", err)
	}
	if len(p.Experiences) != 1 || p.Experiences[0].Company != "Acme" {
		t.Fatalf("expected one Acme experience, got %+v", p.Experiences)
	}
	if len(p.Experiences[0].Achievements) != 2 || p.Experiences[0].Achievements[0] != "Delivered project A" {
		t.Fatalf("expected experience achievements, got %+v", p.Experiences[0].Achievements)
	}
	if len(p.Skills) != 2 || p.Skills[0].Level != "expert" {
		t.Fatalf("expected 2 skills, got %d (%v)", len(p.Skills), p.Skills)
	}
}
