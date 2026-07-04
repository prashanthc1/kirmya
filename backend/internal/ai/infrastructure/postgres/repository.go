// Package postgres implements ai/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/ai/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateThread(ctx context.Context, t *domain.CoachThread) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO coach_threads (user_id, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		t.UserID, t.Title).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) GetThread(ctx context.Context, id string) (*domain.CoachThread, error) {
	var t domain.CoachThread
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM coach_threads WHERE id = $1`, id).
		Scan(&t.ID, &t.UserID, &t.Title, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListThreads(ctx context.Context, userID string) ([]domain.CoachThread, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM coach_threads
		 WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.CoachThread{}
	for rows.Next() {
		var t domain.CoachThread
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) AddMessage(ctx context.Context, m *domain.CoachMessage) error {
	if err := r.db.QueryRowContext(ctx,
		`INSERT INTO coach_messages (thread_id, role, content) VALUES ($1, $2, $3) RETURNING id, created_at`,
		m.ThreadID, m.Role, m.Content).Scan(&m.ID, &m.CreatedAt); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE coach_threads SET updated_at = now() WHERE id = $1`, m.ThreadID)
	return err
}

func (r *Repository) ListMessages(ctx context.Context, threadID string) ([]domain.CoachMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, thread_id, role, content, created_at FROM coach_messages
		 WHERE thread_id = $1 ORDER BY created_at`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.CoachMessage{}
	for rows.Next() {
		var m domain.CoachMessage
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repository) LogInteraction(ctx context.Context, userID, kind string, c domain.Completion) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ai_interactions (user_id, kind, model, input_tokens, output_tokens)
		 VALUES (NULLIF($1,'')::uuid, $2, $3, $4, $5)`,
		userID, kind, c.Model, c.InputTokens, c.OutputTokens)
	return err
}
