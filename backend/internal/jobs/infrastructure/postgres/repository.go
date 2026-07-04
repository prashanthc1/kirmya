// Package postgres implements jobs/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"workspace-app/internal/jobs/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const jobCols = `id, title, company, COALESCE(location,''), COALESCE(description,''),
                 COALESCE(salary,''), COALESCE(job_type,''), posted_by, created_at, updated_at`

func scanJob(s interface {
	Scan(...any) error
}) (domain.Job, error) {
	var j domain.Job
	err := s.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Description,
		&j.Salary, &j.JobType, &j.PostedBy, &j.CreatedAt, &j.UpdatedAt)
	return j, err
}

func (r *Repository) CreateJob(ctx context.Context, j *domain.Job) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO jobs (title, company, location, description, salary, job_type, posted_by)
		VALUES ($1,$2,$3,NULLIF($4,''),$5,NULLIF($6,''),$7)
		RETURNING id, created_at, updated_at`,
		j.Title, j.Company, j.Location, j.Description, j.Salary, j.JobType, j.PostedBy).
		Scan(&j.ID, &j.CreatedAt, &j.UpdatedAt)
}

func (r *Repository) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	j, err := scanJob(r.db.QueryRowContext(ctx, `SELECT `+jobCols+` FROM jobs WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *Repository) ListJobs(ctx context.Context, f domain.Filter) ([]domain.Job, error) {
	q := `SELECT ` + jobCols + ` FROM jobs WHERE 1=1`
	args := []any{}
	n := 0
	add := func(v any) string {
		n++
		args = append(args, v)
		return "$" + strconv.Itoa(n)
	}
	if f.Keyword != "" {
		p := add("%" + f.Keyword + "%")
		q += ` AND (title ILIKE ` + p + ` OR company ILIKE ` + p + ` OR description ILIKE ` + p + `)`
	}
	if f.Location != "" {
		q += ` AND location ILIKE ` + add("%"+f.Location+"%")
	}
	if f.JobType != "" {
		q += ` AND job_type = ` + add(f.JobType)
	}
	if f.PostedBy != "" {
		q += ` AND posted_by = ` + add(f.PostedBy)
	}
	q += ` ORDER BY created_at DESC`
	if f.Limit > 0 {
		q += ` LIMIT ` + add(f.Limit)
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Job{}
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

func (r *Repository) UpdateJob(ctx context.Context, j *domain.Job) error {
	return r.db.QueryRowContext(ctx, `
		UPDATE jobs SET title=$2, company=$3, location=NULLIF($4,''), description=NULLIF($5,''),
		    salary=$6, job_type=NULLIF($7,''), updated_at=now()
		WHERE id=$1
		RETURNING updated_at`,
		j.ID, j.Title, j.Company, j.Location, j.Description, j.Salary, j.JobType).Scan(&j.UpdatedAt)
}

func (r *Repository) DeleteJob(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, id)
	return err
}

// --- applications ---

const appCols = `id, job_id, user_id, status, COALESCE(cover_letter,''), created_at, updated_at`

func scanApp(s interface {
	Scan(...any) error
}) (domain.Application, error) {
	var a domain.Application
	err := s.Scan(&a.ID, &a.JobID, &a.UserID, &a.Status, &a.CoverLetter, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *Repository) CreateApplication(ctx context.Context, a *domain.Application) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO job_applications (job_id, user_id, status, cover_letter)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, updated_at`,
		a.JobID, a.UserID, a.Status, a.CoverLetter).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *Repository) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	a, err := scanApp(r.db.QueryRowContext(ctx, `SELECT `+appCols+` FROM job_applications WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrApplicationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) HasApplied(ctx context.Context, jobID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM job_applications WHERE job_id=$1 AND user_id=$2)`, jobID, userID).Scan(&exists)
	return exists, err
}

func (r *Repository) ListApplicationsByUser(ctx context.Context, userID string) ([]domain.Application, error) {
	return r.queryApps(ctx, `SELECT `+appCols+` FROM job_applications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
}

func (r *Repository) ListApplicationsByJob(ctx context.Context, jobID string) ([]domain.Application, error) {
	return r.queryApps(ctx, `SELECT `+appCols+` FROM job_applications WHERE job_id=$1 ORDER BY created_at DESC`, jobID)
}

func (r *Repository) queryApps(ctx context.Context, q, arg string) ([]domain.Application, error) {
	rows, err := r.db.QueryContext(ctx, q, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Application{}
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repository) UpdateApplicationStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE job_applications SET status=$2, updated_at=now() WHERE id=$1`, id, status)
	return err
}

// --- saved jobs ---

func (r *Repository) SaveJob(ctx context.Context, userID, jobID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO saved_jobs (user_id, job_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, jobID)
	return err
}

func (r *Repository) UnsaveJob(ctx context.Context, userID, jobID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM saved_jobs WHERE user_id=$1 AND job_id=$2`, userID, jobID)
	return err
}

func (r *Repository) IsSaved(ctx context.Context, userID, jobID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM saved_jobs WHERE user_id=$1 AND job_id=$2)`, userID, jobID).Scan(&exists)
	return exists, err
}

func (r *Repository) ListSavedJobs(ctx context.Context, userID string) ([]domain.Job, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+jobCols+` FROM jobs j
		JOIN saved_jobs s ON s.job_id = j.id
		WHERE s.user_id = $1 ORDER BY s.saved_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Job{}
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}
