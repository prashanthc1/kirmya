package domain

import "context"

// Repository is the persistence port for the Mentorship context. The concrete
// adapter lives in infrastructure/postgres.
type Repository interface {
	UpsertMentorProfile(ctx context.Context, p *MentorProfile) error
	GetMentorByID(ctx context.Context, id string) (*MentorProfile, error)
	GetMentorByUserID(ctx context.Context, userID string) (*MentorProfile, error)
	ListMentors(ctx context.Context) ([]MentorProfile, error)

	CreateSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, id string) (*Session, error)
	ListSessionsForMentee(ctx context.Context, menteeID string) ([]Session, error)
	ListSessionsForMentor(ctx context.Context, mentorID string) ([]Session, error)
	UpdateSessionStatus(ctx context.Context, id, status string) error

	CreateReview(ctx context.Context, rv *Review) error
	ListReviewsForMentor(ctx context.Context, mentorID string) ([]Review, error)

	AddAvailability(ctx context.Context, slot *AvailabilitySlot) error
	ListAvailability(ctx context.Context, mentorID string, openOnly bool) ([]AvailabilitySlot, error)
	GetSlot(ctx context.Context, id string) (*AvailabilitySlot, error)

	// CreateSessionWithSlot atomically claims an open availability slot and
	// inserts the session in a single transaction. It returns ErrSlotUnavailable
	// (and writes nothing) if the slot is no longer open — i.e. another booking
	// claimed it first. This is the authoritative guard against double-booking.
	CreateSessionWithSlot(ctx context.Context, s *Session, slotID string) error
}
