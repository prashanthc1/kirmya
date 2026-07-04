//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/notifications/domain"
	"workspace-app/internal/testsupport"
)

// TestNotificationsReadTracking verifies create, per-user listing/isolation,
// unread counting, and the mark-read transitions against a real PostgreSQL.
func TestNotificationsReadTracking(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	alice := testsupport.InsertUser(t, db, "alice@cb.test", "Alice")
	bob := testsupport.InsertUser(t, db, "bob@cb.test", "Bob")

	// Create two notifications for Alice and one for Bob.
	n1 := &domain.Notification{UserID: alice, Type: "job_match", Title: "New match", Body: "A role fits you", Link: "/jobs/1"}
	n2 := &domain.Notification{UserID: alice, Type: "mentorship", Title: "Session confirmed"}
	nBob := &domain.Notification{UserID: bob, Type: "job_match", Title: "Bob's match"}
	for _, n := range []*domain.Notification{n1, n2, nBob} {
		if err := repo.Create(ctx, n); err != nil {
			t.Fatalf("create notification: %v", err)
		}
		if n.ID == "" || n.CreatedAt.IsZero() {
			t.Fatalf("expected id + created_at populated, got %+v", n)
		}
	}

	// Listing is scoped per user.
	aliceList, err := repo.ListByUser(ctx, alice, 50, 0)
	if err != nil {
		t.Fatalf("list alice: %v", err)
	}
	if len(aliceList) != 2 {
		t.Fatalf("expected 2 notifications for alice, got %d", len(aliceList))
	}

	// Both start unread.
	count, err := repo.UnreadCount(ctx, alice)
	if err != nil {
		t.Fatalf("unread count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 unread, got %d", count)
	}

	// Marking one read drops the unread count by one.
	if err := repo.MarkRead(ctx, alice, n1.ID); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	if count, _ = repo.UnreadCount(ctx, alice); count != 1 {
		t.Fatalf("expected 1 unread after MarkRead, got %d", count)
	}

	// A user cannot mark another user's notification read.
	if err := repo.MarkRead(ctx, bob, n2.ID); err != nil {
		t.Fatalf("cross-user mark read should be a no-op, not an error: %v", err)
	}
	if count, _ = repo.UnreadCount(ctx, alice); count != 1 {
		t.Fatalf("cross-user MarkRead must not affect alice; unread=%d", count)
	}

	// MarkAllRead clears the rest.
	if err := repo.MarkAllRead(ctx, alice); err != nil {
		t.Fatalf("mark all read: %v", err)
	}
	if count, _ = repo.UnreadCount(ctx, alice); count != 0 {
		t.Fatalf("expected 0 unread after MarkAllRead, got %d", count)
	}

	// Bob's notification is untouched throughout.
	if count, _ = repo.UnreadCount(ctx, bob); count != 1 {
		t.Fatalf("expected bob to still have 1 unread, got %d", count)
	}
}
