// Package application implements the Profile use cases over the domain
// repository port.
package application

import (
	"context"
	"encoding/json"
	"time"

	"workspace-app/internal/profile/domain"
)

// EventPublisher publishes domain events (the platform bus satisfies this).
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

// Cache is the cache-aside port (the platform cache satisfies this). A nil cache
// disables caching; the platform's no-op cache also makes every call a no-op.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	Delete(ctx context.Context, keys ...string)
}

const (
	eventProfileUpdated = "ProfileUpdated"
	profileCacheTTL     = 10 * time.Minute
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

// Get returns the full profile aggregate (created lazily if missing). Reads are
// served cache-aside: hit → return cached; miss → load, then populate the cache.
func (s *Service) Get(ctx context.Context, userID string) (*domain.Profile, error) {
	if s.cache != nil {
		if b, ok := s.cache.Get(ctx, profileKey(userID)); ok {
			var p domain.Profile
			if json.Unmarshal(b, &p) == nil {
				return &p, nil
			}
		}
	}
	p, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.put(ctx, p)
	return p, nil
}

// put writes the profile to the cache (best-effort).
func (s *Service) put(ctx context.Context, p *domain.Profile) {
	if s.cache == nil || p == nil {
		return
	}
	if b, err := json.Marshal(p); err == nil {
		s.cache.Set(ctx, profileKey(p.UserID), b, profileCacheTTL)
	}
}

func (s *Service) UpdateScalars(ctx context.Context, userID string, sc domain.Scalars) (*domain.Profile, error) {
	// Validate
	temp := &domain.Profile{
		SalaryMin:        sc.SalaryMin,
		SalaryMax:        sc.SalaryMax,
		TransitionReason: sc.TransitionReason,
		CareerStatus:     sc.CareerStatus,
	}
	if err := temp.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateScalars(ctx, userID, sc); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) AddExperience(ctx context.Context, userID string, e *domain.WorkExperience) (*domain.Profile, error) {
	temp := &domain.Profile{
		Experiences: []domain.WorkExperience{*e},
	}
	if err := temp.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.AddExperience(ctx, userID, e); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) UpdateExperience(ctx context.Context, userID string, e domain.WorkExperience) (*domain.Profile, error) {
	temp := &domain.Profile{
		Experiences: []domain.WorkExperience{e},
	}
	if err := temp.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateExperience(ctx, userID, e); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) DeleteExperience(ctx context.Context, userID, id string) (*domain.Profile, error) {
	if err := s.repo.DeleteExperience(ctx, userID, id); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) AddEducation(ctx context.Context, userID string, e *domain.Education) (*domain.Profile, error) {
	temp := &domain.Profile{
		Educations: []domain.Education{*e},
	}
	if err := temp.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.AddEducation(ctx, userID, e); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) UpdateEducation(ctx context.Context, userID string, e domain.Education) (*domain.Profile, error) {
	temp := &domain.Profile{
		Educations: []domain.Education{e},
	}
	if err := temp.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateEducation(ctx, userID, e); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) DeleteEducation(ctx context.Context, userID, id string) (*domain.Profile, error) {
	if err := s.repo.DeleteEducation(ctx, userID, id); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) AddCertification(ctx context.Context, userID string, c *domain.Certification) (*domain.Profile, error) {
	if err := s.repo.AddCertification(ctx, userID, c); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) UpdateCertification(ctx context.Context, userID string, c domain.Certification) (*domain.Profile, error) {
	if err := s.repo.UpdateCertification(ctx, userID, c); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) DeleteCertification(ctx context.Context, userID, id string) (*domain.Profile, error) {
	if err := s.repo.DeleteCertification(ctx, userID, id); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) SetSkills(ctx context.Context, userID string, skills []domain.ProfileSkill) (*domain.Profile, error) {
	if err := s.repo.SetSkills(ctx, userID, skills); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) SetLanguages(ctx context.Context, userID string, langs []domain.Language) (*domain.Profile, error) {
	if err := s.repo.SetLanguages(ctx, userID, langs); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) SetPortfolio(ctx context.Context, userID string, links []domain.PortfolioLink) (*domain.Profile, error) {
	if err := s.repo.SetPortfolio(ctx, userID, links); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

// --- new features use cases ---

func (s *Service) AddEndorsement(ctx context.Context, toUserID string, e *domain.Endorsement) (*domain.Profile, error) {
	if err := s.repo.AddEndorsement(ctx, toUserID, e); err != nil {
		return nil, err
	}
	return s.reload(ctx, toUserID)
}

func (s *Service) AddReference(ctx context.Context, userID string, rf *domain.Reference) (*domain.Profile, error) {
	if err := s.repo.AddReference(ctx, userID, rf); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) UpdateReference(ctx context.Context, userID string, rf domain.Reference) (*domain.Profile, error) {
	if err := s.repo.UpdateReference(ctx, userID, rf); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) DeleteReference(ctx context.Context, userID, id string) (*domain.Profile, error) {
	if err := s.repo.DeleteReference(ctx, userID, id); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

func (s *Service) AddConsentLog(ctx context.Context, cl *domain.ConsentLog) error {
	if err := s.repo.AddConsentLog(ctx, cl); err != nil {
		return err
	}
	// If it is background_check consent, update the profile aggregate fields
	if cl.ConsentType == "background_check" {
		p, err := s.repo.Get(ctx, cl.UserID)
		if err == nil {
			sc := domain.Scalars{
				Headline: p.Headline, About: p.About, PhotoURL: p.PhotoURL, Bio: p.Bio, Location: p.Location, Website: p.Website,
				Pronouns: p.Pronouns, CareerStatus: p.CareerStatus, TransitionReason: p.TransitionReason, TargetComebackTimeline: p.TargetComebackTimeline,
				OpenToRemote: p.OpenToRemote, OpenToRelocation: p.OpenToRelocation, EmploymentType: p.EmploymentType,
				SalaryMin: p.SalaryMin, SalaryMax: p.SalaryMax, SalaryCurrency: p.SalaryCurrency, SalaryVisible: p.SalaryVisible,
				WorkMode: p.WorkMode, AvailabilityDate: p.AvailabilityDate, NoticePeriod: p.NoticePeriod,
				ReferralEligible: p.ReferralEligible, CareerNarrative: p.CareerNarrative, CoachingMetadata: p.CoachingMetadata,
				WorkAuthStatus: p.WorkAuthStatus, PassportNationality: p.PassportNationality, DrivingLicenseBool: p.DrivingLicenseBool, DrivingLicenseType: p.DrivingLicenseType,
				PreferredContactChannel: p.PreferredContactChannel, AccessibilityNeeds: p.AccessibilityNeeds, VideoIntroURL: p.VideoIntroURL,
				WillingToMentor:   p.WillingToMentor,
				JobAlertFrequency: p.JobAlertFrequency, JobAlertChannel: p.JobAlertChannel,
				VisibilityProfile: p.VisibilityProfile, VisibilitySalary: p.VisibilitySalary, VisibilityTransitionReason: p.VisibilityTransitionReason,
				VisibilityExperience: p.VisibilityExperience, VisibilityEducation: p.VisibilityEducation, VisibilityCertifications: p.VisibilityCertifications,
				VisibilitySkills: p.VisibilitySkills, VisibilityPortfolio: p.VisibilityPortfolio, VisibilityReferences: p.VisibilityReferences,
				SupportsNeeded: p.SupportsNeeded, RelocationLocations: p.RelocationLocations, DesiredRoles: p.DesiredRoles, DesiredIndustries: p.DesiredIndustries,
				// Consent fields
				BackgroundCheckConsent:   cl.Consented,
				BackgroundCheckConsentAt: cl.CreatedAt,
			}
			_ = s.repo.UpdateScalars(ctx, cl.UserID, sc)
			_, _ = s.reload(ctx, cl.UserID)
		}
	}
	return nil
}

func (s *Service) SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) (*domain.Profile, error) {
	if err := s.repo.SetVerificationStatus(ctx, userID, field, verified); err != nil {
		return nil, err
	}
	return s.reload(ctx, userID)
}

// reload re-reads the aggregate after a write, recalculates completeness score,
// refreshes the cache write-through with the fresh value, and emits ProfileUpdated.
func (s *Service) reload(ctx context.Context, userID string) (*domain.Profile, error) {
	p, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Recalculate completeness score
	score := calculateCompletenessScore(p)

	// Mock or calculate average response time
	avgResponse := p.AvgResponseTimeHours
	if avgResponse == 0 {
		avgResponse = 2.5 // default/mock
	}

	// Update calculated fields in repository
	nowStr := time.Now().UTC().Format(time.RFC3339)
	_ = s.repo.UpdateCalculatedFields(ctx, userID, score, avgResponse, nowStr)

	// Fetch again to get the updated calculated fields
	p, err = s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.put(ctx, p)
	if s.events != nil {
		_ = s.events.Publish(ctx, eventProfileUpdated, userID, nil)
	}
	return p, nil
}

func calculateCompletenessScore(p *domain.Profile) int {
	score := 0
	if p.Headline != "" {
		score += 10
	}
	if p.About != "" || p.Bio != "" {
		score += 10
	}
	if p.Location != "" {
		score += 10
	}
	if p.CareerStatus != "" {
		score += 10
	}
	if len(p.Experiences) > 0 {
		score += 20
	}
	if len(p.Educations) > 0 {
		score += 15
	}
	if len(p.Skills) > 0 {
		score += 15
	}
	if p.PreferredContactChannel != "" {
		score += 10
	}
	return score
}
