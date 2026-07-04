//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/referrals/domain"
	"workspace-app/internal/testsupport"
)

func TestReferralsRepository_RequestClaimAndOutcome(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	seeker := testsupport.InsertUser(t, db, "asha@cb.test", "Asha Rao")
	referrer := testsupport.InsertUser(t, db, "carla@cb.test", "Carla Mendes")

	// Open request (no referrer assigned yet).
	ref := &domain.Referral{SeekerID: seeker, Company: "BuildCo", Message: "please refer me", Status: domain.StatusRequested}
	if err := repo.Create(ctx, ref); err != nil {
		t.Fatalf("create: %v", err)
	}
	if ref.ID == "" || ref.Version != 1 {
		t.Fatalf("expected id + version=1, got %+v", ref)
	}

	got, err := repo.Get(ctx, ref.ID)
	if err != nil || !got.IsOpen() || got.Company != "BuildCo" {
		t.Fatalf("get: %+v err=%v", got, err)
	}

	// Seeker sees it in their outgoing list.
	out, err := repo.ListBySeeker(ctx, seeker)
	if err != nil || len(out) != 1 {
		t.Fatalf("by seeker: len=%d err=%v", len(out), err)
	}

	// A referrer claims + accepts it.
	if err := repo.Decide(ctx, ref.ID, referrer, domain.StatusAccepted); err != nil {
		t.Fatalf("decide: %v", err)
	}
	got, _ = repo.Get(ctx, ref.ID)
	if got.Status != domain.StatusAccepted || got.ReferrerID != referrer || got.DecidedAt == nil {
		t.Fatalf("expected accepted+claimed+decided, got %+v", got)
	}

	// Now it appears in the referrer's incoming list.
	in, err := repo.ListByReferrer(ctx, referrer)
	if err != nil || len(in) != 1 {
		t.Fatalf("by referrer: len=%d err=%v", len(in), err)
	}

	// Outcome tracking.
	if err := repo.SetOutcome(ctx, ref.ID, domain.OutcomeHired); err != nil {
		t.Fatalf("outcome: %v", err)
	}
	got, _ = repo.Get(ctx, ref.ID)
	if got.Outcome != domain.OutcomeHired {
		t.Fatalf("expected hired, got %q", got.Outcome)
	}
}
