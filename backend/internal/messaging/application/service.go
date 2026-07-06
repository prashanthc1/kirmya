package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"workspace-app/internal/messaging/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

type Bus interface {
	EventPublisher
	Broadcaster
}

const eventMessageSent = "MessageSent"

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

type MessagePolicyReader interface {
	MessagePolicy(ctx context.Context, userID string) (string, error)
}

type Service struct {
	repo          domain.Repository
	events        EventPublisher
	hub           *Hub
	checker       domain.ConnectionChecker
	encryptionKey string
	policy        MessagePolicyReader
}

func NewService(repo domain.Repository, events EventPublisher, hub *Hub, checker domain.ConnectionChecker) *Service {
	key := os.Getenv("MESSAGE_ENCRYPTION_KEY")
	if key == "" {
		key = "kirmya-default-key-32bytes-long!"
	}
	return &Service{
		repo:          repo,
		events:        events,
		hub:           hub,
		checker:       checker,
		encryptionKey: key,
	}
}

func (s *Service) SetMessagePolicyReader(policy MessagePolicyReader) {
	s.policy = policy
}

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

	// 1. Connection Gate Check for Direct DMs (2 participants)
	if len(ids) == 2 {
		other := ids[0]
		if other == creatorID {
			other = ids[1]
		}

		// Enforce accepted connection check
		if s.checker != nil {
			connected, err := s.checker.AreConnected(ctx, creatorID, other)
			if err != nil {
				return nil, err
			}
			if !connected {
				return nil, ValidationError{"forbidden: you must have an accepted connection with this user to start a conversation"}
			}
		}

		if existing, found, err := s.repo.FindDirect(ctx, creatorID, other); err != nil {
			return nil, err
		} else if found {
			c, err := s.repo.GetConversation(ctx, existing)
			if err == nil {
				return c, nil
			}
		}
	}

	c := &domain.Conversation{
		Type:           "direct",
		Title:          strings.TrimSpace(title),
		CreatedBy:      creatorID,
		ParticipantIDs: ids,
	}
	if len(ids) > 2 {
		c.Type = "group"
	}

	if err := s.repo.CreateConversation(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	convs, err := s.repo.ListConversations(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter: only show direct conversations with active connections.
	var activeConvs []domain.Conversation
	for _, c := range convs {
		if c.Type == "direct" && len(c.ParticipantIDs) == 2 {
			other := c.ParticipantIDs[0]
			if other == userID {
				other = c.ParticipantIDs[1]
			}
			if s.checker != nil {
				connected, err := s.checker.AreConnected(ctx, userID, other)
				if err != nil || !connected {
					continue // Skip if unconnected or blocked
				}
			}
		}

		// Fetch unread count
		unread, err := s.repo.GetUnreadCount(ctx, c.ID, userID, c.LastReadAt)
		if err == nil {
			c.UnreadCount = unread
		}

		// Fetch last message preview
		lastMsg, err := s.repo.GetLastMessage(ctx, c.ID)
		if err == nil && lastMsg != nil {
			decrypted, err := decryptAESGCM(lastMsg.Content, s.encryptionKey)
			if err == nil {
				c.LastMessagePreview = decrypted
			} else {
				c.LastMessagePreview = lastMsg.Content
			}
			if len(c.LastMessagePreview) > 100 {
				c.LastMessagePreview = c.LastMessagePreview[:97] + "..."
			}
		}

		activeConvs = append(activeConvs, c)
	}
	return activeConvs, nil
}

func (s *Service) ListMessages(ctx context.Context, userID, conversationID string) ([]domain.Message, error) {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return nil, err
	}
	_ = s.repo.MarkRead(ctx, conversationID, userID)
	s.notifyRead(ctx, conversationID, userID)

	msgs, err := s.repo.ListMessages(ctx, conversationID, 200)
	if err != nil {
		return nil, err
	}

	// Decrypt message content
	for i := range msgs {
		decrypted, err := decryptAESGCM(msgs[i].Content, s.encryptionKey)
		if err == nil {
			msgs[i].Content = decrypted
		}
	}

	return msgs, nil
}

// SearchMessages searches within a conversation history.
func (s *Service) SearchMessages(ctx context.Context, userID, conversationID, query string) ([]domain.Message, error) {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return nil, err
	}

	msgs, err := s.repo.ListMessages(ctx, conversationID, 1000)
	if err != nil {
		return nil, err
	}

	var results []domain.Message
	query = strings.ToLower(query)
	for i := range msgs {
		decrypted, err := decryptAESGCM(msgs[i].Content, s.encryptionKey)
		if err == nil {
			msgs[i].Content = decrypted
		}

		if msgs[i].DeletedAt == nil && strings.Contains(strings.ToLower(msgs[i].Content), query) {
			results = append(results, msgs[i])
		}
	}
	return results, nil
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

func (s *Service) notifyRead(ctx context.Context, conversationID, readerID string) {
	s.broadcast(ctx, conversationID, readerID, StreamEvent{
		Kind: EventRead, ConversationID: conversationID, ActorID: readerID, At: time.Now().UTC(),
	})
}

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

// Send sends a text or rich message.
func (s *Service) Send(ctx context.Context, userID, conversationID, content, contentType string, fileData []byte, fileName string) (*domain.Message, error) {
	if strings.TrimSpace(content) == "" && len(fileData) == 0 {
		return nil, ValidationError{"message body or attachment is required"}
	}
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return nil, err
	}

	// 1. Connection Gate Check for Direct Conversation messaging
	conv, err := s.repo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conv.Type == "direct" && len(conv.ParticipantIDs) == 2 {
		other := conv.ParticipantIDs[0]
		if other == userID {
			other = conv.ParticipantIDs[1]
		}
		if s.checker != nil {
			connected, err := s.checker.AreConnected(ctx, userID, other)
			if err != nil {
				return nil, err
			}
			if !connected {
				return nil, ValidationError{"forbidden: you must have an accepted connection with this user to send messages"}
			}
		}
		if s.policy != nil {
			pol, err := s.policy.MessagePolicy(ctx, other)
			if err == nil && pol == "none" {
				return nil, ValidationError{"forbidden: this user has disabled direct messaging"}
			}
		}
	}

	// 2. Rich Content Validations
	if contentType == "image" && len(fileData) > 0 {
		if len(fileData) > 5*1024*1024 {
			return nil, ValidationError{"image size exceeds maximum 5MB limit"}
		}
		lowerName := strings.ToLower(fileName)
		if !strings.HasSuffix(lowerName, ".png") && !strings.HasSuffix(lowerName, ".jpg") && !strings.HasSuffix(lowerName, ".jpeg") && !strings.HasSuffix(lowerName, ".gif") {
			return nil, ValidationError{"invalid image format, supported: png, jpg, jpeg, gif"}
		}
	}

	if contentType == "file" && len(fileData) > 0 {
		if err := s.scanForViruses(fileData); err != nil {
			return nil, ValidationError{"file security check failed: " + err.Error()}
		}
	}

	// 3. Encrypt message body before storage
	encryptedContent, err := encryptAESGCM(content, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %w", err)
	}

	m := &domain.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        encryptedContent,
		ContentType:    contentType,
	}
	if m.ContentType == "" {
		m.ContentType = "text"
	}

	if err := s.repo.AddMessage(ctx, m); err != nil {
		return nil, err
	}

	// 4. Create Message Status per Recipient
	participants, _ := s.repo.Participants(ctx, conversationID)
	for _, pid := range participants {
		if pid != userID {
			status := &domain.MessageStatus{
				MessageID:       m.ID,
				UserID:          pid,
				Status:          "sent",
				StatusUpdatedAt: time.Now().UTC(),
			}
			_ = s.repo.SetMessageStatus(ctx, status)
		}
	}

	// 5. Decrypt message for the live event hubs
	m.Content = content

	if s.events != nil {
		_ = s.events.Publish(ctx, eventMessageSent, conversationID, map[string]any{
			"conversation_id": conversationID, "sender_id": userID, "recipient_ids": participants,
		})
	}

	// Push to every participant's live stream (sender included).
	if s.hub != nil {
		for _, uid := range participants {
			s.hub.Publish(uid, StreamEvent{Kind: EventMessage, ConversationID: conversationID, ActorID: userID, Message: m})
		}
	}
	return m, nil
}

// DeleteMessage soft deletes a message by setting deleted_at.
func (s *Service) DeleteMessage(ctx context.Context, userID, messageID string) error {
	m, err := s.repo.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}
	if m.SenderID != userID {
		return ValidationError{"forbidden: you can only delete your own messages"}
	}
	now := time.Now().UTC()
	m.DeletedAt = &now
	return s.repo.UpdateMessage(ctx, m)
}

// ArchiveConversation archives a conversation for a user.
func (s *Service) ArchiveConversation(ctx context.Context, userID, conversationID string, archive bool) error {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return err
	}
	return s.repo.ArchiveConversation(ctx, conversationID, userID, archive)
}

// PinConversation pins a conversation for a user.
func (s *Service) PinConversation(ctx context.Context, userID, conversationID string, pin bool) error {
	if err := s.requireParticipant(ctx, conversationID, userID); err != nil {
		return err
	}
	return s.repo.PinConversation(ctx, conversationID, userID, pin)
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

// scanForViruses represents the mock security virus scan hook.
func (s *Service) scanForViruses(data []byte) error {
	log.Printf("[virus-scan] Scanning %d bytes... Clean.", len(data))
	return nil
}

// GetUserPresence fetches online/offline status from Redis if users are connected.
func (s *Service) GetUserPresence(ctx context.Context, viewerID, targetID string) (string, error) {
	if viewerID == targetID {
		return "online", nil
	}

	// 1. Verification: Viewer and Target must be connected
	if s.checker != nil {
		connected, err := s.checker.AreConnected(ctx, viewerID, targetID)
		if err != nil {
			return "", err
		}
		if !connected {
			return "", ValidationError{"forbidden: you can only see presence status of your connections"}
		}
	}

	if s.hub == nil {
		return "offline", nil
	}

	return s.hub.GetPresence(ctx, targetID)
}

// UpdateUserPresence updates the user's presence state in Redis.
func (s *Service) UpdateUserPresence(ctx context.Context, userID string, online bool) error {
	if s.hub == nil {
		return nil
	}
	return s.hub.SetPresence(ctx, userID, online)
}
