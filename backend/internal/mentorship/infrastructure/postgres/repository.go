// Package postgres implements mentorship/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/mentorship/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) UpsertMentorProfile(ctx context.Context, p *domain.MentorProfile) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO mentor_profiles (user_id, headline, bio, expertise, is_active)
		VALUES ($1,$2,NULLIF($3,''),NULLIF($4,''),true)
		ON CONFLICT (user_id) DO UPDATE
		  SET headline = EXCLUDED.headline, bio = EXCLUDED.bio, expertise = EXCLUDED.expertise,
		      is_active = true, updated_at = now()
		RETURNING id, created_at, updated_at`,
		p.UserID, p.Headline, p.Bio, p.Expertise).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

const mentorCols = `id, user_id, headline, COALESCE(bio,''), COALESCE(expertise,''), is_active, created_at, updated_at`

func scanMentor(s interface{ Scan(...any) error }) (domain.MentorProfile, error) {
	var p domain.MentorProfile
	err := s.Scan(&p.ID, &p.UserID, &p.Headline, &p.Bio, &p.Expertise, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetMentorByID(ctx context.Context, id string) (*domain.MentorProfile, error) {
	p, err := scanMentor(r.db.QueryRowContext(ctx, `SELECT `+mentorCols+` FROM mentor_profiles WHERE id = $1`, id))
	return mentorOrErr(p, err)
}

func (r *Repository) GetMentorByUserID(ctx context.Context, userID string) (*domain.MentorProfile, error) {
	p, err := scanMentor(r.db.QueryRowContext(ctx, `SELECT `+mentorCols+` FROM mentor_profiles WHERE user_id = $1`, userID))
	return mentorOrErr(p, err)
}

func mentorOrErr(p domain.MentorProfile, err error) (*domain.MentorProfile, error) {
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListMentors(ctx context.Context) ([]domain.MentorProfile, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+mentorCols+` FROM mentor_profiles WHERE is_active ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.MentorProfile{}
	for rows.Next() {
		p, err := scanMentor(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

const sessionCols = `id, mentor_id, mentee_id, COALESCE(topic,''), status, scheduled_at, created_at, updated_at`

func scanSession(s interface{ Scan(...any) error }) (domain.Session, error) {
	var x domain.Session
	err := s.Scan(&x.ID, &x.MentorID, &x.MenteeID, &x.Topic, &x.Status, &x.ScheduledAt, &x.CreatedAt, &x.UpdatedAt)
	return x, err
}

func (r *Repository) CreateSession(ctx context.Context, s *domain.Session) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO mentorship_sessions (mentor_id, mentee_id, topic, status, scheduled_at)
		VALUES ($1,$2,NULLIF($3,''),$4,$5)
		RETURNING id, created_at, updated_at`,
		s.MentorID, s.MenteeID, s.Topic, s.Status, s.ScheduledAt).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *Repository) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	x, err := scanSession(r.db.QueryRowContext(ctx, `SELECT `+sessionCols+` FROM mentorship_sessions WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *Repository) ListSessionsForMentee(ctx context.Context, menteeID string) ([]domain.Session, error) {
	return r.querySessions(ctx, `SELECT `+sessionCols+` FROM mentorship_sessions WHERE mentee_id=$1 ORDER BY scheduled_at DESC`, menteeID)
}

func (r *Repository) ListSessionsForMentor(ctx context.Context, mentorID string) ([]domain.Session, error) {
	return r.querySessions(ctx, `SELECT `+sessionCols+` FROM mentorship_sessions WHERE mentor_id=$1 ORDER BY scheduled_at DESC`, mentorID)
}

func (r *Repository) querySessions(ctx context.Context, q, arg string) ([]domain.Session, error) {
	rows, err := r.db.QueryContext(ctx, q, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Session{}
	for rows.Next() {
		x, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repository) UpdateSessionStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE mentorship_sessions SET status=$2, updated_at=now() WHERE id=$1`, id, status)
	return err
}

func (r *Repository) CreateReview(ctx context.Context, rv *domain.Review) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO mentorship_reviews (session_id, rating, comment)
		VALUES ($1,$2,NULLIF($3,''))
		ON CONFLICT (session_id) DO UPDATE SET rating = EXCLUDED.rating, comment = EXCLUDED.comment
		RETURNING id, created_at`,
		rv.SessionID, rv.Rating, rv.Comment).Scan(&rv.ID, &rv.CreatedAt)
}

func (r *Repository) ListReviewsForMentor(ctx context.Context, mentorID string) ([]domain.Review, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rv.id, rv.session_id, rv.rating, COALESCE(rv.comment,''), rv.created_at
		FROM mentorship_reviews rv
		JOIN mentorship_sessions s ON s.id = rv.session_id
		WHERE s.mentor_id = $1
		ORDER BY rv.created_at DESC`, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Review{}
	for rows.Next() {
		var rv domain.Review
		if err := rows.Scan(&rv.ID, &rv.SessionID, &rv.Rating, &rv.Comment, &rv.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rv)
	}
	return out, rows.Err()
}

const slotCols = `id, mentor_id, starts_at, ends_at, is_booked, created_at`

func scanSlot(s interface{ Scan(...any) error }) (domain.AvailabilitySlot, error) {
	var sl domain.AvailabilitySlot
	err := s.Scan(&sl.ID, &sl.MentorID, &sl.StartsAt, &sl.EndsAt, &sl.IsBooked, &sl.CreatedAt)
	return sl, err
}

func (r *Repository) AddAvailability(ctx context.Context, slot *domain.AvailabilitySlot) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO mentor_availability (mentor_id, starts_at, ends_at)
		VALUES ($1,$2,$3)
		RETURNING id, is_booked, created_at`,
		slot.MentorID, slot.StartsAt, slot.EndsAt).Scan(&slot.ID, &slot.IsBooked, &slot.CreatedAt)
}

func (r *Repository) ListAvailability(ctx context.Context, mentorID string, openOnly bool) ([]domain.AvailabilitySlot, error) {
	q := `SELECT ` + slotCols + ` FROM mentor_availability WHERE mentor_id=$1`
	if openOnly {
		q += ` AND is_booked = false`
	}
	q += ` ORDER BY starts_at`
	rows, err := r.db.QueryContext(ctx, q, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.AvailabilitySlot{}
	for rows.Next() {
		sl, err := scanSlot(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sl)
	}
	return out, rows.Err()
}

func (r *Repository) GetSlot(ctx context.Context, id string) (*domain.AvailabilitySlot, error) {
	sl, err := scanSlot(r.db.QueryRowContext(ctx, `SELECT `+slotCols+` FROM mentor_availability WHERE id=$1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSlotNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

// CreateSessionWithSlot claims the slot and creates the session atomically. The
// conditional UPDATE (is_booked = false) is the race guard: only one concurrent
// caller flips the flag, so RowsAffected == 0 means the slot was already taken
// (or removed) and the whole transaction is rolled back — no orphan session.
func (r *Repository) CreateSessionWithSlot(ctx context.Context, s *domain.Session, slotID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx,
		`UPDATE mentor_availability SET is_booked = true WHERE id=$1 AND is_booked = false`, slotID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrSlotUnavailable
	}

	if err := tx.QueryRowContext(ctx, `
		INSERT INTO mentorship_sessions (mentor_id, mentee_id, topic, status, scheduled_at)
		VALUES ($1,$2,NULLIF($3,''),$4,$5)
		RETURNING id, created_at, updated_at`,
		s.MentorID, s.MenteeID, s.Topic, s.Status, s.ScheduledAt).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return err
	}
	return tx.Commit()
}
