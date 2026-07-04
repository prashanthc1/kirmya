package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/settings/domain"
)

func newSvc() (*Service, *fakeRepo, *recordingEvents) {
	repo := newFakeRepo()
	ev := &recordingEvents{}
	return NewService(repo, ev), repo, ev
}

func TestGetMaterialisesDefaults(t *testing.T) {
	svc, _, _ := newSvc()
	s, err := svc.Get(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if s.Theme != domain.ThemeSystem || s.EmailDigest != domain.DigestWeekly {
		t.Fatalf("unexpected general defaults: %+v", s)
	}
	if s.ProfileVisibility != domain.VisibilityPublic || !s.Discoverable {
		t.Fatalf("unexpected privacy defaults: %+v", s)
	}
	if !s.Notifications.EmailJobs || !s.LoginAlerts {
		t.Fatalf("unexpected notification/security defaults: %+v", s)
	}
}

func TestUpdateGeneralPersistsAndBumpsVersion(t *testing.T) {
	svc, _, ev := newSvc()
	s, err := svc.UpdateGeneral(context.Background(), "user-1", GeneralInput{
		Language: "fr", Timezone: "Europe/Paris", Theme: domain.ThemeDark, EmailDigest: domain.DigestDaily,
	})
	if err != nil {
		t.Fatalf("UpdateGeneral: %v", err)
	}
	if s.Language != "fr" || s.Theme != domain.ThemeDark || s.Version != 2 {
		t.Fatalf("update not applied: %+v", s)
	}
	// Read-back reflects the change.
	got, _ := svc.Get(context.Background(), "user-1")
	if got.Timezone != "Europe/Paris" {
		t.Fatalf("read-back mismatch: %+v", got)
	}
	if len(ev.events) != 1 || ev.events[0] != domain.EventSettingsUpdated {
		t.Fatalf("expected SettingsUpdated event, got %v", ev.events)
	}
}

func TestUpdateGeneralRejectsBadEnums(t *testing.T) {
	svc, _, _ := newSvc()
	cases := []GeneralInput{
		{Language: "en", Timezone: "UTC", Theme: "neon", EmailDigest: domain.DigestOff},
		{Language: "en", Timezone: "UTC", Theme: domain.ThemeLight, EmailDigest: "hourly"},
		{Language: "", Timezone: "UTC", Theme: domain.ThemeLight, EmailDigest: domain.DigestOff},
		{Language: "en", Timezone: "", Theme: domain.ThemeLight, EmailDigest: domain.DigestOff},
	}
	for i, in := range cases {
		if _, err := svc.UpdateGeneral(context.Background(), "u", in); err == nil {
			t.Fatalf("case %d: expected validation error", i)
		} else {
			var ve ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("case %d: expected ValidationError, got %T", i, err)
			}
		}
	}
}

func TestUpdatePrivacyValidatesEnums(t *testing.T) {
	svc, _, ev := newSvc()
	if _, err := svc.UpdatePrivacy(context.Background(), "u", PrivacyInput{ProfileVisibility: "secret", AllowMessages: domain.MessagesNone}); err == nil {
		t.Fatal("expected error for bad visibility")
	}
	s, err := svc.UpdatePrivacy(context.Background(), "u", PrivacyInput{
		ProfileVisibility: domain.VisibilityPrivate, ShowEmail: true, Discoverable: false, AllowMessages: domain.MessagesNetwork,
	})
	if err != nil {
		t.Fatalf("UpdatePrivacy: %v", err)
	}
	if !s.ShowEmail || s.Discoverable || s.AllowMessages != domain.MessagesNetwork {
		t.Fatalf("privacy not applied: %+v", s)
	}
	if len(ev.events) != 1 || ev.events[0] != domain.EventPrivacyChanged {
		t.Fatalf("expected PrivacySettingsChanged, got %v", ev.events)
	}
}

func TestUpdateNotificationsReplacesSection(t *testing.T) {
	svc, _, ev := newSvc()
	in := domain.NotificationPrefs{EmailJobs: false, EmailMessages: true}
	s, err := svc.UpdateNotifications(context.Background(), "u", in)
	if err != nil {
		t.Fatalf("UpdateNotifications: %v", err)
	}
	if s.Notifications.EmailJobs || !s.Notifications.EmailMessages {
		t.Fatalf("notifications not applied: %+v", s.Notifications)
	}
	if len(ev.events) != 1 || ev.events[0] != domain.EventNotificationsChanged {
		t.Fatalf("expected NotificationSettingsChanged, got %v", ev.events)
	}
}

func TestUpdateSecurityTogglesLoginAlerts(t *testing.T) {
	svc, _, _ := newSvc()
	s, err := svc.UpdateSecurity(context.Background(), "u", SecurityInput{LoginAlerts: false})
	if err != nil {
		t.Fatalf("UpdateSecurity: %v", err)
	}
	if s.LoginAlerts {
		t.Fatal("login_alerts should be false")
	}
}

func TestUpdatePropagatesOptimisticLock(t *testing.T) {
	svc, repo, _ := newSvc()
	if _, err := svc.Get(context.Background(), "u"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	repo.failNext = domain.ErrOptimisticLock
	_, err := svc.UpdateSecurity(context.Background(), "u", SecurityInput{LoginAlerts: false})
	if !errors.Is(err, domain.ErrOptimisticLock) {
		t.Fatalf("expected ErrOptimisticLock, got %v", err)
	}
}

func TestNilEventBusIsSafe(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	if _, err := svc.UpdateSecurity(context.Background(), "u", SecurityInput{LoginAlerts: true}); err != nil {
		t.Fatalf("nil bus should be safe: %v", err)
	}
}

func TestReadAPIDefaultsAndReflectsUpdates(t *testing.T) {
	svc, _, _ := newSvc()
	ctx := context.Background()

	// Defaults are returned without materialising a row.
	if ok, err := svc.WantsInApp(ctx, "u", "messages"); err != nil || !ok {
		t.Fatalf("default WantsInApp(messages) = %v, %v", ok, err)
	}
	if pol, err := svc.MessagePolicy(ctx, "u"); err != nil || pol != domain.MessagesEveryone {
		t.Fatalf("default MessagePolicy = %v, %v", pol, err)
	}
	if vis, err := svc.ProfileVisibility(ctx, "u"); err != nil || vis != domain.VisibilityPublic {
		t.Fatalf("default ProfileVisibility = %v, %v", vis, err)
	}

	// Reads reflect persisted updates.
	if _, err := svc.UpdatePrivacy(ctx, "u", PrivacyInput{ProfileVisibility: domain.VisibilityPrivate, AllowMessages: domain.MessagesNone}); err != nil {
		t.Fatalf("UpdatePrivacy: %v", err)
	}
	if vis, _ := svc.ProfileVisibility(ctx, "u"); vis != domain.VisibilityPrivate {
		t.Fatalf("ProfileVisibility after update = %v", vis)
	}
	if pol, _ := svc.MessagePolicy(ctx, "u"); pol != domain.MessagesNone {
		t.Fatalf("MessagePolicy after update = %v", pol)
	}
	if _, err := svc.UpdateNotifications(ctx, "u", domain.NotificationPrefs{InAppMessages: false}); err != nil {
		t.Fatalf("UpdateNotifications: %v", err)
	}
	if ok, _ := svc.WantsInApp(ctx, "u", "messages"); ok {
		t.Fatal("WantsInApp(messages) should be false after opt-out")
	}
}
