// Package application implements Messaging use cases.
package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"workspace-app/internal/messaging/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

// Bus is the platform event bus: event publishing plus the cross-instance SSE
// fanout transport. The composition root passes one bus that satisfies both.
type Bus interface {
	EventPublisher
	Broadcaster
}

const eventMessageSent = "MessageSent"

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

// MessagePolicyReader reports a user's "who can message me" policy. A nil reader
// (the default) disables the check.
type MessagePolicyReader interface {
	MessagePolicy(ctx context.Context, userID string) (string, error)
}

type Service struct {
	repo   domain.Repository
	events EventPublisher
	hub    *Hub
	policy MessagePolicyReader
}

func NewService(repo domain.Repository, events EventPublisher, hub *Hub) *Service {
	return &Service{repo: repo, events: events, hub: hub}
}

// SetMessagePolicyReader injects the reader used to honour recipients' message policy.
func (s *Service) SetMessagePolicyReader(r MessagePolicyReader) { s.policy = r }

// Subscribe registers a real-time subscriber for the user's conversation events.
func (s *Service) Subscribe(userID string) (<-chan StreamEvent, func()) {
	return s.hub.Subscribe(userID)
}

// Start creates a conversation (reusing an existing 1:1 if it already exists).
func (s *Service) Start(ctx context.Context, creatorID string, participantIDs []string, title string) (*domain.Conversation, error) {
	set := map[string]bool{creatorID: true}
	for _, id := range participantIDs {
		if id != "" {
			set[id] = true
		}
	}
	if len(set) < 2 {
		return nil, ValidationError{"a conversation needs at least one other participant"}
	}
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}

	if len(ids) == 2 {
		other := ids[0]
		if other == creatorID {
			other = ids[1]
		}
		if existing, found, err := s.repo.FindDirect(ctx, creatorID, other); err != nil {
			return nil, err
		} else if found {
			return s.repo.GetConversation(ctx, existing)
		}
	}

	// Honour each recipient's "who can message me" policy before opening a new
	// conversation. ("network" needs a connections graph that does not exist yet,
	// so only an explicit "none" is enforced here.)
	if s.policy != nil {
		for _, id := range ids {
			if id == creatorID {
				continue
			}
			if pol, err := s.policy.MessagePolicy(ctx, id); err == nil && pol == "none" {
				return nil, ValidationError{"this person is not accepting new messages"}
			}
		}
	}

	c := &domain.Conversation{
		IsGroup: len(ids) > 2, Title: strings.TrimSpace(title), CreatedBy: creatorID, ParticipantIDs: ids,
	}
	if err := s.repo.CreateConversation(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	return s.repo.ListConversations(ctx, userID)
}

func (s *Service) ListMessages(ctx context.Context, userID, conversationID string) ([]domain.Message, error) {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return nil, err
	}
	_ = s.repo.MarkRead(ctx, conversationID, userID)
	s.notifyRead(ctx, conversationID, userID)
	return s.repo.ListMessages(ctx, conversationID, 200)
}

// Typing pushes an ephemeral "typing" indicator to the other participants.
func (s *Service) Typing(ctx context.Context, userID, conversationID string) error {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return err
	}
	s.broadcast(ctx, conversationID, userID, StreamEvent{
		Kind: EventTyping, ConversationID: conversationID, ActorID: userID,
	})
	return nil
}

// notifyRead tells the other participants that readerID has read the conversation.
func (s *Service) notifyRead(ctx context.Context, conversationID, readerID string) {
	s.broadcast(ctx, conversationID, readerID, StreamEvent{
		Kind: EventRead, ConversationID: conversationID, ActorID: readerID, At: time.Now().UTC(),
	})
}

// broadcast publishes an event to every participant except the actor.
func (s *Service) broadcast(ctx context.Context, conversationID, actorID string, ev StreamEvent) {
	if s.hub == nil {
		return
	}
	participants, _ := s.repo.Participants(ctx, conversationID)
	for _, uid := range participants {
		if uid != actorID {
			s.hub.Publish(uid, ev)
		}
	}
}

func (s *Service) Send(ctx context.Context, userID, conversationID, body string) (*domain.Message, error) {
	if strings.TrimSpace(body) == "" {
		return nil, ValidationError{"message body is required"}
	}
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return nil, err
	}
	m := &domain.Message{ConversationID: conversationID, SenderID: userID, Body: body}
	if err := s.repo.AddMessage(ctx, m); err != nil {
		return nil, err
	}

	participants, _ := s.repo.Participants(ctx, conversationID)
	if s.events != nil {
		_ = s.events.Publish(ctx, eventMessageSent, conversationID, map[string]any{
			"conversation_id": conversationID, "sender_id": userID, "recipient_ids": participants,
		})
	}
	// Push to every participant's live stream (sender included; clients dedupe by id).
	if s.hub != nil {
		for _, uid := range participants {
			s.hub.Publish(uid, StreamEvent{Kind: EventMessage, ConversationID: conversationID, ActorID: userID, Message: m})
		}
	}
	return m, nil
}

func (s *Service) MarkRead(ctx context.Context, userID, conversationID string) error {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return err
	}
	if err := s.repo.MarkRead(ctx, conversationID, userID); err != nil {
		return err
	}
	s.notifyRead(ctx, conversationID, userID)
	return nil
}

func (s *Service) requireParticipant(ctx context.Context, conversationID, userID string) error {
	ok, err := s.repo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	if !ok {
		return domain.ErrNotParticipant
	}
	return nil
}
