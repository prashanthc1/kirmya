package api

import (
	"fmt"
	"strconv"
	"time"

	"workspace-app/internal/profile/domain"
)

type profileResponse struct {
	UserID   string `json:"user_id"`
	Headline string `json:"headline"`
	About    string `json:"about"`
	PhotoURL string `json:"photo_url"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`
	Version  int    `json:"version"`

	// Core Identity
	Pronouns     string `json:"pronouns"`
	CareerStatus string `json:"career_status"`

	// Career Recovery
	TransitionReason       string   `json:"transition_reason,omitempty"`
	TargetComebackTimeline string   `json:"target_comeback_timeline"`
	SupportsNeeded         []string `json:"supports_needed"`

	// Mobility & Preferences
	OpenToRemote        bool     `json:"open_to_remote"`
	OpenToRelocation    bool     `json:"open_to_relocation"`
	RelocationLocations []string `json:"relocation_locations"`
	DesiredRoles        []string `json:"desired_roles"`
	DesiredIndustries   []string `json:"desired_industries"`
	EmploymentType      string   `json:"employment_type"`
	SalaryMin           int      `json:"salary_min,omitempty"`
	SalaryMax           int      `json:"salary_max,omitempty"`
	SalaryCurrency      string   `json:"salary_currency,omitempty"`
	SalaryVisible       bool     `json:"salary_visible"`
	WorkMode            string   `json:"work_mode"`
	AvailabilityDate    string   `json:"availability_date"`
	NoticePeriod        string   `json:"notice_period"`

	// Trust & Verification
	ReferralEligible bool `json:"referral_eligible"`
	EmailVerified    bool `json:"email_verified"`
	PhoneVerified    bool `json:"phone_verified"`
	LinkedinVerified bool `json:"linkedin_verified"`
	IdVerified       bool `json:"id_verified"`

	// AI Coach
	CareerNarrative  string `json:"career_narrative"`
	CoachingMetadata string `json:"coaching_metadata"`

	// Work Auth
	WorkAuthStatus      string `json:"work_auth_status"`
	PassportNationality string `json:"passport_nationality"`
	DrivingLicenseBool  bool   `json:"driving_license_bool"`
	DrivingLicenseType  string `json:"driving_license_type"`

	// Communication & Accessibility
	PreferredContactChannel string `json:"preferred_contact_channel"`
	AccessibilityNeeds      string `json:"accessibility_needs,omitempty"`
	VideoIntroURL           string `json:"video_intro_url"`

	// Mentorship
	WillingToMentor bool `json:"willing_to_mentor"`

	// Calculated Fields
	AvgResponseTimeHours     float64 `json:"avg_response_time_hours"`
	ProfileCompletenessScore int     `json:"profile_completeness_score"`
	LastActiveAt             string  `json:"last_active_at"`

	// Background Check Consent
	BackgroundCheckConsent   bool   `json:"background_check_consent"`
	BackgroundCheckConsentAt string `json:"background_check_consent_at"`

	// Job Alerts
	JobAlertFrequency string `json:"job_alert_frequency"`
	JobAlertChannel   string `json:"job_alert_channel"`

	// Privacy settings
	VisibilityProfile          string `json:"visibility_profile"`
	VisibilitySalary           string `json:"visibility_salary"`
	VisibilityTransitionReason string `json:"visibility_transition_reason"`
	VisibilityExperience       string `json:"visibility_experience"`
	VisibilityEducation        string `json:"visibility_education"`
	VisibilityCertifications   string `json:"visibility_certifications"`
	VisibilitySkills           string `json:"visibility_skills"`
	VisibilityPortfolio        string `json:"visibility_portfolio"`
	VisibilityReferences       string `json:"visibility_references"`

	// Collections
	Experiences    []experienceDTO    `json:"experiences"`
	Educations     []educationDTO     `json:"educations"`
	Certifications []certificationDTO `json:"certifications"`
	Skills         []skillDTO         `json:"skills"`
	Languages      []languageDTO      `json:"languages"`
	Portfolio      []portfolioLinkDTO `json:"portfolio"`
	Endorsements   []endorsementDTO   `json:"endorsements,omitempty"`
	References     []referenceDTO     `json:"references,omitempty"`
}

type experienceDTO struct {
	ID             string   `json:"id,omitempty"`
	Title          string   `json:"title"`
	Company        string   `json:"company"`
	Location       string   `json:"location"`
	EmploymentType string   `json:"employment_type"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	IsCurrent      bool     `json:"is_current"`
	Description    string   `json:"description"`
	Achievements   []string `json:"achievements"`
}

type educationDTO struct {
	ID           string `json:"id,omitempty"`
	School       string `json:"school"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Grade        string `json:"grade"`
	Description  string `json:"description"`
}

type certificationDTO struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Issuer        string `json:"issuer"`
	IssueDate     string `json:"issue_date"`
	ExpiryDate    string `json:"expiry_date"`
	CredentialID  string `json:"credential_id"`
	CredentialURL string `json:"credential_url"`
}

type skillDTO struct {
	Name             string `json:"name"`
	ProficiencyLevel string `json:"proficiency_level"`
	EndorsedCount    int    `json:"endorsed_count"`
}

type languageDTO struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

type portfolioLinkDTO struct {
	ID       string `json:"id,omitempty"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type endorsementDTO struct {
	ID           string `json:"id,omitempty"`
	ToUserID     string `json:"to_user_id"`
	FromUserID   string `json:"from_user_id"`
	Relationship string `json:"relationship"`
	Text         string `json:"text"`
	CreatedAt    string `json:"created_at"`
}

type referenceDTO struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	Relationship        string `json:"relationship"`
	ContactInfo         string `json:"contact_info"`
	PermissionToContact bool   `json:"permission_to_contact"`
}

// Request DTOs
type updateProfileRequest struct {
	Headline string `json:"headline"`
	About    string `json:"about"`
	PhotoURL string `json:"photo_url"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`

	Pronouns     string `json:"pronouns"`
	CareerStatus string `json:"career_status"`

	TransitionReason       string   `json:"transition_reason"`
	TargetComebackTimeline string   `json:"target_comeback_timeline"`
	SupportsNeeded         []string `json:"supports_needed"`

	OpenToRemote        bool     `json:"open_to_remote"`
	OpenToRelocation    bool     `json:"open_to_relocation"`
	RelocationLocations []string `json:"relocation_locations"`
	DesiredRoles        []string `json:"desired_roles"`
	DesiredIndustries   []string `json:"desired_industries"`
	EmploymentType      string   `json:"employment_type"`
	SalaryMin           int      `json:"salary_min"`
	SalaryMax           int      `json:"salary_max"`
	SalaryCurrency      string   `json:"salary_currency"`
	SalaryVisible       bool     `json:"salary_visible"`
	WorkMode            string   `json:"work_mode"`
	AvailabilityDate    string   `json:"availability_date"`
	NoticePeriod        string   `json:"notice_period"`

	ReferralEligible bool `json:"referral_eligible"`

	CareerNarrative  string `json:"career_narrative"`
	CoachingMetadata string `json:"coaching_metadata"`

	WorkAuthStatus      string `json:"work_auth_status"`
	PassportNationality string `json:"passport_nationality"`
	DrivingLicenseBool  bool   `json:"driving_license_bool"`
	DrivingLicenseType  string `json:"driving_license_type"`

	PreferredContactChannel string `json:"preferred_contact_channel"`
	AccessibilityNeeds      string `json:"accessibility_needs"`
	VideoIntroURL           string `json:"video_intro_url"`

	WillingToMentor bool `json:"willing_to_mentor"`

	BackgroundCheckConsent   bool   `json:"background_check_consent"`
	BackgroundCheckConsentAt string `json:"background_check_consent_at"`

	JobAlertFrequency string `json:"job_alert_frequency"`
	JobAlertChannel   string `json:"job_alert_channel"`

	VisibilityProfile          string `json:"visibility_profile"`
	VisibilitySalary           string `json:"visibility_salary"`
	VisibilityTransitionReason string `json:"visibility_transition_reason"`
	VisibilityExperience       string `json:"visibility_experience"`
	VisibilityEducation        string `json:"visibility_education"`
	VisibilityCertifications   string `json:"visibility_certifications"`
	VisibilitySkills           string `json:"visibility_skills"`
	VisibilityPortfolio        string `json:"visibility_portfolio"`
	VisibilityReferences       string `json:"visibility_references"`
}

type skillsRequest struct {
	Skills []skillDTO `json:"skills"`
}

type languagesRequest struct {
	Languages []languageDTO `json:"languages"`
}

type portfolioRequest struct {
	Portfolio []portfolioLinkDTO `json:"portfolio"`
}

// --- mapping ---

func toResponse(p *domain.Profile) profileResponse {
	r := profileResponse{
		UserID: p.UserID, Headline: p.Identity.Headline, About: p.Identity.Bio, PhotoURL: p.Identity.PhotoURL,
		Bio: p.Identity.Bio, Location: p.Identity.Location, Website: p.Identity.SocialLinks.Website, Version: p.Version,
		Pronouns: p.Identity.VisaStatus, CareerStatus: p.Identity.Availability,
		TransitionReason: p.Identity.Bio, TargetComebackTimeline: p.Identity.Availability,
		SupportsNeeded: []string{}, OpenToRemote: p.Preferences.OpenToRelocation, OpenToRelocation: p.Preferences.OpenToRelocation,
		RelocationLocations: []string{}, DesiredRoles: p.Preferences.DesiredRoles, DesiredIndustries: p.Preferences.DesiredIndustries,
		EmploymentType: p.Preferences.NoticePeriod, SalaryMin: p.Preferences.SalaryMin, SalaryMax: p.Preferences.SalaryMax,
		SalaryCurrency: p.Preferences.SalaryCurrency, SalaryVisible: p.Preferences.OpenToRelocation, WorkMode: p.Preferences.RemotePreference,
		AvailabilityDate: "", NoticePeriod: p.Preferences.NoticePeriod,
		ReferralEligible: p.Verification.IdentityVerified, EmailVerified: p.Verification.EmailVerified, PhoneVerified: p.Verification.PhoneVerified,
		LinkedinVerified: p.Verification.IdentityVerified, IdVerified: p.Verification.IdentityVerified,
		CareerNarrative: string(p.AICareerAssistant.GapAnalysis), CoachingMetadata: string(p.AICareerAssistant.InterviewPrep),
		WorkAuthStatus: p.Identity.WorkAuthorization, PassportNationality: p.Identity.Nationality,
		DrivingLicenseBool: p.Verification.IdentityVerified, DrivingLicenseType: p.Identity.VisaStatus,
		PreferredContactChannel: p.Identity.PreferredContactChannel, AccessibilityNeeds: p.Identity.VisaStatus,
		VideoIntroURL: p.Identity.CoverURL, WillingToMentor: p.Verification.IdentityVerified,
		AvgResponseTimeHours: float64(p.Analytics.ProfileViews), ProfileCompletenessScore: p.ProfileCompletenessScore,
		LastActiveAt: p.LastActiveAt.Format(time.RFC3339), BackgroundCheckConsent: p.Verification.IdentityVerified,
		BackgroundCheckConsentAt: "", JobAlertFrequency: p.Identity.VisaStatus,
		JobAlertChannel: p.Identity.VisaStatus, VisibilityProfile: p.Privacy.FieldVisibility["profile"],
		VisibilitySalary: p.Privacy.FieldVisibility["salary"], VisibilityTransitionReason: p.Privacy.FieldVisibility["transition_reason"],
		VisibilityExperience: p.Privacy.FieldVisibility["experience"], VisibilityEducation: p.Privacy.FieldVisibility["education"],
		VisibilityCertifications: p.Privacy.FieldVisibility["certifications"], VisibilitySkills: p.Privacy.FieldVisibility["skills"],
		VisibilityPortfolio: p.Privacy.FieldVisibility["portfolio"], VisibilityReferences: p.Privacy.FieldVisibility["references"],
		Experiences:    make([]experienceDTO, 0, len(p.Experiences)),
		Educations:     make([]educationDTO, 0, len(p.Educations)),
		Certifications: make([]certificationDTO, 0, len(p.Certifications)),
		Skills:         make([]skillDTO, 0, len(p.Skills)),
		Languages:      make([]languageDTO, 0, len(p.Identity.Languages)),
		Portfolio:      make([]portfolioLinkDTO, 0, len(p.Projects)),
		Endorsements:   make([]endorsementDTO, 0, len(p.Networking.Recommendations)),
		References:     make([]referenceDTO, 0),
	}

	for _, e := range p.Experiences {
		achs := e.Achievements
		if achs == nil {
			achs = []string{}
		}
		r.Experiences = append(r.Experiences, experienceDTO{
			ID:             e.ID,
			Title:          e.Position,
			Company:        e.Company,
			Location:       e.Location,
			EmploymentType: e.EmploymentType,
			StartDate:      e.StartDate.Format("2006-01-02"),
			EndDate:        e.EndDate.Format("2006-01-02"),
			IsCurrent:      e.IsCurrent,
			Description:    e.Responsibilities,
			Achievements:   achs,
		})
	}
	for _, e := range p.Educations {
		r.Educations = append(r.Educations, educationDTO{
			ID:           e.ID,
			School:       e.Institution,
			Degree:       e.Degree,
			FieldOfStudy: e.Major,
			StartDate:    e.GraduationDate.Format("2006-01-02"),
			EndDate:      e.GraduationDate.Format("2006-01-02"),
			Grade:        fmt.Sprintf("%.2f", e.GPA),
			Description:  e.Thesis,
		})
	}
	for _, c := range p.Certifications {
		r.Certifications = append(r.Certifications, certificationDTO{
			ID:            c.ID,
			Name:          c.Name,
			Issuer:        c.Issuer,
			IssueDate:     c.IssueDate.Format("2006-01-02"),
			ExpiryDate:    c.ExpirationDate.Format("2006-01-02"),
			CredentialID:  c.CredentialID,
			CredentialURL: c.VerificationURL,
		})
	}
	for _, s := range p.Skills {
		r.Skills = append(r.Skills, skillDTO{
			Name:             s.Name,
			ProficiencyLevel: s.Level,
			EndorsedCount:    0,
		})
	}
	for _, l := range p.Identity.Languages {
		r.Languages = append(r.Languages, languageDTO{
			Name:        l.Name,
			Proficiency: l.Proficiency,
		})
	}
	for _, pr := range p.Projects {
		r.Portfolio = append(r.Portfolio, portfolioLinkDTO{
			ID:       pr.ID,
			Platform: "custom",
			URL:      pr.LiveDemoURL,
		})
	}
	for _, rec := range p.Networking.Recommendations {
		r.Endorsements = append(r.Endorsements, endorsementDTO{
			FromUserID:   rec.FromUserName,
			Relationship: rec.Relationship,
			Text:         rec.Text,
			CreatedAt:    rec.CreatedAt.Format(time.RFC3339),
		})
	}

	return r
}

type consentDTO struct {
	ID           string `json:"id,omitempty"`
	ConsentType  string `json:"consent_type"`
	TargetEntity string `json:"target_entity"`
	Consented    bool   `json:"consented"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	CreatedAt    string `json:"created_at"`
}

func (d consentDTO) toDomain() domain.ConsentLog {
	return domain.ConsentLog{
		ID:           d.ID,
		ConsentType:  d.ConsentType,
		TargetEntity: d.TargetEntity,
		Consented:    d.Consented,
		IPAddress:    d.IPAddress,
		UserAgent:    d.UserAgent,
		CreatedAt:    d.CreatedAt,
	}
}

func (d experienceDTO) toDomain() domain.WorkExperience {
	var start, end time.Time
	if d.StartDate != "" {
		start, _ = time.Parse("2006-01-02", d.StartDate)
	}
	if d.EndDate != "" {
		end, _ = time.Parse("2006-01-02", d.EndDate)
	}
	return domain.WorkExperience{
		ID:               d.ID,
		Position:         d.Title,
		Company:          d.Company,
		Location:         d.Location,
		EmploymentType:   d.EmploymentType,
		StartDate:        start,
		EndDate:          end,
		IsCurrent:        d.IsCurrent,
		Responsibilities: d.Description,
		Achievements:     d.Achievements,
	}
}

func (d educationDTO) toDomain() domain.Education {
	var grad time.Time
	if d.EndDate != "" {
		grad, _ = time.Parse("2006-01-02", d.EndDate)
	}
	gpaVal, _ := strconv.ParseFloat(d.Grade, 64)
	return domain.Education{
		ID:                 d.ID,
		Institution:        d.School,
		Degree:             d.Degree,
		Major:              d.FieldOfStudy,
		GraduationDate:     grad,
		GPA:                gpaVal,
		Thesis:             d.Description,
		VerificationStatus: "unverified",
	}
}

func (d certificationDTO) toDomain() domain.CertificationItem {
	var issue, expiry time.Time
	if d.IssueDate != "" {
		issue, _ = time.Parse("2006-01-02", d.IssueDate)
	}
	if d.ExpiryDate != "" {
		expiry, _ = time.Parse("2006-01-02", d.ExpiryDate)
	}
	return domain.CertificationItem{
		ID:              d.ID,
		Name:            d.Name,
		Issuer:          d.Issuer,
		CredentialID:    d.CredentialID,
		VerificationURL: d.CredentialURL,
		IssueDate:       issue,
		ExpirationDate:  expiry,
		Status:          "active",
	}
}

func (d skillDTO) toDomain() domain.SkillItem {
	return domain.SkillItem{
		Name:  d.Name,
		Level: d.ProficiencyLevel,
	}
}

func (d languageDTO) toDomain() domain.LanguageItem {
	return domain.LanguageItem{
		Name:        d.Name,
		Proficiency: d.Proficiency,
	}
}

func (d portfolioLinkDTO) toDomain() domain.ProjectItem {
	return domain.ProjectItem{
		ID:          d.ID,
		LiveDemoURL: d.URL,
	}
}

func (d endorsementDTO) toDomain() domain.Endorsement {
	return domain.Endorsement{
		ID:           d.ID,
		ToUserID:     d.ToUserID,
		FromUserID:   d.FromUserID,
		Relationship: d.Relationship,
		Text:         d.Text,
		CreatedAt:    d.CreatedAt,
	}
}

func (d referenceDTO) toDomain() domain.Reference {
	return domain.Reference{
		ID:                  d.ID,
		Name:                d.Name,
		Relationship:        d.Relationship,
		ContactInfo:         d.ContactInfo,
		PermissionToContact: d.PermissionToContact,
	}
}

func (r updateProfileRequest) toDomain() domain.AggregateUpdate {
	isDraft := true
	return domain.AggregateUpdate{
		Identity: &domain.IdentitySection{
			Headline:                r.Headline,
			Bio:                     r.Bio,
			PhotoURL:                r.PhotoURL,
			Location:                r.Location,
			PreferredContactChannel: r.PreferredContactChannel,
			WorkAuthorization:       r.WorkAuthStatus,
			Nationality:             r.PassportNationality,
		},
		Summary: &domain.SummarySection{
			ExecutiveSummary: r.About,
		},
		Preferences: &domain.CareerPreferences{
			DesiredRoles:      r.DesiredRoles,
			DesiredIndustries: r.DesiredIndustries,
			OpenToRelocation:  r.OpenToRelocation,
			NoticePeriod:      r.NoticePeriod,
			RemotePreference:  r.WorkMode,
			SalaryMin:         r.SalaryMin,
			SalaryMax:         r.SalaryMax,
			SalaryCurrency:    r.SalaryCurrency,
		},
		Privacy: &domain.PrivacySecuritySettings{
			FieldVisibility: map[string]string{
				"profile":           r.VisibilityProfile,
				"salary":            r.VisibilitySalary,
				"transition_reason": r.VisibilityTransitionReason,
				"experience":        r.VisibilityExperience,
				"education":         r.VisibilityEducation,
				"certifications":    r.VisibilityCertifications,
				"skills":            r.VisibilitySkills,
				"portfolio":         r.VisibilityPortfolio,
				"references":        r.VisibilityReferences,
			},
		},
		IsDraft: &isDraft,
	}
}
