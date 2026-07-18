package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"workspace-app/internal/profile/domain"
)

// fakeRepo is an in-memory domain.Repository for unit tests.
type fakeRepo struct {
	mu        sync.Mutex
	seq       int
	store     map[string]*domain.Profile
	snapshots map[string]map[int]*domain.Profile
	auditLogs []*domain.AuditLogEntry
	analytics map[string][]string
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		store:     map[string]*domain.Profile{},
		snapshots: map[string]map[int]*domain.Profile{},
		analytics: map[string][]string{},
	}
}

func (r *fakeRepo) get(userID string) *domain.Profile {
	p, ok := r.store[userID]
	if !ok {
		p = &domain.Profile{UserID: userID}
		r.store[userID] = p
	}
	return p
}

func clone(p *domain.Profile) *domain.Profile {
	b, _ := json.Marshal(p)
	var c domain.Profile
	_ = json.Unmarshal(b, &c)
	return &c
}

func (r *fakeRepo) Get(_ context.Context, userID string, includeDraft bool) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !includeDraft {
		snaps, ok := r.snapshots[userID]
		if !ok || len(snaps) == 0 {
			return nil, errors.New("no published version")
		}
		maxVer := 0
		for v := range snaps {
			if v > maxVer {
				maxVer = v
			}
		}
		return clone(snaps[maxVer]), nil
	}

	return clone(r.get(userID)), nil
}

func (r *fakeRepo) id(prefix string) string {
	r.seq++
	return fmt.Sprintf("%s-%d", prefix, r.seq)
}

func (r *fakeRepo) UpdateAggregate(ctx context.Context, userID string, expectedVersion int, u domain.AggregateUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.get(userID)
	if expectedVersion > 0 && p.Version != expectedVersion {
		return domain.ErrOptimisticLock
	}

	if u.Identity != nil {
		p.Identity = *u.Identity
	}
	if u.Summary != nil {
		p.Summary = *u.Summary
	}
	if u.Experiences != nil {
		for i, e := range *u.Experiences {
			if e.ID == "" {
				(*u.Experiences)[i].ID = r.id("exp")
			}
		}
		p.Experiences = *u.Experiences
	}
	if u.Educations != nil {
		for i, e := range *u.Educations {
			if e.ID == "" {
				(*u.Educations)[i].ID = r.id("edu")
			}
		}
		p.Educations = *u.Educations
	}
	if u.Skills != nil {
		p.Skills = *u.Skills
	}
	if u.Projects != nil {
		for i, pr := range *u.Projects {
			if pr.ID == "" {
				(*u.Projects)[i].ID = r.id("proj")
			}
		}
		p.Projects = *u.Projects
	}
	if u.Certifications != nil {
		for i, c := range *u.Certifications {
			if c.ID == "" {
				(*u.Certifications)[i].ID = r.id("cert")
			}
		}
		p.Certifications = *u.Certifications
	}
	if u.Achievements != nil {
		for i, a := range *u.Achievements {
			if a.ID == "" {
				(*u.Achievements)[i].ID = r.id("ach")
			}
		}
		p.Achievements = *u.Achievements
	}
	if u.Preferences != nil {
		p.Preferences = *u.Preferences
	}
	if u.Privacy != nil {
		p.Privacy = *u.Privacy
	}
	if u.IsDraft != nil {
		p.IsDraft = *u.IsDraft
	}

	p.Version++
	return nil
}

func (r *fakeRepo) CreateVersionSnapshot(ctx context.Context, userID string, version int, p *domain.Profile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.snapshots[userID] == nil {
		r.snapshots[userID] = map[int]*domain.Profile{}
	}
	r.snapshots[userID][version] = clone(p)
	return nil
}

func (r *fakeRepo) GetVersionSnapshot(ctx context.Context, userID string, version int) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	snaps, ok := r.snapshots[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	p, ok := snaps[version]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return clone(p), nil
}

func (r *fakeRepo) ListVersions(ctx context.Context, userID string) ([]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []int
	for v := range r.snapshots[userID] {
		out = append(out, v)
	}
	return out, nil
}

func (r *fakeRepo) WriteAuditLog(ctx context.Context, log *domain.AuditLogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.auditLogs = append(r.auditLogs, log)
	return nil
}

func (r *fakeRepo) RecordAnalyticsEvent(ctx context.Context, profileID string, eventType string, actorID *string, ip, ua string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.analytics[profileID] = append(r.analytics[profileID], eventType)
	return nil
}

func (r *fakeRepo) GetAnalytics(ctx context.Context, profileID string) (*domain.AnalyticsSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sum domain.AnalyticsSummary
	for _, t := range r.analytics[profileID] {
		switch t {
		case "view":
			sum.ProfileViews++
		case "search_appearance":
			sum.SearchAppearances++
		case "recruiter_view":
			sum.RecruiterViews++
		case "resume_download":
			sum.ResumeDownloads++
		}
	}
	return &sum, nil
}

func (r *fakeRepo) SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.get(userID)
	switch field {
	case "email_verified":
		p.Verification.EmailVerified = verified
	case "phone_verified":
		p.Verification.PhoneVerified = verified
	case "identity_verified":
		p.Verification.IdentityVerified = verified
	case "employment_verified":
		p.Verification.EmploymentVerified = verified
	case "education_verified":
		p.Verification.EducationVerified = verified
	case "certification_verified":
		p.Verification.CertificationVerified = verified
	}
	return nil
}

func (r *fakeRepo) UpdateCompletenessScore(ctx context.Context, userID string, score int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.get(userID)
	p.ProfileCompletenessScore = score
	return nil
}
