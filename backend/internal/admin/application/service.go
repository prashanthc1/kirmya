// Package application implements the Admin use cases over the domain ports.
// Every mutating action is recorded in the audit log with the acting admin's id.
package application

import (
	"context"
	"strings"

	"workspace-app/internal/admin/domain"
)

// ValidationError is returned for invalid input (mapped to HTTP 400).
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

// ----- Users -------------------------------------------------------------

// ListUsers returns a filtered, paginated page of users plus the total count.
func (s *Service) ListUsers(ctx context.Context, f domain.UserFilter) ([]domain.UserSummary, int, error) {
	if f.Status != "" && !domain.ValidUserStatuses[f.Status] {
		return nil, 0, ValidationError{"invalid status filter"}
	}
	if f.Role != "" && !domain.ValidRoles[f.Role] {
		return nil, 0, ValidationError{"invalid role filter"}
	}
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 25
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return s.repo.ListUsers(ctx, f)
}

// GetUser returns one user's admin projection.
func (s *Service) GetUser(ctx context.Context, id string) (*domain.UserSummary, error) {
	return s.repo.GetUser(ctx, id)
}

// SetUserStatus suspends/reactivates/deactivates an account. An admin may not
// change their own status (guards against self-lockout).
func (s *Service) SetUserStatus(ctx context.Context, adminID, userID, status string) (*domain.UserSummary, error) {
	if !domain.ValidUserStatuses[status] {
		return nil, ValidationError{"status must be one of active, suspended, deactivated"}
	}
	if adminID == userID {
		return nil, ValidationError{"you cannot change your own account status"}
	}
	if _, err := s.repo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.repo.SetUserStatus(ctx, userID, status); err != nil {
		return nil, err
	}
	s.audit(ctx, adminID, "admin.user.status", "user", userID, map[string]any{"status": status})
	return s.repo.GetUser(ctx, userID)
}

// AssignRole grants an RBAC role to a user.
func (s *Service) AssignRole(ctx context.Context, adminID, userID, role string) (*domain.UserSummary, error) {
	if !domain.ValidRoles[role] {
		return nil, ValidationError{"unknown role"}
	}
	if _, err := s.repo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.repo.AssignRole(ctx, userID, role); err != nil {
		return nil, err
	}
	s.audit(ctx, adminID, "admin.user.role.grant", "user", userID, map[string]any{"role": role})
	return s.repo.GetUser(ctx, userID)
}

// RevokeRole removes an RBAC role from a user. An admin may not revoke their own
// admin role (guards against self-lockout).
func (s *Service) RevokeRole(ctx context.Context, adminID, userID, role string) (*domain.UserSummary, error) {
	if !domain.ValidRoles[role] {
		return nil, ValidationError{"unknown role"}
	}
	if adminID == userID && role == "admin" {
		return nil, ValidationError{"you cannot revoke your own admin role"}
	}
	if _, err := s.repo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.repo.RevokeRole(ctx, userID, role); err != nil {
		return nil, err
	}
	s.audit(ctx, adminID, "admin.user.role.revoke", "user", userID, map[string]any{"role": role})
	return s.repo.GetUser(ctx, userID)
}

// ----- Moderation --------------------------------------------------------

// RemovePost deletes a community post as a moderation action.
func (s *Service) RemovePost(ctx context.Context, adminID, postID string) error {
	if err := s.repo.DeletePost(ctx, postID); err != nil {
		return err
	}
	s.audit(ctx, adminID, "admin.post.remove", "post", postID, nil)
	return nil
}

// RemoveComment deletes a community comment as a moderation action.
func (s *Service) RemoveComment(ctx context.Context, adminID, commentID string) error {
	if err := s.repo.DeleteComment(ctx, commentID); err != nil {
		return err
	}
	s.audit(ctx, adminID, "admin.comment.remove", "comment", commentID, nil)
	return nil
}

// ----- Reports -----------------------------------------------------------

// ReportInput is the file-a-report command (open to any authenticated user).
type ReportInput struct {
	TargetType string
	TargetID   string
	Reason     string
}

// FileReport lets any authenticated user report content for moderation.
func (s *Service) FileReport(ctx context.Context, reporterID string, in ReportInput) (*domain.Report, error) {
	if !domain.ValidTargetTypes[in.TargetType] {
		return nil, ValidationError{"target_type must be one of post, comment, user, message"}
	}
	if strings.TrimSpace(in.TargetID) == "" {
		return nil, ValidationError{"target_id is required"}
	}
	if strings.TrimSpace(in.Reason) == "" {
		return nil, ValidationError{"a reason is required"}
	}
	r := &domain.Report{
		ReporterID: reporterID,
		TargetType: in.TargetType,
		TargetID:   strings.TrimSpace(in.TargetID),
		Reason:     strings.TrimSpace(in.Reason),
		Status:     domain.ReportOpen,
	}
	if err := s.repo.CreateReport(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

// ListReports returns the moderation queue, optionally filtered by status.
func (s *Service) ListReports(ctx context.Context, status string) ([]domain.Report, error) {
	if status != "" && status != domain.ReportOpen && status != domain.ReportReviewing &&
		!domain.IsTerminalReportStatus(status) {
		return nil, ValidationError{"invalid status filter"}
	}
	return s.repo.ListReports(ctx, status)
}

// ResolveReport triages a report to resolved or dismissed.
func (s *Service) ResolveReport(ctx context.Context, adminID, id, status, actionTaken string) (*domain.Report, error) {
	if !domain.IsTerminalReportStatus(status) {
		return nil, ValidationError{"status must be resolved or dismissed"}
	}
	rep, err := s.repo.GetReport(ctx, id)
	if err != nil {
		return nil, err
	}
	if domain.IsTerminalReportStatus(rep.Status) {
		return nil, ValidationError{"this report has already been triaged"}
	}
	if err := s.repo.ResolveReport(ctx, id, status, strings.TrimSpace(actionTaken), adminID); err != nil {
		return nil, err
	}
	s.audit(ctx, adminID, "admin.report."+status, "report", id, map[string]any{"action_taken": actionTaken})
	return s.repo.GetReport(ctx, id)
}

// ----- Analytics ---------------------------------------------------------

// Analytics returns the platform overview metrics.
func (s *Service) Analytics(ctx context.Context) (*domain.Analytics, error) {
	return s.repo.Analytics(ctx)
}

func (s *Service) audit(ctx context.Context, actorID, action, targetType, targetID string, metadata map[string]any) {
	_ = s.repo.WriteAudit(ctx, actorID, action, targetType, targetID, metadata)
}
