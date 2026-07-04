package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/admin/domain"
)

func seedAdmin(repo *fakeRepo) {
	repo.addUser(&domain.UserSummary{ID: "admin-1", Email: "admin@cb.io", Status: domain.StatusActive, Roles: []string{"admin"}})
	repo.addUser(&domain.UserSummary{ID: "user-1", Email: "u1@cb.io", Status: domain.StatusActive, Roles: []string{"job_seeker"}})
}

func TestSetUserStatus(t *testing.T) {
	repo := newFakeRepo()
	seedAdmin(repo)
	svc := NewService(repo)
	ctx := context.Background()

	// Invalid status rejected.
	if _, err := svc.SetUserStatus(ctx, "admin-1", "user-1", "bogus"); err == nil {
		t.Fatal("expected validation error for bad status")
	}
	// Self-status change blocked.
	if _, err := svc.SetUserStatus(ctx, "admin-1", "admin-1", domain.StatusSuspended); err == nil {
		t.Fatal("expected error changing own status")
	}
	// Unknown user not found.
	if _, err := svc.SetUserStatus(ctx, "admin-1", "ghost", domain.StatusSuspended); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	// Happy path suspends.
	u, err := svc.SetUserStatus(ctx, "admin-1", "user-1", domain.StatusSuspended)
	if err != nil {
		t.Fatalf("suspend: %v", err)
	}
	if u.Status != domain.StatusSuspended {
		t.Fatalf("expected suspended, got %s", u.Status)
	}
	if len(repo.audits) == 0 {
		t.Fatal("expected an audit log entry")
	}
}

func TestRoleManagement(t *testing.T) {
	repo := newFakeRepo()
	seedAdmin(repo)
	svc := NewService(repo)
	ctx := context.Background()

	// Unknown role rejected.
	if _, err := svc.AssignRole(ctx, "admin-1", "user-1", "wizard"); err == nil {
		t.Fatal("expected validation error for unknown role")
	}
	// Grant mentor.
	u, err := svc.AssignRole(ctx, "admin-1", "user-1", "mentor")
	if err != nil {
		t.Fatalf("assign: %v", err)
	}
	if !contains(u.Roles, "mentor") {
		t.Fatalf("expected mentor role, got %v", u.Roles)
	}
	// Admin cannot revoke their own admin role.
	if _, err := svc.RevokeRole(ctx, "admin-1", "admin-1", "admin"); err == nil {
		t.Fatal("expected error revoking own admin role")
	}
	// Revoke mentor.
	u, err = svc.RevokeRole(ctx, "admin-1", "user-1", "mentor")
	if err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if contains(u.Roles, "mentor") {
		t.Fatalf("expected mentor removed, got %v", u.Roles)
	}
}

func TestFileAndResolveReport(t *testing.T) {
	repo := newFakeRepo()
	seedAdmin(repo)
	svc := NewService(repo)
	ctx := context.Background()

	// Bad target type rejected.
	if _, err := svc.FileReport(ctx, "user-1", ReportInput{TargetType: "spaceship", TargetID: "x", Reason: "nope"}); err == nil {
		t.Fatal("expected validation error for bad target type")
	}
	// Missing reason rejected.
	if _, err := svc.FileReport(ctx, "user-1", ReportInput{TargetType: "post", TargetID: "p1"}); err == nil {
		t.Fatal("expected validation error for missing reason")
	}
	// File a valid report.
	rep, err := svc.FileReport(ctx, "user-1", ReportInput{TargetType: "post", TargetID: "p1", Reason: "spam"})
	if err != nil {
		t.Fatalf("file: %v", err)
	}
	if rep.Status != domain.ReportOpen {
		t.Fatalf("expected open, got %s", rep.Status)
	}
	// Non-terminal status rejected on resolve.
	if _, err := svc.ResolveReport(ctx, "admin-1", rep.ID, domain.ReportReviewing, ""); err == nil {
		t.Fatal("expected validation error for non-terminal resolve status")
	}
	// Resolve it.
	resolved, err := svc.ResolveReport(ctx, "admin-1", rep.ID, domain.ReportResolved, "removed post")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolved.Status != domain.ReportResolved || resolved.ResolvedBy != "admin-1" {
		t.Fatalf("unexpected resolved report %+v", resolved)
	}
	// Re-resolving is rejected.
	if _, err := svc.ResolveReport(ctx, "admin-1", rep.ID, domain.ReportDismissed, ""); err == nil {
		t.Fatal("expected error re-triaging a resolved report")
	}
}

func TestRemovePost(t *testing.T) {
	repo := newFakeRepo()
	seedAdmin(repo)
	repo.posts["p1"] = true
	svc := NewService(repo)
	ctx := context.Background()

	if err := svc.RemovePost(ctx, "admin-1", "p1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if err := svc.RemovePost(ctx, "admin-1", "p1"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found on second remove, got %v", err)
	}
}

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
