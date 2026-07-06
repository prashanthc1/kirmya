// Package domain holds the Profile bounded context's entities and ports.
package domain

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("profile not found")

// Profile is the aggregate: scalar fields plus child collections.
type Profile struct {
	UserID   string
	Headline string
	About    string
	PhotoURL string
	Bio      string
	Location string
	Website  string
	Version  int

	// Core Identity
	Pronouns     string
	CareerStatus string

	// Career Recovery
	TransitionReason       string // Encrypted at rest
	TargetComebackTimeline string
	SupportsNeeded         []string

	// Mobility & Preferences
	OpenToRemote        bool
	OpenToRelocation    bool
	RelocationLocations []string
	DesiredRoles        []string
	DesiredIndustries   []string
	EmploymentType      string
	SalaryMin           int    // Encrypted at rest
	SalaryMax           int    // Encrypted at rest
	SalaryCurrency      string // Encrypted at rest
	SalaryVisible       bool
	WorkMode            string
	AvailabilityDate    string
	NoticePeriod        string

	// Trust & Verification
	ReferralEligible bool
	EmailVerified    bool
	PhoneVerified    bool
	LinkedinVerified bool
	IdVerified       bool

	// AI Coach
	CareerNarrative  string
	CoachingMetadata string

	// Work Auth
	WorkAuthStatus      string
	PassportNationality string
	DrivingLicenseBool  bool
	DrivingLicenseType  string

	// Communication & Accessibility
	PreferredContactChannel string
	AccessibilityNeeds      string
	VideoIntroURL           string

	// Mentorship
	WillingToMentor bool

	// Calculated Fields
	AvgResponseTimeHours     float64
	ProfileCompletenessScore int
	LastActiveAt             string

	// Background Check Consent
	BackgroundCheckConsent   bool
	BackgroundCheckConsentAt string

	// Job Alerts
	JobAlertFrequency string
	JobAlertChannel   string

	// Privacy settings
	VisibilityProfile          string
	VisibilitySalary           string
	VisibilityTransitionReason string
	VisibilityExperience       string
	VisibilityEducation        string
	VisibilityCertifications   string
	VisibilitySkills           string
	VisibilityPortfolio        string
	VisibilityReferences       string

	// Child collections
	Experiences    []WorkExperience
	Educations     []Education
	Certifications []Certification
	Skills         []ProfileSkill
	Languages      []Language
	Portfolio      []PortfolioLink
	Endorsements   []Endorsement
	References     []Reference
}

// Validate checks profile business rules.
func (p *Profile) Validate() error {
	// Salary range validation
	if p.SalaryMin > 0 && p.SalaryMax > 0 && p.SalaryMin > p.SalaryMax {
		return errors.New("minimum salary cannot be greater than maximum salary")
	}

	// Transition reason condition
	if p.TransitionReason != "" && p.CareerStatus != "career_break" && p.CareerStatus != "open_to_opportunities" && p.CareerStatus != "actively_looking" {
		return errors.New("transition reason requires career status to be career_break, open_to_opportunities, or actively_looking")
	}

	// Date range validations in experiences
	for _, exp := range p.Experiences {
		if exp.StartDate != "" && exp.EndDate != "" && exp.StartDate > exp.EndDate {
			return fmt.Errorf("experience start date %s cannot be after end date %s", exp.StartDate, exp.EndDate)
		}
	}

	// Date range validations in educations
	for _, edu := range p.Educations {
		if edu.StartDate != "" && edu.EndDate != "" && edu.StartDate > edu.EndDate {
			return fmt.Errorf("education start date %s cannot be after end date %s", edu.StartDate, edu.EndDate)
		}
	}

	return nil
}

// Scalars carries the editable top-level fields.
type Scalars struct {
	Headline string
	About    string
	PhotoURL string
	Bio      string
	Location string
	Website  string

	// Core Identity
	Pronouns     string
	CareerStatus string

	// Career Recovery
	TransitionReason       string
	TargetComebackTimeline string
	SupportsNeeded         []string

	// Mobility & Preferences
	OpenToRemote        bool
	OpenToRelocation    bool
	RelocationLocations []string
	DesiredRoles        []string
	DesiredIndustries   []string
	EmploymentType      string
	SalaryMin           int
	SalaryMax           int
	SalaryCurrency      string
	SalaryVisible       bool
	WorkMode            string
	AvailabilityDate    string
	NoticePeriod        string

	// Trust
	ReferralEligible bool

	// AI Coach
	CareerNarrative  string
	CoachingMetadata string

	// Work Auth
	WorkAuthStatus      string
	PassportNationality string
	DrivingLicenseBool  bool
	DrivingLicenseType  string

	// Communication & Accessibility
	PreferredContactChannel string
	AccessibilityNeeds      string
	VideoIntroURL           string

	// Mentorship
	WillingToMentor bool

	// Background Check Consent
	BackgroundCheckConsent   bool
	BackgroundCheckConsentAt string

	// Job Alerts
	JobAlertFrequency string
	JobAlertChannel   string

	// Privacy settings
	VisibilityProfile          string
	VisibilitySalary           string
	VisibilityTransitionReason string
	VisibilityExperience       string
	VisibilityEducation        string
	VisibilityCertifications   string
	VisibilitySkills           string
	VisibilityPortfolio        string
	VisibilityReferences       string
}

type WorkExperience struct {
	ID             string
	Title          string
	Company        string
	Location       string
	EmploymentType string
	StartDate      string
	EndDate        string
	IsCurrent      bool
	Description    string
	Achievements   []string
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

type ProfileSkill struct {
	Name             string
	ProficiencyLevel string
	EndorsedCount    int
}

type Language struct {
	Name        string
	Proficiency string
}

type PortfolioLink struct {
	ID       string
	Platform string
	URL      string
}

type Endorsement struct {
	ID           string
	ToUserID     string
	FromUserID   string
	Relationship string
	Text         string
	CreatedAt    string
}

type Reference struct {
	ID                  string
	Name                string
	Relationship        string
	ContactInfo         string
	PermissionToContact bool
}

type ConsentLog struct {
	ID           string
	UserID       string
	ConsentType  string // e.g. background_check, data_sharing
	TargetEntity string // e.g. recruiter/employer, empty for background check
	Consented    bool
	IPAddress    string
	UserAgent    string
	CreatedAt    string
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

	SetSkills(ctx context.Context, userID string, skills []ProfileSkill) error
	SetLanguages(ctx context.Context, userID string, langs []Language) error
	SetPortfolio(ctx context.Context, userID string, links []PortfolioLink) error

	AddEndorsement(ctx context.Context, toUserID string, e *Endorsement) error
	AddReference(ctx context.Context, userID string, r *Reference) error
	UpdateReference(ctx context.Context, userID string, r Reference) error
	DeleteReference(ctx context.Context, userID, id string) error
	AddConsentLog(ctx context.Context, cl *ConsentLog) error

	// Verification updates (internally triggered)
	SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) error

	// Calculated Fields updates
	UpdateCalculatedFields(ctx context.Context, userID string, completeness int, avgResponse float64, lastActive string) error
}
