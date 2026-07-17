//go:build integration

package connections

import (
	"context"
	"fmt"
	"testing"

	"workspace-app/internal/testsupport"
)

// mockBus implements EventPublisher for testing
type mockBus struct {
	events []string
}

func (m *mockBus) Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	m.events = append(m.events, eventType)
	return nil
}

func TestNormalizePair(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)

	u1 := "11111111-1111-1111-1111-111111111111"
	u2 := "22222222-2222-2222-2222-222222222222"

	// Check that order is always smaller, larger
	a1, b1 := repo.NormalizePair(u1, u2)
	a2, b2 := repo.NormalizePair(u2, u1)

	if a1 != u1 || b1 != u2 {
		t.Fatalf("expected order %s, %s; got %s, %s", u1, u2, a1, b1)
	}
	if a2 != u1 || b2 != u2 {
		t.Fatalf("expected order %s, %s; got %s, %s", u1, u2, a2, b2)
	}
}

func TestConnectionsLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	bus := &mockBus{}
	svc := NewService(db, repo, bus)
	ctx := context.Background()

	// Create test users
	userA := testsupport.InsertUser(t, db, "usera@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "userb@kirmya.test", "User B")

	// 1. SendRequest rejects self-request
	err := svc.SendRequest(ctx, userA, userA, nil, nil)
	if err != ErrSelfRequest {
		t.Fatalf("expected ErrSelfRequest, got %v", err)
	}

	// 2. SendRequest blocks duplicate pending requests
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != nil {
		t.Fatalf("failed to send first request: %v", err)
	}
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != ErrAlreadyPending {
		t.Fatalf("expected ErrAlreadyPending, got %v", err)
	}

	// 3. AcceptRequest rejects if responder == requester (sender cannot accept their own request)
	// First let's find the connection ID
	conn, err := repo.GetConnection(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	err = svc.AcceptRequest(ctx, conn.ID, userA)
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden for sender accepting request, got %v", err)
	}

	// 4. AcceptRequest works correctly when accepted by receiver
	err = svc.AcceptRequest(ctx, conn.ID, userB)
	if err != nil {
		t.Fatalf("failed to accept connection request: %v", err)
	}
	
	// Verify they can message
	ok, err := CanMessage(ctx, db, userA, userB)
	if err != nil || !ok {
		t.Fatalf("expected CanMessage to be true after accept, got %t (err: %v)", ok, err)
	}

	// 5. SendRequest rejects duplicate when already connected
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != ErrAlreadyConnected {
		t.Fatalf("expected ErrAlreadyConnected, got %v", err)
	}

	// 6. BlockUser cancels existing accepted connection and sets CanMessage() = false
	err = svc.BlockUser(ctx, userA, userB, "spamming")
	if err != nil {
		t.Fatalf("failed to block user: %v", err)
	}
	ok, err = CanMessage(ctx, db, userA, userB)
	if err != nil || ok {
		t.Fatalf("expected CanMessage to be false after block, got %t (err: %v)", ok, err)
	}

	// 7. SendRequest rejects when blocked in either direction
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != ErrBlocked {
		t.Fatalf("expected ErrBlocked for requester, got %v", err)
	}
	err = svc.SendRequest(ctx, userB, userA, nil, nil)
	if err != ErrBlocked {
		t.Fatalf("expected ErrBlocked for blockee, got %v", err)
	}

	// 8. UnblockUser allows clean state (CanMessage is still false, but re-request is allowed)
	err = svc.UnblockUser(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to unblock user: %v", err)
	}
	ok, err = CanMessage(ctx, db, userA, userB)
	if err != nil || ok {
		t.Fatalf("expected CanMessage to stay false after unblock, got %t (err: %v)", ok, err)
	}

	// Send request again should now succeed
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != nil {
		t.Fatalf("expected send request to succeed after unblock, got %v", err)
	}
}

func TestCooldownAfterDecline(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "usera@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "userb@kirmya.test", "User B")

	// Request -> Decline
	err := svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != nil {
		t.Fatalf("failed request: %v", err)
	}
	conn, err := repo.GetConnection(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed connection get: %v", err)
	}
	err = svc.DeclineRequest(ctx, conn.ID, userB)
	if err != nil {
		t.Fatalf("failed decline: %v", err)
	}

	// Try re-requesting immediately should fail with ErrCooldown
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != ErrCooldown {
		t.Fatalf("expected ErrCooldown, got %v", err)
	}

	// Artificially change responded_at to 31 days ago in database to bypass cooldown
	_, err = db.Exec("UPDATE connections SET responded_at = now() - interval '31 days' WHERE id = $1", conn.ID)
	if err != nil {
		t.Fatalf("failed to backdate responded_at: %v", err)
	}

	// Re-request should now succeed!
	err = svc.SendRequest(ctx, userA, userB, nil, nil)
	if err != nil {
		t.Fatalf("expected request to succeed after 31 days, got %v", err)
	}
}

func TestRateLimiting(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)
	// Force nil redis to use local memory limiter for testing
	svc.SetRedisClient(nil)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "usera@kirmya.test", "User A")

	// Trigger 20 requests (rate limit local check should pass)
	for i := 0; i < 20; i++ {
		// Mock unique target users so it doesn't fail on duplicate checks
		targetUser := testsupport.InsertUser(t, db, fmt.Sprintf("target%d@kirmya.test", i), fmt.Sprintf("Target %d", i))
		err := svc.SendRequest(ctx, userA, targetUser, nil, nil)
		if err != nil {
			t.Fatalf("failed request at %d: %v", i, err)
		}
	}

	// 21st request should be rejected with ErrRateLimited
	target21 := testsupport.InsertUser(t, db, "target21@kirmya.test", "Target 21")
	err := svc.SendRequest(ctx, userA, target21, nil, nil)
	if err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited on 21st request, got %v", err)
	}
}

func TestGetMutualConnections(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)
	ctx := context.Background()

	// Insert users
	userA := testsupport.InsertUser(t, db, "a@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "b@kirmya.test", "User B")
	userM1 := testsupport.InsertUser(t, db, "m1@kirmya.test", "Mutual 1")
	userM2 := testsupport.InsertUser(t, db, "m2@kirmya.test", "Mutual 2")

	// Connect A <-> M1
	_ = svc.SendRequest(ctx, userA, userM1, nil, nil)
	c, _ := repo.GetConnection(ctx, userA, userM1)
	_ = svc.AcceptRequest(ctx, c.ID, userM1)

	// Connect B <-> M1
	_ = svc.SendRequest(ctx, userB, userM1, nil, nil)
	c, _ = repo.GetConnection(ctx, userB, userM1)
	_ = svc.AcceptRequest(ctx, c.ID, userM1)

	// Connect A <-> M2
	_ = svc.SendRequest(ctx, userA, userM2, nil, nil)
	c, _ = repo.GetConnection(ctx, userA, userM2)
	_ = svc.AcceptRequest(ctx, c.ID, userM2)

	// Connect B <-> M2
	_ = svc.SendRequest(ctx, userB, userM2, nil, nil)
	c, _ = repo.GetConnection(ctx, userB, userM2)
	_ = svc.AcceptRequest(ctx, c.ID, userM2)

	// Get Mutual
	mutuals, total, err := repo.GetMutualConnections(ctx, userA, userB, 10)
	if err != nil {
		t.Fatalf("failed to get mutual connections: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 mutual connections, got %d", total)
	}
	if len(mutuals) != 2 {
		t.Fatalf("expected 2 hydrated mutual connections, got %d", len(mutuals))
	}
}

func TestGetSuggestions(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "a@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "b@kirmya.test", "User B")
	userC := testsupport.InsertUser(t, db, "c@kirmya.test", "User C") // blocked user
	userD := testsupport.InsertUser(t, db, "d@kirmya.test", "User D") // pending user

	// 1. Block user C
	_ = svc.BlockUser(ctx, userA, userC, "spam")

	// 2. Send pending request to user D
	_ = svc.SendRequest(ctx, userA, userD, nil, nil)

	// Get Suggestions for user A
	sugs, err := repo.GetSuggestions(ctx, userA, 10)
	if err != nil {
		t.Fatalf("failed to get suggestions: %v", err)
	}

	// Should suggest user B, but NOT user C (blocked), D (pending), or user A itself (self)
	foundB := false
	for _, s := range sugs {
		if s.User.ID == userB {
			foundB = true
		}
		if s.User.ID == userC || s.User.ID == userD || s.User.ID == userA {
			t.Fatalf("unexpected user %s suggested", s.User.ID)
		}
	}

	if !foundB {
		t.Fatalf("expected user B to be suggested")
	}
}

func TestConnectionCountsReconciliation(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "usera@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "userb@kirmya.test", "User B")

	// Helper to get counts
	getCounts := func(userID string) ConnectionCounts {
		var cc ConnectionCounts
		err := db.QueryRowContext(ctx, "SELECT user_id, connection_count, pending_incoming_count, pending_outgoing_count FROM connection_counts WHERE user_id = $1", userID).
			Scan(&cc.UserID, &cc.ConnectionCount, &cc.PendingIncomingCount, &cc.PendingOutgoingCount)
		if err != nil {
			return ConnectionCounts{}
		}
		return cc
	}

	// 1. Send connection request A -> B
	_ = svc.SendRequest(ctx, userA, userB, nil, nil)

	ccA := getCounts(userA)
	if ccA.PendingOutgoingCount != 1 || ccA.ConnectionCount != 0 {
		t.Fatalf("unexpected A counts after request: %+v", ccA)
	}
	ccB := getCounts(userB)
	if ccB.PendingIncomingCount != 1 || ccB.ConnectionCount != 0 {
		t.Fatalf("unexpected B counts after request: %+v", ccB)
	}

	// 2. Accept connection
	conn, _ := repo.GetConnection(ctx, userA, userB)
	_ = svc.AcceptRequest(ctx, conn.ID, userB)

	ccA = getCounts(userA)
	if ccA.PendingOutgoingCount != 0 || ccA.ConnectionCount != 1 {
		t.Fatalf("unexpected A counts after accept: %+v", ccA)
	}
	ccB = getCounts(userB)
	if ccB.PendingIncomingCount != 0 || ccB.ConnectionCount != 1 {
		t.Fatalf("unexpected B counts after accept: %+v", ccB)
	}

	// 3. Remove connection
	_ = svc.RemoveConnection(ctx, conn.ID, userA)

	ccA = getCounts(userA)
	if ccA.ConnectionCount != 0 {
		t.Fatalf("unexpected A counts after remove: %+v", ccA)
	}
	ccB = getCounts(userB)
	if ccB.ConnectionCount != 0 {
		t.Fatalf("unexpected B counts after remove: %+v", ccB)
	}
}
