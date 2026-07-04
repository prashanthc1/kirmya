package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workspace-app/internal/resume/domain"
)

type fakeRepo struct {
	mu       sync.Mutex
	seq      int
	resumes  map[string]*domain.Resume
	versions map[string][]domain.Version // resumeID -> versions
	scores   map[string]domain.Score     // versionID -> score
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		resumes:  map[string]*domain.Resume{},
		versions: map[string][]domain.Version{},
		scores:   map[string]domain.Score{},
	}
}

func (r *fakeRepo) id(p string) string { r.seq++; return fmt.Sprintf("%s-%d", p, r.seq) }

func (r *fakeRepo) CreateResume(_ context.Context, res *domain.Resume) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res.ID = r.id("resume")
	res.CreatedAt = time.Now()
	res.UpdatedAt = res.CreatedAt
	c := *res
	r.resumes[res.ID] = &c
	return nil
}

func (r *fakeRepo) GetResume(_ context.Context, id string) (*domain.Resume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, ok := r.resumes[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *res
	return &c, nil
}

func (r *fakeRepo) ListResumesByUser(_ context.Context, userID string) ([]domain.Resume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Resume{}
	for _, res := range r.resumes {
		if res.UserID == userID {
			out = append(out, *res)
		}
	}
	return out, nil
}

func (r *fakeRepo) SoftDeleteResume(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.resumes, id)
	return nil
}

func (r *fakeRepo) NextVersionNo(_ context.Context, resumeID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.versions[resumeID]) + 1, nil
}

func (r *fakeRepo) AddVersion(_ context.Context, v *domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v.ID = r.id("ver")
	v.CreatedAt = time.Now()
	r.versions[v.ResumeID] = append(r.versions[v.ResumeID], *v)
	return nil
}

func (r *fakeRepo) ListVersions(_ context.Context, resumeID string) ([]domain.Version, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vs := r.versions[resumeID]
	out := make([]domain.Version, len(vs))
	copy(out, vs)
	return out, nil
}

func (r *fakeRepo) LatestVersion(_ context.Context, resumeID string) (*domain.Version, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vs := r.versions[resumeID]
	if len(vs) == 0 {
		return nil, domain.ErrNotFound
	}
	v := vs[len(vs)-1]
	return &v, nil
}

func (r *fakeRepo) SaveScore(_ context.Context, s *domain.Score) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.CreatedAt = time.Now()
	r.scores[s.VersionID] = *s
	return nil
}

func (r *fakeRepo) LatestScore(_ context.Context, resumeID string) (*domain.Score, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vs := r.versions[resumeID]
	if len(vs) == 0 {
		return nil, domain.ErrNotFound
	}
	s, ok := r.scores[vs[len(vs)-1].ID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &s, nil
}

// --- fake storage & parser ---

type fakeStorage struct {
	mu    sync.Mutex
	saved map[string][]byte
}

func newFakeStorage() *fakeStorage { return &fakeStorage{saved: map[string][]byte{}} }

func (s *fakeStorage) Save(_ context.Context, key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.saved[key] = data
	return nil
}

type fakeParser struct{ text string }

func (p fakeParser) ExtractText(_, _ string, data []byte) (string, error) {
	if len(data) == 0 {
		return "", domain.ErrEmptyUpload
	}
	return p.text, nil
}
