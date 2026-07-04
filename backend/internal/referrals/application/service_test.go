package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/referrals/domain"
)

func TestRequestValidation(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()

	var ve ValidationError
	if _, err := svc.Request(ctx, "seeker", RequestInput{}); !errors.As(err, &ve) {
		t.Fatalf("expected validation error when no job/company, got %v", err)
	}
	if _, err := svc.Request(ctx, "seeker", RequestInput{ReferrerID: "seeker", Company: "Acme"}); !errors.As(err, &ve) {
		t.Fatalf("expected validation error for self-referral, got %v", err)
	}
	ref, err := svc.Request(ctx, "seeker", RequestInput{Company: "Acme", Message: "please refer me"})
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if ref.Status != domain.StatusRequested || ref.ID == "" {
		t.Fatalf("unexpected referral %+v", ref)
	}
}

func TestAcceptOpenRequestClaimsIt(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()
	ref, _ := svc.Request(ctx, "seeker", RequestInput{Company: "Acme"})

	accepted, err := svc.Accept(ctx, "employee", ref.ID)
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	if accepted.Status != domain.StatusAccepted || accepted.ReferrerID != "employee" {
		t.Fatalf("expected claimed+accepted, got %+v", accepted)
	}

	// Re-deciding is a conflict.
	if _, err := svc.Decline(ctx, "employee", ref.ID); !errors.Is(err, domain.ErrAlreadyDecided) {
		t.Fatalf("expected ErrAlreadyDecided, got %v", err)
	}
}

func TestAcceptDirectedToSomeoneElseForbidden(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()
	ref, _ := svc.Request(ctx, "seeker", RequestInput{ReferrerID: "alice", Company: "Acme"})

	if _, err := svc.Accept(ctx, "bob", ref.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden for non-target, got %v", err)
	}
	if _, err := svc.Accept(ctx, "seeker", ref.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden for self-review, got %v", err)
	}
	if _, err := svc.Accept(ctx, "alice", ref.ID); err != nil {
		t.Fatalf("target should accept: %v", err)
	}
}

func TestUpdateOutcomeRules(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()
	ref, _ := svc.Request(ctx, "seeker", RequestInput{Company: "Acme"})

	// Cannot set outcome before acceptance.
	if _, err := svc.UpdateOutcome(ctx, "seeker", ref.ID, domain.OutcomeHired); !errors.Is(err, domain.ErrNotAccepted) {
		t.Fatalf("expected ErrNotAccepted, got %v", err)
	}
	if _, err := svc.Accept(ctx, "employee", ref.ID); err != nil {
		t.Fatalf("accept: %v", err)
	}
	// Invalid outcome rejected.
	if _, err := svc.UpdateOutcome(ctx, "seeker", ref.ID, "bogus"); err == nil {
		t.Fatal("expected validation error for bad outcome")
	}
	// Non-participant forbidden.
	if _, err := svc.UpdateOutcome(ctx, "stranger", ref.ID, domain.OutcomeHired); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden for non-participant, got %v", err)
	}
	// Participant succeeds.
	updated, err := svc.UpdateOutcome(ctx, "seeker", ref.ID, domain.OutcomeHired)
	if err != nil {
		t.Fatalf("outcome: %v", err)
	}
	if updated.Outcome != domain.OutcomeHired {
		t.Fatalf("expected hired, got %s", updated.Outcome)
	}
}
