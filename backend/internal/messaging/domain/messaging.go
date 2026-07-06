// Package domain holds the Messaging bounded context.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("conversation not found")
	ErrNotParticipant = errors.New("not a participant in this conversation")
)

type Conversation struct {
	ID             string     `json:"id"`
	Type           string     `json:"type"` // direct | group
	Title          string     `json:"title,omitempty"`
	CreatedBy      string     `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastMessageAt  *time.Time `json:"last_message_at,omitempty"`
	ParticipantIDs []string   `json:"participant_ids"`

	// User-specific details populated in lists
	IsPinned           bool       `json:"is_pinned"`
	IsArchived         bool       `json:"is_archived"`
	LastReadAt         *time.Time `json:"last_read_at,omitempty"`
	UnreadCount        int        `json:"unread_count"`
	LastMessagePreview string     `json:"last_message_preview,omitempty"`
}

type Participant struct {
	ConversationID string     `json:"conversation_id"`
	UserID         string     `json:"user_id"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	Role           string     `json:"role"`
	IsArchived     bool       `json:"is_archived"`
	IsPinned       bool       `json:"is_pinned"`
}

type Message struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Content        string     `json:"content"`
	ContentType    string     `json:"content_type"` // text | image | file | system
	CreatedAt      time.Time  `json:"created_at"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type MessageStatus struct {
	MessageID       string    `json:"message_id"`
	UserID          string    `json:"user_id"`
	Status          string    `json:"status"` // sent | delivered | read
	StatusUpdatedAt time.Time `json:"status_updated_at"`
}

type ConnectionChecker interface {
	AreConnected(ctx context.Context, userA, userB string) (bool, error)
}

type Repository interface {
	CreateConversation(ctx context.Context, c *Conversation) error
	GetConversation(ctx context.Context, id string) (*Conversation, error)
	ListConversations(ctx context.Context, userID string) ([]Conversation, error)
	FindDirect(ctx context.Context, userA, userB string) (string, bool, error)
	IsParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	Participants(ctx context.Context, conversationID string) ([]string, error)
	GetParticipantDetail(ctx context.Context, conversationID, userID string) (*Participant, error)
	AddMessage(ctx context.Context, m *Message) error
	GetMessage(ctx context.Context, id string) (*Message, error)
	ListMessages(ctx context.Context, conversationID string, limit int) ([]Message, error)
	UpdateMessage(ctx context.Context, m *Message) error
	MarkRead(ctx context.Context, conversationID, userID string) error
	SetMessageStatus(ctx context.Context, status *MessageStatus) error
	GetMessageStatuses(ctx context.Context, messageID string) ([]MessageStatus, error)
	ArchiveConversation(ctx context.Context, conversationID, userID string, archive bool) error
	PinConversation(ctx context.Context, conversationID, userID string, pin bool) error
	GetUnreadCount(ctx context.Context, conversationID, userID string, lastReadAt *time.Time) (int, error)
	GetLastMessage(ctx context.Context, conversationID string) (*Message, error)
}
