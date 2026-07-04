// Package postgres implements admin/domain.Repository on PostgreSQL. It reads
// across context-owned tables (users, profiles, posts, comments) and owns the
// content_reports + audit_logs writes.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"workspace-app/internal/admin/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

// ----- Users -------------------------------------------------------------

// userCols selects the admin user projection. Roles are aggregated into a single
// comma-separated string (string_agg) to avoid driver-specific array scanning.
const userCols = `u.id, u.email, u.full_name, COALESCE(p.headline,''), u.status,
                  u.email_verified, u.mfa_enabled, u.last_login_at, u.created_at,
                  COALESCE((SELECT string_agg(ro.name, ',' ORDER BY ro.name)
                            FROM user_roles ur JOIN roles ro ON ro.id = ur.role_id
                            WHERE ur.user_id = u.id), '')`

func scanUser(s interface{ Scan(...any) error }) (domain.UserSummary, error) {
	var u domain.UserSummary
	var lastLogin sql.NullTime
	var roles string
	err := s.Scan(&u.ID, &u.Email, &u.FullName, &u.Headline, &u.Status,
		&u.EmailVerified, &u.MFAEnabled, &lastLogin, &u.CreatedAt, &roles)
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	u.Roles = splitRoles(roles)
	return u, err
}

// splitRoles turns the string_agg output into a non-nil slice.
func splitRoles(agg string) []string {
	if agg == "" {
		return []string{}
	}
	return strings.Split(agg, ",")
}

func (r *Repository) ListUsers(ctx context.Context, f domain.UserFilter) ([]domain.UserSummary, int, error) {
	var where []string
	var args []any
	where = append(where, "u.deleted_at IS NULL")

	if f.Query != "" {
		args = append(args, "%"+f.Query+"%")
		where = append(where, fmt.Sprintf("(u.full_name ILIKE $%d OR u.email ILIKE $%d)", len(args), len(args)))
	}
	if f.Status != "" {
		args = append(args, f.Status)
		where = append(where, fmt.Sprintf("u.status = $%d", len(args)))
	}
	if f.Role != "" {
		args = append(args, f.Role)
		where = append(where, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM user_roles ur JOIN roles ro ON ro.id = ur.role_id WHERE ur.user_id = u.id AND ro.name = $%d)",
			len(args)))
	}
	whereSQL := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM users u WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, f.Offset)
	q := `SELECT ` + userCols + `
		FROM users u LEFT JOIN profiles p ON p.user_id = u.id
		WHERE ` + whereSQL + `
		ORDER BY u.created_at DESC
		LIMIT $` + fmt.Sprint(len(args)-1) + ` OFFSET $` + fmt.Sprint(len(args))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []domain.UserSummary{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repository) GetUser(ctx context.Context, id string) (*domain.UserSummary, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx, `SELECT `+userCols+`
		FROM users u LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.id = $1 AND u.deleted_at IS NULL`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) SetUserStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET status = $2, updated_at = now(), version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`, id, status)
	return err
}

func (r *Repository) AssignRole(ctx context.Context, userID, role string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
		ON CONFLICT DO NOTHING`, userID, role)
	return err
}

func (r *Repository) RevokeRole(ctx context.Context, userID, role string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = (SELECT id FROM roles WHERE name = $2)`, userID, role)
	return err
}

// ----- Moderation --------------------------------------------------------

func (r *Repository) DeletePost(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = $1`, id)
	return affected(res, err)
}

func (r *Repository) DeleteComment(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)
	return affected(res, err)
}

func affected(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ----- Reports -----------------------------------------------------------

const reportCols = `id, reporter_id, target_type, target_id, reason, status,
                    COALESCE(action_taken,''), COALESCE(resolved_by::text,''),
                    resolved_at, created_at, updated_at`

func scanReport(s interface{ Scan(...any) error }) (domain.Report, error) {
	var r domain.Report
	var resolvedAt sql.NullTime
	err := s.Scan(&r.ID, &r.ReporterID, &r.TargetType, &r.TargetID, &r.Reason, &r.Status,
		&r.ActionTaken, &r.ResolvedBy, &resolvedAt, &r.CreatedAt, &r.UpdatedAt)
	if resolvedAt.Valid {
		r.ResolvedAt = &resolvedAt.Time
	}
	return r, err
}

func (r *Repository) CreateReport(ctx context.Context, rep *domain.Report) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO content_reports (reporter_id, target_type, target_id, reason, status)
		VALUES ($1, $2, $3::uuid, $4, $5)
		RETURNING id, created_at, updated_at`,
		rep.ReporterID, rep.TargetType, rep.TargetID, rep.Reason, rep.Status).
		Scan(&rep.ID, &rep.CreatedAt, &rep.UpdatedAt)
}

func (r *Repository) ListReports(ctx context.Context, status string) ([]domain.Report, error) {
	q := `SELECT ` + reportCols + ` FROM content_reports`
	var args []any
	if status != "" {
		q += ` WHERE status = $1`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Report{}
	for rows.Next() {
		rep, err := scanReport(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rep)
	}
	return out, rows.Err()
}

func (r *Repository) GetReport(ctx context.Context, id string) (*domain.Report, error) {
	rep, err := scanReport(r.db.QueryRowContext(ctx,
		`SELECT `+reportCols+` FROM content_reports WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *Repository) ResolveReport(ctx context.Context, id, status, actionTaken, resolvedBy string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE content_reports
		SET status = $2, action_taken = NULLIF($3,''), resolved_by = $4::uuid,
		    resolved_at = now(), updated_at = now()
		WHERE id = $1`, id, status, actionTaken, resolvedBy)
	return affected(res, err)
}

// ----- Analytics + audit -------------------------------------------------

func (r *Repository) Analytics(ctx context.Context) (*domain.Analytics, error) {
	var a domain.Analytics
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL),
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL AND status = 'active'),
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL AND status = 'suspended'),
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL AND status = 'deactivated'),
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL AND created_at >= now() - interval '7 days'),
		  (SELECT count(*) FROM users WHERE deleted_at IS NULL AND created_at >= now() - interval '30 days'),
		  (SELECT count(*) FROM jobs),
		  (SELECT count(*) FROM job_applications),
		  (SELECT count(*) FROM referral_requests WHERE deleted_at IS NULL),
		  (SELECT count(*) FROM referral_requests WHERE deleted_at IS NULL AND status = 'accepted'),
		  (SELECT count(*) FROM referral_requests WHERE deleted_at IS NULL AND outcome = 'hired'),
		  (SELECT count(*) FROM communities),
		  (SELECT count(*) FROM posts),
		  (SELECT count(*) FROM comments),
		  (SELECT count(*) FROM content_reports WHERE status IN ('open','reviewing'))`).
		Scan(
			&a.Users.Total, &a.Users.Active, &a.Users.Suspended, &a.Users.Deactivated,
			&a.Users.New7d, &a.Users.New30d,
			&a.Jobs.Total, &a.Jobs.Applications,
			&a.Referrals.Total, &a.Referrals.Accepted, &a.Referrals.Hired,
			&a.Communities.Total, &a.Communities.Posts, &a.Communities.Comments,
			&a.Reports.Open,
		)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) WriteAudit(ctx context.Context, actorID, action, targetType, targetID string, metadata map[string]any) error {
	meta := []byte("{}")
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = b
		}
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1::uuid, $2, $3, NULLIF($4,'')::uuid, $5)`,
		actorID, action, targetType, targetID, meta)
	return err
}
