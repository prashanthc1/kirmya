package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workspace-app/internal/referrals/domain"
)

type fakeRepo struct {
	mu    sync.Mutex
	seq   int
	store map[string]*domain.Referral
}

func newFakeRepo() *fakeRepo { return &fakeRepo{store: map[string]*domain.Referral{}} }

func (r *fakeRepo) Create(_ context.Context, ref *domain.Referral) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	ref.ID = fmt.Sprintf("ref-%d", r.seq)
	ref.Version = 1
	ref.CreatedAt = time.Now()
	ref.UpdatedAt = ref.CreatedAt
	c := *ref
	r.store[ref.ID] = &c
	return nil
}

func (r *fakeRepo) Get(_ context.Context, id string) (*domain.Referral, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ref, ok := r.store[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *ref
	return &c, nil
}

func (r *fakeRepo) ListByReferrer(_ context.Context, referrerID string) ([]domain.Referral, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Referral{}
	for _, ref := range r.store {
		if ref.ReferrerID == referrerID {
			out = append(out, *ref)
		}
	}
	return out, nil
}

func (r *fakeRepo) ListBySeeker(_ context.Context, seekerID string) ([]domain.Referral, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Referral{}
	for _, ref := range r.store {
		if ref.SeekerID == seekerID {
			out = append(out, *ref)
		}
	}
	return out, nil
}

func (r *fakeRepo) Decide(_ context.Context, id, referrerID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ref, ok := r.store[id]
	if !ok {
		return domain.ErrNotFound
	}
	now := time.Now()
	ref.ReferrerID = referrerID
	ref.Status = status
	ref.DecidedAt = &now
	ref.Version++
	return nil
}

func (r *fakeRepo) SetOutcome(_ context.Context, id, outcome string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ref, ok := r.store[id]
	if !ok {
		return domain.ErrNotFound
	}
	ref.Outcome = outcome
	ref.Version++
	return nil
}
