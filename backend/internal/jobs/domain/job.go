// Package domain holds the Jobs bounded context entities and ports.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrJobNotFound         = errors.New("job not found")
	ErrApplicationNotFound = errors.New("application not found")
	ErrAlreadyApplied      = errors.New("already applied to this job")
)

// Application status values (job-seeker pipeline).
const (
	StatusPending      = "pending"
	StatusReviewed     = "reviewed"
	StatusInterviewing = "interviewing"
	StatusOffer        = "offer"
	StatusHired        = "hired"
	StatusRejected     = "rejected"
	StatusWithdrawn    = "withdrawn"
)

// ValidStatuses is the set of allowed application statuses.
var ValidStatuses = map[string]bool{
	StatusPending: true, StatusReviewed: true, StatusInterviewing: true,
	StatusOffer: true, StatusHired: true, StatusRejected: true, StatusWithdrawn: true,
}

type Job struct {
	ID          string
	Title       string
	Company     string
	Location    string
	Description string
	Salary      string
	JobType     string
	PostedBy    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Application struct {
	ID          string
	JobID       string
	UserID      string
	Status      string
	CoverLetter string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Filter parameterizes job search.
type Filter struct {
	Keyword  string
	Location string
	JobType  string
	// PostedBy, when set, restricts results to jobs posted by the given user.
	// Used by the recruiter dashboard to list a recruiter's own postings.
	PostedBy string
	Limit    int
}

// Repository is the persistence port for jobs and applications.
type Repository interface {
	CreateJob(ctx context.Context, j *Job) error
	GetJob(ctx context.Context, id string) (*Job, error)
	ListJobs(ctx context.Context, f Filter) ([]Job, error)
	UpdateJob(ctx context.Context, j *Job) error
	DeleteJob(ctx context.Context, id string) error

	CreateApplication(ctx context.Context, a *Application) error
	GetApplication(ctx context.Context, id string) (*Application, error)
	HasApplied(ctx context.Context, jobID, userID string) (bool, error)
	ListApplicationsByUser(ctx context.Context, userID string) ([]Application, error)
	ListApplicationsByJob(ctx context.Context, jobID string) ([]Application, error)
	UpdateApplicationStatus(ctx context.Context, id, status string) error

	SaveJob(ctx context.Context, userID, jobID string) error
	UnsaveJob(ctx context.Context, userID, jobID string) error
	IsSaved(ctx context.Context, userID, jobID string) (bool, error)
	ListSavedJobs(ctx context.Context, userID string) ([]Job, error)
}
