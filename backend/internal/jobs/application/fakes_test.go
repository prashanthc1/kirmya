package application

import (
	"context"
	"fmt"
	"sync"

	"workspace-app/internal/jobs/domain"
)

type fakeRepo struct {
	mu    sync.Mutex
	seq   int
	jobs  map[string]domain.Job
	apps  map[string]domain.Application
	saved map[string]bool // key: userID+"|"+jobID
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{jobs: map[string]domain.Job{}, apps: map[string]domain.Application{}, saved: map[string]bool{}}
}

func (r *fakeRepo) id(p string) string { r.seq++; return fmt.Sprintf("%s-%d", p, r.seq) }

func (r *fakeRepo) CreateJob(_ context.Context, j *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	j.ID = r.id("job")
	r.jobs[j.ID] = *j
	return nil
}

func (r *fakeRepo) GetJob(_ context.Context, id string) (*domain.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return &j, nil
}

func (r *fakeRepo) ListJobs(_ context.Context, f domain.Filter) ([]domain.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Job{}
	for _, j := range r.jobs {
		if f.PostedBy != "" && j.PostedBy != f.PostedBy {
			continue
		}
		out = append(out, j)
	}
	return out, nil
}

func (r *fakeRepo) UpdateJob(_ context.Context, j *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = *j
	return nil
}

func (r *fakeRepo) DeleteJob(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.jobs, id)
	return nil
}

func (r *fakeRepo) CreateApplication(_ context.Context, a *domain.Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = r.id("app")
	r.apps[a.ID] = *a
	return nil
}

func (r *fakeRepo) GetApplication(_ context.Context, id string) (*domain.Application, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.apps[id]
	if !ok {
		return nil, domain.ErrApplicationNotFound
	}
	return &a, nil
}

func (r *fakeRepo) HasApplied(_ context.Context, jobID, userID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, a := range r.apps {
		if a.JobID == jobID && a.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeRepo) ListApplicationsByUser(_ context.Context, userID string) ([]domain.Application, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Application{}
	for _, a := range r.apps {
		if a.UserID == userID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *fakeRepo) ListApplicationsByJob(_ context.Context, jobID string) ([]domain.Application, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Application{}
	for _, a := range r.apps {
		if a.JobID == jobID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *fakeRepo) UpdateApplicationStatus(_ context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a := r.apps[id]
	a.Status = status
	r.apps[id] = a
	return nil
}

func savedKey(userID, jobID string) string { return userID + "|" + jobID }

func (r *fakeRepo) SaveJob(_ context.Context, userID, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.saved[savedKey(userID, jobID)] = true
	return nil
}

func (r *fakeRepo) UnsaveJob(_ context.Context, userID, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.saved, savedKey(userID, jobID))
	return nil
}

func (r *fakeRepo) IsSaved(_ context.Context, userID, jobID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.saved[savedKey(userID, jobID)], nil
}

func (r *fakeRepo) ListSavedJobs(_ context.Context, userID string) ([]domain.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Job{}
	for k := range r.saved {
		if len(k) > len(userID) && k[:len(userID)] == userID {
			jobID := k[len(userID)+1:]
			if j, ok := r.jobs[jobID]; ok {
				out = append(out, j)
			}
		}
	}
	return out, nil
}
