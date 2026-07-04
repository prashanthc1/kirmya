// Package application implements Notification use cases.
package application

import (
	"context"
	"strings"

	"workspace-app/internal/notifications/domain"
)

// PrefChecker lets the notifier honour a user's in-app notification toggles. A
// nil checker (the default) delivers every notification.
type PrefChecker interface {
	WantsInApp(ctx context.Context, userID, category string) (bool, error)
}

type Service struct {
	repo  domain.Repository
	hub   *Hub
	prefs PrefChecker
}

func NewService(repo domain.Repository, hub *Hub) *Service {
	return &Service{repo: repo, hub: hub}
}

// SetPrefChecker injects the preference checker used to honour in-app toggles.
func (s *Service) SetPrefChecker(p PrefChecker) { s.prefs = p }

// Notify creates a notification for a user (used by event subscribers) and
// pushes it to any of the user's live SSE subscribers.
func (s *Service) Notify(ctx context.Context, userID, typ, title, body, link string) error {
	if userID == "" {
		return nil
	}
	// Honour the recipient's in-app toggle for this category.
	if s.prefs != nil {
		if want, err := s.prefs.WantsInApp(ctx, userID, categoryForType(typ)); err == nil && !want {
			return nil
		}
	}
	n := &domain.Notification{UserID: userID, Type: typ, Title: title, Body: body, Link: link}
	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}
	if s.hub != nil {
		s.hub.Publish(*n)
	}
	return nil
}

// Subscribe registers a real-time subscriber for the user's notifications.
func (s *Service) Subscribe(userID string) (<-chan domain.Notification, func()) {
	return s.hub.Subscribe(userID)
}

// List returns a page of the user's notifications (unread first, then newest).
func (s *Service) List(ctx context.Context, userID string, limit, offset int) ([]domain.Notification, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *Service) MarkRead(ctx context.Context, userID, id string) error {
	return s.repo.MarkRead(ctx, userID, id)
}

func (s *Service) MarkAllRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *Service) UnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.UnreadCount(ctx, userID)
}

// categoryForType maps a notification type to a settings notification category.
func categoryForType(typ string) string {
	switch {
	case strings.HasPrefix(typ, "referral"):
		return "referrals"
	case strings.HasPrefix(typ, "mentorship"):
		return "mentorship"
	case typ == "message":
		return "messages"
	case strings.HasPrefix(typ, "job"):
		return "jobs"
	default:
		return ""
	}
}
