// Package postgres contains PostgreSQL adapters implementing the identity
// domain ports.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"workspace-app/internal/identity/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	const q = `
		INSERT INTO users (email, password_hash, full_name, email_verified, status)
		VALUES ($1, NULLIF($2,''), $3, $4, $5)
		RETURNING id, created_at, updated_at, version`
	err := r.db.QueryRowContext(ctx, q, u.Email, u.PasswordHash, u.FullName, u.EmailVerified, statusOrDefault(u.Status)).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Version)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailTaken
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.getOne(ctx, "id = $1", id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.getOne(ctx, "lower(email) = lower($1)", email)
}

func (r *UserRepository) getOne(ctx context.Context, where, arg string) (*domain.User, error) {
	q := `
		SELECT id, email, COALESCE(password_hash,''), full_name, email_verified, status,
		       mfa_enabled, last_login_at, created_at, updated_at, version
		FROM users
		WHERE ` + where + ` AND deleted_at IS NULL`
	var (
		u           domain.User
		lastLogin   sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, q, arg).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.EmailVerified, &u.Status,
		&u.MFAEnabled, &lastLogin, &u.CreatedAt, &u.UpdatedAt, &u.Version,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	roles, err := r.GetRoles(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	const q = `
		UPDATE users
		SET full_name = $2, status = $3, email_verified = $4, mfa_enabled = $5,
		    updated_at = now(), version = version + 1
		WHERE id = $1 AND version = $6 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, u.ID, u.FullName, statusOrDefault(u.Status), u.EmailVerified, u.MFAEnabled, u.Version)
	if err != nil {
		return err
	}
	if err := ensureRowAffected(res); err != nil {
		return err
	}
	u.Version++
	return nil
}

func (r *UserRepository) AssignRole(ctx context.Context, userID, roleName string) error {
	const q = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
		ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, q, userID, roleName)
	return err
}

func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleName string) error {
	const q = `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = (SELECT id FROM roles WHERE name = $2)`
	_, err := r.db.ExecContext(ctx, q, userID, roleName)
	return err
}

func (r *UserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	const q = `
		SELECT r.name FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	return roles, rows.Err()
}

func (r *UserRepository) SetEmailVerified(ctx context.Context, userID string) error {
	return r.simpleUpdate(ctx, `UPDATE users SET email_verified = true, updated_at = now(), version = version + 1 WHERE id = $1`, userID)
}

func (r *UserRepository) SetPasswordHash(ctx context.Context, userID, hash string) error {
	const q = `UPDATE users SET password_hash = $2, updated_at = now(), version = version + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, userID, hash)
	return err
}

func (r *UserRepository) SetMFAEnabled(ctx context.Context, userID string, enabled bool) error {
	const q = `UPDATE users SET mfa_enabled = $2, updated_at = now(), version = version + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, userID, enabled)
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	return r.simpleUpdate(ctx, `UPDATE users SET last_login_at = now() WHERE id = $1`, userID)
}

func (r *UserRepository) simpleUpdate(ctx context.Context, q, userID string) error {
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

const directoryCols = `u.id, u.full_name, u.email, COALESCE(p.headline,''), COALESCE(p.photo_url,'')`

func scanDirectory(s interface{ Scan(...any) error }) (domain.DirectoryEntry, error) {
	var d domain.DirectoryEntry
	err := s.Scan(&d.ID, &d.FullName, &d.Email, &d.Headline, &d.PhotoURL)
	return d, err
}

func (r *UserRepository) Search(ctx context.Context, query string, limit int) ([]domain.DirectoryEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+directoryCols+`
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
		LEFT JOIN user_settings us ON us.user_id = u.id
		WHERE u.deleted_at IS NULL
		  AND COALESCE(us.discoverable, true)
		  AND (u.full_name ILIKE $1 OR u.email ILIKE $1)
		ORDER BY u.full_name NULLS LAST, u.email
		LIMIT $2`, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.DirectoryEntry{}
	for rows.Next() {
		d, err := scanDirectory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *UserRepository) GetDirectory(ctx context.Context, id string) (*domain.DirectoryEntry, error) {
	d, err := scanDirectory(r.db.QueryRowContext(ctx, `
		SELECT `+directoryCols+`
		FROM users u LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.id = $1 AND u.deleted_at IS NULL`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func statusOrDefault(s string) string {
	if strings.TrimSpace(s) == "" {
		return domain.StatusActive
	}
	return s
}
