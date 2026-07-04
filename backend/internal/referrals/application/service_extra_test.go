package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/referrals/domain"
)

// spyEvents records published event types.
type spyEvents struct{ types []string }

func (e *spyEvents) Publish(_ context.Context, eventType, _ string, _ map[string]any) error {
	e.types = append(e.types, eventType)
	return nil
}

func TestIncomingAndOutgoingFilter(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()

	// seeker requests, directed at alice.
	ref, _ := svc.Request(ctx, "seeker", RequestInput{ReferrerID: "alice", Company: "Acme"})
	// another open request by the same seeker.
	_, _ = svc.Request(ctx, "seeker", RequestInput{Company: "Globex"})

	out, err := svc.Outgoing(ctx, "seeker")
	if err != nil {
		t.Fatalf("outgoing: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 outgoing for seeker, got %d", len(out))
	}

	// alice accepts her directed referral, so it becomes incoming for her.
	if _, err := svc.Accept(ctx, "alice", ref.ID); err != nil {
		t.Fatalf("accept: %v", err)
	}
	in, err := svc.Incoming(ctx, "alice")
	if err != nil {
		t.Fatalf("incoming: %v", err)
	}
	if len(in) != 1 || in[0].ReferrerID != "alice" {
		t.Fatalf("expected 1 incoming for alice, got %+v", in)
	}
}

func TestDeclineOpenRequest(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()
	ref, _ := svc.Request(ctx, "seeker", RequestInput{Company: "Acme"})

	declined, err := svc.Decline(ctx, "employee", ref.ID)
	if err != nil {
		t.Fatalf("decline: %v", err)
	}
	if declined.Status != domain.StatusDeclined || declined.ReferrerID != "employee" {
		t.Fatalf("expected declined+claimed, got %+v", declined)
	}
}

func TestDecideOnMissingReferral(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	if _, err := svc.Accept(context.Background(), "employee", "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEventsEmittedOnLifecycle(t *testing.T) {
	ev := &spyEvents{}
	svc := NewService(newFakeRepo(), ev)
	ctx := context.Background()

	ref, _ := svc.Request(ctx, "seeker", RequestInput{Company: "Acme"})
	if _, err := svc.Accept(ctx, "employee", ref.ID); err != nil {
		t.Fatalf("accept: %v", err)
	}

	if len(ev.types) != 2 {
		t.Fatalf("expected 2 events, got %v", ev.types)
	}
	if ev.types[0] != eventReferralRequested || ev.types[1] != eventReferralAccepted {
		t.Fatalf("unexpected event sequence: %v", ev.types)
	}
}
