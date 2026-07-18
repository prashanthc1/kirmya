package connections

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"workspace-app/internal/platform/tx"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// NormalizePair orders user IDs so user_a_id < user_b_id holds
func (r *Repository) NormalizePair(userA, userB string) (string, string) {
	if userA < userB {
		return userA, userB
	}
	return userB, userA
}

// GetConnectionStatus retrieves the connection row if it exists for a pair
func (r *Repository) GetConnection(ctx context.Context, userA, userB string) (*Connection, error) {
	uA, uB := r.NormalizePair(userA, userB)
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		SELECT c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.created_at, c.responded_at, c.updated_at,
		       meta.note, meta.source
		FROM connections c
		LEFT JOIN connection_requests_meta meta ON c.id = meta.connection_id
		WHERE c.user_a_id = $1 AND c.user_b_id = $2
	`

	var conn Connection
	var note sql.NullString
	var source sql.NullString
	var respondedAt sql.NullTime

	err := exec.QueryRowContext(ctx, query, uA, uB).Scan(
		&conn.ID, &conn.UserAID, &conn.UserBID, &conn.Status, &conn.RequestedBy,
		&conn.CreatedAt, &respondedAt, &conn.UpdatedAt, &note, &source,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if respondedAt.Valid {
		conn.RespondedAt = &respondedAt.Time
	}
	if note.Valid {
		conn.Note = &note.String
	}
	if source.Valid {
		src := ConnectionSource(source.String)
		conn.Source = &src
	}

	return &conn, nil
}

// GetConnectionByID retrieves a connection by its ID
func (r *Repository) GetConnectionByID(ctx context.Context, id string) (*Connection, error) {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		SELECT c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.created_at, c.responded_at, c.updated_at,
		       meta.note, meta.source
		FROM connections c
		LEFT JOIN connection_requests_meta meta ON c.id = meta.connection_id
		WHERE c.id = $1
	`

	var conn Connection
	var note sql.NullString
	var source sql.NullString
	var respondedAt sql.NullTime

	err := exec.QueryRowContext(ctx, query, id).Scan(
		&conn.ID, &conn.UserAID, &conn.UserBID, &conn.Status, &conn.RequestedBy,
		&conn.CreatedAt, &respondedAt, &conn.UpdatedAt, &note, &source,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if respondedAt.Valid {
		conn.RespondedAt = &respondedAt.Time
	}
	if note.Valid {
		conn.Note = &note.String
	}
	if source.Valid {
		src := ConnectionSource(source.String)
		conn.Source = &src
	}

	return &conn, nil
}

// CreateConnection inserts a new connection and its optional metadata
func (r *Repository) CreateConnection(ctx context.Context, fromUser, toUser string, status ConnectionStatus, note *string, source *ConnectionSource) (*Connection, error) {
	uA, uB := r.NormalizePair(fromUser, toUser)
	exec := tx.GetExecutor(ctx, r.db)

	const insertConn = `
		INSERT INTO connections (user_a_id, user_b_id, status, requested_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		RETURNING id, created_at, updated_at
	`

	var conn Connection
	conn.UserAID = uA
	conn.UserBID = uB
	conn.Status = status
	conn.RequestedBy = fromUser

	err := exec.QueryRowContext(ctx, insertConn, uA, uB, status, fromUser).Scan(&conn.ID, &conn.CreatedAt, &conn.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if note != nil || source != nil {
		const insertMeta = `
			INSERT INTO connection_requests_meta (connection_id, note, source)
			VALUES ($1, $2, $3)
		`
		var nVal sql.NullString
		var sVal sql.NullString
		if note != nil {
			nVal = sql.NullString{String: *note, Valid: true}
			conn.Note = note
		}
		if source != nil {
			sVal = sql.NullString{String: string(*source), Valid: true}
			conn.Source = source
		}
		_, err = exec.ExecContext(ctx, insertMeta, conn.ID, nVal, sVal)
		if err != nil {
			return nil, err
		}
	}

	return &conn, nil
}

// UpdateConnectionStatus updates status, responded_at, and updated_at
func (r *Repository) UpdateConnectionStatus(ctx context.Context, connectionID string, status ConnectionStatus, respondedAt *time.Time) error {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		UPDATE connections
		SET status = $1, responded_at = $2, updated_at = now()
		WHERE id = $3
	`
	var respVal sql.NullTime
	if respondedAt != nil {
		respVal = sql.NullTime{Time: *respondedAt, Valid: true}
	}

	res, err := exec.ExecContext(ctx, query, status, respVal, connectionID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteConnection deletes a connection row (and cascaded meta)
func (r *Repository) DeleteConnection(ctx context.Context, connectionID string) error {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `DELETE FROM connections WHERE id = $1`
	res, err := exec.ExecContext(ctx, query, connectionID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// InsertBlock records a block in the blocks table idempotently
func (r *Repository) InsertBlock(ctx context.Context, blocker, blocked, reason string) error {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		INSERT INTO blocks (blocker_id, blocked_id, reason, created_at)
		VALUES ($1, $2, NULLIF($3, ''), now())
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
	`
	_, err := exec.ExecContext(ctx, query, blocker, blocked, reason)
	return err
}

// DeleteBlock removes a block
func (r *Repository) DeleteBlock(ctx context.Context, blocker, blocked string) error {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		DELETE FROM blocks
		WHERE blocker_id = $1 AND blocked_id = $2
	`
	res, err := exec.ExecContext(ctx, query, blocker, blocked)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// IsBlocked checks if a block exists between two users in either direction
func (r *Repository) IsBlocked(ctx context.Context, userA, userB string) (bool, error) {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		SELECT EXISTS (
			SELECT 1 FROM blocks
			WHERE (blocker_id = $1 AND blocked_id = $2)
			   OR (blocker_id = $2 AND blocked_id = $1)
		)
	`
	var blocked bool
	err := exec.QueryRowContext(ctx, query, userA, userB).Scan(&blocked)
	return blocked, err
}

// ReconcileCounts recalculates counts and updates connection_counts table
func (r *Repository) ReconcileCounts(ctx context.Context, userID string) error {
	exec := tx.GetExecutor(ctx, r.db)

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
	_, err := exec.ExecContext(ctx, query, userID)
	return err
}

// GetConnections retrieves accepted connections for a user
func (r *Repository) GetConnections(ctx context.Context, userID string, page, limit int) ([]Connection, error) {
	exec := tx.GetExecutor(ctx, r.db)
	offset := (page - 1) * limit

	const query = `
		SELECT c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.created_at, c.responded_at, c.updated_at,
		       meta.note, meta.source,
		       other.id AS other_id, other.full_name AS other_name,
		       COALESCE(p.headline, '') AS other_headline, COALESCE(p.photo_url, '') AS other_avatar_url
		FROM connections c
		LEFT JOIN connection_requests_meta meta ON c.id = meta.connection_id
		JOIN users other ON (other.id = CASE WHEN c.user_a_id = $1 THEN c.user_b_id ELSE c.user_a_id END)
		LEFT JOIN profiles p ON other.id = p.user_id
		WHERE (c.user_a_id = $1 OR c.user_b_id = $1) AND c.status = 'accepted'
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := exec.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conns := []Connection{}
	for rows.Next() {
		var conn Connection
		var note sql.NullString
		var source sql.NullString
		var respondedAt sql.NullTime
		var other PublicProfileSummary

		err := rows.Scan(
			&conn.ID, &conn.UserAID, &conn.UserBID, &conn.Status, &conn.RequestedBy,
			&conn.CreatedAt, &respondedAt, &conn.UpdatedAt, &note, &source,
			&other.ID, &other.Name, &other.Headline, &other.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		if respondedAt.Valid {
			conn.RespondedAt = &respondedAt.Time
		}
		if note.Valid {
			conn.Note = &note.String
		}
		if source.Valid {
			src := ConnectionSource(source.String)
			conn.Source = &src
		}
		conn.User = &other
		conns = append(conns, conn)
	}

	return conns, rows.Err()
}

// GetPendingRequests retrieves pending connection requests (incoming or outgoing)
func (r *Repository) GetPendingRequests(ctx context.Context, userID string, direction string) ([]Connection, error) {
	exec := tx.GetExecutor(ctx, r.db)

	var filter string
	if direction == "incoming" {
		filter = "c.requested_by != $1"
	} else {
		filter = "c.requested_by = $1"
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.user_a_id, c.user_b_id, c.status, c.requested_by, c.created_at, c.responded_at, c.updated_at,
		       meta.note, meta.source,
		       other.id AS other_id, other.full_name AS other_name,
		       COALESCE(p.headline, '') AS other_headline, COALESCE(p.photo_url, '') AS other_avatar_url
		FROM connections c
		LEFT JOIN connection_requests_meta meta ON c.id = meta.connection_id
		JOIN users other ON (other.id = CASE WHEN c.user_a_id = $1 THEN c.user_b_id ELSE c.user_a_id END)
		LEFT JOIN profiles p ON other.id = p.user_id
		WHERE (c.user_a_id = $1 OR c.user_b_id = $1) AND c.status = 'pending' AND %s
		ORDER BY c.created_at DESC
	`, filter)

	rows, err := exec.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conns := []Connection{}
	for rows.Next() {
		var conn Connection
		var note sql.NullString
		var source sql.NullString
		var respondedAt sql.NullTime
		var other PublicProfileSummary

		err := rows.Scan(
			&conn.ID, &conn.UserAID, &conn.UserBID, &conn.Status, &conn.RequestedBy,
			&conn.CreatedAt, &respondedAt, &conn.UpdatedAt, &note, &source,
			&other.ID, &other.Name, &other.Headline, &other.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		if respondedAt.Valid {
			conn.RespondedAt = &respondedAt.Time
		}
		if note.Valid {
			conn.Note = &note.String
		}
		if source.Valid {
			src := ConnectionSource(source.String)
			conn.Source = &src
		}
		conn.User = &other
		conns = append(conns, conn)
	}

	return conns, rows.Err()
}

// GetMutualConnections intersects accepted connection partner IDs of both users
func (r *Repository) GetMutualConnections(ctx context.Context, userA, userB string, limit int) ([]PublicProfileSummary, int, error) {
	exec := tx.GetExecutor(ctx, r.db)

	const queryCount = `
		SELECT COUNT(*)
		FROM users other
		WHERE other.id IN (
			SELECT CASE WHEN user_a_id = $1 THEN user_b_id ELSE user_a_id END
			FROM connections
			WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'accepted'
		)
		AND other.id IN (
			SELECT CASE WHEN user_a_id = $2 THEN user_b_id ELSE user_a_id END
			FROM connections
			WHERE (user_a_id = $2 OR user_b_id = $2) AND status = 'accepted'
		)
	`
	var totalCount int
	err := exec.QueryRowContext(ctx, queryCount, userA, userB).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	const queryList = `
		SELECT other.id, other.full_name, COALESCE(p.headline, ''), COALESCE(p.photo_url, '')
		FROM users other
		LEFT JOIN profiles p ON other.id = p.user_id
		WHERE other.id IN (
			SELECT CASE WHEN user_a_id = $1 THEN user_b_id ELSE user_a_id END
			FROM connections
			WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'accepted'
		)
		AND other.id IN (
			SELECT CASE WHEN user_a_id = $2 THEN user_b_id ELSE user_a_id END
			FROM connections
			WHERE (user_a_id = $2 OR user_b_id = $2) AND status = 'accepted'
		)
		ORDER BY other.full_name ASC
		LIMIT $3
	`
	rows, err := exec.QueryContext(ctx, queryList, userA, userB, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	summaries := []PublicProfileSummary{}
	for rows.Next() {
		var s PublicProfileSummary
		if err := rows.Scan(&s.ID, &s.Name, &s.Headline, &s.AvatarURL); err != nil {
			return nil, 0, err
		}
		summaries = append(summaries, s)
	}

	return summaries, totalCount, rows.Err()
}

// GetSuggestions recommends connections excluding existing/pending/blocked pairs
func (r *Repository) GetSuggestions(ctx context.Context, userID string, limit int) ([]Suggestion, error) {
	exec := tx.GetExecutor(ctx, r.db)

	const query = `
		SELECT 
			other.id,
			other.full_name,
			COALESCE(p.headline, '') AS headline,
			COALESCE(p.photo_url, '') AS photo_url,
			-- mutual connection count
			(
				SELECT COUNT(*)
				FROM connections c1
				JOIN connections c2 ON (
					(c1.user_a_id = $1 AND c1.user_b_id = CASE WHEN c2.user_a_id = other.id THEN c2.user_b_id ELSE c2.user_a_id END) OR
					(c1.user_b_id = $1 AND c1.user_a_id = CASE WHEN c2.user_a_id = other.id THEN c2.user_b_id ELSE c2.user_a_id END)
				)
				WHERE c1.status = 'accepted' AND c2.status = 'accepted'
				  AND (c2.user_a_id = other.id OR c2.user_b_id = other.id)
			) AS mutual_connection_count,
			-- same industry or location match
			(
				(p.location IS NOT NULL AND p.location = (SELECT location FROM profiles WHERE user_id = $1))
				OR
				EXISTS (
					SELECT 1 FROM profile_desired_industries pdi1
					JOIN profile_desired_industries pdi2 ON pdi1.industry = pdi2.industry
					WHERE pdi1.user_id = $1 AND pdi2.user_id = other.id
				)
			) AS same_match
		FROM users other
		LEFT JOIN profiles p ON other.id = p.user_id
		WHERE other.id != $1
		  AND other.deleted_at IS NULL
		  -- Exclude existing connections or pending/declined/blocked rows
		  AND NOT EXISTS (
			  SELECT 1 FROM connections 
			  WHERE (user_a_id = $1 AND user_b_id = other.id) OR (user_a_id = other.id AND user_b_id = $1)
		  )
		  -- Exclude blocked users (both directions)
		  AND NOT EXISTS (
			  SELECT 1 FROM blocks 
			  WHERE (blocker_id = $1 AND blocked_id = other.id) OR (blocker_id = other.id AND blocked_id = $1)
		  )
		ORDER BY 
			mutual_connection_count DESC,
			same_match DESC,
			p.last_active_at DESC NULLS LAST
		LIMIT $2
	`

	rows, err := exec.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suggestions := []Suggestion{}
	for rows.Next() {
		var s Suggestion
		var sameMatch bool
		err := rows.Scan(
			&s.User.ID, &s.User.Name, &s.User.Headline, &s.User.AvatarURL,
			&s.MutualConnectionCount, &sameMatch,
		)
		if err != nil {
			return nil, err
		}

		// Determine reason
		if s.MutualConnectionCount > 0 {
			s.Reason = fmt.Sprintf("%d mutual connections", s.MutualConnectionCount)
		} else if sameMatch {
			s.Reason = "Similar profile or location"
		} else {
			s.Reason = "Recommended for you"
		}

		suggestions = append(suggestions, s)
	}

	return suggestions, rows.Err()
}
