// Package application implements the Jobs use cases over the domain ports.
package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"workspace-app/internal/jobs/domain"
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
	eventJobPosted  = "JobPosted"
	eventJobApplied = "JobApplied"

	jobCacheTTL = 10 * time.Minute
)

func jobKey(id string) string { return "job:" + id }

// ValidationError is returned for invalid input (mapped to HTTP 400).
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

// ErrForbidden indicates the caller does not own the target resource.
var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo   domain.Repository
	events EventPublisher
	cache  Cache
}

func NewService(repo domain.Repository, events EventPublisher, cache Cache) *Service {
	return &Service{repo: repo, events: events, cache: cache}
}

// PostJobInput is the create-job command.
type PostJobInput struct {
	Title, Company, Location, Description, Salary, JobType string
}

func (s *Service) PostJob(ctx context.Context, userID string, in PostJobInput) (*domain.Job, error) {
	in.Title, in.Company = strings.TrimSpace(in.Title), strings.TrimSpace(in.Company)
	in.Description = strings.TrimSpace(in.Description)
	if in.Title == "" || in.Company == "" {
		return nil, ValidationError{"title and company are required"}
	}
	if len(in.Description) < 20 {
		return nil, ValidationError{"description must be at least 20 characters"}
	}
	job := &domain.Job{
		Title: in.Title, Company: in.Company, Location: in.Location,
		Description: in.Description, Salary: in.Salary, JobType: in.JobType, PostedBy: userID,
	}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}
	s.publish(ctx, eventJobPosted, job.ID, map[string]any{"posted_by": userID, "title": job.Title})
	return job, nil
}

// GetJob returns a single job, served cache-aside (hit → cached; miss → load and
// populate). Write paths below read from the repo directly for authoritative data.
func (s *Service) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	if s.cache != nil {
		if b, ok := s.cache.Get(ctx, jobKey(id)); ok {
			var j domain.Job
			if json.Unmarshal(b, &j) == nil {
				return &j, nil
			}
		}
	}
	j, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		if b, err := json.Marshal(j); err == nil {
			s.cache.Set(ctx, jobKey(id), b, jobCacheTTL)
		}
	}
	return j, nil
}

// invalidateJob drops a job from the cache after a mutation.
func (s *Service) invalidateJob(ctx context.Context, id string) {
	if s.cache != nil {
		s.cache.Delete(ctx, jobKey(id))
	}
}

func (s *Service) SearchJobs(ctx context.Context, f domain.Filter) ([]domain.Job, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	return s.repo.ListJobs(ctx, f)
}

func (s *Service) UpdateJob(ctx context.Context, userID, jobID string, in PostJobInput) (*domain.Job, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job.PostedBy != userID {
		return nil, ErrForbidden
	}
	// Patch only provided fields.
	if in.Title != "" {
		job.Title = in.Title
	}
	if in.Company != "" {
		job.Company = in.Company
	}
	if in.Location != "" {
		job.Location = in.Location
	}
	if in.Description != "" {
		job.Description = in.Description
	}
	if in.Salary != "" {
		job.Salary = in.Salary
	}
	if in.JobType != "" {
		job.JobType = in.JobType
	}
	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return nil, err
	}
	s.invalidateJob(ctx, jobID)
	return job, nil
}

func (s *Service) DeleteJob(ctx context.Context, userID, jobID string) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	if job.PostedBy != userID {
		return ErrForbidden
	}
	if err := s.repo.DeleteJob(ctx, jobID); err != nil {
		return err
	}
	s.invalidateJob(ctx, jobID)
	return nil
}

func (s *Service) Apply(ctx context.Context, userID, jobID, coverLetter string) (*domain.Application, error) {
	if len(strings.TrimSpace(coverLetter)) < 10 {
		return nil, ValidationError{"cover letter must be at least 10 characters"}
	}
	if _, err := s.repo.GetJob(ctx, jobID); err != nil {
		return nil, err
	}
	applied, err := s.repo.HasApplied(ctx, jobID, userID)
	if err != nil {
		return nil, err
	}
	if applied {
		return nil, domain.ErrAlreadyApplied
	}
	app := &domain.Application{JobID: jobID, UserID: userID, Status: domain.StatusPending, CoverLetter: coverLetter}
	if err := s.repo.CreateApplication(ctx, app); err != nil {
		return nil, err
	}
	s.publish(ctx, eventJobApplied, app.ID, map[string]any{"job_id": jobID, "user_id": userID})
	return app, nil
}

func (s *Service) MyApplications(ctx context.Context, userID string) ([]domain.Application, error) {
	return s.repo.ListApplicationsByUser(ctx, userID)
}

// JobApplicants returns applicants for a job the caller owns.
func (s *Service) JobApplicants(ctx context.Context, userID, jobID string) ([]domain.Application, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job.PostedBy != userID {
		return nil, ErrForbidden
	}
	return s.repo.ListApplicationsByJob(ctx, jobID)
}

// UpdateApplicationStatus lets the owning recruiter advance an application.
func (s *Service) UpdateApplicationStatus(ctx context.Context, userID, appID, status string) (*domain.Application, error) {
	if !domain.ValidStatuses[status] {
		return nil, ValidationError{"invalid status"}
	}
	app, err := s.repo.GetApplication(ctx, appID)
	if err != nil {
		return nil, err
	}
	job, err := s.repo.GetJob(ctx, app.JobID)
	if err != nil {
		return nil, err
	}
	if job.PostedBy != userID {
		return nil, ErrForbidden
	}
	if err := s.repo.UpdateApplicationStatus(ctx, appID, status); err != nil {
		return nil, err
	}
	return s.repo.GetApplication(ctx, appID)
}

// ToggleSave saves or unsaves a job for the user, returning the new saved state.
func (s *Service) ToggleSave(ctx context.Context, userID, jobID string) (bool, error) {
	if _, err := s.repo.GetJob(ctx, jobID); err != nil {
		return false, err
	}
	saved, err := s.repo.IsSaved(ctx, userID, jobID)
	if err != nil {
		return false, err
	}
	if saved {
		return false, s.repo.UnsaveJob(ctx, userID, jobID)
	}
	return true, s.repo.SaveJob(ctx, userID, jobID)
}

func (s *Service) SavedJobs(ctx context.Context, userID string) ([]domain.Job, error) {
	return s.repo.ListSavedJobs(ctx, userID)
}

func (s *Service) publish(ctx context.Context, evt, aggID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, evt, aggID, payload)
	}
}
