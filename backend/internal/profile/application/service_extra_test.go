package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/profile/domain"
)

// spyEvents records published event types for assertions.
type spyEvents struct {
	types []string
}

func (e *spyEvents) Publish(_ context.Context, eventType, _ string, _ map[string]any) error {
	e.types = append(e.types, eventType)
	return nil
}

func TestEducationLifecycle(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	ed := domain.Education{Institution: "State University", Degree: "BSc"}
	p, err := svc.AddEducation(ctx, "u1", &ed)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(p.Educations) != 1 || ed.ID == "" {
		t.Fatalf("expected one education with assigned id, got %+v", p.Educations)
	}

	ed.Degree = "MSc"
	if _, err := svc.UpdateEducation(ctx, "u1", ed); err != nil {
		t.Fatalf("update: %v", err)
	}
	if _, err := svc.UpdateEducation(ctx, "intruder", ed); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-owner, got %v", err)
	}

	p, err = svc.DeleteEducation(ctx, "u1", ed.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(p.Educations) != 0 {
		t.Fatalf("expected no educations after delete, got %d", len(p.Educations))
	}
}

func TestCertificationLifecycle(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	c := domain.CertificationItem{Name: "PMP", Issuer: "PMI"}
	p, err := svc.AddCertification(ctx, "u1", &c)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(p.Certifications) != 1 || c.ID == "" {
		t.Fatalf("expected one certification with assigned id, got %+v", p.Certifications)
	}

	c.Issuer = "Project Management Institute"
	if _, err := svc.UpdateCertification(ctx, "u1", c); err != nil {
		t.Fatalf("update: %v", err)
	}
	if _, err := svc.UpdateCertification(ctx, "intruder", c); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-owner, got %v", err)
	}

	p, err = svc.DeleteCertification(ctx, "u1", c.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(p.Certifications) != 0 {
		t.Fatalf("expected no certifications after delete, got %d", len(p.Certifications))
	}
}

func TestSetLanguagesAndPortfolioReplace(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	if _, err := svc.SetLanguages(ctx, "u1", []domain.LanguageItem{{Name: "English", Proficiency: "native"}}); err != nil {
		t.Fatalf("set languages: %v", err)
	}
	p, err := svc.SetLanguages(ctx, "u1", []domain.LanguageItem{{Name: "Arabic", Proficiency: "professional"}})
	if err != nil {
		t.Fatalf("set languages 2: %v", err)
	}
	if len(p.Identity.Languages) != 1 || p.Identity.Languages[0].Name != "Arabic" {
		t.Fatalf("languages should be replaced, got %+v", p.Identity.Languages)
	}

	p, err = svc.SetPortfolio(ctx, "u1", []domain.ProjectItem{{LiveDemoURL: "https://example.com"}})
	if err != nil {
		t.Fatalf("set portfolio: %v", err)
	}
	if len(p.Projects) != 1 || p.Projects[0].LiveDemoURL != "https://example.com" {
		t.Fatalf("portfolio not persisted, got %+v", p.Projects)
	}
}

func TestWritesEmitProfileUpdated(t *testing.T) {
	ev := &spyEvents{}
	svc := NewService(newFakeRepo(), ev, nil)
	ctx := context.Background()

	isDraft := true
	if _, err := svc.UpdateProfile(ctx, "u1", 0, domain.AggregateUpdate{
		Identity: &domain.IdentitySection{Headline: "Engineer"},
		IsDraft:  &isDraft,
	}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if _, err := svc.SetSkills(ctx, "u1", []domain.SkillItem{{Name: "Go"}}); err != nil {
		t.Fatalf("skills: %v", err)
	}
	if len(ev.types) != 2 {
		t.Fatalf("expected 2 events emitted, got %d (%v)", len(ev.types), ev.types)
	}
	for _, ty := range ev.types {
		if ty != eventProfileUpdated {
			t.Fatalf("expected %q events, got %q", eventProfileUpdated, ty)
		}
	}
}

func TestGetReturnsLazyEmptyProfile(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	p, err := svc.Get(context.Background(), "newuser")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p == nil || p.UserID != "newuser" {
		t.Fatalf("expected lazily-created profile for newuser, got %+v", p)
	}
}
