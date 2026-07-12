package domain

import "time"

// CookiePreferences represents the user's granular cookie choices.
type CookiePreferences struct {
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
	IPAddress       string    `json:"ip_address"`
	Country         string    `json:"country"`
	UserAgent       string    `json:"user_agent"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DefaultConsent returns a baseline of essential-only cookies.
func DefaultConsent(userID *string, anonymousID *string) CookiePreferences {
	now := time.Now()
	return CookiePreferences{
		Essential:       true,
		Functional:      false,
		Analytics:       false,
		Marketing:       false,
		Performance:     false,
		Personalization: false,
		AIPreferences:   false,
		ConsentVersion:  "1.0",
		AcceptedAt:      now,
		CreatedAt:       now,
		UpdatedAt:       now,
		UserID:          userID,
		AnonymousID:     anonymousID,
	}
}
