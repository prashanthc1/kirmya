// Package domain holds the Notifications bounded context.
package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Title     string
	Body      string
	Link      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error)
	MarkRead(ctx context.Context, userID, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	UnreadCount(ctx context.Context, userID string) (int, error)
}
