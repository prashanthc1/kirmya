package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"workspace-app/internal/messaging/domain"
)

// in-memory repo
type fakeRepo struct {
	seq   int
	convs map[string]*domain.Conversation
	parts map[string][]string
	msgs  map[string][]domain.Message
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{convs: map[string]*domain.Conversation{}, parts: map[string][]string{}, msgs: map[string][]domain.Message{}}
}

func (r *fakeRepo) CreateConversation(_ context.Context, c *domain.Conversation) error {
	r.seq++
	c.ID = fmt.Sprintf("c-%d", r.seq)
	cp := *c
	r.convs[c.ID] = &cp
	r.parts[c.ID] = append([]string(nil), c.ParticipantIDs...)
	return nil
}

func (r *fakeRepo) GetConversation(_ context.Context, id string) (*domain.Conversation, error) {
	c, ok := r.convs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	cp.ParticipantIDs = r.parts[id]
	return &cp, nil
}

func (r *fakeRepo) ListConversations(_ context.Context, userID string) ([]domain.Conversation, error) {
	out := []domain.Conversation{}
	for id, ps := range r.parts {
		for _, p := range ps {
			if p == userID {
				out = append(out, *r.convs[id])
			}
		}
	}
	return out, nil
}

func (r *fakeRepo) FindDirect(_ context.Context, a, b string) (string, bool, error) {
	for id, ps := range r.parts {
		if len(ps) == 2 && r.convs[id].Type == "direct" && contains(ps, a) && contains(ps, b) {
			return id, true, nil
		}
	}
	return "", false, nil
}

func (r *fakeRepo) IsParticipant(_ context.Context, id, userID string) (bool, error) {
	return contains(r.parts[id], userID), nil
}

func (r *fakeRepo) Participants(_ context.Context, id string) ([]string, error) {
	return r.parts[id], nil
}

func (r *fakeRepo) GetParticipantDetail(_ context.Context, conversationID, userID string) (*domain.Participant, error) {
	return &domain.Participant{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           "member",
	}, nil
}

func (r *fakeRepo) AddMessage(_ context.Context, m *domain.Message) error {
	r.seq++
	m.ID = fmt.Sprintf("m-%d", r.seq)
	r.msgs[m.ConversationID] = append(r.msgs[m.ConversationID], *m)
	return nil
}

func (r *fakeRepo) GetMessage(_ context.Context, id string) (*domain.Message, error) {
	for _, ms := range r.msgs {
		for _, m := range ms {
			if m.ID == id {
				return &m, nil
			}
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakeRepo) ListMessages(_ context.Context, id string, _ int) ([]domain.Message, error) {
	return r.msgs[id], nil
}

func (r *fakeRepo) UpdateMessage(_ context.Context, m *domain.Message) error {
	for cid, ms := range r.msgs {
		for i, existing := range ms {
			if existing.ID == m.ID {
				r.msgs[cid][i] = *m
				return nil
			}
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepo) MarkRead(_ context.Context, _, _ string) error { return nil }

func (r *fakeRepo) SetMessageStatus(_ context.Context, _ *domain.MessageStatus) error { return nil }

func (r *fakeRepo) GetMessageStatuses(_ context.Context, _ string) ([]domain.MessageStatus, error) {
	return nil, nil
}

func (r *fakeRepo) ArchiveConversation(_ context.Context, _, _ string, _ bool) error { return nil }

func (r *fakeRepo) PinConversation(_ context.Context, _, _ string, _ bool) error { return nil }

func (r *fakeRepo) GetUnreadCount(_ context.Context, _, _ string, _ *time.Time) (int, error) {
	return 0, nil
}

func (r *fakeRepo) GetLastMessage(_ context.Context, _ string) (*domain.Message, error) {
	return nil, nil
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

type recordingPublisher struct{ lastPayload map[string]any }

func (p *recordingPublisher) Publish(_ context.Context, _, _ string, payload map[string]any) error {
	p.lastPayload = payload
	return nil
}

type fakeConnectionChecker struct{}

func (c *fakeConnectionChecker) AreConnected(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

func TestStartReusesDirectConversation(t *testing.T) {
	svc := NewService(newFakeRepo(), nil, NewHub(nil), &fakeConnectionChecker{})
	ctx := context.Background()
	c1, err := svc.Start(ctx, "alice", []string{"bob"}, "")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	c1.Type = "direct"
	svc.repo.(*fakeRepo).convs[c1.ID].Type = "direct"

	c2, err := svc.Start(ctx, "bob", []string{"alice"}, "")
	if err != nil {
		t.Fatalf("start2: %v", err)
	}
	if c1.ID != c2.ID {
		t.Fatalf("expected the same direct conversation, got %s vs %s", c1.ID, c2.ID)
	}
}

func TestSendRequiresParticipantAndPublishes(t *testing.T) {
	pub := &recordingPublisher{}
	svc := NewService(newFakeRepo(), pub, NewHub(nil), &fakeConnectionChecker{})
	ctx := context.Background()
	c, _ := svc.Start(ctx, "alice", []string{"bob"}, "")
	c.Type = "direct"
	svc.repo.(*fakeRepo).convs[c.ID].Type = "direct"

	if _, err := svc.Send(ctx, "intruder", c.ID, "hi", "text", nil, ""); !errors.Is(err, domain.ErrNotParticipant) {
		t.Fatalf("expected ErrNotParticipant, got %v", err)
	}
	if _, err := svc.Send(ctx, "alice", c.ID, "hello bob", "text", nil, ""); err != nil {
		t.Fatalf("send: %v", err)
	}
	recips, _ := pub.lastPayload["recipient_ids"].([]string)
	if !contains(recips, "bob") {
		t.Fatalf("expected bob in recipients, got %v", recips)
	}
}
