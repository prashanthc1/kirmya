//go:build integration

package postgres

import (
	"context"
	"testing"

	"workspace-app/internal/network/domain"
	"workspace-app/internal/testsupport"
)

func TestConnectionLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "user_a@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "user_b@kirmya.test", "User B")

	// 1. Create connection request (A -> B)
	c, err := repo.Create(ctx, userA, userB, domain.OriginManualRequest)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	if c.RequesterID != userA || c.ReceiverID != userB || c.Status != domain.StatusPending {
		t.Fatalf("unexpected connection fields: %+v", c)
	}

	// 2. Try sending duplicate connection should fail
	_, err = repo.Create(ctx, userA, userB, domain.OriginManualRequest)
	if err == nil {
		t.Fatal("expected duplicate request error, got nil")
	}

	// 3. Get connection status
	status, reqID, err := repo.GetConnectionStatus(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to get connection status: %v", err)
	}
	if status != domain.StatusPending || reqID != userA {
		t.Fatalf("unexpected connection status: %s / %s", status, reqID)
	}

	// 4. Get incoming requests for B
	reqs, err := repo.GetIncomingRequests(ctx, userB)
	if err != nil {
		t.Fatalf("failed to get incoming requests: %v", err)
	}
	if len(reqs) != 1 || reqs[0].ID != c.ID {
		t.Fatalf("expected 1 incoming request for B, got: %+v", reqs)
	}

	// 5. Update connection status to accepted
	err = repo.UpdateStatus(ctx, c.ID, domain.StatusAccepted)
	if err != nil {
		t.Fatalf("failed to accept request: %v", err)
	}

	// 6. Get accepted connections for A
	conns, err := repo.GetConnections(ctx, userA)
	if err != nil {
		t.Fatalf("failed to get connections: %v", err)
	}
	if len(conns) != 1 || conns[0].ID != c.ID {
		t.Fatalf("expected 1 connection for A, got: %+v", conns)
	}

	// 7. Delete connection
	err = repo.Delete(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to delete connection: %v", err)
	}

	// 8. Connection status should now be unconnected
	status, _, err = repo.GetConnectionStatus(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to get connection status: %v", err)
	}
	if status != "" {
		t.Fatalf("expected unconnected status, got %s", status)
	}
}
