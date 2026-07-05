package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"workspace-app/internal/network/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{db: db}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}

func (r *repository) Create(ctx context.Context, requesterID, receiverID string) (*domain.Connection, error) {
	var c domain.Connection
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO user_connections (requester_id, receiver_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, requester_id, receiver_id, status, created_at, updated_at
	`, requesterID, receiverID).Scan(&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &createdAt, &updatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrDuplicateRequest
		}
		return nil, err
	}
	c.CreatedAt = createdAt.Format(time.RFC3339)
	c.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &c, nil
}

func (r *repository) UpdateStatus(ctx context.Context, connectionID string, status domain.ConnectionStatus) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE user_connections
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, string(status), connectionID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *repository) GetConnections(ctx context.Context, userID string) ([]domain.Connection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			c.id, c.requester_id, c.receiver_id, c.status, c.created_at, c.updated_at,
			u1.full_name as req_name, COALESCE(p1.headline, '') as req_hl, COALESCE(p1.photo_url, '') as req_photo,
			u2.full_name as rec_name, COALESCE(p2.headline, '') as rec_hl, COALESCE(p2.photo_url, '') as rec_photo
		FROM user_connections c
		JOIN users u1 ON c.requester_id = u1.id
		LEFT JOIN profiles p1 ON c.requester_id = p1.user_id
		JOIN users u2 ON c.receiver_id = u2.id
		LEFT JOIN profiles p2 ON c.receiver_id = p2.user_id
		WHERE (c.requester_id = $1 OR c.receiver_id = $1) AND c.status = 'accepted'
		ORDER BY c.updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Connection
	for rows.Next() {
		var c domain.Connection
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &createdAt, &updatedAt,
			&c.RequesterName, &c.RequesterHeadline, &c.RequesterPhotoURL,
			&c.ReceiverName, &c.ReceiverHeadline, &c.ReceiverPhotoURL,
		)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) GetIncomingRequests(ctx context.Context, userID string) ([]domain.Connection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			c.id, c.requester_id, c.receiver_id, c.status, c.created_at, c.updated_at,
			u1.full_name as req_name, COALESCE(p1.headline, '') as req_hl, COALESCE(p1.photo_url, '') as req_photo,
			u2.full_name as rec_name, COALESCE(p2.headline, '') as rec_hl, COALESCE(p2.photo_url, '') as rec_photo
		FROM user_connections c
		JOIN users u1 ON c.requester_id = u1.id
		LEFT JOIN profiles p1 ON c.requester_id = p1.user_id
		JOIN users u2 ON c.receiver_id = u2.id
		LEFT JOIN profiles p2 ON c.receiver_id = p2.user_id
		WHERE c.receiver_id = $1 AND c.status = 'pending'
		ORDER BY c.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Connection
	for rows.Next() {
		var c domain.Connection
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &createdAt, &updatedAt,
			&c.RequesterName, &c.RequesterHeadline, &c.RequesterPhotoURL,
			&c.ReceiverName, &c.ReceiverHeadline, &c.ReceiverPhotoURL,
		)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) GetConnectionStatus(ctx context.Context, userA, userB string) (domain.ConnectionStatus, string, error) {
	var status domain.ConnectionStatus
	var requesterID string
	err := r.db.QueryRowContext(ctx, `
		SELECT status, requester_id
		FROM user_connections
		WHERE (requester_id = $1 AND receiver_id = $2) OR (requester_id = $2 AND receiver_id = $1)
	`, userA, userB).Scan(&status, &requesterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		return "", "", err
	}
	return status, requesterID, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Connection, error) {
	var c domain.Connection
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT id, requester_id, receiver_id, status, created_at, updated_at
		FROM user_connections
		WHERE id = $1
	`, id).Scan(&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	c.CreatedAt = createdAt.Format(time.RFC3339)
	c.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &c, nil
}

func (r *repository) Delete(ctx context.Context, requesterID, receiverID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_connections
		WHERE (requester_id = $1 AND receiver_id = $2) OR (requester_id = $2 AND receiver_id = $1)
	`, requesterID, receiverID)
	return err
}

