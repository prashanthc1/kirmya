package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/profile/domain"
)

func TestUpdateScalarsAndReload(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	p, err := svc.UpdateScalars(ctx, "u1", domain.Scalars{Headline: "Operations Leader", Location: "Dubai"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if p.Headline != "Operations Leader" || p.Location != "Dubai" {
		t.Fatalf("scalars not persisted: %+v", p)
	}
}

func TestExperienceLifecycle(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	e := domain.WorkExperience{Title: "Coordinator", Company: "Acme"}
	p, err := svc.AddExperience(ctx, "u1", &e)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(p.Experiences) != 1 || e.ID == "" {
		t.Fatalf("expected one experience with assigned id, got %+v", p.Experiences)
	}

	e.Title = "Senior Coordinator"
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

	if _, err := svc.SetSkills(ctx, "u1", []string{"Go", "PostgreSQL"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	p, err := svc.SetSkills(ctx, "u1", []string{"Leadership"})
	if err != nil {
		t.Fatalf("set2: %v", err)
	}
	if len(p.Skills) != 1 || p.Skills[0] != "Leadership" {
		t.Fatalf("skills should be replaced, got %v", p.Skills)
	}
}
