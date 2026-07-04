// Package postgres implements referrals/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/referrals/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const cols = `id, seeker_id, COALESCE(referrer_id::text,''), COALESCE(job_id::text,''),
              COALESCE(company,''), COALESCE(message,''), status, COALESCE(outcome,''),
              decided_at, created_at, updated_at, version`

func scan(s interface {
	Scan(...any) error
}) (domain.Referral, error) {
	var r domain.Referral
	var decided sql.NullTime
	err := s.Scan(&r.ID, &r.SeekerID, &r.ReferrerID, &r.JobID, &r.Company, &r.Message,
		&r.Status, &r.Outcome, &decided, &r.CreatedAt, &r.UpdatedAt, &r.Version)
	if decided.Valid {
		r.DecidedAt = &decided.Time
	}
	return r, err
}

func (r *Repository) Create(ctx context.Context, ref *domain.Referral) error {
	const q = `
		INSERT INTO referral_requests (seeker_id, referrer_id, job_id, company, message, status)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid, NULLIF($4,''), NULLIF($5,''), $6)
		RETURNING id, created_at, updated_at, version`
	return r.db.QueryRowContext(ctx, q,
		ref.SeekerID, ref.ReferrerID, ref.JobID, ref.Company, ref.Message, ref.Status).
		Scan(&ref.ID, &ref.CreatedAt, &ref.UpdatedAt, &ref.Version)
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.Referral, error) {
	ref, err := scan(r.db.QueryRowContext(ctx,
		`SELECT `+cols+` FROM referral_requests WHERE id = $1 AND deleted_at IS NULL`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *Repository) ListByReferrer(ctx context.Context, referrerID string) ([]domain.Referral, error) {
	return r.list(ctx, `SELECT `+cols+` FROM referral_requests
		WHERE referrer_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`, referrerID)
}

func (r *Repository) ListBySeeker(ctx context.Context, seekerID string) ([]domain.Referral, error) {
	return r.list(ctx, `SELECT `+cols+` FROM referral_requests
		WHERE seeker_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`, seekerID)
}

func (r *Repository) list(ctx context.Context, q, arg string) ([]domain.Referral, error) {
	rows, err := r.db.QueryContext(ctx, q, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Referral{}
	for rows.Next() {
		ref, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ref)
	}
	return out, rows.Err()
}

func (r *Repository) Decide(ctx context.Context, id, referrerID, status string) error {
	const q = `
		UPDATE referral_requests
		SET referrer_id = $2::uuid, status = $3, decided_at = now(), updated_at = now(), version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id, referrerID, status)
	return err
}

func (r *Repository) SetOutcome(ctx context.Context, id, outcome string) error {
	const q = `
		UPDATE referral_requests
		SET outcome = $2, updated_at = now(), version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id, outcome)
	return err
}
