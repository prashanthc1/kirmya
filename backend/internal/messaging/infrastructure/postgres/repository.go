// Package postgres implements messaging/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/messaging/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateConversation(ctx context.Context, c *domain.Conversation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.QueryRowContext(ctx,
		`INSERT INTO conversations (is_group, title, created_by) VALUES ($1, NULLIF($2,''), $3)
		 RETURNING id, created_at, updated_at`,
		c.IsGroup, c.Title, c.CreatedBy).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return err
	}
	for _, uid := range c.ParticipantIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1,$2)
			 ON CONFLICT DO NOTHING`, c.ID, uid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.QueryRowContext(ctx,
		`SELECT id, is_group, COALESCE(title,''), COALESCE(created_by::text,''), created_at, updated_at
		 FROM conversations WHERE id = $1`, id).
		Scan(&c.ID, &c.IsGroup, &c.Title, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if c.ParticipantIDs, err = r.Participants(ctx, id); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.is_group, COALESCE(c.title,''), COALESCE(c.created_by::text,''), c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants p ON p.conversation_id = c.id
		WHERE p.user_id = $1
		ORDER BY c.updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Conversation{}
	for rows.Next() {
		var c domain.Conversation
		if err := rows.Scan(&c.ID, &c.IsGroup, &c.Title, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		if out[i].ParticipantIDs, err = r.Participants(ctx, out[i].ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *Repository) FindDirect(ctx context.Context, userA, userB string) (string, bool, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		SELECT cp.conversation_id
		FROM conversation_participants cp
		JOIN conversations c ON c.id = cp.conversation_id AND c.is_group = false
		GROUP BY cp.conversation_id
		HAVING COUNT(*) = 2 AND COUNT(*) FILTER (WHERE cp.user_id IN ($1,$2)) = 2
		LIMIT 1`, userA, userB).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (r *Repository) IsParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id=$1 AND user_id=$2)`,
		conversationID, userID).Scan(&exists)
	return exists, err
}

func (r *Repository) Participants(ctx context.Context, conversationID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id FROM conversation_participants WHERE conversation_id = $1`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (r *Repository) AddMessage(ctx context.Context, m *domain.Message) error {
	if err := r.db.QueryRowContext(ctx,
		`INSERT INTO messages (conversation_id, sender_id, body) VALUES ($1,$2,$3) RETURNING id, created_at`,
		m.ConversationID, m.SenderID, m.Body).Scan(&m.ID, &m.CreatedAt); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE conversations SET updated_at = now() WHERE id = $1`, m.ConversationID)
	return err
}

func (r *Repository) ListMessages(ctx context.Context, conversationID string, limit int) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, conversation_id, sender_id, body, created_at FROM messages
		 WHERE conversation_id = $1 ORDER BY created_at LIMIT $2`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Message{}
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repository) MarkRead(ctx context.Context, conversationID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE conversation_participants SET last_read_at = now() WHERE conversation_id=$1 AND user_id=$2`,
		conversationID, userID)
	return err
}
