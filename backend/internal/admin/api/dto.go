package api

import (
	"time"

	"workspace-app/internal/admin/domain"
)

type userResponse struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	FullName      string   `json:"full_name"`
	Headline      string   `json:"headline,omitempty"`
	Status        string   `json:"status"`
	EmailVerified bool     `json:"email_verified"`
	MFAEnabled    bool     `json:"mfa_enabled"`
	Roles         []string `json:"roles"`
	LastLoginAt   string   `json:"last_login_at,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

type reportResponse struct {
	ID          string `json:"id"`
	ReporterID  string `json:"reporter_id"`
	TargetType  string `json:"target_type"`
	TargetID    string `json:"target_id"`
	Reason      string `json:"reason"`
	Status      string `json:"status"`
	ActionTaken string `json:"action_taken,omitempty"`
	ResolvedBy  string `json:"resolved_by,omitempty"`
	ResolvedAt  string `json:"resolved_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type statusRequest struct {
	Status string `json:"status"`
}

type roleRequest struct {
	Role string `json:"role"`
}

type fileReportRequest struct {
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Reason     string `json:"reason"`
}

type resolveReportRequest struct {
	Status      string `json:"status"`
	ActionTaken string `json:"action_taken"`
}

func toUser(u *domain.UserSummary) userResponse {
	resp := userResponse{
		ID: u.ID, Email: u.Email, FullName: u.FullName, Headline: u.Headline,
		Status: u.Status, EmailVerified: u.EmailVerified, MFAEnabled: u.MFAEnabled,
		Roles: u.Roles, CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	}
	if u.LastLoginAt != nil {
		resp.LastLoginAt = u.LastLoginAt.UTC().Format(time.RFC3339)
	}
	return resp
}

func toUsers(users []domain.UserSummary) []userResponse {
	out := make([]userResponse, 0, len(users))
	for i := range users {
		out = append(out, toUser(&users[i]))
	}
	return out
}

func toReport(r *domain.Report) reportResponse {
	resp := reportResponse{
		ID: r.ID, ReporterID: r.ReporterID, TargetType: r.TargetType, TargetID: r.TargetID,
		Reason: r.Reason, Status: r.Status, ActionTaken: r.ActionTaken, ResolvedBy: r.ResolvedBy,
		CreatedAt: r.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: r.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if r.ResolvedAt != nil {
		resp.ResolvedAt = r.ResolvedAt.UTC().Format(time.RFC3339)
	}
	return resp
}

func toReports(reports []domain.Report) []reportResponse {
	out := make([]reportResponse, 0, len(reports))
	for i := range reports {
		out = append(out, toReport(&reports[i]))
	}
	return out
}
