package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"workspace-app/internal/identity/domain"
	"workspace-app/internal/platform/tx"
)

// RefreshTokenRepository implements domain.RefreshTokenRepository.
type RefreshTokenRepository struct{ db tx.DBTX }

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: tx.NewTxDB(db)}
}

func (r *RefreshTokenRepository) Store(ctx context.Context, t *domain.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, family_id, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	return r.db.QueryRowContext(ctx, q, t.UserID, t.TokenHash, t.FamilyID, t.ExpiresAt).Scan(&t.ID)
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, family_id, expires_at, revoked_at, replaced_by
		FROM refresh_tokens WHERE token_hash = $1`
	var (
		t          domain.RefreshToken
		revokedAt  sql.NullTime
		replacedBy sql.NullString
	)
	err := r.db.QueryRowContext(ctx, q, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.FamilyID, &t.ExpiresAt, &revokedAt, &replacedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	if revokedAt.Valid {
		t.RevokedAt = &revokedAt.Time
	}
	if replacedBy.Valid {
		t.ReplacedBy = &replacedBy.String
	}
	return &t, nil
}

func (r *RefreshTokenRepository) MarkReplaced(ctx context.Context, id, replacedBy string) error {
	const q = `UPDATE refresh_tokens SET replaced_by = $2, revoked_at = now() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id, replacedBy)
	return err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	const q = `UPDATE refresh_tokens SET revoked_at = now() WHERE id = $1 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *RefreshTokenRepository) RevokeFamily(ctx context.Context, familyID string) error {
	const q = `UPDATE refresh_tokens SET revoked_at = now() WHERE family_id = $1 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, familyID)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	const q = `UPDATE refresh_tokens SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

// VerificationRepository implements domain.VerificationRepository.
type VerificationRepository struct{ db tx.DBTX }

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: tx.NewTxDB(db)}
}

func (r *VerificationRepository) StoreEmailToken(ctx context.Context, userID, hash string, expiresAt time.Time) error {
	const q = `INSERT INTO email_verification_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, userID, hash, expiresAt)
	return err
}

func (r *VerificationRepository) ConsumeEmailToken(ctx context.Context, hash string) (string, error) {
	const q = `
		UPDATE email_verification_tokens SET used_at = now()
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
		RETURNING user_id`
	return r.consume(ctx, q, hash)
}

func (r *VerificationRepository) StorePasswordToken(ctx context.Context, userID, hash string, expiresAt time.Time) error {
	const q = `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, userID, hash, expiresAt)
	return err
}

func (r *VerificationRepository) ConsumePasswordToken(ctx context.Context, hash string) (string, error) {
	const q = `
		UPDATE password_reset_tokens SET used_at = now()
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
		RETURNING user_id`
	return r.consume(ctx, q, hash)
}

func (r *VerificationRepository) consume(ctx context.Context, q, hash string) (string, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, q, hash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrTokenNotFound
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}
