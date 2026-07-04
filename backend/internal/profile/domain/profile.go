// Package domain holds the Profile bounded context's entities and ports.
package domain

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("profile not found")

// Profile is the aggregate: scalar fields plus child collections.
type Profile struct {
	UserID         string
	Headline       string
	About          string
	PhotoURL       string
	Bio            string
	Location       string
	Website        string
	Version        int
	Experiences    []WorkExperience
	Educations     []Education
	Certifications []Certification
	Skills         []string
	Languages      []Language
	Portfolio      []PortfolioLink
}

// Scalars carries the editable top-level fields.
type Scalars struct {
	Headline string
	About    string
	PhotoURL string
	Bio      string
	Location string
	Website  string
}

type WorkExperience struct {
	ID             string
	Title          string
	Company        string
	Location       string
	EmploymentType string
	StartDate      string // ISO date "2006-01-02" or empty
	EndDate        string
	IsCurrent      bool
	Description    string
}

type Education struct {
	ID           string
	School       string
	Degree       string
	FieldOfStudy string
	StartDate    string
	EndDate      string
	Grade        string
	Description  string
}

type Certification struct {
	ID            string
	Name          string
	Issuer        string
	IssueDate     string
	ExpiryDate    string
	CredentialID  string
	CredentialURL string
}

type Language struct {
	Name        string
	Proficiency string
}

type PortfolioLink struct {
	ID    string
	Label string
	URL   string
}

// Repository is the persistence port for the profile aggregate.
type Repository interface {
	// Get returns the full aggregate. A profile row is created lazily if absent.
	Get(ctx context.Context, userID string) (*Profile, error)
	UpdateScalars(ctx context.Context, userID string, s Scalars) error

	AddExperience(ctx context.Context, userID string, e *WorkExperience) error
	UpdateExperience(ctx context.Context, userID string, e WorkExperience) error
	DeleteExperience(ctx context.Context, userID, id string) error

	AddEducation(ctx context.Context, userID string, e *Education) error
	UpdateEducation(ctx context.Context, userID string, e Education) error
	DeleteEducation(ctx context.Context, userID, id string) error

	AddCertification(ctx context.Context, userID string, c *Certification) error
	UpdateCertification(ctx context.Context, userID string, c Certification) error
	DeleteCertification(ctx context.Context, userID, id string) error

	SetSkills(ctx context.Context, userID string, skills []string) error
	SetLanguages(ctx context.Context, userID string, langs []Language) error
	SetPortfolio(ctx context.Context, userID string, links []PortfolioLink) error
}
