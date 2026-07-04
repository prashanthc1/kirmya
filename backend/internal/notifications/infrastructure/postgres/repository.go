// Package postgres implements notifications/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"

	"workspace-app/internal/notifications/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, n *domain.Notification) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO notifications (user_id, type, title, body, link)
		VALUES ($1,$2,$3,NULLIF($4,''),NULLIF($5,''))
		RETURNING id, created_at`,
		n.UserID, n.Type, n.Title, n.Body, n.Link).Scan(&n.ID, &n.CreatedAt)
}

func (r *Repository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]domain.Notification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, type, title, COALESCE(body,''), COALESCE(link,''), read_at, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY (read_at IS NULL) DESC, created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Notification{}
	for rows.Next() {
		var n domain.Notification
		var read sql.NullTime
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Link, &read, &n.CreatedAt); err != nil {
			return nil, err
		}
		if read.Valid {
			n.ReadAt = &read.Time
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *Repository) MarkRead(ctx context.Context, userID, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET read_at = now() WHERE id = $1 AND user_id = $2 AND read_at IS NULL`, id, userID)
	return err
}

func (r *Repository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET read_at = now() WHERE user_id = $1 AND read_at IS NULL`, userID)
	return err
}

func (r *Repository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read_at IS NULL`, userID).Scan(&n)
	return n, err
}
