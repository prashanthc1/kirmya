// Package postgres implements profile/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"

	"workspace-app/internal/profile/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const dateLayout = "2006-01-02"

func dateStr(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format(dateLayout)
}

func (r *Repository) ensureRow(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO profiles (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`, userID)
	return err
}

func (r *Repository) Get(ctx context.Context, userID string) (*domain.Profile, error) {
	if err := r.ensureRow(ctx, userID); err != nil {
		return nil, err
	}
	p := &domain.Profile{UserID: userID}
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(headline,''), COALESCE(about,''), COALESCE(photo_url,''),
		       COALESCE(bio,''), COALESCE(location,''), COALESCE(website,''), version
		FROM profiles WHERE user_id = $1 AND deleted_at IS NULL`, userID).
		Scan(&p.Headline, &p.About, &p.PhotoURL, &p.Bio, &p.Location, &p.Website, &p.Version)
	if err != nil {
		return nil, err
	}
	if p.Experiences, err = r.loadExperiences(ctx, userID); err != nil {
		return nil, err
	}
	if p.Educations, err = r.loadEducations(ctx, userID); err != nil {
		return nil, err
	}
	if p.Certifications, err = r.loadCertifications(ctx, userID); err != nil {
		return nil, err
	}
	if p.Skills, err = r.loadSkills(ctx, userID); err != nil {
		return nil, err
	}
	if p.Languages, err = r.loadLanguages(ctx, userID); err != nil {
		return nil, err
	}
	if p.Portfolio, err = r.loadPortfolio(ctx, userID); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) UpdateScalars(ctx context.Context, userID string, s domain.Scalars) error {
	if err := r.ensureRow(ctx, userID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE profiles
		SET headline = $2, about = $3, photo_url = $4, bio = $5, location = $6, website = $7,
		    updated_at = now(), version = version + 1
		WHERE user_id = $1 AND deleted_at IS NULL`,
		userID, s.Headline, s.About, s.PhotoURL, s.Bio, s.Location, s.Website)
	return err
}

// --- experiences ---

func (r *Repository) loadExperiences(ctx context.Context, userID string) ([]domain.WorkExperience, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, company, COALESCE(location,''), COALESCE(employment_type,''),
		       start_date, end_date, is_current, COALESCE(description,'')
		FROM work_experiences WHERE user_id = $1
		ORDER BY is_current DESC, start_date DESC NULLS LAST`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.WorkExperience{}
	for rows.Next() {
		var e domain.WorkExperience
		var start, end sql.NullTime
		if err := rows.Scan(&e.ID, &e.Title, &e.Company, &e.Location, &e.EmploymentType,
			&start, &end, &e.IsCurrent, &e.Description); err != nil {
			return nil, err
		}
		e.StartDate, e.EndDate = dateStr(start), dateStr(end)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *Repository) AddExperience(ctx context.Context, userID string, e *domain.WorkExperience) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO work_experiences (user_id, title, company, location, employment_type, start_date, end_date, is_current, description)
		VALUES ($1,$2,$3,$4,$5,NULLIF($6,'')::date,NULLIF($7,'')::date,$8,$9)
		RETURNING id`,
		userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description).
		Scan(&e.ID)
}

func (r *Repository) UpdateExperience(ctx context.Context, userID string, e domain.WorkExperience) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE work_experiences
		SET title=$3, company=$4, location=$5, employment_type=$6,
		    start_date=NULLIF($7,'')::date, end_date=NULLIF($8,'')::date, is_current=$9, description=$10, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description)
	return owned(res, err)
}

func (r *Repository) DeleteExperience(ctx context.Context, userID, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM work_experiences WHERE id=$1 AND user_id=$2`, id, userID)
	return owned(res, err)
}

// --- educations ---

func (r *Repository) loadEducations(ctx context.Context, userID string) ([]domain.Education, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, school, COALESCE(degree,''), COALESCE(field_of_study,''),
		       start_date, end_date, COALESCE(grade,''), COALESCE(description,'')
		FROM educations WHERE user_id = $1 ORDER BY start_date DESC NULLS LAST`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Education{}
	for rows.Next() {
		var e domain.Education
		var start, end sql.NullTime
		if err := rows.Scan(&e.ID, &e.School, &e.Degree, &e.FieldOfStudy, &start, &end, &e.Grade, &e.Description); err != nil {
			return nil, err
		}
		e.StartDate, e.EndDate = dateStr(start), dateStr(end)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *Repository) AddEducation(ctx context.Context, userID string, e *domain.Education) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO educations (user_id, school, degree, field_of_study, start_date, end_date, grade, description)
		VALUES ($1,$2,$3,$4,NULLIF($5,'')::date,NULLIF($6,'')::date,$7,$8)
		RETURNING id`,
		userID, e.School, e.Degree, e.FieldOfStudy, e.StartDate, e.EndDate, e.Grade, e.Description).Scan(&e.ID)
}

func (r *Repository) UpdateEducation(ctx context.Context, userID string, e domain.Education) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE educations
		SET school=$3, degree=$4, field_of_study=$5, start_date=NULLIF($6,'')::date,
		    end_date=NULLIF($7,'')::date, grade=$8, description=$9, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.School, e.Degree, e.FieldOfStudy, e.StartDate, e.EndDate, e.Grade, e.Description)
	return owned(res, err)
}

func (r *Repository) DeleteEducation(ctx context.Context, userID, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM educations WHERE id=$1 AND user_id=$2`, id, userID)
	return owned(res, err)
}

// --- certifications ---

func (r *Repository) loadCertifications(ctx context.Context, userID string) ([]domain.Certification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(issuer,''), issue_date, expiry_date,
		       COALESCE(credential_id,''), COALESCE(credential_url,'')
		FROM certifications WHERE user_id = $1 ORDER BY issue_date DESC NULLS LAST`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Certification{}
	for rows.Next() {
		var c domain.Certification
		var issue, expiry sql.NullTime
		if err := rows.Scan(&c.ID, &c.Name, &c.Issuer, &issue, &expiry, &c.CredentialID, &c.CredentialURL); err != nil {
			return nil, err
		}
		c.IssueDate, c.ExpiryDate = dateStr(issue), dateStr(expiry)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) AddCertification(ctx context.Context, userID string, c *domain.Certification) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO certifications (user_id, name, issuer, issue_date, expiry_date, credential_id, credential_url)
		VALUES ($1,$2,$3,NULLIF($4,'')::date,NULLIF($5,'')::date,$6,$7)
		RETURNING id`,
		userID, c.Name, c.Issuer, c.IssueDate, c.ExpiryDate, c.CredentialID, c.CredentialURL).Scan(&c.ID)
}

func (r *Repository) UpdateCertification(ctx context.Context, userID string, c domain.Certification) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE certifications
		SET name=$3, issuer=$4, issue_date=NULLIF($5,'')::date, expiry_date=NULLIF($6,'')::date,
		    credential_id=$7, credential_url=$8, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		c.ID, userID, c.Name, c.Issuer, c.IssueDate, c.ExpiryDate, c.CredentialID, c.CredentialURL)
	return owned(res, err)
}

func (r *Repository) DeleteCertification(ctx context.Context, userID, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM certifications WHERE id=$1 AND user_id=$2`, id, userID)
	return owned(res, err)
}

// --- skills ---

func (r *Repository) loadSkills(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.name FROM profile_skills ps JOIN skills s ON s.id = ps.skill_id
		WHERE ps.user_id = $1 ORDER BY s.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

func (r *Repository) SetSkills(ctx context.Context, userID string, skills []string) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_skills WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, name := range skills {
			if name == "" {
				continue
			}
			var skillID string
			if err := tx.QueryRowContext(ctx,
				`INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				name).Scan(&skillID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO profile_skills (user_id, skill_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
				userID, skillID); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- languages ---

func (r *Repository) loadLanguages(ctx context.Context, userID string) ([]domain.Language, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.name, pl.proficiency FROM profile_languages pl JOIN languages l ON l.id = pl.language_id
		WHERE pl.user_id = $1 ORDER BY l.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Language{}
	for rows.Next() {
		var l domain.Language
		if err := rows.Scan(&l.Name, &l.Proficiency); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *Repository) SetLanguages(ctx context.Context, userID string, langs []domain.Language) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_languages WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, l := range langs {
			if l.Name == "" {
				continue
			}
			prof := l.Proficiency
			if prof == "" {
				prof = "professional"
			}
			var langID string
			if err := tx.QueryRowContext(ctx,
				`INSERT INTO languages (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				l.Name).Scan(&langID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO profile_languages (user_id, language_id, proficiency) VALUES ($1,$2,$3)
				 ON CONFLICT (user_id, language_id) DO UPDATE SET proficiency = EXCLUDED.proficiency`,
				userID, langID, prof); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- portfolio ---

func (r *Repository) loadPortfolio(ctx context.Context, userID string) ([]domain.PortfolioLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, COALESCE(label,''), url FROM portfolio_links WHERE user_id = $1 ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.PortfolioLink{}
	for rows.Next() {
		var l domain.PortfolioLink
		if err := rows.Scan(&l.ID, &l.Label, &l.URL); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *Repository) SetPortfolio(ctx context.Context, userID string, links []domain.PortfolioLink) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM portfolio_links WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, l := range links {
			if l.URL == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO portfolio_links (user_id, label, url) VALUES ($1,$2,$3)`, userID, l.Label, l.URL); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- helpers ---

func (r *Repository) tx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func owned(res sql.Result, err error) error {
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
