package connections

import (
	"time"
)

type ConnectionStatus string

const (
	StatusPending  ConnectionStatus = "pending"
	StatusAccepted ConnectionStatus = "accepted"
	StatusDeclined ConnectionStatus = "declined"
	StatusBlocked  ConnectionStatus = "blocked"
)

type ConnectionSource string

const (
	SourceSearch      ConnectionSource = "search"
	SourceProfileView ConnectionSource = "profile_view"
	SourceSuggested   ConnectionSource = "suggested"
	SourceImport      ConnectionSource = "import"
)

// PublicProfileSummary is the other party's public profile summary
type PublicProfileSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Headline  string `json:"headline"`
	AvatarURL string `json:"avatar_url"`
}

// Connection represents a relationship between two users
type Connection struct {
	ID          string               `json:"id"`
	UserAID     string               `json:"user_a_id"`
	UserBID     string               `json:"user_b_id"`
	Status      ConnectionStatus     `json:"status"`
	RequestedBy string               `json:"requested_by"`
	CreatedAt   time.Time            `json:"created_at"`
	RespondedAt *time.Time           `json:"responded_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Note        *string              `json:"note,omitempty"`
	Source      *ConnectionSource    `json:"source,omitempty"`
	User        *PublicProfileSummary `json:"user,omitempty"`
}

// Suggestion represents a recommended connection
type Suggestion struct {
	User                  PublicProfileSummary `json:"user"`
	MutualConnectionCount int                  `json:"mutual_connection_count"`
	Reason                string               `json:"reason"`
}

// Block represents a user block
type Block struct {
	ID        string    `json:"id"`
	BlockerID string    `json:"blocker_id"`
	BlockedID string    `json:"blocked_id"`
	Reason    *string   `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ConnectionCounts tracks a user's connections metrics
type ConnectionCounts struct {
	UserID               string    `json:"user_id"`
	ConnectionCount      int       `json:"connection_count"`
	PendingIncomingCount int       `json:"pending_incoming_count"`
	PendingOutgoingCount int       `json:"pending_outgoing_count"`
	UpdatedAt            time.Time `json:"updated_at"`
}
