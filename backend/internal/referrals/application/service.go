// Package application implements the Referrals use cases over the domain ports.
package application

import (
	"context"
	"errors"
	"strings"

	"workspace-app/internal/referrals/domain"
)

// EventPublisher publishes domain events (the platform bus satisfies this).
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

const (
	eventReferralRequested = "ReferralRequested"
	eventReferralAccepted  = "ReferralAccepted"
	eventReferralDeclined  = "ReferralDeclined"
)

// ValidationError is returned for invalid input (mapped to HTTP 400).
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

// ErrForbidden indicates the caller may not act on this referral.
var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo   domain.Repository
	events EventPublisher
}

func NewService(repo domain.Repository, events EventPublisher) *Service {
	return &Service{repo: repo, events: events}
}

// RequestInput is the create-referral command.
type RequestInput struct {
	ReferrerID string // optional: directed referral
	JobID      string // optional
	Company    string // optional
	Message    string
}

// Request creates a referral request for the seeker.
func (s *Service) Request(ctx context.Context, seekerID string, in RequestInput) (*domain.Referral, error) {
	if strings.TrimSpace(in.JobID) == "" && strings.TrimSpace(in.Company) == "" {
		return nil, ValidationError{"a job or company is required"}
	}
	if in.ReferrerID != "" && in.ReferrerID == seekerID {
		return nil, ValidationError{"you cannot request a referral from yourself"}
	}
	ref := &domain.Referral{
		SeekerID: seekerID, ReferrerID: in.ReferrerID, JobID: in.JobID,
		Company: strings.TrimSpace(in.Company), Message: strings.TrimSpace(in.Message),
		Status: domain.StatusRequested,
	}
	if err := s.repo.Create(ctx, ref); err != nil {
		return nil, err
	}
	s.publish(ctx, eventReferralRequested, ref.ID, map[string]any{
		"seeker_id": seekerID, "referrer_id": in.ReferrerID, "job_id": in.JobID,
	})
	return ref, nil
}

// Incoming lists referrals directed at (or claimed by) the referrer.
func (s *Service) Incoming(ctx context.Context, referrerID string) ([]domain.Referral, error) {
	return s.repo.ListByReferrer(ctx, referrerID)
}

// Outgoing lists the seeker's own requests.
func (s *Service) Outgoing(ctx context.Context, seekerID string) ([]domain.Referral, error) {
	return s.repo.ListBySeeker(ctx, seekerID)
}

// Accept lets a referrer accept (claiming the request if it was open).
func (s *Service) Accept(ctx context.Context, referrerID, id string) (*domain.Referral, error) {
	return s.decide(ctx, referrerID, id, domain.StatusAccepted, eventReferralAccepted)
}

// Decline lets a referrer decline (claiming the request if it was open).
func (s *Service) Decline(ctx context.Context, referrerID, id string) (*domain.Referral, error) {
	return s.decide(ctx, referrerID, id, domain.StatusDeclined, eventReferralDeclined)
}

func (s *Service) decide(ctx context.Context, referrerID, id, status, event string) (*domain.Referral, error) {
	ref, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ref.SeekerID == referrerID {
		return nil, ErrForbidden // a seeker cannot review their own request
	}
	if !ref.IsOpen() && ref.ReferrerID != referrerID {
		return nil, ErrForbidden // directed at someone else
	}
	if ref.Status == domain.StatusAccepted || ref.Status == domain.StatusDeclined {
		return nil, domain.ErrAlreadyDecided
	}
	if err := s.repo.Decide(ctx, id, referrerID, status); err != nil {
		return nil, err
	}
	s.publish(ctx, event, id, map[string]any{"referrer_id": referrerID, "seeker_id": ref.SeekerID})
	return s.repo.Get(ctx, id)
}

// UpdateOutcome records the hiring-pipeline outcome. Either participant may set it.
func (s *Service) UpdateOutcome(ctx context.Context, userID, id, outcome string) (*domain.Referral, error) {
	if !domain.ValidOutcomes[outcome] {
		return nil, ValidationError{"invalid outcome"}
	}
	ref, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if userID != ref.SeekerID && userID != ref.ReferrerID {
		return nil, ErrForbidden
	}
	if ref.Status != domain.StatusAccepted {
		return nil, domain.ErrNotAccepted
	}
	if err := s.repo.SetOutcome(ctx, id, outcome); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) publish(ctx context.Context, evt, aggID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, evt, aggID, payload)
	}
}
