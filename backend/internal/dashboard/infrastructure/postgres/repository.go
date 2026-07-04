// Package postgres implements dashboard/domain.Repository on PostgreSQL. It runs
// read-only COUNT aggregates against the tables owned by the jobs, referrals,
// mentorship and notifications contexts. As a backend-for-frontend read model it
// is permitted to query those tables directly; it never writes.
package postgres

import (
	"context"
	"database/sql"

	"workspace-app/internal/dashboard/domain"
)

// Repository is the PostgreSQL dashboard read adapter.
type Repository struct{ db *sql.DB }

// NewRepository constructs the dashboard repository.
func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

// count runs a single-row COUNT(*) query and returns the integer result.
func (r *Repository) count(ctx context.Context, query string, args ...any) (int, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// Summary aggregates the per-user dashboard counts across modules.
func (r *Repository) Summary(ctx context.Context, userID string) (domain.Summary, error) {
	var s domain.Summary
	var err error

	if s.UnreadNotifications, err = r.count(ctx,
		`SELECT count(*) FROM notifications WHERE user_id = $1 AND read_at IS NULL`, userID); err != nil {
		return s, err
	}

	// Job-seeker.
	if s.JobSeeker.Applications, err = r.count(ctx,
		`SELECT count(*) FROM job_applications WHERE user_id = $1`, userID); err != nil {
		return s, err
	}
	if s.JobSeeker.SavedJobs, err = r.count(ctx,
		`SELECT count(*) FROM saved_jobs WHERE user_id = $1`, userID); err != nil {
		return s, err
	}
	if s.JobSeeker.OutgoingReferrals, err = r.count(ctx,
		`SELECT count(*) FROM referral_requests WHERE seeker_id = $1 AND deleted_at IS NULL`, userID); err != nil {
		return s, err
	}

	// Recruiter.
	if s.Recruiter.PostedJobs, err = r.count(ctx,
		`SELECT count(*) FROM jobs WHERE posted_by = $1`, userID); err != nil {
		return s, err
	}
	if s.Recruiter.TotalApplicants, err = r.count(ctx,
		`SELECT count(*) FROM job_applications a JOIN jobs j ON j.id = a.job_id WHERE j.posted_by = $1`, userID); err != nil {
		return s, err
	}
	if s.Recruiter.IncomingReferrals, err = r.count(ctx,
		`SELECT count(*) FROM referral_requests WHERE referrer_id = $1 AND deleted_at IS NULL`, userID); err != nil {
		return s, err
	}

	// Mentor (sessions are joined through the caller's mentor_profiles row).
	const mentorJoin = `FROM mentorship_sessions s JOIN mentor_profiles m ON m.id = s.mentor_id WHERE m.user_id = $1`
	if s.Mentor.UpcomingSessions, err = r.count(ctx,
		`SELECT count(*) `+mentorJoin+` AND s.status = 'confirmed' AND s.scheduled_at > now()`, userID); err != nil {
		return s, err
	}
	if s.Mentor.PendingRequests, err = r.count(ctx,
		`SELECT count(*) `+mentorJoin+` AND s.status = 'requested'`, userID); err != nil {
		return s, err
	}
	if s.Mentor.CompletedSessions, err = r.count(ctx,
		`SELECT count(*) `+mentorJoin+` AND s.status = 'completed'`, userID); err != nil {
		return s, err
	}

	return s, nil
}
