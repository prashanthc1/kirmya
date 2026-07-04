package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNotifyPersistsAndPushesToSubscriber(t *testing.T) {
	repo := newFakeRepo()
	hub := NewHub(nil)
	svc := NewService(repo, hub)
	ctx := context.Background()

	ch, cancel := svc.Subscribe("u1")
	defer cancel()

	if err := svc.Notify(ctx, "u1", "job_posted", "New job", "A role you may like", "/jobs/1"); err != nil {
		t.Fatalf("notify: %v", err)
	}

	// Persisted with an ID.
	list, err := svc.List(ctx, "u1", 100, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].ID == "" || list[0].Title != "New job" {
		t.Fatalf("unexpected persisted notifications: %+v", list)
	}

	// Pushed live to the subscriber.
	select {
	case n := <-ch:
		if n.UserID != "u1" || n.Title != "New job" {
			t.Fatalf("unexpected pushed notification: %+v", n)
		}
	case <-time.After(time.Second):
		t.Fatal("expected the notification to be pushed to the subscriber")
	}
}

func TestNotifyEmptyUserIsNoop(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, NewHub(nil))

	if err := svc.Notify(context.Background(), "", "t", "title", "body", ""); err != nil {
		t.Fatalf("notify with empty user should be a no-op, got %v", err)
	}
	if got := repo.seq; got != 0 {
		t.Fatalf("expected nothing persisted, got seq=%d", got)
	}
}

func TestNotifyWithoutHubDoesNotPanic(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	if err := svc.Notify(context.Background(), "u1", "t", "title", "body", ""); err != nil {
		t.Fatalf("notify without hub: %v", err)
	}
}

func TestNotifyPropagatesRepoError(t *testing.T) {
	repo := newFakeRepo()
	boom := errors.New("db down")
	repo.createErr = boom
	svc := NewService(repo, NewHub(nil))

	if err := svc.Notify(context.Background(), "u1", "t", "title", "body", ""); !errors.Is(err, boom) {
		t.Fatalf("expected repo error to propagate, got %v", err)
	}
}

func TestListOnlyReturnsOwnNotifications(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, NewHub(nil))
	ctx := context.Background()

	_ = svc.Notify(ctx, "u1", "t", "mine", "", "")
	_ = svc.Notify(ctx, "u2", "t", "theirs", "", "")

	list, err := svc.List(ctx, "u1", 100, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].Title != "mine" {
		t.Fatalf("expected only u1's notification, got %+v", list)
	}
}

func TestMarkReadUpdatesUnreadCount(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, NewHub(nil))
	ctx := context.Background()

	_ = svc.Notify(ctx, "u1", "t", "one", "", "")
	_ = svc.Notify(ctx, "u1", "t", "two", "", "")

	count, err := svc.UnreadCount(ctx, "u1")
	if err != nil {
		t.Fatalf("unread count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 unread, got %d", count)
	}

	list, _ := svc.List(ctx, "u1", 100, 0)
	if err := svc.MarkRead(ctx, "u1", list[0].ID); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	if count, _ := svc.UnreadCount(ctx, "u1"); count != 1 {
		t.Fatalf("expected 1 unread after MarkRead, got %d", count)
	}
}

func TestMarkAllReadClearsUnread(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, NewHub(nil))
	ctx := context.Background()

	_ = svc.Notify(ctx, "u1", "t", "one", "", "")
	_ = svc.Notify(ctx, "u1", "t", "two", "", "")
	_ = svc.Notify(ctx, "u2", "t", "other", "", "")

	if err := svc.MarkAllRead(ctx, "u1"); err != nil {
		t.Fatalf("mark all read: %v", err)
	}
	if count, _ := svc.UnreadCount(ctx, "u1"); count != 0 {
		t.Fatalf("expected 0 unread for u1, got %d", count)
	}
	// u2 is untouched.
	if count, _ := svc.UnreadCount(ctx, "u2"); count != 1 {
		t.Fatalf("expected u2 unread untouched, got %d", count)
	}
}
