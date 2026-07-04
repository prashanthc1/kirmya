package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/jobs/domain"
)

func validJob() PostJobInput {
	return PostJobInput{Title: "Operations Manager", Company: "Acme", Location: "Dubai",
		Description: "Lead operations across the region with full P&L ownership.", JobType: "full-time"}
}

func TestPostJobSuccessAndValidation(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	job, err := svc.PostJob(ctx, "recruiter-1", validJob())
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if job.ID == "" || job.PostedBy != "recruiter-1" {
		t.Fatalf("unexpected job %+v", job)
	}

	var ve ValidationError
	if _, err := svc.PostJob(ctx, "r", PostJobInput{Title: "", Company: "x", Description: "long enough description here"}); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for empty title, got %v", err)
	}
	if _, err := svc.PostJob(ctx, "r", PostJobInput{Title: "t", Company: "c", Description: "short"}); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for short description, got %v", err)
	}
}

func TestSearchJobsPostedByFilter(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()

	if _, err := svc.PostJob(ctx, "recruiter-1", validJob()); err != nil {
		t.Fatalf("post: %v", err)
	}
	if _, err := svc.PostJob(ctx, "recruiter-2", validJob()); err != nil {
		t.Fatalf("post: %v", err)
	}

	all, err := svc.SearchJobs(ctx, domain.Filter{})
	if err != nil {
		t.Fatalf("search all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 jobs unfiltered, got %d", len(all))
	}

	mine, err := svc.SearchJobs(ctx, domain.Filter{PostedBy: "recruiter-1"})
	if err != nil {
		t.Fatalf("search mine: %v", err)
	}
	if len(mine) != 1 || mine[0].PostedBy != "recruiter-1" {
		t.Fatalf("expected only recruiter-1's job, got %+v", mine)
	}
}

func TestUpdateAndDeleteOwnership(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()
	job, _ := svc.PostJob(ctx, "owner", validJob())

	if _, err := svc.UpdateJob(ctx, "intruder", job.ID, PostJobInput{Title: "Hacked"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden on update, got %v", err)
	}
	if err := svc.DeleteJob(ctx, "intruder", job.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden on delete, got %v", err)
	}
	if _, err := svc.UpdateJob(ctx, "owner", job.ID, PostJobInput{Title: "Senior Operations Manager"}); err != nil {
		t.Fatalf("owner update: %v", err)
	}
}

func TestApplyFlowAndDuplicate(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()
	job, _ := svc.PostJob(ctx, "owner", validJob())

	app, err := svc.Apply(ctx, "seeker", job.ID, "I am genuinely interested in this role.")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if app.Status != domain.StatusPending {
		t.Fatalf("expected pending, got %s", app.Status)
	}
	if _, err := svc.Apply(ctx, "seeker", job.ID, "applying again with enough text"); !errors.Is(err, domain.ErrAlreadyApplied) {
		t.Fatalf("expected ErrAlreadyApplied, got %v", err)
	}
	if _, err := svc.Apply(ctx, "seeker", "missing", "enough text here please"); !errors.Is(err, domain.ErrJobNotFound) {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}
}

func TestUpdateApplicationStatusOwnership(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()
	job, _ := svc.PostJob(ctx, "owner", validJob())
	app, _ := svc.Apply(ctx, "seeker", job.ID, "I am genuinely interested in this role.")

	if _, err := svc.UpdateApplicationStatus(ctx, "owner", app.ID, "bogus"); err == nil {
		t.Fatal("expected validation error for bad status")
	}
	if _, err := svc.UpdateApplicationStatus(ctx, "intruder", app.ID, domain.StatusHired); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	updated, err := svc.UpdateApplicationStatus(ctx, "owner", app.ID, domain.StatusInterviewing)
	if err != nil {
		t.Fatalf("status update: %v", err)
	}
	if updated.Status != domain.StatusInterviewing {
		t.Fatalf("expected interviewing, got %s", updated.Status)
	}
}

func TestToggleSave(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, nil)
	ctx := context.Background()
	job, _ := svc.PostJob(ctx, "owner", validJob())

	saved, err := svc.ToggleSave(ctx, "seeker", job.ID)
	if err != nil || !saved {
		t.Fatalf("expected saved=true, got %v err=%v", saved, err)
	}
	list, _ := svc.SavedJobs(ctx, "seeker")
	if len(list) != 1 {
		t.Fatalf("expected 1 saved job, got %d", len(list))
	}
	saved, err = svc.ToggleSave(ctx, "seeker", job.ID)
	if err != nil || saved {
		t.Fatalf("expected saved=false after toggle, got %v err=%v", saved, err)
	}
}
