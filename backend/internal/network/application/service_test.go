package application

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"workspace-app/internal/network/domain"
)

type fakeNetworkRepo struct {
	mu          sync.Mutex
	seq         int
	connections map[string]*domain.Connection
}

func newFakeNetworkRepo() *fakeNetworkRepo {
	return &fakeNetworkRepo{connections: make(map[string]*domain.Connection)}
}

func (r *fakeNetworkRepo) Create(ctx context.Context, requesterID, receiverID string, origin domain.ConnectionOrigin) (*domain.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.connections {
		if (c.RequesterID == requesterID && c.ReceiverID == receiverID) || (c.RequesterID == receiverID && c.ReceiverID == requesterID) {
			return nil, domain.ErrDuplicateRequest
		}
	}

	r.seq++
	id := fmt.Sprintf("conn-%d", r.seq)
	c := &domain.Connection{
		ID:          id,
		RequesterID: requesterID,
		ReceiverID:  receiverID,
		Status:      domain.StatusPending,
		Origin:      origin,
		CreatedAt:   "2026-07-06T00:00:00Z",
		UpdatedAt:   "2026-07-06T00:00:00Z",
	}
	r.connections[id] = c
	return c, nil
}

func (r *fakeNetworkRepo) CreateAccepted(ctx context.Context, requesterID, receiverID string, origin domain.ConnectionOrigin) (*domain.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	id := fmt.Sprintf("conn-%d", r.seq)
	c := &domain.Connection{
		ID:          id,
		RequesterID: requesterID,
		ReceiverID:  receiverID,
		Status:      domain.StatusAccepted,
		Origin:      origin,
		RespondedAt: "2026-07-06T00:00:00Z",
		CreatedAt:   "2026-07-06T00:00:00Z",
		UpdatedAt:   "2026-07-06T00:00:00Z",
	}
	r.connections[id] = c
	return c, nil
}

func (r *fakeNetworkRepo) UpdateStatus(ctx context.Context, connectionID string, status domain.ConnectionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.connections[connectionID]
	if !ok {
		return domain.ErrNotFound
	}
	c.Status = status
	c.RespondedAt = "2026-07-06T00:00:00Z"
	c.UpdatedAt = "2026-07-06T00:00:00Z"
	return nil
}

func (r *fakeNetworkRepo) Block(ctx context.Context, blockerID, blockedID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.connections {
		if (c.RequesterID == blockerID && c.ReceiverID == blockedID) || (c.RequesterID == blockedID && c.ReceiverID == blockerID) {
			c.Status = domain.StatusBlocked
			c.RespondedAt = "2026-07-06T00:00:00Z"
			c.UpdatedAt = "2026-07-06T00:00:00Z"
			return nil
		}
	}

	r.seq++
	id := fmt.Sprintf("conn-%d", r.seq)
	c := &domain.Connection{
		ID:          id,
		RequesterID: blockerID,
		ReceiverID:  blockedID,
		Status:      domain.StatusBlocked,
		Origin:      domain.OriginManualRequest,
		RespondedAt: "2026-07-06T00:00:00Z",
		CreatedAt:   "2026-07-06T00:00:00Z",
		UpdatedAt:   "2026-07-06T00:00:00Z",
	}
	r.connections[id] = c
	return nil
}

func (r *fakeNetworkRepo) Unconnect(ctx context.Context, userA, userB string) error {
	return r.Delete(ctx, userA, userB)
}

func (r *fakeNetworkRepo) GetConnections(ctx context.Context, userID string) ([]domain.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []domain.Connection
	for _, c := range r.connections {
		if (c.RequesterID == userID || c.ReceiverID == userID) && c.Status == domain.StatusAccepted {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (r *fakeNetworkRepo) GetIncomingRequests(ctx context.Context, userID string) ([]domain.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []domain.Connection
	for _, c := range r.connections {
		if c.ReceiverID == userID && c.Status == domain.StatusPending {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (r *fakeNetworkRepo) GetConnectionStatus(ctx context.Context, userA, userB string) (domain.ConnectionStatus, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.connections {
		if (c.RequesterID == userA && c.ReceiverID == userB) || (c.RequesterID == userB && c.ReceiverID == userA) {
			return c.Status, c.RequesterID, nil
		}
	}
	return "", "", nil
}

func (r *fakeNetworkRepo) GetByID(ctx context.Context, id string) (*domain.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.connections[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (r *fakeNetworkRepo) Delete(ctx context.Context, requesterID, receiverID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, c := range r.connections {
		if (c.RequesterID == requesterID && c.ReceiverID == receiverID) || (c.RequesterID == receiverID && c.ReceiverID == requesterID) {
			delete(r.connections, id)
		}
	}
	return nil
}

func TestNetworkConnectionFlow(t *testing.T) {
	repo := newFakeNetworkRepo()
	svc := NewService(repo)
	ctx := context.Background()

	userA := "user-a"
	userB := "user-b"

	// 1. Cannot connect to self
	_, err := svc.SendRequest(ctx, userA, userA, domain.OriginManualRequest)
	if !errors.Is(err, domain.ErrSelfConnection) {
		t.Fatalf("expected ErrSelfConnection, got %v", err)
	}

	// 2. Send request from A to B
	c, err := svc.SendRequest(ctx, userA, userB, domain.OriginManualRequest)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	if c.Status != domain.StatusPending || c.RequesterID != userA || c.ReceiverID != userB {
		t.Fatalf("invalid connection state: %+v", c)
	}

	// 3. Duplicate request should fail
	_, err = svc.SendRequest(ctx, userA, userB, domain.OriginManualRequest)
	if !errors.Is(err, domain.ErrDuplicateRequest) {
		t.Fatalf("expected ErrDuplicateRequest, got %v", err)
	}

	// 4. Reverse duplicate request should also fail
	_, err = svc.SendRequest(ctx, userB, userA, domain.OriginManualRequest)
	if !errors.Is(err, domain.ErrDuplicateRequest) {
		t.Fatalf("expected ErrDuplicateRequest, got %v", err)
	}

	// 5. Incoming requests for B should have 1 item
	reqs, err := svc.GetIncomingRequests(ctx, userB)
	if err != nil {
		t.Fatalf("failed to get incoming requests: %v", err)
	}
	if len(reqs) != 1 || reqs[0].ID != c.ID {
		t.Fatalf("expected 1 incoming request for B, got: %+v", reqs)
	}

	// 6. Incoming requests for A should have 0 items
	reqsA, err := svc.GetIncomingRequests(ctx, userA)
	if err != nil {
		t.Fatalf("failed to get incoming requests for A: %v", err)
	}
	if len(reqsA) != 0 {
		t.Fatalf("expected 0 incoming requests for A, got: %+v", reqsA)
	}

	// 7. Accept request from A (wrong receiver should fail)
	err = svc.AcceptRequest(ctx, userA, c.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for wrong receiver accepting, got: %v", err)
	}

	// 8. Accept request from B (correct receiver)
	err = svc.AcceptRequest(ctx, userB, c.ID)
	if err != nil {
		t.Fatalf("failed to accept request: %v", err)
	}

	// 9. Status should be accepted
	status, reqID, err := svc.GetConnectionStatus(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to get connection status: %v", err)
	}
	if status != domain.StatusAccepted || reqID != userA {
		t.Fatalf("expected accepted status with requester A, got %s / %s", status, reqID)
	}

	// 10. List connections for A should have 1 item
	conns, err := svc.GetConnections(ctx, userA)
	if err != nil {
		t.Fatalf("failed to get connections: %v", err)
	}
	if len(conns) != 1 || conns[0].ID != c.ID {
		t.Fatalf("expected 1 connection for A, got: %+v", conns)
	}

	// 11. Block user B
	err = svc.BlockUser(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to block user: %v", err)
	}
	status, _, _ = svc.GetConnectionStatus(ctx, userA, userB)
	if status != domain.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", status)
	}

	// 12. Unconnect user A and B (removes blocked/existing connection)
	err = svc.Unconnect(ctx, userA, userB)
	if err != nil {
		t.Fatalf("failed to unconnect: %v", err)
	}
	status, _, _ = svc.GetConnectionStatus(ctx, userA, userB)
	if status != "" {
		t.Fatalf("expected empty status after unconnecting, got %s", status)
	}
}

func TestAutoCreateConnection(t *testing.T) {
	repo := newFakeNetworkRepo()
	svc := NewService(repo)
	ctx := context.Background()

	userA := "user-a"
	userB := "user-b"

	c, err := svc.AutoCreateConnection(ctx, userA, userB, domain.OriginMentorshipMatch)
	if err != nil {
		t.Fatalf("failed to auto-create connection: %v", err)
	}
	if c.Status != domain.StatusAccepted || c.Origin != domain.OriginMentorshipMatch {
		t.Fatalf("invalid auto connection state: %+v", c)
	}
}

func TestRejectRequest(t *testing.T) {
	repo := newFakeNetworkRepo()
	svc := NewService(repo)
	ctx := context.Background()

	userA := "user-a"
	userB := "user-b"

	c, _ := svc.SendRequest(ctx, userA, userB, domain.OriginManualRequest)

	err := svc.RejectRequest(ctx, userA, c.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for wrong receiver rejecting, got: %v", err)
	}

	err = svc.RejectRequest(ctx, userB, c.ID)
	if err != nil {
		t.Fatalf("failed to reject request: %v", err)
	}

	status, _, _ := svc.GetConnectionStatus(ctx, userA, userB)
	if status != domain.StatusDeclined {
		t.Fatalf("expected declined status, got %s", status)
	}
}
