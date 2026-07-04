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
		Version:       s.Version,
		UpdatedAt:     s.UpdatedAt.UTC().Format(time.RFC3339),
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
