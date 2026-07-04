// Package application implements Mentorship use cases.
package application

import (
	"context"
	"math"
	"strings"
	"time"

	"workspace-app/internal/mentorship/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

type Service struct {
	repo   domain.Repository
	events EventPublisher
}

func NewService(repo domain.Repository, events EventPublisher) *Service {
	return &Service{repo: repo, events: events}
}

type MentorInput struct {
	Headline  string
	Bio       string
	Expertise string
}

// BecomeMentor creates or updates the caller's mentor profile.
func (s *Service) BecomeMentor(ctx context.Context, userID string, in MentorInput) (*domain.MentorProfile, error) {
	if strings.TrimSpace(in.Headline) == "" {
		return nil, ValidationError{"headline is required"}
	}
	p := &domain.MentorProfile{UserID: userID, Headline: in.Headline, Bio: in.Bio, Expertise: in.Expertise, IsActive: true}
	if err := s.repo.UpsertMentorProfile(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListMentors(ctx context.Context) ([]domain.MentorProfile, error) {
	return s.repo.ListMentors(ctx)
}

func (s *Service) GetMentor(ctx context.Context, id string) (*domain.MentorProfile, error) {
	return s.repo.GetMentorByID(ctx, id)
}

// BookInput captures an optional availability slot to consume when booking.
type BookInput struct {
	MentorID    string
	Topic       string
	ScheduledAt time.Time
	SlotID      string // optional: consume this availability slot
}

// Book requests a session with a mentor. If a SlotID is supplied it is validated
// (must belong to the mentor and be open) and consumed; its window then sets the
// scheduled time. Booking without a slot stays supported for backward compatibility.
func (s *Service) Book(ctx context.Context, menteeID string, in BookInput) (*domain.Session, error) {
	mentor, err := s.repo.GetMentorByID(ctx, in.MentorID)
	if err != nil {
		return nil, err
	}
	if mentor.UserID == menteeID {
		return nil, ValidationError{"you cannot book a session with yourself"}
	}

	scheduledAt := in.ScheduledAt
	if in.SlotID != "" {
		slot, err := s.repo.GetSlot(ctx, in.SlotID)
		if err != nil {
			return nil, err
		}
		if slot.MentorID != in.MentorID {
			return nil, ValidationError{"slot does not belong to this mentor"}
		}
		if slot.IsBooked {
			return nil, domain.ErrSlotUnavailable
		}
		scheduledAt = slot.StartsAt
	}
	if scheduledAt.IsZero() {
		return nil, ValidationError{"a scheduled time is required"}
	}

	// When a slot is involved, claim it and create the session in one atomic
	// step so the two can never diverge (no double-booked slot, no orphan
	// session). The earlier GetSlot check gives friendly errors; this is the
	// authoritative guard for the concurrent case.
	sess := &domain.Session{MentorID: in.MentorID, MenteeID: menteeID, Topic: in.Topic, Status: domain.StatusRequested, ScheduledAt: scheduledAt}
	if in.SlotID != "" {
		if err := s.repo.CreateSessionWithSlot(ctx, sess, in.SlotID); err != nil {
			return nil, err
		}
	} else if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	s.publish(ctx, domain.EventMentorshipBooked, sess.ID, map[string]any{
		"session_id": sess.ID, "mentor_user_id": mentor.UserID, "mentee_id": menteeID,
	})
	return sess, nil
}

// AddAvailability lets a mentor open a booking window. The caller must already
// have a mentor profile; the slot is attached to it.
func (s *Service) AddAvailability(ctx context.Context, userID string, startsAt, endsAt time.Time) (*domain.AvailabilitySlot, error) {
	if startsAt.IsZero() || endsAt.IsZero() {
		return nil, ValidationError{"start and end times are required"}
	}
	if !endsAt.After(startsAt) {
		return nil, ValidationError{"end time must be after start time"}
	}
	mentor, err := s.repo.GetMentorByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	slot := &domain.AvailabilitySlot{MentorID: mentor.ID, StartsAt: startsAt, EndsAt: endsAt}
	if err := s.repo.AddAvailability(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

// MentorAvailability lists a mentor's open (unbooked) slots.
func (s *Service) MentorAvailability(ctx context.Context, mentorID string) ([]domain.AvailabilitySlot, error) {
	if _, err := s.repo.GetMentorByID(ctx, mentorID); err != nil {
		return nil, err
	}
	return s.repo.ListAvailability(ctx, mentorID, true)
}

func (s *Service) publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, eventType, aggregateID, payload)
	}
}

// Sessions returns the caller's sessions both as mentee and (if applicable) as mentor.
func (s *Service) Sessions(ctx context.Context, userID string) (asMentee, asMentor []domain.Session, err error) {
	if asMentee, err = s.repo.ListSessionsForMentee(ctx, userID); err != nil {
		return nil, nil, err
	}
	if mentor, mErr := s.repo.GetMentorByUserID(ctx, userID); mErr == nil {
		if asMentor, err = s.repo.ListSessionsForMentor(ctx, mentor.ID); err != nil {
			return nil, nil, err
		}
	}
	return asMentee, asMentor, nil
}

// UpdateStatus lets the owning mentor advance a session.
func (s *Service) UpdateStatus(ctx context.Context, userID, sessionID, status string) (*domain.Session, error) {
	if !domain.ValidStatuses[status] {
		return nil, ValidationError{"invalid status"}
	}
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	mentor, err := s.repo.GetMentorByID(ctx, sess.MentorID)
	if err != nil {
		return nil, err
	}
	if mentor.UserID != userID {
		return nil, domain.ErrNotMentor
	}
	if err := s.repo.UpdateSessionStatus(ctx, sessionID, status); err != nil {
		return nil, err
	}
	switch status {
	case domain.StatusConfirmed:
		s.publish(ctx, domain.EventSessionConfirmed, sessionID, map[string]any{"session_id": sessionID})
	case domain.StatusCompleted:
		s.publish(ctx, domain.EventSessionCompleted, sessionID, map[string]any{"session_id": sessionID})
	}
	return s.repo.GetSession(ctx, sessionID)
}

// Review lets the mentee rate a completed session.
func (s *Service) Review(ctx context.Context, userID, sessionID string, rating int, comment string) (*domain.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ValidationError{"rating must be between 1 and 5"}
	}
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if sess.MenteeID != userID {
		return nil, domain.ErrNotMentee
	}
	if sess.Status != domain.StatusCompleted {
		return nil, domain.ErrNotComplete
	}
	rv := &domain.Review{SessionID: sessionID, Rating: rating, Comment: comment}
	if err := s.repo.CreateReview(ctx, rv); err != nil {
		return nil, err
	}
	s.publish(ctx, domain.EventReviewLeft, rv.ID, map[string]any{
		"session_id": sessionID, "mentor_id": sess.MentorID, "rating": rating,
	})
	return rv, nil
}

// MentorReviewStats is the aggregate view of a mentor's reviews.
type MentorReviewStats struct {
	Reviews       []domain.Review
	AverageRating float64
	Count         int
}

// MentorReviews returns a mentor's reviews together with their average rating.
// Returns domain.ErrNotFound if the mentor does not exist.
func (s *Service) MentorReviews(ctx context.Context, mentorID string) (*MentorReviewStats, error) {
	if _, err := s.repo.GetMentorByID(ctx, mentorID); err != nil {
		return nil, err
	}
	reviews, err := s.repo.ListReviewsForMentor(ctx, mentorID)
	if err != nil {
		return nil, err
	}
	stats := &MentorReviewStats{Reviews: reviews, Count: len(reviews)}
	if stats.Count > 0 {
		sum := 0
		for _, rv := range reviews {
			sum += rv.Rating
		}
		stats.AverageRating = math.Round((float64(sum)/float64(stats.Count))*100) / 100
	}
	return stats, nil
}
