package application

import (
	"context"
	"errors"
	"time"

	"workspace-app/internal/cookies/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// GetPreferences fetches cookie consent settings by User ID (authenticated) or Anonymous ID.
func (s *Service) GetPreferences(ctx context.Context, userID, anonymousID string) (*domain.CookiePreferences, error) {
	if userID != "" {
		p, err := s.repo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return p, nil
		}
	}

	if anonymousID != "" {
		p, err := s.repo.GetByAnonymousID(ctx, anonymousID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return p, nil
		}
	}

	// Fallback to default essential-only preferences
	var uID *string
	if userID != "" {
		uID = &userID
	}
	var aID *string
	if anonymousID != "" {
		aID = &anonymousID
	}
	d := domain.DefaultConsent(uID, aID)
	return &d, nil
}

type SaveInput struct {
	UserID          *string
	AnonymousID     *string
	Functional      bool
	Analytics       bool
	Marketing       bool
	Performance     bool
	Personalization bool
	AIPreferences   bool
	ConsentVersion  string
	IPAddress       string
	Country         string
	UserAgent       string
}

// SavePreferences persists cookie consent choices and auto-merges if needed.
func (s *Service) SavePreferences(ctx context.Context, in SaveInput) (*domain.CookiePreferences, error) {
	p := &domain.CookiePreferences{
		UserID:          in.UserID,
		AnonymousID:     in.AnonymousID,
		Essential:       true, // Hardcoded safety check
		Functional:      in.Functional,
		Analytics:       in.Analytics,
		Marketing:       in.Marketing,
		Performance:     in.Performance,
		Personalization: in.Personalization,
		AIPreferences:   in.AIPreferences,
		ConsentVersion:  in.ConsentVersion,
		AcceptedAt:      time.Now(),
		IPAddress:       in.IPAddress,
		Country:         in.Country,
		UserAgent:       in.UserAgent,
	}

	if p.ConsentVersion == "" {
		p.ConsentVersion = "1.0"
	}

	err := s.repo.Save(ctx, p)
	if err != nil {
		return nil, err
	}

	// Trigger async or sync merge if both IDs exist
	if in.UserID != nil && in.AnonymousID != nil && *in.AnonymousID != "" {
		_ = s.repo.Merge(ctx, *in.AnonymousID, *in.UserID)
	}

	return p, nil
}

// DeletePreferences clears saved cookie choices.
func (s *Service) DeletePreferences(ctx context.Context, userID, anonymousID string) error {
	if userID == "" && anonymousID == "" {
		return errors.New("must provide either user_id or anonymous_id to delete")
	}
	return s.repo.Delete(ctx, userID, anonymousID)
}

// MergePreferences consolidates anonymous cookie selections into a logged-in user profile.
func (s *Service) MergePreferences(ctx context.Context, anonymousID, userID string) error {
	if anonymousID == "" || userID == "" {
		return errors.New("both anonymous_id and user_id are required to merge")
	}
	return s.repo.Merge(ctx, anonymousID, userID)
}
