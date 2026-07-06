// Package postgres implements messaging/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"workspace-app/internal/messaging/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateConversation(ctx context.Context, c *domain.Conversation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := tx.QueryRowContext(ctx,
		`INSERT INTO conversations (type, title, created_by) VALUES ($1, NULLIF($2,''), $3)
		 RETURNING id, created_at, updated_at`,
		c.Type, c.Title, c.CreatedBy).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return err
	}
	for _, uid := range c.ParticipantIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO conversation_participants (conversation_id, user_id, role) VALUES ($1,$2,$3)
			 ON CONFLICT DO NOTHING`, c.ID, uid, "member"); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	var c domain.Conversation
	var lastMessageAt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT id, type, COALESCE(title,''), COALESCE(created_by::text,''), created_at, updated_at, last_message_at
		 FROM conversations WHERE id = $1`, id).
		Scan(&c.ID, &c.Type, &c.Title, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt, &lastMessageAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if lastMessageAt.Valid {
		c.LastMessageAt = &lastMessageAt.Time
	}
	if c.ParticipantIDs, err = r.Participants(ctx, id); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.type, COALESCE(c.title,''), COALESCE(c.created_by::text,''), c.created_at, c.updated_at, c.last_message_at,
		       p.is_pinned, p.is_archived, p.last_read_at
		FROM conversations c
		JOIN conversation_participants p ON p.conversation_id = c.id
		WHERE p.user_id = $1 AND p.left_at IS NULL
		ORDER BY p.is_pinned DESC, COALESCE(c.last_message_at, c.updated_at) DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Conversation{}
	for rows.Next() {
		var c domain.Conversation
		var lastMessageAt sql.NullTime
		var lastReadAt sql.NullTime
		err := rows.Scan(
			&c.ID, &c.Type, &c.Title, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt, &lastMessageAt,
			&c.IsPinned, &c.IsArchived, &lastReadAt,
		)
		if err != nil {
			return nil, err
		}
		if lastMessageAt.Valid {
			c.LastMessageAt = &lastMessageAt.Time
		}
		if lastReadAt.Valid {
			c.LastReadAt = &lastReadAt.Time
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
		JOIN conversations c ON c.id = cp.conversation_id AND c.type = 'direct'
		WHERE cp.left_at IS NULL
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
		`SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id=$1 AND user_id=$2 AND left_at IS NULL)`,
		conversationID, userID).Scan(&exists)
	return exists, err
}

func (r *Repository) Participants(ctx context.Context, conversationID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id FROM conversation_participants WHERE conversation_id = $1 AND left_at IS NULL`, conversationID)
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

func (r *Repository) GetParticipantDetail(ctx context.Context, conversationID, userID string) (*domain.Participant, error) {
	var p domain.Participant
	var leftAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT conversation_id, user_id, joined_at, left_at, role, is_archived, is_pinned
		FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID).Scan(&p.ConversationID, &p.UserID, &p.JoinedAt, &leftAt, &p.Role, &p.IsArchived, &p.IsPinned)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if leftAt.Valid {
		p.LeftAt = &leftAt.Time
	}
	return &p, nil
}

func (r *Repository) AddMessage(ctx context.Context, m *domain.Message) error {
	if err := r.db.QueryRowContext(ctx,
		`INSERT INTO messages (conversation_id, sender_id, content, content_type) VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		m.ConversationID, m.SenderID, m.Content, m.ContentType).Scan(&m.ID, &m.CreatedAt); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE conversations SET last_message_at = now(), updated_at = now() WHERE id = $1`, m.ConversationID)
	return err
}

func (r *Repository) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	var m domain.Message
	var editedAt, deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, conversation_id, sender_id, content, content_type, created_at, edited_at, deleted_at
		FROM messages
		WHERE id = $1
	`, id).Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.ContentType, &m.CreatedAt, &editedAt, &deletedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if editedAt.Valid {
		m.EditedAt = &editedAt.Time
	}
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}
	return &m, nil
}

func (r *Repository) ListMessages(ctx context.Context, conversationID string, limit int) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, conversation_id, sender_id, content, content_type, created_at, edited_at, deleted_at FROM messages
		 WHERE conversation_id = $1 ORDER BY created_at LIMIT $2`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Message{}
	for rows.Next() {
		var m domain.Message
		var editedAt, deletedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.ContentType, &m.CreatedAt, &editedAt, &deletedAt); err != nil {
			return nil, err
		}
		if editedAt.Valid {
			m.EditedAt = &editedAt.Time
		}
		if deletedAt.Valid {
			m.DeletedAt = &deletedAt.Time
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repository) UpdateMessage(ctx context.Context, m *domain.Message) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE messages
		SET content = $1, content_type = $2, edited_at = $3, deleted_at = $4
		WHERE id = $5
	`, m.Content, m.ContentType, m.EditedAt, m.DeletedAt, m.ID)
	return err
}

func (r *Repository) MarkRead(ctx context.Context, conversationID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE conversation_participants SET last_read_at = now() WHERE conversation_id=$1 AND user_id=$2`,
		conversationID, userID)
	return err
}

func (r *Repository) SetMessageStatus(ctx context.Context, status *domain.MessageStatus) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO message_statuses (message_id, user_id, status, status_updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (message_id, user_id) DO UPDATE
		SET status = $3, status_updated_at = $4
	`, status.MessageID, status.UserID, status.Status, status.StatusUpdatedAt)
	return err
}

func (r *Repository) GetMessageStatuses(ctx context.Context, messageID string) ([]domain.MessageStatus, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT message_id, user_id, status, status_updated_at
		FROM message_statuses
		WHERE message_id = $1
	`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.MessageStatus
	for rows.Next() {
		var s domain.MessageStatus
		if err := rows.Scan(&s.MessageID, &s.UserID, &s.Status, &s.StatusUpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) ArchiveConversation(ctx context.Context, conversationID, userID string, archive bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE conversation_participants
		SET is_archived = $1
		WHERE conversation_id = $2 AND user_id = $3
	`, archive, conversationID, userID)
	return err
}

func (r *Repository) PinConversation(ctx context.Context, conversationID, userID string, pin bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE conversation_participants
		SET is_pinned = $1
		WHERE conversation_id = $2 AND user_id = $3
	`, pin, conversationID, userID)
	return err
}

func (r *Repository) GetUnreadCount(ctx context.Context, conversationID, userID string, lastReadAt *time.Time) (int, error) {
	var count int
	var query string
	var args []any
	if lastReadAt != nil {
		query = `SELECT COUNT(*) FROM messages WHERE conversation_id = $1 AND sender_id != $2 AND created_at > $3 AND deleted_at IS NULL`
		args = []any{conversationID, userID, *lastReadAt}
	} else {
		query = `SELECT COUNT(*) FROM messages WHERE conversation_id = $1 AND sender_id != $2 AND deleted_at IS NULL`
		args = []any{conversationID, userID}
	}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *Repository) GetLastMessage(ctx context.Context, conversationID string) (*domain.Message, error) {
	var m domain.Message
	var editedAt, deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, conversation_id, sender_id, content, content_type, created_at, edited_at, deleted_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1
	`, conversationID).Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.ContentType, &m.CreatedAt, &editedAt, &deletedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // No message yet
	}
	if err != nil {
		return nil, err
	}
	if editedAt.Valid {
		m.EditedAt = &editedAt.Time
	}
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}
	return &m, nil
}
