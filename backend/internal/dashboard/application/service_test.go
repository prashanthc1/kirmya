package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/dashboard/domain"
)

func TestSummaryReturnsRepoProjection(t *testing.T) {
	want := domain.Summary{
		UnreadNotifications: 3,
		JobSeeker:           domain.JobSeekerStats{Applications: 5, SavedJobs: 2, OutgoingReferrals: 1},
		Recruiter:           domain.RecruiterStats{PostedJobs: 4, TotalApplicants: 12},
		Mentor:              domain.MentorStats{UpcomingSessions: 1, PendingRequests: 2, CompletedSessions: 7},
	}
	repo := &fakeRepo{summary: want}
	svc := NewService(repo)

	got, err := svc.Summary(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if repo.gotUser != "user-1" {
		t.Fatalf("expected userID passthrough, got %q", repo.gotUser)
	}
}

func TestSummaryPropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	svc := NewService(&fakeRepo{err: sentinel})

	if _, err := svc.Summary(context.Background(), "user-1"); !errors.Is(err, sentinel) {
		t.Fatalf("expected repo error to propagate, got %v", err)
	}
}
