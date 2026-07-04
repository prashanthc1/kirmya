// Package postgres implements settings/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/settings/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const settingsCols = `
	user_id, language, timezone, theme, email_digest,
	profile_visibility, show_email, discoverable, allow_messages,
	notif_email_jobs, notif_email_mentorship, notif_email_messages, notif_email_referrals,
	notif_inapp_jobs, notif_inapp_mentorship, notif_inapp_messages, notif_inapp_referrals,
	login_alerts, version, created_at, updated_at`

func scanSettings(s interface{ Scan(...any) error }) (domain.UserSettings, error) {
	var x domain.UserSettings
	n := &x.Notifications
	err := s.Scan(
		&x.UserID, &x.Language, &x.Timezone, &x.Theme, &x.EmailDigest,
		&x.ProfileVisibility, &x.ShowEmail, &x.Discoverable, &x.AllowMessages,
		&n.EmailJobs, &n.EmailMentorship, &n.EmailMessages, &n.EmailReferrals,
		&n.InAppJobs, &n.InAppMentorship, &n.InAppMessages, &n.InAppReferrals,
		&x.LoginAlerts, &x.Version, &x.CreatedAt, &x.UpdatedAt,
	)
	return x, err
}

func (r *Repository) Get(ctx context.Context, userID string) (*domain.UserSettings, error) {
	x, err := scanSettings(r.db.QueryRowContext(ctx, `SELECT `+settingsCols+` FROM user_settings WHERE user_id = $1`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// EnsureDefaults inserts the default row if missing and returns the current
// settings. The ON CONFLICT ... DO UPDATE (no-op self assignment) guarantees a
// row is always RETURNED, whether it was just created or already existed.
func (r *Repository) EnsureDefaults(ctx context.Context, userID string) (*domain.UserSettings, error) {
	x, err := scanSettings(r.db.QueryRowContext(ctx, `
		INSERT INTO user_settings (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING `+settingsCols, userID))
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// Update applies a version-checked write. RETURNING version confirms the row was
// updated; sql.ErrNoRows means the version was stale (concurrent write).
func (r *Repository) Update(ctx context.Context, s *domain.UserSettings) error {
	n := s.Notifications
	err := r.db.QueryRowContext(ctx, `
		UPDATE user_settings SET
			language = $2, timezone = $3, theme = $4, email_digest = $5,
			profile_visibility = $6, show_email = $7, discoverable = $8, allow_messages = $9,
			notif_email_jobs = $10, notif_email_mentorship = $11, notif_email_messages = $12, notif_email_referrals = $13,
			notif_inapp_jobs = $14, notif_inapp_mentorship = $15, notif_inapp_messages = $16, notif_inapp_referrals = $17,
			login_alerts = $18, version = version + 1, updated_at = now()
		WHERE user_id = $1 AND version = $19
		RETURNING version, updated_at`,
		s.UserID, s.Language, s.Timezone, s.Theme, s.EmailDigest,
		s.ProfileVisibility, s.ShowEmail, s.Discoverable, s.AllowMessages,
		n.EmailJobs, n.EmailMentorship, n.EmailMessages, n.EmailReferrals,
		n.InAppJobs, n.InAppMentorship, n.InAppMessages, n.InAppReferrals,
		s.LoginAlerts, s.Version,
	).Scan(&s.Version, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrOptimisticLock
	}
	return err
}
