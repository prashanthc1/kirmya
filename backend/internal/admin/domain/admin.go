// Package domain holds the Admin bounded context entities and ports. The admin
// context is a read/operate layer over data owned by other contexts (identity
// users, community posts/comments) plus its own content-report queue; it has no
// dependency on any other layer — infrastructure implements its ports.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidRole   = errors.New("invalid role")
)

// Account status values (mirrors the identity users.status check constraint).
const (
	StatusActive      = "active"
	StatusSuspended   = "suspended"
	StatusDeactivated = "deactivated"
)

// ValidUserStatuses are the statuses an admin may set on a user.
var ValidUserStatuses = map[string]bool{
	StatusActive: true, StatusSuspended: true, StatusDeactivated: true,
}

// Assignable RBAC role names (mirrors the seeded roles table).
var ValidRoles = map[string]bool{
	"job_seeker": true, "referrer": true, "mentor": true, "recruiter": true, "admin": true,
}

// Report status values.
const (
	ReportOpen      = "open"
	ReportReviewing = "reviewing"
	ReportResolved  = "resolved"
	ReportDismissed = "dismissed"
)

// Reportable target types.
var ValidTargetTypes = map[string]bool{
	"post": true, "comment": true, "user": true, "message": true,
}

// terminalReportStatus reports whether a report has been triaged to completion.
func IsTerminalReportStatus(s string) bool {
	return s == ReportResolved || s == ReportDismissed
}

// UserSummary is the admin read projection of an account (joined with profile
// headline and role names).
type UserSummary struct {
	ID            string
	Email         string
	FullName      string
	Headline      string
	Status        string
	EmailVerified bool
	MFAEnabled    bool
	Roles         []string
	LastLoginAt   *time.Time
	CreatedAt     time.Time
}

// UserFilter narrows a user listing. Empty fields are ignored.
type UserFilter struct {
	Query  string // matches full_name/email (ILIKE)
	Status string
	Role   string
	Limit  int
	Offset int
}

// Report is the aggregate root of the moderation queue.
type Report struct {
	ID          string
	ReporterID  string
	TargetType  string
	TargetID    string
	Reason      string
	Status      string
	ActionTaken string
	ResolvedBy  string
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Analytics is the platform overview shown on the admin dashboard.
type Analytics struct {
	Users struct {
		Total       int `json:"total"`
		Active      int `json:"active"`
		Suspended   int `json:"suspended"`
		Deactivated int `json:"deactivated"`
		New7d       int `json:"new_7d"`
		New30d      int `json:"new_30d"`
	} `json:"users"`
	Jobs struct {
		Total        int `json:"total"`
		Applications int `json:"applications"`
	} `json:"jobs"`
	Referrals struct {
		Total    int `json:"total"`
		Accepted int `json:"accepted"`
		Hired    int `json:"hired"`
	} `json:"referrals"`
	Communities struct {
		Total    int `json:"total"`
		Posts    int `json:"posts"`
		Comments int `json:"comments"`
	} `json:"communities"`
	Reports struct {
		Open int `json:"open"`
	} `json:"reports"`
}

// Repository is the persistence port for the admin context.
type Repository interface {
	// Users.
	ListUsers(ctx context.Context, f UserFilter) ([]UserSummary, int, error)
	GetUser(ctx context.Context, id string) (*UserSummary, error)
	SetUserStatus(ctx context.Context, id, status string) error
	AssignRole(ctx context.Context, userID, role string) error
	RevokeRole(ctx context.Context, userID, role string) error

	// Content moderation (data owned by the community context; admin removes it).
	DeletePost(ctx context.Context, id string) error
	DeleteComment(ctx context.Context, id string) error

	// Report queue.
	CreateReport(ctx context.Context, r *Report) error
	ListReports(ctx context.Context, status string) ([]Report, error)
	GetReport(ctx context.Context, id string) (*Report, error)
	ResolveReport(ctx context.Context, id, status, actionTaken, resolvedBy string) error

	// Analytics + audit.
	Analytics(ctx context.Context) (*Analytics, error)
	WriteAudit(ctx context.Context, actorID, action, targetType, targetID string, metadata map[string]any) error
}
