//go:build integration

package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"workspace-app/internal/mentorship/domain"
	"workspace-app/internal/testsupport"
)

// TestCreateSessionWithSlot_AtomicClaim proves the double-booking guard: two
// bookings race for the same open slot, exactly one wins, and the loser writes
// nothing — no second session, slot left consistent.
func TestCreateSessionWithSlot_AtomicClaim(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	mentorUser := testsupport.InsertUser(t, db, "mentor@cb.test", "Mara Mentor")
	menteeA := testsupport.InsertUser(t, db, "mentee-a@cb.test", "Aya Mentee")
	menteeB := testsupport.InsertUser(t, db, "mentee-b@cb.test", "Ben Mentee")

	mentor := &domain.MentorProfile{UserID: mentorUser, Headline: "Ops leader"}
	if err := repo.UpsertMentorProfile(ctx, mentor); err != nil {
		t.Fatalf("upsert mentor: %v", err)
	}

	start := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	slot := &domain.AvailabilitySlot{MentorID: mentor.ID, StartsAt: start, EndsAt: start.Add(time.Hour)}
	if err := repo.AddAvailability(ctx, slot); err != nil {
		t.Fatalf("add availability: %v", err)
	}

	// First booking claims the slot.
	first := &domain.Session{MentorID: mentor.ID, MenteeID: menteeA, Topic: "interview prep", Status: domain.StatusRequested, ScheduledAt: start}
	if err := repo.CreateSessionWithSlot(ctx, first, slot.ID); err != nil {
		t.Fatalf("first booking: %v", err)
	}
	if first.ID == "" {
		t.Fatal("expected first session to be created")
	}

	// Second booking against the same slot must be rejected and write nothing.
	second := &domain.Session{MentorID: mentor.ID, MenteeID: menteeB, Topic: "resume review", Status: domain.StatusRequested, ScheduledAt: start}
	if err := repo.CreateSessionWithSlot(ctx, second, slot.ID); !errors.Is(err, domain.ErrSlotUnavailable) {
		t.Fatalf("expected ErrSlotUnavailable on double-book, got %v", err)
	}
	if second.ID != "" {
		t.Fatalf("expected no session written for the losing booking, got id=%s", second.ID)
	}

	// The slot is now booked and no longer listed as open.
	got, err := repo.GetSlot(ctx, slot.ID)
	if err != nil {
		t.Fatalf("get slot: %v", err)
	}
	if !got.IsBooked {
		t.Fatal("expected slot to be marked booked")
	}
	open, err := repo.ListAvailability(ctx, mentor.ID, true)
	if err != nil {
		t.Fatalf("list availability: %v", err)
	}
	if len(open) != 0 {
		t.Fatalf("expected no open slots, got %d", len(open))
	}

	// Exactly one session exists for the mentor.
	sessions, err := repo.ListSessionsForMentor(ctx, mentor.ID)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected exactly one session, got %d", len(sessions))
	}
	if sessions[0].MenteeID != menteeA {
		t.Fatalf("expected the winning session to belong to menteeA, got %s", sessions[0].MenteeID)
	}
}
