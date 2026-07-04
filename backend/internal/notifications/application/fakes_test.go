package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workspace-app/internal/notifications/domain"
)

// fakeRepo is an in-memory domain.Repository for service unit tests.
type fakeRepo struct {
	mu        sync.Mutex
	seq       int
	store     map[string]*domain.Notification
	createErr error
}

func newFakeRepo() *fakeRepo { return &fakeRepo{store: map[string]*domain.Notification{}} }

func (r *fakeRepo) Create(_ context.Context, n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.createErr != nil {
		return r.createErr
	}
	r.seq++
	n.ID = fmt.Sprintf("ntf-%d", r.seq)
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	c := *n
	r.store[n.ID] = &c
	return nil
}

func (r *fakeRepo) ListByUser(_ context.Context, userID string, limit, offset int) ([]domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	all := []domain.Notification{}
	for _, n := range r.store {
		if n.UserID == userID {
			all = append(all, *n)
		}
	}
	if offset > len(all) {
		offset = len(all)
	}
	all = all[offset:]
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func (r *fakeRepo) MarkRead(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.store[id]
	if !ok || n.UserID != userID {
		return nil
	}
	now := time.Now()
	n.ReadAt = &now
	return nil
}

func (r *fakeRepo) MarkAllRead(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, n := range r.store {
		if n.UserID == userID && n.ReadAt == nil {
			n.ReadAt = &now
		}
	}
	return nil
}

func (r *fakeRepo) UnreadCount(_ context.Context, userID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, n := range r.store {
		if n.UserID == userID && n.ReadAt == nil {
			count++
		}
	}
	return count, nil
}
