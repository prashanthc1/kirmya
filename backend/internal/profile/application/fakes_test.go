package application

import (
	"context"
	"fmt"
	"sync"

	"workspace-app/internal/profile/domain"
)

// fakeRepo is an in-memory domain.Repository for unit tests.
type fakeRepo struct {
	mu    sync.Mutex
	seq   int
	store map[string]*domain.Profile
}

func newFakeRepo() *fakeRepo { return &fakeRepo{store: map[string]*domain.Profile{}} }

func (r *fakeRepo) get(userID string) *domain.Profile {
	p, ok := r.store[userID]
	if !ok {
		p = &domain.Profile{UserID: userID}
		r.store[userID] = p
	}
	return p
}

func (r *fakeRepo) id(prefix string) string {
	r.seq++
	return fmt.Sprintf("%s-%d", prefix, r.seq)
}

func clone(p *domain.Profile) *domain.Profile {
	c := *p
	c.Experiences = append([]domain.WorkExperience(nil), p.Experiences...)
	c.Educations = append([]domain.Education(nil), p.Educations...)
	c.Certifications = append([]domain.Certification(nil), p.Certifications...)
	c.Skills = append([]string(nil), p.Skills...)
	c.Languages = append([]domain.Language(nil), p.Languages...)
	c.Portfolio = append([]domain.PortfolioLink(nil), p.Portfolio...)
	return &c
}

func (r *fakeRepo) Get(_ context.Context, userID string) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return clone(r.get(userID)), nil
}

func (r *fakeRepo) UpdateScalars(_ context.Context, userID string, s domain.Scalars) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	p.Headline, p.About, p.PhotoURL = s.Headline, s.About, s.PhotoURL
	p.Bio, p.Location, p.Website = s.Bio, s.Location, s.Website
	return nil
}

func (r *fakeRepo) AddExperience(_ context.Context, userID string, e *domain.WorkExperience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.ID = r.id("exp")
	p := r.get(userID)
	p.Experiences = append(p.Experiences, *e)
	return nil
}

func (r *fakeRepo) UpdateExperience(_ context.Context, userID string, e domain.WorkExperience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Experiences {
		if p.Experiences[i].ID == e.ID {
			p.Experiences[i] = e
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) DeleteExperience(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Experiences {
		if p.Experiences[i].ID == id {
			p.Experiences = append(p.Experiences[:i], p.Experiences[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) AddEducation(_ context.Context, userID string, e *domain.Education) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.ID = r.id("edu")
	p := r.get(userID)
	p.Educations = append(p.Educations, *e)
	return nil
}

func (r *fakeRepo) UpdateEducation(_ context.Context, userID string, e domain.Education) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Educations {
		if p.Educations[i].ID == e.ID {
			p.Educations[i] = e
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) DeleteEducation(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Educations {
		if p.Educations[i].ID == id {
			p.Educations = append(p.Educations[:i], p.Educations[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) AddCertification(_ context.Context, userID string, c *domain.Certification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = r.id("cert")
	p := r.get(userID)
	p.Certifications = append(p.Certifications, *c)
	return nil
}

func (r *fakeRepo) UpdateCertification(_ context.Context, userID string, c domain.Certification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Certifications {
		if p.Certifications[i].ID == c.ID {
			p.Certifications[i] = c
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) DeleteCertification(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.Certifications {
		if p.Certifications[i].ID == id {
			p.Certifications = append(p.Certifications[:i], p.Certifications[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) SetSkills(_ context.Context, userID string, skills []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.get(userID).Skills = append([]string(nil), skills...)
	return nil
}

func (r *fakeRepo) SetLanguages(_ context.Context, userID string, langs []domain.Language) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.get(userID).Languages = append([]domain.Language(nil), langs...)
	return nil
}

func (r *fakeRepo) SetPortfolio(_ context.Context, userID string, links []domain.PortfolioLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.get(userID).Portfolio = append([]domain.PortfolioLink(nil), links...)
	return nil
}
