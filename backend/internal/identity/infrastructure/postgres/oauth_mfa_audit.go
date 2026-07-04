package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"workspace-app/internal/identity/domain"
)

// OAuthRepository implements domain.OAuthRepository.
type OAuthRepository struct{ db *sql.DB }

func NewOAuthRepository(db *sql.DB) *OAuthRepository { return &OAuthRepository{db: db} }

func (r *OAuthRepository) FindUserIDByProvider(ctx context.Context, provider, providerUID string) (string, bool, error) {
	const q = `SELECT user_id FROM oauth_accounts WHERE provider = $1 AND provider_uid = $2`
	var userID string
	err := r.db.QueryRowContext(ctx, q, provider, providerUID).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return userID, true, nil
}

func (r *OAuthRepository) Link(ctx context.Context, userID, provider, providerUID string) error {
	const q = `
		INSERT INTO oauth_accounts (user_id, provider, provider_uid)
		VALUES ($1, $2, $3)
		ON CONFLICT (provider, provider_uid) DO NOTHING`
	_, err := r.db.ExecContext(ctx, q, userID, provider, providerUID)
	return err
}

// MFARepository implements domain.MFARepository.
type MFARepository struct{ db *sql.DB }

func NewMFARepository(db *sql.DB) *MFARepository { return &MFARepository{db: db} }

func (r *MFARepository) Upsert(ctx context.Context, c *domain.MFACredential) error {
	const q = `
		INSERT INTO mfa_credentials (user_id, type, secret_enc)
		VALUES ($1, 'totp', $2)
		ON CONFLICT (id) DO NOTHING`
	// One credential per user in MVP: replace any unconfirmed prior secret.
	if _, err := r.db.ExecContext(ctx, `DELETE FROM mfa_credentials WHERE user_id = $1 AND confirmed_at IS NULL`, c.UserID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, q, c.UserID, c.SecretEnc)
	return err
}

func (r *MFARepository) Get(ctx context.Context, userID string) (*domain.MFACredential, error) {
	const q = `
		SELECT user_id, secret_enc, (confirmed_at IS NOT NULL)
		FROM mfa_credentials WHERE user_id = $1
		ORDER BY created_at DESC LIMIT 1`
	var c domain.MFACredential
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&c.UserID, &c.SecretEnc, &c.Confirmed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MFARepository) Confirm(ctx context.Context, userID string) error {
	const q = `UPDATE mfa_credentials SET confirmed_at = now() WHERE user_id = $1 AND confirmed_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

// AuditRepository implements domain.AuditRepository.
type AuditRepository struct{ db *sql.DB }

func NewAuditRepository(db *sql.DB) *AuditRepository { return &AuditRepository{db: db} }

func (r *AuditRepository) Record(ctx context.Context, actorID, action, targetType, targetID string, metadata map[string]any, ip string) error {
	meta, err := json.Marshal(metadata)
	if err != nil {
		meta = []byte("{}")
	}
	const q = `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata, ip)
		VALUES (NULLIF($1,'')::uuid, $2, NULLIF($3,''), NULLIF($4,'')::uuid, $5::jsonb, NULLIF($6,''))`
	_, err = r.db.ExecContext(ctx, q, actorID, action, targetType, targetID, string(meta), ip)
	return err
}
