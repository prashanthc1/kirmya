// Package postgres implements resume/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"workspace-app/internal/resume/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateResume(ctx context.Context, res *domain.Resume) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO resumes (user_id, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		res.UserID, res.Title).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

func (r *Repository) GetResume(ctx context.Context, id string) (*domain.Resume, error) {
	var res domain.Resume
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM resumes WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&res.ID, &res.UserID, &res.Title, &res.CreatedAt, &res.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *Repository) ListResumesByUser(ctx context.Context, userID string) ([]domain.Resume, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM resumes
		 WHERE user_id = $1 AND deleted_at IS NULL ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Resume{}
	for rows.Next() {
		var res domain.Resume
		if err := rows.Scan(&res.ID, &res.UserID, &res.Title, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, rows.Err()
}

func (r *Repository) SoftDeleteResume(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE resumes SET deleted_at = now() WHERE id = $1`, id)
	return err
}

func (r *Repository) NextVersionNo(ctx context.Context, resumeID string) (int, error) {
	var n sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT MAX(version_no) FROM resume_versions WHERE resume_id = $1`, resumeID).Scan(&n)
	if err != nil {
		return 0, err
	}
	return int(n.Int64) + 1, nil
}

func (r *Repository) AddVersion(ctx context.Context, v *domain.Version) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO resume_versions (resume_id, version_no, filename, content_type, size_bytes, storage_key, extracted_text)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at`,
		v.ResumeID, v.VersionNo, v.Filename, v.ContentType, v.SizeBytes, v.StorageKey, v.ExtractedText).
		Scan(&v.ID, &v.CreatedAt)
}

const versionCols = `id, resume_id, version_no, filename, content_type, size_bytes, storage_key, extracted_text, created_at`

func scanVersion(s interface{ Scan(...any) error }) (domain.Version, error) {
	var v domain.Version
	err := s.Scan(&v.ID, &v.ResumeID, &v.VersionNo, &v.Filename, &v.ContentType,
		&v.SizeBytes, &v.StorageKey, &v.ExtractedText, &v.CreatedAt)
	return v, err
}

func (r *Repository) ListVersions(ctx context.Context, resumeID string) ([]domain.Version, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+versionCols+` FROM resume_versions WHERE resume_id = $1 ORDER BY version_no DESC`, resumeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Version{}
	for rows.Next() {
		v, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *Repository) LatestVersion(ctx context.Context, resumeID string) (*domain.Version, error) {
	v, err := scanVersion(r.db.QueryRowContext(ctx,
		`SELECT `+versionCols+` FROM resume_versions WHERE resume_id = $1 ORDER BY version_no DESC LIMIT 1`, resumeID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) SaveScore(ctx context.Context, s *domain.Score) error {
	sugg, err := json.Marshal(s.Suggestions)
	if err != nil {
		sugg = []byte("[]")
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO resume_scores (version_id, overall, formatting, keywords, ats, suggestions)
		VALUES ($1,$2,$3,$4,$5,$6::jsonb)
		ON CONFLICT (version_id) DO UPDATE
		  SET overall=EXCLUDED.overall, formatting=EXCLUDED.formatting, keywords=EXCLUDED.keywords,
		      ats=EXCLUDED.ats, suggestions=EXCLUDED.suggestions
		RETURNING created_at`,
		s.VersionID, s.Overall, s.Formatting, s.Keywords, s.ATS, string(sugg)).Scan(&s.CreatedAt)
}

func (r *Repository) LatestScore(ctx context.Context, resumeID string) (*domain.Score, error) {
	var s domain.Score
	var sugg []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT sc.version_id, sc.overall, sc.formatting, sc.keywords, sc.ats, sc.suggestions, sc.created_at
		FROM resume_scores sc
		JOIN resume_versions v ON v.id = sc.version_id
		WHERE v.resume_id = $1
		ORDER BY v.version_no DESC LIMIT 1`, resumeID).
		Scan(&s.VersionID, &s.Overall, &s.Formatting, &s.Keywords, &s.ATS, &sugg, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(sugg, &s.Suggestions)
	if s.Suggestions == nil {
		s.Suggestions = []string{}
	}
	return &s, nil
}
