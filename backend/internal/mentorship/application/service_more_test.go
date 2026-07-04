package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"workspace-app/internal/mentorship/domain"
)

// Covers the mentor-profile, listing, and session-read use cases plus the Book
// error branches not exercised by service_test.go, reusing its in-package fakes.

func TestBecomeMentorValidationAndUpsert(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	if _, err := svc.BecomeMentor(ctx, "u1", MentorInput{Headline: "   "}); err == nil {
		t.Fatal("expected validation error for blank headline")
	}

	first, err := svc.BecomeMentor(ctx, "u1", MentorInput{Headline: "Ops leader", Bio: "b", Expertise: "ops"})
	if err != nil {
		t.Fatalf("become mentor: %v", err)
	}
	if first.ID == "" || !first.IsActive {
		t.Fatalf("expected an active mentor profile, got %+v", first)
	}

	// Calling again for the same user upserts onto the same profile id.
	second, err := svc.BecomeMentor(ctx, "u1", MentorInput{Headline: "Ops director"})
	if err != nil {
		t.Fatalf("re-become mentor: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected upsert to keep id %q, got %q", first.ID, second.ID)
	}
}

func TestListAndGetMentor(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()

	mentors, err := svc.ListMentors(ctx)
	if err != nil {
		t.Fatalf("list mentors: %v", err)
	}
	if len(mentors) != 1 || mentors[0].ID != mentorID {
		t.Fatalf("expected the one seeded mentor, got %+v", mentors)
	}

	m, err := svc.GetMentor(ctx, mentorID)
	if err != nil {
		t.Fatalf("get mentor: %v", err)
	}
	if m.UserID != "mentor-user" {
		t.Fatalf("unexpected mentor %+v", m)
	}

	if _, err := svc.GetMentor(ctx, "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSessionsAsMenteeAndMentor(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()

	// mentor-user books with another mentor? Simpler: a mentee books this mentor.
	if _, err := svc.Book(ctx, "mentee", BookInput{MentorID: mentorID, Topic: "t", ScheduledAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("book: %v", err)
	}

	// As the mentee: one session as mentee, none as mentor.
	asMentee, asMentor, err := svc.Sessions(ctx, "mentee")
	if err != nil {
		t.Fatalf("sessions (mentee): %v", err)
	}
	if len(asMentee) != 1 || len(asMentor) != 0 {
		t.Fatalf("expected 1 mentee / 0 mentor session, got %d / %d", len(asMentee), len(asMentor))
	}

	// As the mentor-user: zero as mentee, one as mentor.
	asMentee, asMentor, err = svc.Sessions(ctx, "mentor-user")
	if err != nil {
		t.Fatalf("sessions (mentor): %v", err)
	}
	if len(asMentee) != 0 || len(asMentor) != 1 {
		t.Fatalf("expected 0 mentee / 1 mentor session, got %d / %d", len(asMentee), len(asMentor))
	}
}

func TestMentorAvailabilityUnknownMentor(t *testing.T) {
	svc, _, _ := setup(t)
	if _, err := svc.MentorAvailability(context.Background(), "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown mentor, got %v", err)
	}
}

func TestBookErrorBranches(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()

	// Unknown mentor.
	if _, err := svc.Book(ctx, "mentee", BookInput{MentorID: "missing", Topic: "t", ScheduledAt: time.Now().Add(time.Hour)}); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound booking unknown mentor, got %v", err)
	}

	// No slot and no scheduled time.
	if _, err := svc.Book(ctx, "mentee", BookInput{MentorID: mentorID, Topic: "t"}); err == nil {
		t.Fatal("expected validation error when no scheduled time is given")
	}

	// A slot that belongs to mentorID, booked against a different mentor.
	start := time.Now().Add(time.Hour)
	slot, err := svc.AddAvailability(ctx, "mentor-user", start, start.Add(time.Hour))
	if err != nil {
		t.Fatalf("add availability: %v", err)
	}
	other, err := svc.BecomeMentor(ctx, "other-user", MentorInput{Headline: "Another"})
	if err != nil {
		t.Fatalf("become other mentor: %v", err)
	}
	if _, err := svc.Book(ctx, "mentee", BookInput{MentorID: other.ID, Topic: "t", SlotID: slot.ID}); err == nil {
		t.Fatal("expected error booking a slot that belongs to a different mentor")
	}
}
