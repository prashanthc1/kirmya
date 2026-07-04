package domain

import "testing"

func TestIsTerminalReportStatus(t *testing.T) {
	terminal := map[string]bool{
		ReportOpen:      false,
		ReportReviewing: false,
		ReportResolved:  true,
		ReportDismissed: true,
		"":              false,
	}
	for status, want := range terminal {
		if got := IsTerminalReportStatus(status); got != want {
			t.Errorf("IsTerminalReportStatus(%q) = %v, want %v", status, got, want)
		}
	}
}

func TestValidRolesAndTargetTypes(t *testing.T) {
	for _, r := range []string{"job_seeker", "referrer", "mentor", "recruiter", "admin"} {
		if !ValidRoles[r] {
			t.Errorf("expected %q to be a valid role", r)
		}
	}
	if ValidRoles["superuser"] {
		t.Error("did not expect an unknown role to be valid")
	}

	for _, ttype := range []string{"post", "comment", "user", "message"} {
		if !ValidTargetTypes[ttype] {
			t.Errorf("expected %q to be a valid target type", ttype)
		}
	}
	if ValidTargetTypes["planet"] {
		t.Error("did not expect an unknown target type to be valid")
	}
}
