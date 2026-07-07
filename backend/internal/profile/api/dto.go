package api

import "workspace-app/internal/profile/domain"

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
	CreatedAt    string `json:"created_at,omitempty"`
}

type referenceDTO struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	Relationship        string `json:"relationship"`
	ContactInfo         string `json:"contact_info"`
	PermissionToContact bool   `json:"permission_to_contact"`
}

type consentDTO struct {
	ID           string `json:"id,omitempty"`
	ConsentType  string `json:"consent_type"`
	TargetEntity string `json:"target_entity"`
	Consented    bool   `json:"consented"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type updateScalarsRequest struct {
	// Version is the client's last-known aggregate version, used for optimistic
	// concurrency. Zero (omitted) skips the check for backward compatibility.
	Version int `json:"version"`

	Headline string `json:"headline"`
	About    string `json:"about"`
	PhotoURL string `json:"photo_url"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`

	// Core Identity
	Pronouns     string `json:"pronouns"`
	CareerStatus string `json:"career_status"`

	// Career Recovery
	TransitionReason       string   `json:"transition_reason"`
	TargetComebackTimeline string   `json:"target_comeback_timeline"`
	SupportsNeeded         []string `json:"supports_needed"`

	// Mobility & Preferences
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

	// Trust
	ReferralEligible bool `json:"referral_eligible"`

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
	AccessibilityNeeds      string `json:"accessibility_needs"`
	VideoIntroURL           string `json:"video_intro_url"`

	// Mentorship
	WillingToMentor bool `json:"willing_to_mentor"`

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

	// Collections for synchronization
	Experiences    *[]experienceDTO    `json:"experiences"`
	Educations     *[]educationDTO     `json:"educations"`
	Certifications *[]certificationDTO `json:"certifications"`
	Skills         *[]skillDTO         `json:"skills"`
	Languages      *[]languageDTO      `json:"languages"`
	Portfolio      *[]portfolioLinkDTO `json:"portfolio"`
	References     *[]referenceDTO     `json:"references"`
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
		UserID: p.UserID, Headline: p.Headline, About: p.About, PhotoURL: p.PhotoURL,
		Bio: p.Bio, Location: p.Location, Website: p.Website, Version: p.Version,
		Pronouns: p.Pronouns, CareerStatus: p.CareerStatus,
		TransitionReason: p.TransitionReason, TargetComebackTimeline: p.TargetComebackTimeline,
		SupportsNeeded: p.SupportsNeeded, OpenToRemote: p.OpenToRemote, OpenToRelocation: p.OpenToRelocation,
		RelocationLocations: p.RelocationLocations, DesiredRoles: p.DesiredRoles, DesiredIndustries: p.DesiredIndustries,
		EmploymentType: p.EmploymentType, SalaryMin: p.SalaryMin, SalaryMax: p.SalaryMax,
		SalaryCurrency: p.SalaryCurrency, SalaryVisible: p.SalaryVisible, WorkMode: p.WorkMode,
		AvailabilityDate: p.AvailabilityDate, NoticePeriod: p.NoticePeriod,
		ReferralEligible: p.ReferralEligible, EmailVerified: p.EmailVerified, PhoneVerified: p.PhoneVerified,
		LinkedinVerified: p.LinkedinVerified, IdVerified: p.IdVerified,
		CareerNarrative: p.CareerNarrative, CoachingMetadata: p.CoachingMetadata,
		WorkAuthStatus: p.WorkAuthStatus, PassportNationality: p.PassportNationality,
		DrivingLicenseBool: p.DrivingLicenseBool, DrivingLicenseType: p.DrivingLicenseType,
		PreferredContactChannel: p.PreferredContactChannel, AccessibilityNeeds: p.AccessibilityNeeds,
		VideoIntroURL: p.VideoIntroURL, WillingToMentor: p.WillingToMentor,
		AvgResponseTimeHours: p.AvgResponseTimeHours, ProfileCompletenessScore: p.ProfileCompletenessScore,
		LastActiveAt: p.LastActiveAt, BackgroundCheckConsent: p.BackgroundCheckConsent,
		BackgroundCheckConsentAt: p.BackgroundCheckConsentAt, JobAlertFrequency: p.JobAlertFrequency,
		JobAlertChannel: p.JobAlertChannel, VisibilityProfile: p.VisibilityProfile,
		VisibilitySalary: p.VisibilitySalary, VisibilityTransitionReason: p.VisibilityTransitionReason,
		VisibilityExperience: p.VisibilityExperience, VisibilityEducation: p.VisibilityEducation,
		VisibilityCertifications: p.VisibilityCertifications, VisibilitySkills: p.VisibilitySkills,
		VisibilityPortfolio: p.VisibilityPortfolio, VisibilityReferences: p.VisibilityReferences,
		Experiences:    make([]experienceDTO, 0, len(p.Experiences)),
		Educations:     make([]educationDTO, 0, len(p.Educations)),
		Certifications: make([]certificationDTO, 0, len(p.Certifications)),
		Skills:         make([]skillDTO, 0, len(p.Skills)),
		Languages:      make([]languageDTO, 0, len(p.Languages)),
		Portfolio:      make([]portfolioLinkDTO, 0, len(p.Portfolio)),
		Endorsements:   make([]endorsementDTO, 0, len(p.Endorsements)),
		References:     make([]referenceDTO, 0, len(p.References)),
	}

	if r.SupportsNeeded == nil {
		r.SupportsNeeded = []string{}
	}
	if r.RelocationLocations == nil {
		r.RelocationLocations = []string{}
	}
	if r.DesiredRoles == nil {
		r.DesiredRoles = []string{}
	}
	if r.DesiredIndustries == nil {
		r.DesiredIndustries = []string{}
	}

	for _, e := range p.Experiences {
		achs := e.Achievements
		if achs == nil {
			achs = []string{}
		}
		r.Experiences = append(r.Experiences, experienceDTO{
			ID:             e.ID,
			Title:          e.Title,
			Company:        e.Company,
			Location:       e.Location,
			EmploymentType: e.EmploymentType,
			StartDate:      e.StartDate,
			EndDate:        e.EndDate,
			IsCurrent:      e.IsCurrent,
			Description:    e.Description,
			Achievements:   achs,
		})
	}
	for _, e := range p.Educations {
		r.Educations = append(r.Educations, educationDTO(e))
	}
	for _, c := range p.Certifications {
		r.Certifications = append(r.Certifications, certificationDTO(c))
	}
	for _, s := range p.Skills {
		r.Skills = append(r.Skills, skillDTO(s))
	}
	for _, l := range p.Languages {
		r.Languages = append(r.Languages, languageDTO(l))
	}
	for _, l := range p.Portfolio {
		r.Portfolio = append(r.Portfolio, portfolioLinkDTO(l))
	}
	for _, en := range p.Endorsements {
		r.Endorsements = append(r.Endorsements, endorsementDTO(en))
	}
	for _, ref := range p.References {
		r.References = append(r.References, referenceDTO(ref))
	}
	return r
}

// Convert DTOs to domain
func (d experienceDTO) toDomain() domain.WorkExperience {
	achs := d.Achievements
	if achs == nil {
		achs = []string{}
	}
	return domain.WorkExperience{
		ID:             d.ID,
		Title:          d.Title,
		Company:        d.Company,
		Location:       d.Location,
		EmploymentType: d.EmploymentType,
		StartDate:      d.StartDate,
		EndDate:        d.EndDate,
		IsCurrent:      d.IsCurrent,
		Description:    d.Description,
		Achievements:   achs,
	}
}

func (d educationDTO) toDomain() domain.Education         { return domain.Education(d) }
func (d certificationDTO) toDomain() domain.Certification { return domain.Certification(d) }
func (d skillDTO) toDomain() domain.ProfileSkill          { return domain.ProfileSkill(d) }
func (d endorsementDTO) toDomain() domain.Endorsement     { return domain.Endorsement(d) }
func (d referenceDTO) toDomain() domain.Reference         { return domain.Reference(d) }
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

func toScalars(r updateScalarsRequest) domain.Scalars {
	return domain.Scalars{
		Headline: r.Headline, About: r.About, PhotoURL: r.PhotoURL, Bio: r.Bio, Location: r.Location, Website: r.Website,
		Pronouns: r.Pronouns, CareerStatus: r.CareerStatus, TransitionReason: r.TransitionReason, TargetComebackTimeline: r.TargetComebackTimeline,
		OpenToRemote: r.OpenToRemote, OpenToRelocation: r.OpenToRelocation, EmploymentType: r.EmploymentType,
		SalaryMin: r.SalaryMin, SalaryMax: r.SalaryMax, SalaryCurrency: r.SalaryCurrency, SalaryVisible: r.SalaryVisible,
		WorkMode: r.WorkMode, AvailabilityDate: r.AvailabilityDate, NoticePeriod: r.NoticePeriod,
		ReferralEligible: r.ReferralEligible, CareerNarrative: r.CareerNarrative, CoachingMetadata: r.CoachingMetadata,
		WorkAuthStatus: r.WorkAuthStatus, PassportNationality: r.PassportNationality, DrivingLicenseBool: r.DrivingLicenseBool, DrivingLicenseType: r.DrivingLicenseType,
		PreferredContactChannel: r.PreferredContactChannel, AccessibilityNeeds: r.AccessibilityNeeds, VideoIntroURL: r.VideoIntroURL,
		WillingToMentor: r.WillingToMentor, BackgroundCheckConsent: r.BackgroundCheckConsent, BackgroundCheckConsentAt: r.BackgroundCheckConsentAt,
		JobAlertFrequency: r.JobAlertFrequency, JobAlertChannel: r.JobAlertChannel,
		VisibilityProfile: r.VisibilityProfile, VisibilitySalary: r.VisibilitySalary, VisibilityTransitionReason: r.VisibilityTransitionReason,
		VisibilityExperience: r.VisibilityExperience, VisibilityEducation: r.VisibilityEducation, VisibilityCertifications: r.VisibilityCertifications,
		VisibilitySkills: r.VisibilitySkills, VisibilityPortfolio: r.VisibilityPortfolio, VisibilityReferences: r.VisibilityReferences,
		SupportsNeeded: r.SupportsNeeded, RelocationLocations: r.RelocationLocations, DesiredRoles: r.DesiredRoles, DesiredIndustries: r.DesiredIndustries,
	}
}
