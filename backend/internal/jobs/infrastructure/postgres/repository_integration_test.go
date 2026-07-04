//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/jobs/domain"
	"workspace-app/internal/testsupport"
)

func TestJobsRepository_PostSearchApplySave(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	recruiter := testsupport.InsertUser(t, db, "rita@cb.test", "Rita Shah")
	seeker := testsupport.InsertUser(t, db, "asha@cb.test", "Asha Rao")

	job := &domain.Job{
		Title: "Operations Manager", Company: "Acme Logistics", Location: "Dubai",
		Description: "Lead the operations team across two warehouses.", JobType: "full_time", PostedBy: recruiter,
	}
	if err := repo.CreateJob(ctx, job); err != nil {
		t.Fatalf("create job: %v", err)
	}
	if job.ID == "" {
		t.Fatal("expected job id from RETURNING")
	}

	// Get round-trips.
	got, err := repo.GetJob(ctx, job.ID)
	if err != nil || got.Title != "Operations Manager" || got.PostedBy != recruiter {
		t.Fatalf("get job: %+v err=%v", got, err)
	}

	// ILIKE search matches title; non-match returns empty.
	hits, err := repo.ListJobs(ctx, domain.Filter{Keyword: "operations", Limit: 10})
	if err != nil || len(hits) != 1 {
		t.Fatalf("search hit: %v len=%d", err, len(hits))
	}
	none, _ := repo.ListJobs(ctx, domain.Filter{Keyword: "zzz-nomatch", Limit: 10})
	if len(none) != 0 {
		t.Fatalf("expected no matches, got %d", len(none))
	}

	// Apply: first succeeds, dup-guard via HasApplied.
	app := &domain.Application{JobID: job.ID, UserID: seeker, Status: domain.StatusPending, CoverLetter: "please consider me"}
	if err := repo.CreateApplication(ctx, app); err != nil {
		t.Fatalf("apply: %v", err)
	}
	applied, err := repo.HasApplied(ctx, job.ID, seeker)
	if err != nil || !applied {
		t.Fatalf("expected HasApplied true, got %v err=%v", applied, err)
	}

	// Applicants list visible to the recruiter side.
	apps, err := repo.ListApplicationsByJob(ctx, job.ID)
	if err != nil || len(apps) != 1 || apps[0].UserID != seeker {
		t.Fatalf("applicants: %+v err=%v", apps, err)
	}

	// Status transition.
	if err := repo.UpdateApplicationStatus(ctx, app.ID, domain.StatusInterviewing); err != nil {
		t.Fatalf("status: %v", err)
	}
	reloaded, _ := repo.GetApplication(ctx, app.ID)
	if reloaded.Status != domain.StatusInterviewing {
		t.Fatalf("expected interviewing, got %s", reloaded.Status)
	}

	// Save / unsave toggle.
	if err := repo.SaveJob(ctx, seeker, job.ID); err != nil {
		t.Fatalf("save: %v", err)
	}
	saved, _ := repo.IsSaved(ctx, seeker, job.ID)
	if !saved {
		t.Fatal("expected saved=true")
	}
	if err := repo.UnsaveJob(ctx, seeker, job.ID); err != nil {
		t.Fatalf("unsave: %v", err)
	}
	if saved, _ := repo.IsSaved(ctx, seeker, job.ID); saved {
		t.Fatal("expected saved=false after unsave")
	}
}
