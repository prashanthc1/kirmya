package api

import (
	"time"

	"workspace-app/internal/settings/domain"
)

// settingsDTO is the wire representation of a user's full settings object.
type settingsDTO struct {
	General       generalDTO       `json:"general"`
	Privacy       privacyDTO       `json:"privacy"`
	Notifications notificationsDTO `json:"notifications"`
	Security      securityDTO      `json:"security"`
	Accessibility accessibilityDTO `json:"accessibility"`
	AI            aiDTO            `json:"ai"`
	Learning      learningDTO      `json:"learning"`
	Version       int              `json:"version"`
	UpdatedAt     string           `json:"updated_at"`
}

type generalDTO struct {
	Language    string `json:"language"`
	Timezone    string `json:"timezone"`
	Theme       string `json:"theme"`
	EmailDigest string `json:"email_digest"`
}

type privacyDTO struct {
	ProfileVisibility string `json:"profile_visibility"`
	ShowEmail         bool   `json:"show_email"`
	Discoverable      bool   `json:"discoverable"`
	AllowMessages     string `json:"allow_messages"`
}

type notificationsDTO struct {
	EmailJobs       bool `json:"email_jobs"`
	EmailMentorship bool `json:"email_mentorship"`
	EmailMessages   bool `json:"email_messages"`
	EmailReferrals  bool `json:"email_referrals"`
	InAppJobs       bool `json:"inapp_jobs"`
	InAppMentorship bool `json:"inapp_mentorship"`
	InAppMessages   bool `json:"inapp_messages"`
	InAppReferrals  bool `json:"inapp_referrals"`
}

type securityDTO struct {
	LoginAlerts bool `json:"login_alerts"`
}

type accessibilityDTO struct {
	FontSize                        string `json:"font_size"`
	HighContrast                    bool   `json:"high_contrast"`
	ReducedMotion                   bool   `json:"reduced_motion"`
	CompactMode                     bool   `json:"compact_mode"`
	DefaultLandingPage              string `json:"default_landing_page"`
	AccessibilityKeyboardNavigation bool   `json:"accessibility_keyboard_navigation"`
	AccessibilityScreenReader       bool   `json:"accessibility_screen_reader"`
	AccessibilityFocusIndicators    bool   `json:"accessibility_focus_indicators"`
}

type aiDTO struct {
	EnableAIAssistant         bool `json:"enable_ai_assistant"`
	AIJobRecommendations      bool `json:"ai_job_recommendations"`
	AIResumeSuggestions       bool `json:"ai_resume_suggestions"`
	AIRoadmapSuggestions      bool `json:"ai_roadmap_suggestions"`
	AISkillGapAnalysis        bool `json:"ai_skill_gap_analysis"`
	AIInterviewPrep           bool `json:"ai_interview_prep"`
	AILearningRecommendations bool `json:"ai_learning_recommendations"`
}

type learningDTO struct {
	LearningGoals          []string `json:"learning_goals"`
	TechnologiesOfInterest []string `json:"technologies_of_interest"`
	CertificationGoals     []string `json:"certification_goals"`
	LearningReminders      bool     `json:"learning_reminders"`
}

type connectedAccountDTO struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	ProviderUID string `json:"provider_uid"`
	CreatedAt   string `json:"created_at"`
}

type cookieConsentDTO struct {
	Functional        bool   `json:"functional"`
	Analytics         bool   `json:"analytics"`
	AIPersonalization bool   `json:"ai_personalization"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

type activeSessionDTO struct {
	ID        string `json:"id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

type securityHistoryDTO struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type unifiedSettingsPatchRequest struct {
	General       *generalDTO       `json:"general,omitempty"`
	Privacy       *privacyDTO       `json:"privacy,omitempty"`
	Notifications *notificationsDTO `json:"notifications,omitempty"`
	Security      *securityDTO      `json:"security,omitempty"`
	Accessibility *accessibilityDTO `json:"accessibility,omitempty"`
	AI            *aiDTO            `json:"ai,omitempty"`
	Learning      *learningDTO      `json:"learning,omitempty"`
}

func toDTO(s *domain.UserSettings) settingsDTO {
	return settingsDTO{
		General: generalDTO{
			Language:    s.Language,
			Timezone:    s.Timezone,
			Theme:       s.Theme,
			EmailDigest: s.EmailDigest,
		},
		Privacy: privacyDTO{
			ProfileVisibility: s.ProfileVisibility,
			ShowEmail:         s.ShowEmail,
			Discoverable:      s.Discoverable,
			AllowMessages:     s.AllowMessages,
		},
		Notifications: toNotificationsDTO(s.Notifications),
		Security:      securityDTO{LoginAlerts: s.LoginAlerts},
		Accessibility: accessibilityDTO{
			FontSize:                        s.FontSize,
			HighContrast:                    s.HighContrast,
			ReducedMotion:                   s.ReducedMotion,
			CompactMode:                     s.CompactMode,
			DefaultLandingPage:              s.DefaultLandingPage,
			AccessibilityKeyboardNavigation: s.AccessibilityKeyboardNavigation,
			AccessibilityScreenReader:       s.AccessibilityScreenReader,
			AccessibilityFocusIndicators:    s.AccessibilityFocusIndicators,
		},
		AI: aiDTO{
			EnableAIAssistant:         s.EnableAIAssistant,
			AIJobRecommendations:      s.AIJobRecommendations,
			AIResumeSuggestions:       s.AIResumeSuggestions,
			AIRoadmapSuggestions:      s.AIRoadmapSuggestions,
			AISkillGapAnalysis:        s.AISkillGapAnalysis,
			AIInterviewPrep:           s.AIInterviewPrep,
			AILearningRecommendations: s.AILearningRecommendations,
		},
		Learning: learningDTO{
			LearningGoals:          s.LearningGoals,
			TechnologiesOfInterest: s.TechnologiesOfInterest,
			CertificationGoals:     s.CertificationGoals,
			LearningReminders:      s.LearningReminders,
		},
		Version:   s.Version,
		UpdatedAt: s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toNotificationsDTO(n domain.NotificationPrefs) notificationsDTO {
	return notificationsDTO{
		EmailJobs:       n.EmailJobs,
		EmailMentorship: n.EmailMentorship,
		EmailMessages:   n.EmailMessages,
		EmailReferrals:  n.EmailReferrals,
		InAppJobs:       n.InAppJobs,
		InAppMentorship: n.InAppMentorship,
		InAppMessages:   n.InAppMessages,
		InAppReferrals:  n.InAppReferrals,
	}
}

func (d notificationsDTO) toDomain() domain.NotificationPrefs {
	return domain.NotificationPrefs{
		EmailJobs:       d.EmailJobs,
		EmailMentorship: d.EmailMentorship,
		EmailMessages:   d.EmailMessages,
		EmailReferrals:  d.EmailReferrals,
		InAppJobs:       d.InAppJobs,
		InAppMentorship: d.InAppMentorship,
		InAppMessages:   d.InAppMessages,
		InAppReferrals:  d.InAppReferrals,
	}
}

type profileSettingsDTO struct {
	Username          string            `json:"username"`
	CustomURL         string            `json:"custom_url"`
	ProfileVisibility string            `json:"profile_visibility"`
	FieldVisibility   map[string]string `json:"field_visibility"`
	OpenToWork        bool              `json:"open_to_work"`
	ReferralEligible  bool              `json:"referral_eligible"`
	WillingToMentor   bool              `json:"willing_to_mentor"`
}
