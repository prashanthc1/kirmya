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
	c.Skills = append([]domain.ProfileSkill(nil), p.Skills...)
	c.Languages = append([]domain.Language(nil), p.Languages...)
	c.Portfolio = append([]domain.PortfolioLink(nil), p.Portfolio...)
	c.SupportsNeeded = append([]string(nil), p.SupportsNeeded...)
	c.RelocationLocations = append([]string(nil), p.RelocationLocations...)
	c.DesiredRoles = append([]string(nil), p.DesiredRoles...)
	c.DesiredIndustries = append([]string(nil), p.DesiredIndustries...)
	c.Endorsements = append([]domain.Endorsement(nil), p.Endorsements...)
	c.References = append([]domain.Reference(nil), p.References...)
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
	p.Pronouns, p.CareerStatus = s.Pronouns, s.CareerStatus
	p.TransitionReason, p.TargetComebackTimeline = s.TransitionReason, s.TargetComebackTimeline
	p.SupportsNeeded = s.SupportsNeeded
	p.OpenToRemote, p.OpenToRelocation = s.OpenToRemote, s.OpenToRelocation
	p.RelocationLocations = s.RelocationLocations
	p.DesiredRoles, p.DesiredIndustries = s.DesiredRoles, s.DesiredIndustries
	p.EmploymentType = s.EmploymentType
	p.SalaryMin, p.SalaryMax, p.SalaryCurrency = s.SalaryMin, s.SalaryMax, s.SalaryCurrency
	p.SalaryVisible = s.SalaryVisible
	p.WorkMode = s.WorkMode
	p.AvailabilityDate, p.NoticePeriod = s.AvailabilityDate, s.NoticePeriod
	p.ReferralEligible = s.ReferralEligible
	p.CareerNarrative, p.CoachingMetadata = s.CareerNarrative, s.CoachingMetadata
	p.WorkAuthStatus, p.PassportNationality = s.WorkAuthStatus, s.PassportNationality
	p.DrivingLicenseBool, p.DrivingLicenseType = s.DrivingLicenseBool, s.DrivingLicenseType
	p.PreferredContactChannel, p.AccessibilityNeeds = s.PreferredContactChannel, s.AccessibilityNeeds
	p.VideoIntroURL = s.VideoIntroURL
	p.WillingToMentor = s.WillingToMentor
	p.BackgroundCheckConsent, p.BackgroundCheckConsentAt = s.BackgroundCheckConsent, s.BackgroundCheckConsentAt
	p.JobAlertFrequency, p.JobAlertChannel = s.JobAlertFrequency, s.JobAlertChannel
	p.VisibilityProfile = s.VisibilityProfile
	p.VisibilitySalary = s.VisibilitySalary
	p.VisibilityTransitionReason = s.VisibilityTransitionReason
	p.VisibilityExperience = s.VisibilityExperience
	p.VisibilityEducation = s.VisibilityEducation
	p.VisibilityCertifications = s.VisibilityCertifications
	p.VisibilitySkills = s.VisibilitySkills
	p.VisibilityPortfolio = s.VisibilityPortfolio
	p.VisibilityReferences = s.VisibilityReferences
	return nil
}

func (r *fakeRepo) UpdateAggregate(ctx context.Context, userID string, expectedVersion int, u domain.AggregateUpdate) error {
	r.mu.Lock()
	if p := r.get(userID); expectedVersion > 0 && p.Version != expectedVersion {
		r.mu.Unlock()
		return domain.ErrOptimisticLock
	}
	r.mu.Unlock()

	if err := r.UpdateScalars(ctx, userID, u.Scalars); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	p.Version++

	if u.Experiences != nil {
		out := make([]domain.WorkExperience, 0, len(*u.Experiences))
		for _, e := range *u.Experiences {
			if e.ID == "" {
				e.ID = r.id("exp")
			}
			out = append(out, e)
		}
		p.Experiences = out
	}
	if u.Educations != nil {
		out := make([]domain.Education, 0, len(*u.Educations))
		for _, e := range *u.Educations {
			if e.ID == "" {
				e.ID = r.id("edu")
			}
			out = append(out, e)
		}
		p.Educations = out
	}
	if u.Certifications != nil {
		out := make([]domain.Certification, 0, len(*u.Certifications))
		for _, c := range *u.Certifications {
			if c.ID == "" {
				c.ID = r.id("cert")
			}
			out = append(out, c)
		}
		p.Certifications = out
	}
	if u.Skills != nil {
		p.Skills = append([]domain.ProfileSkill(nil), *u.Skills...)
	}
	if u.Languages != nil {
		p.Languages = append([]domain.Language(nil), *u.Languages...)
	}
	if u.Portfolio != nil {
		p.Portfolio = append([]domain.PortfolioLink(nil), *u.Portfolio...)
	}
	if u.References != nil {
		out := make([]domain.Reference, 0, len(*u.References))
		for _, rf := range *u.References {
			if rf.ID == "" {
				rf.ID = r.id("ref")
			}
			out = append(out, rf)
		}
		p.References = out
	}
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

func (r *fakeRepo) SetSkills(_ context.Context, userID string, skills []domain.ProfileSkill) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.get(userID).Skills = append([]domain.ProfileSkill(nil), skills...)
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

func (r *fakeRepo) AddEndorsement(_ context.Context, toUserID string, e *domain.Endorsement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.ID = r.id("end")
	e.CreatedAt = "2026-07-06T00:00:00Z"
	p := r.get(toUserID)
	p.Endorsements = append(p.Endorsements, *e)
	return nil
}

func (r *fakeRepo) AddReference(_ context.Context, userID string, rf *domain.Reference) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rf.ID = r.id("ref")
	p := r.get(userID)
	p.References = append(p.References, *rf)
	return nil
}

func (r *fakeRepo) UpdateReference(_ context.Context, userID string, rf domain.Reference) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.References {
		if p.References[i].ID == rf.ID {
			p.References[i] = rf
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) DeleteReference(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	for i := range p.References {
		if p.References[i].ID == id {
			p.References = append(p.References[:i], p.References[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) AddConsentLog(_ context.Context, cl *domain.ConsentLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cl.ID = r.id("con")
	cl.CreatedAt = "2026-07-06T00:00:00Z"
	return nil
}

func (r *fakeRepo) SetVerificationStatus(_ context.Context, userID string, field string, verified bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	switch field {
	case "phone_verified":
		p.PhoneVerified = verified
	case "linkedin_verified":
		p.LinkedinVerified = verified
	case "id_verified":
		p.IdVerified = verified
	}
	return nil
}

func (r *fakeRepo) UpdateCalculatedFields(_ context.Context, userID string, completeness int, avgResponse float64, lastActive string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	p.ProfileCompletenessScore = completeness
	p.AvgResponseTimeHours = avgResponse
	p.LastActiveAt = lastActive
	return nil
}
