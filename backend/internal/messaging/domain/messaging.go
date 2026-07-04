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
	ID             string
	IsGroup        bool
	Title          string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ParticipantIDs []string
}

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Body           string
	CreatedAt      time.Time
}

type Repository interface {
	CreateConversation(ctx context.Context, c *Conversation) error
	GetConversation(ctx context.Context, id string) (*Conversation, error)
	ListConversations(ctx context.Context, userID string) ([]Conversation, error)
	FindDirect(ctx context.Context, userA, userB string) (string, bool, error)
	IsParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	Participants(ctx context.Context, conversationID string) ([]string, error)
	AddMessage(ctx context.Context, m *Message) error
	ListMessages(ctx context.Context, conversationID string, limit int) ([]Message, error)
	MarkRead(ctx context.Context, conversationID, userID string) error
}
