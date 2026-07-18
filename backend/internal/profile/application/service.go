// Package application implements the Profile use cases over the domain
// repository port.
package application

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"workspace-app/internal/profile/domain"
)

// EventPublisher publishes domain events (the platform bus satisfies this).
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

// Cache is the cache-aside port (the platform cache satisfies this).
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	Delete(ctx context.Context, keys ...string)
}

const (
	eventProfileUpdated   = "ProfileUpdated"
	eventProfilePublished = "profile.published"
	profileCacheTTL       = 10 * time.Minute
)

func profileKey(userID string) string { return "profile:" + userID }

type Service struct {
	repo   domain.Repository
	events EventPublisher
	cache  Cache
}

func NewService(repo domain.Repository, events EventPublisher, cache Cache) *Service {
	return &Service{repo: repo, events: events, cache: cache}
}

// Get returns the draft (working copy) profile aggregate.
func (s *Service) Get(ctx context.Context, userID string) (*domain.Profile, error) {
	if s.cache != nil {
		if b, ok := s.cache.Get(ctx, profileKey(userID)); ok {
			var p domain.Profile
			if json.Unmarshal(b, &p) == nil {
				return &p, nil
			}
		}
	}
	p, err := s.repo.Get(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	s.put(ctx, p)
	return p, nil
}

// GetPublished returns the latest published snapshot.
func (s *Service) GetPublished(ctx context.Context, userID string) (*domain.Profile, error) {
	return s.repo.Get(ctx, userID, false)
}

// put writes the profile to the cache.
func (s *Service) put(ctx context.Context, p *domain.Profile) {
	if s.cache == nil || p == nil {
		return
	}
	if b, err := json.Marshal(p); err == nil {
		s.cache.Set(ctx, profileKey(p.UserID), b, profileCacheTTL)
	}
}

// UpdateProfile applies updates to the aggregate working copy (draft).
func (s *Service) UpdateProfile(ctx context.Context, userID string, expectedVersion int, upd domain.AggregateUpdate) (*domain.Profile, error) {
	if err := upd.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAggregate(ctx, userID, expectedVersion, upd); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

// Publish commits the current draft state, increments version, creates a historical snapshot, and triggers re-indexing.
func (s *Service) Publish(ctx context.Context, userID string, actorID string, ip, ua string) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	p.IsDraft = false
	p.Version++

	// 1. Create a historical snapshot in the DB
	if err := s.repo.CreateVersionSnapshot(ctx, userID, p.Version, p); err != nil {
		return nil, err
	}

	// 2. Mark draft as false in main tables
	draftVal := false
	err = s.repo.UpdateAggregate(ctx, userID, 0, domain.AggregateUpdate{IsDraft: &draftVal})
	if err != nil {
		return nil, err
	}

	// 3. Write Audit Log
	newVal, _ := json.Marshal(p)
	_ = s.repo.WriteAuditLog(ctx, &domain.AuditLogEntry{
		UserID:    userID,
		Section:   "aggregate",
		Action:    "publish",
		ActorID:   actorID,
		NewValue:  newVal,
		IPAddress: ip,
		UserAgent: ua,
		CreatedAt: time.Now(),
	})

	// 4. Publish Event
	if s.events != nil {
		_ = s.events.Publish(ctx, eventProfilePublished, userID, map[string]any{
			"version": p.Version,
		})
	}

	return s.reload(ctx, userID)
}

// Rollback restores the aggregate state to a historical version.
func (s *Service) Rollback(ctx context.Context, userID string, version int, actorID string, ip, ua string) (*domain.Profile, error) {
	snap, err := s.repo.GetVersionSnapshot(ctx, userID, version)
	if err != nil {
		return nil, err
	}

	// Replace current working draft with historical snapshot values
	upd := domain.AggregateUpdate{
		Identity:       &snap.Identity,
		Summary:        &snap.Summary,
		Experiences:    &snap.Experiences,
		Educations:     &snap.Educations,
		Skills:         &snap.Skills,
		Projects:       &snap.Projects,
		Certifications: &snap.Certifications,
		Achievements:   &snap.Achievements,
		Preferences:    &snap.Preferences,
		Privacy:        &snap.Privacy,
	}

	if err := s.repo.UpdateAggregate(ctx, userID, 0, upd); err != nil {
		return nil, err
	}

	// Write Audit Log
	newVal, _ := json.Marshal(snap)
	_ = s.repo.WriteAuditLog(ctx, &domain.AuditLogEntry{
		UserID:    userID,
		Section:   "aggregate",
		Action:    "rollback",
		ActorID:   actorID,
		NewValue:  newVal,
		IPAddress: ip,
		UserAgent: ua,
		CreatedAt: time.Now(),
	})

	return s.reload(ctx, userID)
}

func (s *Service) ListVersions(ctx context.Context, userID string) ([]int, error) {
	return s.repo.ListVersions(ctx, userID)
}

func (s *Service) GetAnalytics(ctx context.Context, profileID string) (*domain.AnalyticsSummary, error) {
	return s.repo.GetAnalytics(ctx, profileID)
}

func (s *Service) RecordView(ctx context.Context, profileID string, actorID *string, ip, ua string) error {
	return s.repo.RecordAnalyticsEvent(ctx, profileID, "view", actorID, ip, ua)
}

func (s *Service) SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) (*domain.Profile, error) {
	if err := s.repo.SetVerificationStatus(ctx, userID, field, verified); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

// --- Backwards compatible helpers updating slices via UpdateAggregate ---

func (s *Service) AddExperience(ctx context.Context, userID string, e *domain.WorkExperience) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	p.Experiences = append(p.Experiences, *e)
	p, err = s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Experiences: &p.Experiences})
	if err != nil {
		return nil, err
	}
	if len(p.Experiences) > 0 {
		e.ID = p.Experiences[len(p.Experiences)-1].ID
	}
	return p, nil
}

func (s *Service) UpdateExperience(ctx context.Context, userID string, e domain.WorkExperience) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	found := false
	for i, ex := range p.Experiences {
		if ex.ID == e.ID {
			p.Experiences[i] = e
			found = true
			break
		}
	}
	if !found {
		return nil, domain.ErrNotFound
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Experiences: &p.Experiences})
}

func (s *Service) DeleteExperience(ctx context.Context, userID, id string) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	var next []domain.WorkExperience
	for _, ex := range p.Experiences {
		if ex.ID != id {
			next = append(next, ex)
		}
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Experiences: &next})
}

func (s *Service) AddEducation(ctx context.Context, userID string, e *domain.Education) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	p.Educations = append(p.Educations, *e)
	p, err = s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Educations: &p.Educations})
	if err != nil {
		return nil, err
	}
	if len(p.Educations) > 0 {
		e.ID = p.Educations[len(p.Educations)-1].ID
	}
	return p, nil
}

func (s *Service) UpdateEducation(ctx context.Context, userID string, e domain.Education) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	found := false
	for i, ed := range p.Educations {
		if ed.ID == e.ID {
			p.Educations[i] = e
			found = true
			break
		}
	}
	if !found {
		return nil, domain.ErrNotFound
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Educations: &p.Educations})
}

func (s *Service) DeleteEducation(ctx context.Context, userID, id string) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	var next []domain.Education
	for _, ed := range p.Educations {
		if ed.ID != id {
			next = append(next, ed)
		}
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Educations: &next})
}

func (s *Service) AddCertification(ctx context.Context, userID string, c *domain.CertificationItem) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	p.Certifications = append(p.Certifications, *c)
	p, err = s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Certifications: &p.Certifications})
	if err != nil {
		return nil, err
	}
	if len(p.Certifications) > 0 {
		c.ID = p.Certifications[len(p.Certifications)-1].ID
	}
	return p, nil
}

func (s *Service) UpdateCertification(ctx context.Context, userID string, c domain.CertificationItem) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	found := false
	for i, cert := range p.Certifications {
		if cert.ID == c.ID {
			p.Certifications[i] = c
			found = true
			break
		}
	}
	if !found {
		return nil, domain.ErrNotFound
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Certifications: &p.Certifications})
}

func (s *Service) DeleteCertification(ctx context.Context, userID, id string) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	var next []domain.CertificationItem
	for _, cert := range p.Certifications {
		if cert.ID != id {
			next = append(next, cert)
		}
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Certifications: &next})
}

func (s *Service) SetSkills(ctx context.Context, userID string, skills []domain.SkillItem) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Skills: &skills})
}

func (s *Service) SetLanguages(ctx context.Context, userID string, langs []domain.LanguageItem) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	p.Identity.Languages = langs
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Identity: &p.Identity})
}

func (s *Service) SetPortfolio(ctx context.Context, userID string, links []domain.ProjectItem) (*domain.Profile, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateProfile(ctx, userID, p.Version, domain.AggregateUpdate{Projects: &links})
}

func (s *Service) AddEndorsement(ctx context.Context, toUserID string, e *domain.EndorsementSummary) (*domain.Profile, error) {
	// Simple mock insert or trigger
	return s.reload(ctx, toUserID)
}

func (s *Service) AddReference(ctx context.Context, userID string, rf *domain.Reference) (*domain.Profile, error) {
	return s.reload(ctx, userID)
}

func (s *Service) UpdateReference(ctx context.Context, userID string, rf domain.Reference) (*domain.Profile, error) {
	return s.reload(ctx, userID)
}

func (s *Service) DeleteReference(ctx context.Context, userID, id string) (*domain.Profile, error) {
	return s.reload(ctx, userID)
}

func (s *Service) AddConsentLog(ctx context.Context, cl *domain.ConsentLog) error {
	return nil
}

// reload re-reads the aggregate after a write, recalculates completeness score,
// refreshes the cache write-through with the fresh value, and emits ProfileUpdated.
func (s *Service) reload(ctx context.Context, userID string) (*domain.Profile, error) {
	p, err := s.repo.Get(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	score := p.CalculateCompleteness()
	p.ProfileCompletenessScore = score

	if err := s.repo.UpdateCompletenessScore(ctx, userID, score); err != nil {
		log.Printf("failed to update profile completeness score: %v", err)
	}

	s.put(ctx, p)
	if s.events != nil {
		_ = s.events.Publish(ctx, eventProfileUpdated, userID, nil)
	}
	return p, nil
}
