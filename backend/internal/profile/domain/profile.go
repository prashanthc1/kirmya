// Package domain holds the Profile bounded context's entities and ports.
package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("profile not found")
	// ErrOptimisticLock is returned when a version-checked update is applied on
	// top of a stale read (concurrent modification). Mapped to HTTP 409.
	ErrOptimisticLock = errors.New("profile was modified by another request")
)

// Profile represents the full Career Operating System profile aggregate.
type Profile struct {
	UserID                   string                  `json:"user_id"`
	Version                  int                     `json:"version"`
	IsDraft                  bool                    `json:"is_draft"`
	ProfileCompletenessScore int                     `json:"profile_completeness_score"`
	TrustScore               int                     `json:"trust_score"`
	LastActiveAt             time.Time               `json:"last_active_at"`

	// 15 Sections
	Identity          IdentitySection         `json:"identity"`
	Summary           SummarySection          `json:"summary"`
	Experiences       []WorkExperience        `json:"experiences"`
	Educations        []Education             `json:"educations"`
	Skills            []SkillItem             `json:"skills"`
	Projects          []ProjectItem           `json:"projects"`
	Certifications    []CertificationItem     `json:"certifications"`
	Achievements      []AchievementItem       `json:"achievements"`
	Resumes           []ResumeVersion         `json:"resumes"`
	Preferences       CareerPreferences       `json:"preferences"`
	Verification      VerificationStatus      `json:"verification"`
	Networking        NetworkingSummary       `json:"networking"`
	Analytics         AnalyticsSummary        `json:"analytics"`
	Privacy           PrivacySecuritySettings `json:"privacy"`
	AICareerAssistant AICareerState           `json:"ai_career_assistant"`
}

// Section 1 - Identity & Personal Information
type IdentitySection struct {
	PhotoURL                string         `json:"photo_url"`
	CoverURL                string         `json:"cover_url"`
	FullName                string         `json:"full_name"`
	PreferredName           string         `json:"preferred_name"`
	Headline                string         `json:"headline"`
	CurrentTitle            string         `json:"current_title"`
	CurrentEmployer         string         `json:"current_employer"`
	Bio                     string         `json:"bio"`
	Location                string         `json:"location"`
	Country                 string         `json:"country"`
	TimeZone                string         `json:"timezone"`
	Nationality             string         `json:"nationality"`
	Languages               []LanguageItem `json:"languages"`
	Phone                   string         `json:"phone"`
	Email                   string         `json:"email"`
	SocialLinks             SocialLinks    `json:"social_links"`
	Availability            string         `json:"availability"`
	WorkAuthorization       string         `json:"work_authorization"`
	VisaStatus              string         `json:"visa_status"`
	PreferredContactChannel string         `json:"preferred_contact_channel"`
}

type SocialLinks struct {
	Website       string `json:"website"`
	LinkedIn      string `json:"linkedin"`
	GitHub        string `json:"github"`
	Portfolio     string `json:"portfolio"`
	Behance       string `json:"behance"`
	Dribbble      string `json:"dribbble"`
	Medium        string `json:"medium"`
	StackOverflow string `json:"stack_overflow"`
	GoogleScholar string `json:"google_scholar"`
	ResearchGate  string `json:"research_gate"`
	ORCID         string `json:"orcid"`
}

type LanguageItem struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

// Section 2 - Professional Summary
type SummarySection struct {
	ExecutiveSummary       string   `json:"executive_summary"`
	CareerObjectives       string   `json:"career_objectives"`
	CareerHighlights       []string `json:"career_highlights"`
	Industries             []string `json:"industries"`
	FunctionalAreas        []string `json:"functional_areas"`
	PersonalBrandStatement string   `json:"personal_brand_statement"`
	ElevatorPitch          string   `json:"elevator_pitch"`
}

// Section 3 - Work Experience
type WorkExperience struct {
	ID              string    `json:"id"`
	Company         string    `json:"company"`
	CompanyLogo     string    `json:"company_logo"`
	Position        string    `json:"position"`
	EmploymentType  string    `json:"employment_type"`
	Location        string    `json:"location"`
	RemoteType      string    `json:"remote_type"` // remote, hybrid, onsite
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	IsCurrent       bool      `json:"is_current"`
	Responsibilities string   `json:"responsibilities"`
	Achievements    []string  `json:"achievements"`
	KPIs            []string  `json:"kpis"`
	Technologies    []string  `json:"technologies"`
	SkillsUsed      []string  `json:"skills_used"`
	TeamSize        int       `json:"team_size"`
	Attachments     []string  `json:"attachments"`
}

// Section 4 - Education
type Education struct {
	ID                 string    `json:"id"`
	Institution        string    `json:"institution"`
	Degree             string    `json:"degree"`
	FieldOfStudy       string    `json:"field_of_study"`
	Major              string    `json:"major"`
	Minor              string    `json:"minor"`
	GPA                float64   `json:"gpa"`
	Honors             string    `json:"honors"`
	Activities         string    `json:"activities"`
	Projects           string    `json:"projects"`
	Research           string    `json:"research"`
	Thesis             string    `json:"thesis"`
	GraduationDate     time.Time `json:"graduation_date"`
	VerificationStatus string    `json:"verification_status"`
}

// Section 5 - Skills
type SkillItem struct {
	Name                   string  `json:"name"`
	Category               string  `json:"category"`
	Level                  string  `json:"level"` // beginner, intermediate, advanced, expert
	YearsOfExperience      float64 `json:"years_of_experience"`
	LastUsed               int     `json:"last_used"`
	Verified               bool    `json:"verified"`
	RecruiterDemandScore   float64 `json:"recruiter_demand_score"`
	AIRecommendationScore  float64 `json:"ai_recommendation_score"`
}

// Section 6 - Projects
type ProjectItem struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	RepositoryURL  string   `json:"repository_url"`
	LiveDemoURL    string   `json:"live_demo_url"`
	VideoURL       string   `json:"video_url"`
	Images         []string `json:"images"`
	Videos         []string `json:"videos"`
	Documents      []string `json:"documents"`
	Technologies   []string `json:"technologies"`
	Timeline       string   `json:"timeline"`
	TeamMembers    []string `json:"team_members"`
	Metrics        string   `json:"metrics"`
	Awards         string   `json:"awards"`
	BusinessImpact string   `json:"business_impact"`
}

// Section 7 - Certifications
type CertificationItem struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Issuer         string    `json:"issuer"`
	CredentialID   string    `json:"credential_id"`
	VerificationURL string   `json:"verification_url"`
	SkillsCovered  []string  `json:"skills_covered"`
	IssueDate      time.Time `json:"issue_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Status         string    `json:"status"`
}

// Section 8 - Achievements
type AchievementItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	IssuerOrOrg string    `json:"issuer_or_org"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"` // award, patent, publication, conference, etc.
	Description string    `json:"description"`
	EvidenceURL string    `json:"evidence_url"`
}

// Section 9 - Resume
type ResumeVersion struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	FileURL         string    `json:"file_url"`
	FileSize        int       `json:"file_size"`
	ATSScore        int       `json:"ats_score"`
	KeywordAnalysis []byte    `json:"keyword_analysis"`
	IsPrimary       bool      `json:"is_primary"`
	UploadedAt      time.Time `json:"uploaded_at"`
}

// Section 10 - Career Preferences
type CareerPreferences struct {
	DesiredRoles           []string `json:"desired_roles"`
	DesiredIndustries      []string `json:"desired_industries"`
	EmploymentTypes        []string `json:"employment_types"`
	SalaryMin              int      `json:"salary_min"`
	SalaryMax              int      `json:"salary_max"`
	SalaryCurrency         string   `json:"salary_currency"`
	NoticePeriod           string   `json:"notice_period"`
	RemotePreference       string   `json:"remote_preference"`
	OpenToRelocation       bool     `json:"open_to_relocation"`
	PreferredCountries     []string `json:"preferred_countries"`
	PreferredCities        []string `json:"preferred_cities"`
	TravelWillingness      string   `json:"travel_willingness"`
	CompanySizePreferences []string `json:"company_size_preferences"`
}

// Section 11 - Verification & Trust
type VerificationStatus struct {
	EmailVerified         bool `json:"email_verified"`
	PhoneVerified         bool `json:"phone_verified"`
	IdentityVerified      bool `json:"identity_verified"`
	EmploymentVerified   bool `json:"employment_verified"`
	EducationVerified    bool `json:"education_verified"`
	CertificationVerified bool `json:"certification_verified"`
}

// Section 12 - Networking
type NetworkingSummary struct {
	ConnectionsCount int                  `json:"connections_count"`
	FollowersCount   int                  `json:"followers_count"`
	FollowingCount   int                  `json:"following_count"`
	MentorsCount     int                  `json:"mentors_count"`
	ReferralsCount   int                  `json:"referrals_count"`
	Recommendations  []EndorsementSummary `json:"recommendations"`
}

type EndorsementSummary struct {
	FromUserName string    `json:"from_user_name"`
	Relationship string    `json:"relationship"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"created_at"`
}

// Section 13 - Analytics
type AnalyticsSummary struct {
	ProfileViews       int   `json:"profile_views"`
	SearchAppearances  int   `json:"search_appearances"`
	RecruiterViews     int   `json:"recruiter_views"`
	ResumeDownloads    int   `json:"resume_downloads"`
	PortfolioViews     int   `json:"portfolio_views"`
	WeeklyProfileViews []int `json:"weekly_profile_views"`
}

// Section 14 - Privacy & Security
type PrivacySecuritySettings struct {
	FieldVisibility  map[string]string `json:"field_visibility"` // key: section, value: public, recruiter_only, connections_only, private
	TwoFactorEnabled bool              `json:"two_factor_enabled"`
	ActiveSessions   []ActiveSession   `json:"active_sessions"`
}

type ActiveSession struct {
	SessionID  string    `json:"session_id"`
	Device     string    `json:"device"`
	IPAddress  string    `json:"ip_address"`
	LastActive time.Time `json:"last_active"`
}

// Section 15 - AI Career Assistant
type AICareerState struct {
	GapAnalysis     []byte    `json:"gap_analysis"`
	Roadmap         []byte    `json:"roadmap"`
	InterviewPrep   []byte    `json:"interview_prep"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}

// Validate checks profile business rules.
func (p *Profile) Validate() error {
	// Salary range validation
	if p.Preferences.SalaryMin > 0 && p.Preferences.SalaryMax > 0 && p.Preferences.SalaryMin > p.Preferences.SalaryMax {
		return errors.New("minimum salary cannot be greater than maximum salary")
	}

	// Date range validations in experiences
	for _, exp := range p.Experiences {
		if !exp.StartDate.IsZero() && !exp.EndDate.IsZero() && exp.StartDate.After(exp.EndDate) {
			return fmt.Errorf("experience start date %s cannot be after end date %s", exp.StartDate.Format("2006-01-02"), exp.EndDate.Format("2006-01-02"))
		}
	}

	return nil
}

// CalculateCompleteness returns the profile completion percentage based on filled sections.
func (p *Profile) CalculateCompleteness() int {
	score := 0
	totalFields := 15

	if p.Identity.FullName != "" && p.Identity.Email != "" {
		score++
	}
	if p.Summary.ExecutiveSummary != "" {
		score++
	}
	if len(p.Experiences) > 0 {
		score++
	}
	if len(p.Educations) > 0 {
		score++
	}
	if len(p.Skills) > 0 {
		score++
	}
	if len(p.Projects) > 0 {
		score++
	}
	if len(p.Certifications) > 0 {
		score++
	}
	if len(p.Achievements) > 0 {
		score++
	}
	if len(p.Resumes) > 0 {
		score++
	}
	if len(p.Preferences.DesiredRoles) > 0 {
		score++
	}
	if p.Verification.EmailVerified {
		score++
	}
	if p.Networking.ConnectionsCount > 0 {
		score++
	}
	if p.Analytics.ProfileViews > 0 {
		score++
	}
	if len(p.Privacy.FieldVisibility) > 0 {
		score++
	}
	if len(p.AICareerAssistant.Roadmap) > 0 {
		score++
	}

	return int(float64(score) / float64(totalFields) * 100)
}

// AggregateUpdate carries a partial update to the aggregate.
type AggregateUpdate struct {
	Identity       *IdentitySection         `json:"identity,omitempty"`
	Summary        *SummarySection          `json:"summary,omitempty"`
	Experiences    *[]WorkExperience        `json:"experiences,omitempty"`
	Educations     *[]Education             `json:"educations,omitempty"`
	Skills         *[]SkillItem             `json:"skills,omitempty"`
	Projects       *[]ProjectItem           `json:"projects,omitempty"`
	Certifications *[]CertificationItem     `json:"certifications,omitempty"`
	Achievements   *[]AchievementItem       `json:"achievements,omitempty"`
	Preferences    *CareerPreferences       `json:"preferences,omitempty"`
	Privacy        *PrivacySecuritySettings `json:"privacy,omitempty"`
	IsDraft        *bool                    `json:"is_draft,omitempty"`
}

// Validate runs aggregate checks over the update fields.
func (u AggregateUpdate) Validate() error {
	p := &Profile{}
	if u.Preferences != nil {
		p.Preferences = *u.Preferences
	}
	if u.Experiences != nil {
		p.Experiences = *u.Experiences
	}
	return p.Validate()
}

// Repository is the persistence port for the profile aggregate.
type Repository interface {
	Get(ctx context.Context, userID string, includeDraft bool) (*Profile, error)
	UpdateAggregate(ctx context.Context, userID string, expectedVersion int, u AggregateUpdate) error

	// Version Snapshots
	CreateVersionSnapshot(ctx context.Context, userID string, version int, p *Profile) error
	GetVersionSnapshot(ctx context.Context, userID string, version int) (*Profile, error)
	ListVersions(ctx context.Context, userID string) ([]int, error)

	// Auditing
	WriteAuditLog(ctx context.Context, log *AuditLogEntry) error

	// Analytics
	RecordAnalyticsEvent(ctx context.Context, profileID string, eventType string, actorID *string, ip, ua string) error
	GetAnalytics(ctx context.Context, profileID string) (*AnalyticsSummary, error)

	// Verification
	SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) error
}

type AuditLogEntry struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Section   string          `json:"section"`
	Action    string          `json:"action"`
	ActorID   string          `json:"actor_id"`
	OldValue  json.RawMessage `json:"old_value"`
	NewValue  json.RawMessage `json:"new_value"`
	IPAddress string          `json:"ip_address"`
	UserAgent string          `json:"user_agent"`
	CreatedAt time.Time       `json:"created_at"`
}

type Reference struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Relationship        string `json:"relationship"`
	ContactInfo         string `json:"contact_info"`
	PermissionToContact bool   `json:"permission_to_contact"`
}

type ConsentLog struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	ConsentType  string `json:"consent_type"`
	TargetEntity string `json:"target_entity"`
	Consented    bool   `json:"consented"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	CreatedAt    string `json:"created_at"`
}

type Endorsement struct {
	ID           string `json:"id"`
	ToUserID     string `json:"to_user_id"`
	FromUserID   string `json:"from_user_id"`
	Relationship string `json:"relationship"`
	Text         string `json:"text"`
	CreatedAt    string `json:"created_at"`
}
