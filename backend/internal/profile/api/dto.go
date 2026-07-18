package api

import (
	"fmt"
	"strconv"
	"time"

	"workspace-app/internal/profile/domain"
)

type profileResponse struct {
	UserID      string `json:"user_id"`
	FullName    string `json:"full_name"`
	Headline    string `json:"headline"`
	About       string `json:"about"`
	PhotoURL    string `json:"photo_url"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	CoverBanner string `json:"cover_banner"`
	Version     int    `json:"version"`

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

	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`

	// Mentorship
	WillingToMentor bool `json:"willing_to_mentor"`

	// Calculated Fields
	AvgResponseTimeHours     float64 `json:"avg_response_time_hours"`
	ProfileCompletenessScore int     `json:"profile_completeness_score"`
	LastActiveAt             string  `json:"last_active_at"`

	// Background Check Consent
	BackgroundCheckConsent   bool   `json:"background_check_consent"`
	BackgroundCheckConsentAt string `json:"background_check_consent_at"`

	PreferredName string `json:"preferred_name"`
	LinkedInURL   string `json:"linkedin_url"`
	GitHubURL     string `json:"github_url"`
	PersonalBrand string `json:"personal_brand"`
	ElevatorPitch string `json:"elevator_pitch"`
	Industry      string `json:"industry"`
	AnonymousMode bool   `json:"anonymous_mode"`

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
	Projects       []projectDTO       `json:"projects"`
	Achievements   []achievementDTO   `json:"achievements_list"`
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

type projectDTO struct {
	ID             string   `json:"id,omitempty"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	RepositoryURL  string   `json:"repository_url"`
	LiveDemoURL    string   `json:"live_demo_url"`
	VideoURL       string   `json:"video_url"`
	Screenshots    []string `json:"screenshots"`
	Images         []string `json:"images"`
	Technologies   []string `json:"technologies"`
	Timeline       string   `json:"timeline"`
	Metrics        string   `json:"metrics"`
	Awards         string   `json:"awards"`
	BusinessImpact string   `json:"business_impact"`
}

type achievementDTO struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	IssuerOrOrg string `json:"issuer_or_org"`
	Date        string `json:"date"`
	Category    string `json:"category"`
	Description string `json:"description"`
	EvidenceURL string `json:"evidence_url"`
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
	FullName      string `json:"full_name"`
	PreferredName string `json:"preferred_name"`
	Headline      string `json:"headline"`
	About         string `json:"about"`
	PhotoURL      string `json:"photo_url"`
	Bio           string `json:"bio"`
	Location      string `json:"location"`
	Website       string `json:"website"`
	CoverBanner   string `json:"cover_banner"`
	LinkedInURL   string `json:"linkedin_url"`
	GitHubURL     string `json:"github_url"`

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

	PersonalBrand    string   `json:"personal_brand"`
	ElevatorPitch    string   `json:"elevator_pitch"`
	Industry         string   `json:"industry"`
	AnonymousMode    bool     `json:"anonymous_mode"`
	CareerObjectives string   `json:"career_objectives"`
	CareerHighlights []string `json:"career_highlights"`
	FunctionalAreas  []string `json:"functional_areas"`

	// Contact + collections. Nil (JSON key absent) leaves the section untouched;
	// a present array (even empty) replaces it. The editor autosave sends the
	// whole profile, so these carry every section the older subset silently dropped.
	Email   *string `json:"email"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`

	Experiences    []experienceDTO    `json:"experiences"`
	Educations     []educationDTO     `json:"educations"`
	Certifications []certificationDTO `json:"certifications"`
	Skills         []skillDTO         `json:"skills"`
	Languages      []languageDTO      `json:"languages"`
	Portfolio      []portfolioLinkDTO `json:"portfolio"`
	Projects       []projectDTO       `json:"projects"`
	Achievements   []achievementDTO   `json:"achievements_list"`
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
		UserID: p.UserID, FullName: p.Identity.FullName, Headline: p.Identity.Headline, About: p.Identity.About, PhotoURL: p.Identity.PhotoURL,
		Bio: p.Identity.Bio, Location: p.Identity.Location, Website: p.Identity.SocialLinks.Website, Version: p.Version,
		Pronouns: p.Identity.Pronouns, CareerStatus: p.Identity.CareerStatus,
		TransitionReason: "", TargetComebackTimeline: p.Identity.Availability,
		SupportsNeeded: []string{}, OpenToRemote: p.Preferences.OpenToRemote, OpenToRelocation: p.Preferences.OpenToRelocation,
		RelocationLocations: []string{}, DesiredRoles: p.Preferences.DesiredRoles, DesiredIndustries: p.Preferences.DesiredIndustries,
		EmploymentType: p.Preferences.EmploymentType, SalaryMin: p.Preferences.SalaryMin, SalaryMax: p.Preferences.SalaryMax,
		SalaryCurrency: p.Preferences.SalaryCurrency, SalaryVisible: p.Preferences.SalaryVisible, WorkMode: p.Preferences.RemotePreference,
		AvailabilityDate: p.Preferences.AvailabilityDate, NoticePeriod: p.Preferences.NoticePeriod,
		ReferralEligible: p.Preferences.ReferralEligible, EmailVerified: p.Verification.EmailVerified, PhoneVerified: p.Verification.PhoneVerified,
		LinkedinVerified: p.Verification.IdentityVerified, IdVerified: p.Verification.IdentityVerified,
		CareerNarrative: string(p.AICareerAssistant.GapAnalysis), CoachingMetadata: string(p.AICareerAssistant.InterviewPrep),
		WorkAuthStatus: p.Identity.WorkAuthorization, PassportNationality: p.Identity.Nationality,
		DrivingLicenseBool: p.Verification.IdentityVerified, DrivingLicenseType: p.Identity.VisaStatus,
		PreferredContactChannel: p.Identity.PreferredContactChannel, AccessibilityNeeds: p.Identity.VisaStatus,
		VideoIntroURL: p.Identity.VideoIntroURL, WillingToMentor: p.Preferences.WillingToMentor,
		Email: p.Identity.Email, Phone: p.Identity.Phone, Address: p.Identity.Address,
		AvgResponseTimeHours: p.Analytics.AvgResponseTimeHours, ProfileCompletenessScore: p.ProfileCompletenessScore,
		LastActiveAt: p.LastActiveAt.Format(time.RFC3339), BackgroundCheckConsent: p.Verification.IdentityVerified,
		BackgroundCheckConsentAt: "", JobAlertFrequency: p.Identity.VisaStatus,
		JobAlertChannel: p.Identity.VisaStatus, VisibilityProfile: p.Privacy.FieldVisibility["profile"],
		VisibilitySalary: p.Privacy.FieldVisibility["salary"], VisibilityTransitionReason: p.Privacy.FieldVisibility["transition_reason"],
		VisibilityExperience: p.Privacy.FieldVisibility["experience"], VisibilityEducation: p.Privacy.FieldVisibility["education"],
		VisibilityCertifications: p.Privacy.FieldVisibility["certifications"], VisibilitySkills: p.Privacy.FieldVisibility["skills"],
		VisibilityPortfolio: p.Privacy.FieldVisibility["portfolio"], VisibilityReferences: p.Privacy.FieldVisibility["references"],
		PreferredName: p.Identity.PreferredName, LinkedInURL: p.Identity.SocialLinks.LinkedIn, GitHubURL: p.Identity.SocialLinks.GitHub,
		PersonalBrand: p.Summary.PersonalBrandStatement, ElevatorPitch: p.Summary.ElevatorPitch, Industry: p.Summary.Industry,
		AnonymousMode:  p.Privacy.AnonymousMode,
		CoverBanner:    p.Identity.CoverURL,
		Experiences:    make([]experienceDTO, 0, len(p.Experiences)),
		Educations:     make([]educationDTO, 0, len(p.Educations)),
		Certifications: make([]certificationDTO, 0, len(p.Certifications)),
		Skills:         make([]skillDTO, 0, len(p.Skills)),
		Languages:      make([]languageDTO, 0, len(p.Identity.Languages)),
		Portfolio:      make([]portfolioLinkDTO, 0, len(p.Projects)),
		Projects:       make([]projectDTO, 0, len(p.Projects)),
		Achievements:   make([]achievementDTO, 0, len(p.Achievements)),
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
		r.Projects = append(r.Projects, projectDTO{
			ID:             pr.ID,
			Title:          pr.Title,
			Description:    pr.Description,
			RepositoryURL:  pr.RepositoryURL,
			LiveDemoURL:    pr.LiveDemoURL,
			VideoURL:       pr.VideoURL,
			Images:         pr.Images,
			Screenshots:    pr.Images,
			Technologies:   pr.Technologies,
			Timeline:       pr.Timeline,
			Metrics:        pr.Metrics,
			Awards:         pr.Awards,
			BusinessImpact: pr.BusinessImpact,
		})
	}
	for _, ac := range p.Achievements {
		r.Achievements = append(r.Achievements, achievementDTO{
			ID:          ac.ID,
			Title:       ac.Title,
			IssuerOrOrg: ac.IssuerOrOrg,
			Date:        ac.Date.Format("2006-01-02"),
			Category:    ac.Category,
			Description: ac.Description,
			EvidenceURL: ac.EvidenceURL,
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

// parseFlexDate accepts the date formats the profile editor emits: full
// YYYY-MM-DD, month-only YYYY-MM, and year-only YYYY. Returns the zero time for
// unparseable/sentinel values (e.g. "Present").
func parseFlexDate(s string) time.Time {
	for _, layout := range []string{"2006-01-02", "2006-01", "2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func (d experienceDTO) toDomain() domain.WorkExperience {
	start := parseFlexDate(d.StartDate)
	end := parseFlexDate(d.EndDate)
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
	grad := parseFlexDate(d.EndDate)
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
	issue := parseFlexDate(d.IssueDate)
	expiry := parseFlexDate(d.ExpiryDate)
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

func (d projectDTO) toDomain() domain.ProjectItem {
	imgs := d.Images
	if len(d.Screenshots) > 0 {
		imgs = d.Screenshots
	}
	return domain.ProjectItem{
		ID:             d.ID,
		Title:          d.Title,
		Description:    d.Description,
		RepositoryURL:  d.RepositoryURL,
		LiveDemoURL:    d.LiveDemoURL,
		VideoURL:       d.VideoURL,
		Images:         imgs,
		Technologies:   d.Technologies,
		Timeline:       d.Timeline,
		Metrics:        d.Metrics,
		Awards:         d.Awards,
		BusinessImpact: d.BusinessImpact,
	}
}

func (d achievementDTO) toDomain() domain.AchievementItem {
	return domain.AchievementItem{
		ID:          d.ID,
		Title:       d.Title,
		IssuerOrOrg: d.IssuerOrOrg,
		Date:        parseFlexDate(d.Date),
		Category:    d.Category,
		Description: d.Description,
		EvidenceURL: d.EvidenceURL,
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
	identity := &domain.IdentitySection{
		FullName:                r.FullName,
		PreferredName:           r.PreferredName,
		Headline:                r.Headline,
		About:                   r.About,
		Bio:                     r.Bio,
		Pronouns:                r.Pronouns,
		CareerStatus:            r.CareerStatus,
		PhotoURL:                r.PhotoURL,
		CoverURL:                r.CoverBanner,
		VideoIntroURL:           r.VideoIntroURL,
		Location:                r.Location,
		PreferredContactChannel: r.PreferredContactChannel,
		WorkAuthorization:       r.WorkAuthStatus,
		Nationality:             r.PassportNationality,
		SocialLinks: domain.SocialLinks{
			Website:  r.Website,
			LinkedIn: r.LinkedInURL,
			GitHub:   r.GitHubURL,
		},
	}
	if r.Email != nil {
		identity.Email = *r.Email
	}
	if r.Phone != nil {
		identity.Phone = *r.Phone
	}
	if r.Address != nil {
		identity.Address = *r.Address
	}
	if r.Languages != nil {
		langs := make([]domain.LanguageItem, 0, len(r.Languages))
		for _, l := range r.Languages {
			langs = append(langs, l.toDomain())
		}
		identity.Languages = langs
	}

	upd := domain.AggregateUpdate{
		Identity: identity,
		Summary: &domain.SummarySection{
			ExecutiveSummary:       r.About,
			CareerObjectives:       r.CareerObjectives,
			CareerHighlights:       r.CareerHighlights,
			FunctionalAreas:        r.FunctionalAreas,
			PersonalBrandStatement: r.PersonalBrand,
			ElevatorPitch:          r.ElevatorPitch,
			Industry:               r.Industry,
		},
		Preferences: &domain.CareerPreferences{
			DesiredRoles:      r.DesiredRoles,
			DesiredIndustries: r.DesiredIndustries,
			OpenToRelocation:  r.OpenToRelocation,
			OpenToRemote:      r.OpenToRemote,
			NoticePeriod:      r.NoticePeriod,
			RemotePreference:  r.WorkMode,
			SalaryMin:         r.SalaryMin,
			SalaryMax:         r.SalaryMax,
			SalaryCurrency:    r.SalaryCurrency,
			SalaryVisible:     r.SalaryVisible,
			AvailabilityDate:  r.AvailabilityDate,
			ReferralEligible:  r.ReferralEligible,
			WillingToMentor:   r.WillingToMentor,
			EmploymentType:    r.EmploymentType,
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
			AnonymousMode: r.AnonymousMode,
		},
		IsDraft: &isDraft,
	}

	// Collections: only replace a section the client actually sent (non-nil).
	if r.Experiences != nil {
		exps := make([]domain.WorkExperience, 0, len(r.Experiences))
		for _, e := range r.Experiences {
			exps = append(exps, e.toDomain())
		}
		upd.Experiences = &exps
	}
	if r.Educations != nil {
		edus := make([]domain.Education, 0, len(r.Educations))
		for _, e := range r.Educations {
			edus = append(edus, e.toDomain())
		}
		upd.Educations = &edus
	}
	if r.Certifications != nil {
		certs := make([]domain.CertificationItem, 0, len(r.Certifications))
		for _, c := range r.Certifications {
			certs = append(certs, c.toDomain())
		}
		upd.Certifications = &certs
	}
	if r.Skills != nil {
		skills := make([]domain.SkillItem, 0, len(r.Skills))
		for _, s := range r.Skills {
			skills = append(skills, s.toDomain())
		}
		upd.Skills = &skills
	}
	if r.Projects != nil {
		projs := make([]domain.ProjectItem, 0, len(r.Projects))
		for _, p := range r.Projects {
			projs = append(projs, p.toDomain())
		}
		upd.Projects = &projs
	} else if r.Portfolio != nil {
		projs := make([]domain.ProjectItem, 0, len(r.Portfolio))
		for _, p := range r.Portfolio {
			projs = append(projs, p.toDomain())
		}
		upd.Projects = &projs
	}
	if r.Achievements != nil {
		achs := make([]domain.AchievementItem, 0, len(r.Achievements))
		for _, a := range r.Achievements {
			achs = append(achs, a.toDomain())
		}
		upd.Achievements = &achs
	}

	return upd
}
