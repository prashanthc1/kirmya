package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workspace-app/internal/admin/domain"
)

type fakeRepo struct {
	mu      sync.Mutex
	seq     int
	users   map[string]*domain.UserSummary
	reports map[string]*domain.Report
	posts   map[string]bool
	audits  []string
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		users:   map[string]*domain.UserSummary{},
		reports: map[string]*domain.Report{},
		posts:   map[string]bool{},
	}
}

func (r *fakeRepo) addUser(u *domain.UserSummary) {
	r.users[u.ID] = u
}

func (r *fakeRepo) ListUsers(_ context.Context, f domain.UserFilter) ([]domain.UserSummary, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.UserSummary{}
	for _, u := range r.users {
		if f.Status != "" && u.Status != f.Status {
			continue
		}
		out = append(out, *u)
	}
	return out, len(out), nil
}

func (r *fakeRepo) GetUser(_ context.Context, id string) (*domain.UserSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *u
	return &c, nil
}

func (r *fakeRepo) SetUserStatus(_ context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.Status = status
	return nil
}

func (r *fakeRepo) AssignRole(_ context.Context, userID, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[userID]
	if !ok {
		return domain.ErrNotFound
	}
	for _, ex := range u.Roles {
		if ex == role {
			return nil
		}
	}
	u.Roles = append(u.Roles, role)
	return nil
}

func (r *fakeRepo) RevokeRole(_ context.Context, userID, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[userID]
	if !ok {
		return domain.ErrNotFound
	}
	kept := u.Roles[:0]
	for _, ex := range u.Roles {
		if ex != role {
			kept = append(kept, ex)
		}
	}
	u.Roles = kept
	return nil
}

func (r *fakeRepo) DeletePost(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.posts[id] {
		return domain.ErrNotFound
	}
	delete(r.posts, id)
	return nil
}

func (r *fakeRepo) DeleteComment(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.posts[id] {
		return domain.ErrNotFound
	}
	delete(r.posts, id)
	return nil
}

func (r *fakeRepo) CreateReport(_ context.Context, rep *domain.Report) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	rep.ID = fmt.Sprintf("rep-%d", r.seq)
	rep.CreatedAt = time.Now()
	rep.UpdatedAt = rep.CreatedAt
	c := *rep
	r.reports[rep.ID] = &c
	return nil
}

func (r *fakeRepo) ListReports(_ context.Context, status string) ([]domain.Report, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Report{}
	for _, rep := range r.reports {
		if status != "" && rep.Status != status {
			continue
		}
		out = append(out, *rep)
	}
	return out, nil
}

func (r *fakeRepo) GetReport(_ context.Context, id string) (*domain.Report, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rep, ok := r.reports[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *rep
	return &c, nil
}

func (r *fakeRepo) ResolveReport(_ context.Context, id, status, actionTaken, resolvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rep, ok := r.reports[id]
	if !ok {
		return domain.ErrNotFound
	}
	now := time.Now()
	rep.Status = status
	rep.ActionTaken = actionTaken
	rep.ResolvedBy = resolvedBy
	rep.ResolvedAt = &now
	return nil
}

func (r *fakeRepo) Analytics(_ context.Context) (*domain.Analytics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var a domain.Analytics
	a.Users.Total = len(r.users)
	a.Reports.Open = len(r.reports)
	return &a, nil
}

func (r *fakeRepo) WriteAudit(_ context.Context, actorID, action, targetType, targetID string, _ map[string]any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.audits = append(r.audits, fmt.Sprintf("%s:%s:%s:%s", actorID, action, targetType, targetID))
	return nil
}
