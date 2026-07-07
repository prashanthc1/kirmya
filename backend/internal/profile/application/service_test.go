package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"workspace-app/internal/profile/domain"
)

func TestUpdateMeAndReload(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	isDraft := true
	p, err := svc.UpdateProfile(ctx, "u1", 0, domain.AggregateUpdate{
		Identity: &domain.IdentitySection{
			Headline: "Operations Leader",
			Location: "Dubai",
		},
		IsDraft: &isDraft,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if p.Identity.Headline != "Operations Leader" || p.Identity.Location != "Dubai" {
		t.Fatalf("scalars not persisted: %+v", p)
	}
}

func TestExperienceLifecycle(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	e := domain.WorkExperience{
		Position:     "Coordinator",
		Company:      "Acme",
		StartDate:    time.Now().AddDate(-1, 0, 0),
		EndDate:      time.Now(),
		Achievements: []string{"Achievement 1"},
	}
	p, err := svc.AddExperience(ctx, "u1", &e)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(p.Experiences) != 1 || e.ID == "" {
		t.Fatalf("expected one experience with assigned id, got %+v", p.Experiences)
	}
	if len(p.Experiences[0].Achievements) != 1 || p.Experiences[0].Achievements[0] != "Achievement 1" {
		t.Fatalf("expected achievements in saved experience, got %+v", p.Experiences[0].Achievements)
	}

	e.Position = "Senior Coordinator"
	if _, err := svc.UpdateExperience(ctx, "u1", e); err != nil {
		t.Fatalf("update: %v", err)
	}

	// Ownership: another user cannot update this experience.
	if _, err := svc.UpdateExperience(ctx, "intruder", e); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-owner, got %v", err)
	}

	p, err = svc.DeleteExperience(ctx, "u1", e.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(p.Experiences) != 0 {
		t.Fatalf("expected no experiences after delete, got %d", len(p.Experiences))
	}
}

func TestSetSkillsReplaces(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	if _, err := svc.SetSkills(ctx, "u1", []domain.SkillItem{{Name: "Go"}, {Name: "PostgreSQL"}}); err != nil {
		t.Fatalf("set: %v", err)
	}
	p, err := svc.SetSkills(ctx, "u1", []domain.SkillItem{{Name: "Leadership"}})
	if err != nil {
		t.Fatalf("set2: %v", err)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "Leadership" {
		t.Fatalf("skills should be replaced, got %v", p.Skills)
	}
}

func TestValidationRules(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	// Test 1: Salary min > max should fail
	_, err := svc.UpdateProfile(ctx, "u1", 0, domain.AggregateUpdate{
		Preferences: &domain.CareerPreferences{
			SalaryMin: 10000,
			SalaryMax: 5000,
		},
	})
	if err == nil {
		t.Error("expected error for min salary > max salary, got nil")
	}

	// Test 2: Experience start date > end date should fail
	_, err = svc.AddExperience(ctx, "u1", &domain.WorkExperience{
		Position:  "Mgr",
		Company:   "Acme",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(-1, 0, 0),
	})
	if err == nil {
		t.Error("expected error for experience start date > end date, got nil")
	}
}
