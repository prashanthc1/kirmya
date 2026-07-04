// Package domain holds the Dashboard bounded context. Dashboard is a read-only
// projection ("backend-for-frontend"): it aggregates per-user counts that the
// frontend dashboard renders by role. It owns no tables of its own and has no
// dependency on any other layer — infrastructure implements its Repository port.
package domain

import "context"

// JobSeekerStats summarizes a user's job-seeking activity.
type JobSeekerStats struct {
	Applications      int `json:"applications"`
	SavedJobs         int `json:"saved_jobs"`
	OutgoingReferrals int `json:"outgoing_referrals"`
}

// RecruiterStats summarizes a recruiter's hiring activity.
type RecruiterStats struct {
	PostedJobs        int `json:"posted_jobs"`
	TotalApplicants   int `json:"total_applicants"`
	IncomingReferrals int `json:"incoming_referrals"`
}

// MentorStats summarizes a mentor's mentoring activity.
type MentorStats struct {
	UpcomingSessions  int `json:"upcoming_sessions"`
	PendingRequests   int `json:"pending_requests"`
	CompletedSessions int `json:"completed_sessions"`
}

// Summary is the full per-user dashboard projection. The frontend shows the
// section(s) matching the user's active role(s); all sections are always
// populated so a multi-role user sees everything in one payload.
type Summary struct {
	UnreadNotifications int            `json:"unread_notifications"`
	JobSeeker           JobSeekerStats `json:"job_seeker"`
	Recruiter           RecruiterStats `json:"recruiter"`
	Mentor              MentorStats    `json:"mentor"`
}

// Repository is the read port for dashboard aggregates.
type Repository interface {
	Summary(ctx context.Context, userID string) (Summary, error)
}
