// Package application implements the Settings use cases. It depends only on the
// domain ports; infrastructure adapters are injected in module.go.
package application

import (
	"context"
	"errors"
	"strings"

	"workspace-app/internal/settings/domain"
)

// EventPublisher publishes domain events onto the (in-process, NATS-ready) bus.
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

// ValidationError is a user-facing input error mapped to HTTP 400 in api/.
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

// Service is the settings use-case service.
type Service struct {
	repo   domain.Repository
	events EventPublisher
}

func NewService(repo domain.Repository, events EventPublisher) *Service {
	return &Service{repo: repo, events: events}
}

// Get returns the caller's settings, materialising defaults on first access so
// the caller always receives a complete object.
func (s *Service) Get(ctx context.Context, userID string) (*domain.UserSettings, error) {
	return s.repo.EnsureDefaults(ctx, userID)
}

// GeneralInput is the general-settings command (full section replace).
type GeneralInput struct {
	Language    string
	Timezone    string
	Theme       string
	EmailDigest string
}

// UpdateGeneral validates and persists the general section.
func (s *Service) UpdateGeneral(ctx context.Context, userID string, in GeneralInput) (*domain.UserSettings, error) {
	lang := strings.TrimSpace(in.Language)
	tz := strings.TrimSpace(in.Timezone)
	if lang == "" || len(lang) > 10 {
		return nil, ValidationError{"language is required and must be at most 10 characters"}
	}
	if tz == "" || len(tz) > 64 {
		return nil, ValidationError{"timezone is required and must be at most 64 characters"}
	}
	if !domain.ValidThemes[in.Theme] {
		return nil, ValidationError{"theme must be one of: light, dark, system"}
	}
	if !domain.ValidDigests[in.EmailDigest] {
		return nil, ValidationError{"email_digest must be one of: off, daily, weekly"}
	}
	return s.mutate(ctx, userID, func(cur *domain.UserSettings) {
		cur.Language = lang
		cur.Timezone = tz
		cur.Theme = in.Theme
		cur.EmailDigest = in.EmailDigest
	}, domain.EventSettingsUpdated)
}

// PrivacyInput is the privacy-settings command (full section replace).
type PrivacyInput struct {
	ProfileVisibility string
	ShowEmail         bool
	Discoverable      bool
	AllowMessages     string
}

// UpdatePrivacy validates and persists the privacy section.
func (s *Service) UpdatePrivacy(ctx context.Context, userID string, in PrivacyInput) (*domain.UserSettings, error) {
	if !domain.ValidVisibilities[in.ProfileVisibility] {
		return nil, ValidationError{"profile_visibility must be one of: public, network, private"}
	}
	if !domain.ValidMessagePolicy[in.AllowMessages] {
		return nil, ValidationError{"allow_messages must be one of: everyone, network, none"}
	}
	return s.mutate(ctx, userID, func(cur *domain.UserSettings) {
		cur.ProfileVisibility = in.ProfileVisibility
		cur.ShowEmail = in.ShowEmail
		cur.Discoverable = in.Discoverable
		cur.AllowMessages = in.AllowMessages
	}, domain.EventPrivacyChanged)
}

// UpdateNotifications persists the notification toggles (full section replace).
func (s *Service) UpdateNotifications(ctx context.Context, userID string, in domain.NotificationPrefs) (*domain.UserSettings, error) {
	return s.mutate(ctx, userID, func(cur *domain.UserSettings) {
		cur.Notifications = in
	}, domain.EventNotificationsChanged)
}

// SecurityInput is the security-preferences command.
type SecurityInput struct {
	LoginAlerts bool
}

// UpdateSecurity persists the security-preferences section.
func (s *Service) UpdateSecurity(ctx context.Context, userID string, in SecurityInput) (*domain.UserSettings, error) {
	return s.mutate(ctx, userID, func(cur *domain.UserSettings) {
		cur.LoginAlerts = in.LoginAlerts
	}, domain.EventSettingsUpdated)
}

// mutate loads the current settings (creating defaults if absent), applies the
// given change, persists with optimistic locking, and publishes an event.
func (s *Service) mutate(ctx context.Context, userID string, apply func(*domain.UserSettings), eventType string) (*domain.UserSettings, error) {
	cur, err := s.repo.EnsureDefaults(ctx, userID)
	if err != nil {
		return nil, err
	}
	apply(cur)
	if err := s.repo.Update(ctx, cur); err != nil {
		return nil, err
	}
	s.publish(ctx, eventType, userID, map[string]any{"user_id": userID})
	return cur, nil
}

func (s *Service) publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, eventType, aggregateID, payload)
	}
}

// --- Read API consumed by other modules to enforce a user's preferences. ---

// read returns the user's settings WITHOUT materialising a row (no write),
// falling back to defaults when none exists yet. Used by the enforcement reads
// below so a hot read path never performs an upsert.
func (s *Service) read(ctx context.Context, userID string) (*domain.UserSettings, error) {
	cur, err := s.repo.Get(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		d := domain.Defaults(userID)
		return &d, nil
	}
	if err != nil {
		return nil, err
	}
	return cur, nil
}

// WantsInApp reports whether the user wants in-app notifications for a category
// ("jobs"|"mentorship"|"messages"|"referrals"). Unknown categories default to true.
func (s *Service) WantsInApp(ctx context.Context, userID, category string) (bool, error) {
	st, err := s.read(ctx, userID)
	if err != nil {
		return true, err
	}
	n := st.Notifications
	switch category {
	case "jobs":
		return n.InAppJobs, nil
	case "mentorship":
		return n.InAppMentorship, nil
	case "messages":
		return n.InAppMessages, nil
	case "referrals":
		return n.InAppReferrals, nil
	default:
		return true, nil
	}
}

// MessagePolicy returns who may message the user: "everyone"|"network"|"none".
func (s *Service) MessagePolicy(ctx context.Context, userID string) (string, error) {
	st, err := s.read(ctx, userID)
	if err != nil {
		return domain.MessagesEveryone, err
	}
	return st.AllowMessages, nil
}

// ProfileVisibility returns the user's profile visibility:
// "public"|"network"|"private".
func (s *Service) ProfileVisibility(ctx context.Context, userID string) (string, error) {
	st, err := s.read(ctx, userID)
	if err != nil {
		return domain.VisibilityPublic, err
	}
	return st.ProfileVisibility, nil
}

// ShowEmail reports whether the user exposes their email on their profile.
func (s *Service) ShowEmail(ctx context.Context, userID string) (bool, error) {
	st, err := s.read(ctx, userID)
	if err != nil {
		return false, err
	}
	return st.ShowEmail, nil
}
