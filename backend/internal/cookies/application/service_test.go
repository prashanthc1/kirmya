package application_test

import (
	"context"
	"testing"
	"time"

	"workspace-app/internal/cookies/application"
	"workspace-app/internal/cookies/domain"
)

type mockRepo struct {
	byUserID map[string]*domain.CookiePreferences
	byAnonID map[string]*domain.CookiePreferences
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		byUserID: make(map[string]*domain.CookiePreferences),
		byAnonID: make(map[string]*domain.CookiePreferences),
	}
}

func (m *mockRepo) GetByUserID(ctx context.Context, userID string) (*domain.CookiePreferences, error) {
	return m.byUserID[userID], nil
}

func (m *mockRepo) GetByAnonymousID(ctx context.Context, anonymousID string) (*domain.CookiePreferences, error) {
	return m.byAnonID[anonymousID], nil
}

func (m *mockRepo) Save(ctx context.Context, p *domain.CookiePreferences) error {
	p.ID = "test-uuid"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.UserID != nil {
		m.byUserID[*p.UserID] = p
	}
	if p.AnonymousID != nil {
		m.byAnonID[*p.AnonymousID] = p
	}
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, userID, anonymousID string) error {
	delete(m.byUserID, userID)
	delete(m.byAnonID, anonymousID)
	return nil
}

func (m *mockRepo) Merge(ctx context.Context, anonymousID, userID string) error {
	anonVal, ok := m.byAnonID[anonymousID]
	if !ok {
		return nil
	}
	// Copy choices to user
	userVal := *anonVal
	userVal.UserID = &userID
	userVal.AnonymousID = nil
	m.byUserID[userID] = &userVal
	delete(m.byAnonID, anonymousID)
	return nil
}

func TestCookieService_GetPreferences_Default(t *testing.T) {
	repo := newMockRepo()
	svc := application.NewService(repo)

	prefs, err := svc.GetPreferences(context.Background(), "", "anon-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !prefs.Essential {
		t.Error("essential cookies should be enabled by default")
	}
	if prefs.Functional {
		t.Error("non-essential functional cookies should be disabled by default")
	}
}

func TestCookieService_SaveAndGet(t *testing.T) {
	repo := newMockRepo()
	svc := application.NewService(repo)

	in := application.SaveInput{
		AnonymousID:    pointer("anon-123"),
		Functional:     true,
		Analytics:      true,
		ConsentVersion: "1.0",
	}

	saved, err := svc.SavePreferences(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	if !saved.Functional || !saved.Analytics {
		t.Error("expected functional and analytics flags to be saved as true")
	}

	// Fetch it back
	fetched, err := svc.GetPreferences(context.Background(), "", "anon-123")
	if err != nil {
		t.Fatalf("unexpected error fetching: %v", err)
	}

	if fetched.ID != saved.ID || !fetched.Functional {
		t.Error("fetched preferences mismatch")
	}
}

func TestCookieService_MergePreferences(t *testing.T) {
	repo := newMockRepo()
	svc := application.NewService(repo)

	// Save anonymous choices
	_, _ = svc.SavePreferences(context.Background(), application.SaveInput{
		AnonymousID: pointer("anon-123"),
		Functional:  true,
		Analytics:   true,
	})

	// Save logged in user (with link to anonymous ID, should trigger auto merge)
	_, err := svc.SavePreferences(context.Background(), application.SaveInput{
		UserID:      pointer("user-abc"),
		AnonymousID: pointer("anon-123"),
		Functional:  true,
		Analytics:   true,
	})
	if err != nil {
		t.Fatalf("failed saving logged in user preferences: %v", err)
	}

	// Verify merged
	uPref, err := svc.GetPreferences(context.Background(), "user-abc", "")
	if err != nil {
		t.Fatalf("failed getting user preferences: %v", err)
	}

	if uPref == nil || !uPref.Functional || !uPref.Analytics {
		t.Error("expected user preferences to have merged values")
	}

	// Verify anonymous is deleted
	anonPref, err := repo.GetByAnonymousID(context.Background(), "anon-123")
	if err != nil {
		t.Fatalf("failed checking anonymous record: %v", err)
	}
	if anonPref != nil {
		t.Error("expected anonymous record to be deleted after merge")
	}
}

func pointer(s string) *string {
	return &s
}
