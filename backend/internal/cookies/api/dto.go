package api

import (
	"time"

	"workspace-app/internal/cookies/domain"
)

type preferencesDTO struct {
	ID              string    `json:"id"`
	UserID          *string   `json:"user_id,omitempty"`
	AnonymousID     *string   `json:"anonymous_id,omitempty"`
	Essential       bool      `json:"essential"`
	Functional      bool      `json:"functional"`
	Analytics       bool      `json:"analytics"`
	Marketing       bool      `json:"marketing"`
	Performance     bool      `json:"performance"`
	Personalization bool      `json:"personalization"`
	AIPreferences   bool      `json:"ai_preferences"`
	ConsentVersion  string    `json:"consent_version"`
	AcceptedAt      time.Time `json:"accepted_at"`
	IPAddress       string    `json:"ip_address,omitempty"`
	Country         string    `json:"country,omitempty"`
	UserAgent       string    `json:"user_agent,omitempty"`
}

type saveRequest struct {
	AnonymousID     *string `json:"anonymous_id,omitempty"`
	Functional      bool    `json:"functional"`
	Analytics       bool    `json:"analytics"`
	Marketing       bool    `json:"marketing"`
	Performance     bool    `json:"performance"`
	Personalization bool    `json:"personalization"`
	AIPreferences   bool    `json:"ai_preferences"`
	ConsentVersion  string  `json:"consent_version"`
}

func toDTO(p *domain.CookiePreferences) preferencesDTO {
	return preferencesDTO{
		ID:              p.ID,
		UserID:          p.UserID,
		AnonymousID:     p.AnonymousID,
		Essential:       p.Essential,
		Functional:      p.Functional,
		Analytics:       p.Analytics,
		Marketing:       p.Marketing,
		Performance:     p.Performance,
		Personalization: p.Personalization,
		AIPreferences:   p.AIPreferences,
		ConsentVersion:  p.ConsentVersion,
		AcceptedAt:      p.AcceptedAt,
		IPAddress:       p.IPAddress,
		Country:         p.Country,
		UserAgent:       p.UserAgent,
	}
}
