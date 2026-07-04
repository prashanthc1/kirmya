// Package domain holds the Referrals bounded context entities and ports.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("referral not found")
	ErrAlreadyDecided = errors.New("referral has already been decided")
	ErrNotAccepted    = errors.New("referral must be accepted before tracking an outcome")
)

// Request status values.
const (
	StatusRequested   = "requested"
	StatusUnderReview = "under_review"
	StatusAccepted    = "accepted"
	StatusDeclined    = "declined"
)

// Post-acceptance outcome values (hiring pipeline).
const (
	OutcomeApplicationSubmitted = "application_submitted"
	OutcomeInterviewing         = "interviewing"
	OutcomeOffer                = "offer"
	OutcomeHired                = "hired"
	OutcomeRejected             = "rejected"
	OutcomeWithdrawn            = "withdrawn"
)

var ValidOutcomes = map[string]bool{
	OutcomeApplicationSubmitted: true, OutcomeInterviewing: true, OutcomeOffer: true,
	OutcomeHired: true, OutcomeRejected: true, OutcomeWithdrawn: true,
}

// Referral is the aggregate root. Empty ReferrerID/JobID/Company/Outcome mean
// "not set" (stored as NULL).
type Referral struct {
	ID         string
	SeekerID   string
	ReferrerID string
	JobID      string
	Company    string
	Message    string
	Status     string
	Outcome    string
	DecidedAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int
}

// IsOpen reports whether no specific referrer is assigned yet.
func (r *Referral) IsOpen() bool { return r.ReferrerID == "" }

// Repository is the persistence port.
type Repository interface {
	Create(ctx context.Context, r *Referral) error
	Get(ctx context.Context, id string) (*Referral, error)
	ListByReferrer(ctx context.Context, referrerID string) ([]Referral, error)
	ListBySeeker(ctx context.Context, seekerID string) ([]Referral, error)
	// Decide claims the request for referrerID (if open) and sets the status.
	Decide(ctx context.Context, id, referrerID, status string) error
	SetOutcome(ctx context.Context, id, outcome string) error
}
