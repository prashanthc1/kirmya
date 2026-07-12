package application

import (
	"context"
	"sync"

	"workspace-app/internal/settings/domain"
)

// fakeRepo is an in-memory settings repository for unit tests. It mirrors the
// optimistic-locking semantics of the postgres adapter.
type fakeRepo struct {
	mu       sync.Mutex
	byUser   map[string]*domain.UserSettings
	failNext error // when set, the next Update returns this error
}

func newFakeRepo() *fakeRepo { return &fakeRepo{byUser: map[string]*domain.UserSettings{}} }

func (f *fakeRepo) Get(_ context.Context, userID string) (*domain.UserSettings, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s, ok := f.byUser[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	out := *s
	return &out, nil
}

func (f *fakeRepo) EnsureDefaults(_ context.Context, userID string) (*domain.UserSettings, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s, ok := f.byUser[userID]
	if !ok {
		d := domain.Defaults(userID)
		f.byUser[userID] = &d
		s = &d
	}
	out := *s
	return &out, nil
}

func (f *fakeRepo) Update(_ context.Context, s *domain.UserSettings) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failNext != nil {
		err := f.failNext
		f.failNext = nil
		return err
	}
	cur, ok := f.byUser[s.UserID]
	if !ok {
		return domain.ErrNotFound
	}
	if cur.Version != s.Version {
		return domain.ErrOptimisticLock
	}
	clone := *s
	clone.Version = cur.Version + 1
	f.byUser[s.UserID] = &clone
	s.Version = clone.Version
	return nil
}

func (f *fakeRepo) ListConnectedAccounts(_ context.Context, _ string) ([]domain.ConnectedAccount, error) {
	return nil, nil
}
func (f *fakeRepo) DisconnectAccount(_ context.Context, _, _ string) error {
	return nil
}
func (f *fakeRepo) GetCookieConsent(_ context.Context, userID string) (*domain.CookieConsent, error) {
	return &domain.CookieConsent{UserID: userID, Essential: true}, nil
}
func (f *fakeRepo) SaveCookieConsent(_ context.Context, _ *domain.CookieConsent) error {
	return nil
}
func (f *fakeRepo) ListActiveSessions(_ context.Context, _ string) ([]domain.ActiveSession, error) {
	return nil, nil
}
func (f *fakeRepo) RevokeSession(_ context.Context, _, _ string) error {
	return nil
}
func (f *fakeRepo) ListSecurityHistory(_ context.Context, _ string) ([]domain.SecurityHistoryEntry, error) {
	return nil, nil
}
func (f *fakeRepo) WriteSecurityLog(_ context.Context, _, _, _ string) error {
	return nil
}
func (f *fakeRepo) GetProfileSettings(_ context.Context, _ string) (string, string, string, map[string]string, bool, bool, bool, error) {
	return "", "", "", nil, false, false, false, nil
}
func (f *fakeRepo) UpdateProfileSettings(_ context.Context, _ string, _, _, _ string, _ map[string]string, _, _, _ bool) error {
	return nil
}

// recordingEvents captures published events for assertions.
type recordingEvents struct {
	mu     sync.Mutex
	events []string
}

func (e *recordingEvents) Publish(_ context.Context, eventType, _ string, _ map[string]any) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.events = append(e.events, eventType)
	return nil
}
