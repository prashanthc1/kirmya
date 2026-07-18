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

func minMax(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}

func scanConnection(row interface{ Scan(...any) error }, c *domain.Connection) error {
	var createdAt, updatedAt time.Time
	var respondedAt sql.NullTime
	var userAID, userBID, requestedBy, origin string
	err := row.Scan(&c.ID, &userAID, &userBID, &c.Status, &requestedBy, &origin, &respondedAt, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	c.RequesterID = requestedBy
	if requestedBy == userAID {
		c.ReceiverID = userBID
	} else {
		c.ReceiverID = userAID
	}
	c.Origin = domain.ConnectionOrigin(origin)
	c.CreatedAt = createdAt.Format(time.RFC3339)
	c.UpdatedAt = updatedAt.Format(time.RFC3339)
	if respondedAt.Valid {
		c.RespondedAt = respondedAt.Time.Format(time.RFC3339)
	} else {
		c.RespondedAt = ""
	}
	return nil
}

func (r *repository) reconcileCounts(ctx context.Context, userID string) error {
	const query = `
		INSERT INTO connection_counts (
			user_id,
			connection_count,
			pending_incoming_count,
			pending_outgoing_count,
			updated_at
		)
		VALUES (
			$1,
			(SELECT COUNT(*) FROM connections WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'accepted'),
			(SELECT COUNT(*) FROM connections WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'pending' AND requested_by != $1),
			(SELECT COUNT(*) FROM connections WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'pending' AND requested_by = $1),
			now()
		)
		ON CONFLICT (user_id) DO UPDATE SET
			connection_count = EXCLUDED.connection_count,
			pending_incoming_count = EXCLUDED.pending_incoming_count,
			pending_outgoing_count = EXCLUDED.pending_outgoing_count,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) Create(ctx context.Context, requesterID, receiverID string, origin domain.ConnectionOrigin) (*domain.Connection, error) {
	uA, uB := minMax(requesterID, receiverID)
	var c domain.Connection
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO connections (user_a_id, user_b_id, status, requested_by, origin)
		VALUES ($1, $2, 'pending', $3, $4)
		RETURNING id, user_a_id, user_b_id, status, requested_by, origin, responded_at, created_at, updated_at
	`, uA, uB, requesterID, string(origin))

	if err := scanConnection(row, &c); err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrDuplicateRequest
		}
		return nil, err
	}
	_ = r.reconcileCounts(ctx, requesterID)
	_ = r.reconcileCounts(ctx, receiverID)
	return &c, nil
}

func (r *repository) CreateAccepted(ctx context.Context, requesterID, receiverID string, origin domain.ConnectionOrigin) (*domain.Connection, error) {
	uA, uB := minMax(requesterID, receiverID)
	var c domain.Connection
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO connections (user_a_id, user_b_id, status, requested_by, origin, responded_at)
		VALUES ($1, $2, 'accepted', $3, $4, CURRENT_TIMESTAMP)
		RETURNING id, user_a_id, user_b_id, status, requested_by, origin, responded_at, created_at, updated_at
	`, uA, uB, requesterID, string(origin))

	if err := scanConnection(row, &c); err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrDuplicateRequest
		}
		return nil, err
	}
	_ = r.reconcileCounts(ctx, requesterID)
	_ = r.reconcileCounts(ctx, receiverID)
	return &c, nil
}

func (r *repository) UpdateStatus(ctx context.Context, connectionID string, status domain.ConnectionStatus) error {
	var userAID, userBID string
	err := r.db.QueryRowContext(ctx, `
		SELECT user_a_id, user_b_id FROM connections WHERE id = $1
	`, connectionID).Scan(&userAID, &userBID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	res, err := r.db.ExecContext(ctx, `
		UPDATE connections
		SET status = $1::text, responded_at = CASE WHEN $1::text IN ('accepted', 'declined') THEN CURRENT_TIMESTAMP ELSE responded_at END, updated_at = CURRENT_TIMESTAMP
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
	_ = r.reconcileCounts(ctx, userAID)
	_ = r.reconcileCounts(ctx, userBID)
	return nil
}

func (r *repository) Block(ctx context.Context, blockerID, blockedID string) error {
	uA, uB := minMax(blockerID, blockedID)
	// Check if a connection already exists in either direction
	var id string
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM connections
		WHERE user_a_id = $1 AND user_b_id = $2
	`, uA, uB).Scan(&id)

	if err == nil {
		// Update existing connection status to blocked
		_, err = r.db.ExecContext(ctx, `
			UPDATE connections
			SET status = 'blocked', requested_by = $1, responded_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, blockerID, id)
		if err != nil {
			return err
		}
		// Insert into blocks table as well, to be consistent with connections module
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO blocks (blocker_id, blocked_id, created_at)
			VALUES ($1, $2, now())
			ON CONFLICT (blocker_id, blocked_id) DO NOTHING
		`, blockerID, blockedID)

		_ = r.reconcileCounts(ctx, blockerID)
		_ = r.reconcileCounts(ctx, blockedID)
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		// Create new blocked connection
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO connections (user_a_id, user_b_id, status, requested_by, origin, responded_at)
			VALUES ($1, $2, 'blocked', $3, 'manual_request', CURRENT_TIMESTAMP)
		`, uA, uB, blockerID)
		if err != nil {
			return err
		}
		// Insert into blocks table as well, to be consistent with connections module
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO blocks (blocker_id, blocked_id, created_at)
			VALUES ($1, $2, now())
			ON CONFLICT (blocker_id, blocked_id) DO NOTHING
		`, blockerID, blockedID)

		_ = r.reconcileCounts(ctx, blockerID)
		_ = r.reconcileCounts(ctx, blockedID)
		return nil
	}

	return err
}

func (r *repository) Unconnect(ctx context.Context, userA, userB string) error {
	return r.Delete(ctx, userA, userB)
}

func (r *repository) GetConnections(ctx context.Context, userID string) ([]domain.Connection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.origin, c.responded_at, c.created_at, c.updated_at,
			u1.full_name as req_name, COALESCE(p1.headline, '') as req_hl, COALESCE(p1.photo_url, '') as req_photo,
			u2.full_name as rec_name, COALESCE(p2.headline, '') as rec_hl, COALESCE(p2.photo_url, '') as rec_photo
		FROM connections c
		JOIN users u1 ON c.requested_by = u1.id
		LEFT JOIN profiles p1 ON c.requested_by = p1.user_id
		JOIN users u2 ON (u2.id = CASE WHEN c.user_a_id = c.requested_by THEN c.user_b_id ELSE c.user_a_id END)
		LEFT JOIN profiles p2 ON (p2.user_id = CASE WHEN c.user_a_id = c.requested_by THEN c.user_b_id ELSE c.user_a_id END)
		WHERE (c.user_a_id = $1 OR c.user_b_id = $1) AND c.status = 'accepted'
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
		var respondedAt sql.NullTime
		var userAID, userBID, requestedBy, origin string
		err := rows.Scan(
			&c.ID, &userAID, &userBID, &c.Status, &requestedBy, &origin, &respondedAt, &createdAt, &updatedAt,
			&c.RequesterName, &c.RequesterHeadline, &c.RequesterPhotoURL,
			&c.ReceiverName, &c.ReceiverHeadline, &c.ReceiverPhotoURL,
		)
		if err != nil {
			return nil, err
		}
		c.RequesterID = requestedBy
		if requestedBy == userAID {
			c.ReceiverID = userBID
		} else {
			c.ReceiverID = userAID
		}
		c.Origin = domain.ConnectionOrigin(origin)
		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		if respondedAt.Valid {
			c.RespondedAt = respondedAt.Time.Format(time.RFC3339)
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetIncomingRequests(ctx context.Context, userID string) ([]domain.Connection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.origin, c.responded_at, c.created_at, c.updated_at,
			u1.full_name as req_name, COALESCE(p1.headline, '') as req_hl, COALESCE(p1.photo_url, '') as req_photo,
			u2.full_name as rec_name, COALESCE(p2.headline, '') as rec_hl, COALESCE(p2.photo_url, '') as rec_photo
		FROM connections c
		JOIN users u1 ON c.requested_by = u1.id
		LEFT JOIN profiles p1 ON c.requested_by = p1.user_id
		JOIN users u2 ON (u2.id = CASE WHEN c.user_a_id = c.requested_by THEN c.user_b_id ELSE c.user_a_id END)
		LEFT JOIN profiles p2 ON (p2.user_id = CASE WHEN c.user_a_id = c.requested_by THEN c.user_b_id ELSE c.user_a_id END)
		WHERE (c.user_a_id = $1 OR c.user_b_id = $1) AND c.status = 'pending' AND c.requested_by != $1
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
		var respondedAt sql.NullTime
		var userAID, userBID, requestedBy, origin string
		err := rows.Scan(
			&c.ID, &userAID, &userBID, &c.Status, &requestedBy, &origin, &respondedAt, &createdAt, &updatedAt,
			&c.RequesterName, &c.RequesterHeadline, &c.RequesterPhotoURL,
			&c.ReceiverName, &c.ReceiverHeadline, &c.ReceiverPhotoURL,
		)
		if err != nil {
			return nil, err
		}
		c.RequesterID = requestedBy
		if requestedBy == userAID {
			c.ReceiverID = userBID
		} else {
			c.ReceiverID = userAID
		}
		c.Origin = domain.ConnectionOrigin(origin)
		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		if respondedAt.Valid {
			c.RespondedAt = respondedAt.Time.Format(time.RFC3339)
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetConnectionStatus(ctx context.Context, userA, userB string) (domain.ConnectionStatus, string, error) {
	uA, uB := minMax(userA, userB)
	var status domain.ConnectionStatus
	var requesterID string
	err := r.db.QueryRowContext(ctx, `
		SELECT status, requested_by
		FROM connections
		WHERE user_a_id = $1 AND user_b_id = $2
	`, uA, uB).Scan(&status, &requesterID)
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
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_a_id, user_b_id, status, requested_by, origin, responded_at, created_at, updated_at
		FROM connections
		WHERE id = $1
	`, id)
	if err := scanConnection(row, &c); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Delete(ctx context.Context, requesterID, receiverID string) error {
	uA, uB := minMax(requesterID, receiverID)
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM connections
		WHERE user_a_id = $1 AND user_b_id = $2
	`, uA, uB)
	if err != nil {
		return err
	}
	_ = r.reconcileCounts(ctx, requesterID)
	_ = r.reconcileCounts(ctx, receiverID)
	return nil
}
