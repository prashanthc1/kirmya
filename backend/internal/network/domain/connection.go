package domain

import (
	"context"
	"errors"
)

var (
	ErrNotFound          = errors.New("connection not found")
	ErrDuplicateRequest  = errors.New("connection request already exists")
	ErrSelfConnection    = errors.New("cannot connect to yourself")
	ErrInvalidTransition = errors.New("invalid connection status transition")
)

type ConnectionStatus string

const (
	StatusPending  ConnectionStatus = "pending"
	StatusAccepted ConnectionStatus = "accepted"
	StatusDeclined ConnectionStatus = "declined"
	StatusBlocked  ConnectionStatus = "blocked"
)

type ConnectionOrigin string

const (
	OriginMentorshipMatch ConnectionOrigin = "mentorship_match"
	OriginReferralRequest ConnectionOrigin = "referral_request"
	OriginJobApplication  ConnectionOrigin = "job_application"
	OriginManualRequest   ConnectionOrigin = "manual_request"
)

type Connection struct {
	ID          string           `json:"id"`
	RequesterID string           `json:"requester_id"`
	ReceiverID  string           `json:"receiver_id"`
	Status      ConnectionStatus `json:"status"`
	Origin      ConnectionOrigin `json:"origin"`
	RespondedAt string           `json:"responded_at,omitempty"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`

	// Joined user details for UI convenience
	RequesterName     string `json:"requester_name,omitempty"`
	RequesterHeadline string `json:"requester_headline,omitempty"`
	RequesterPhotoURL string `json:"requester_photo_url,omitempty"`
	ReceiverName      string `json:"receiver_name,omitempty"`
	ReceiverHeadline  string `json:"receiver_headline,omitempty"`
	ReceiverPhotoURL  string `json:"receiver_photo_url,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, requesterID, receiverID string, origin ConnectionOrigin) (*Connection, error)
	CreateAccepted(ctx context.Context, requesterID, receiverID string, origin ConnectionOrigin) (*Connection, error)
	UpdateStatus(ctx context.Context, connectionID string, status ConnectionStatus) error
	Block(ctx context.Context, blockerID, blockedID string) error
	Unconnect(ctx context.Context, userA, userB string) error
	GetConnections(ctx context.Context, userID string) ([]Connection, error)
	GetIncomingRequests(ctx context.Context, userID string) ([]Connection, error)
	GetConnectionStatus(ctx context.Context, userA, userB string) (status ConnectionStatus, requesterID string, err error)
	GetByID(ctx context.Context, id string) (*Connection, error)
	Delete(ctx context.Context, requesterID, receiverID string) error
}
