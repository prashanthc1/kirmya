package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"workspace-app/internal/profile/domain"
)

type querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *Repository) UpdateAggregate(ctx context.Context, userID string, expectedVersion int, u domain.AggregateUpdate) error {
	if err := r.ensureRow(ctx, userID); err != nil {
		return err
	}
	return r.tx(ctx, func(tx *sql.Tx) error {
		if expectedVersion > 0 {
			var current int
			err := tx.QueryRowContext(ctx,
				`SELECT version FROM profiles WHERE user_id = $1 AND deleted_at IS NULL FOR UPDATE`, userID).
				Scan(&current)
			if err == sql.ErrNoRows {
				return domain.ErrNotFound
			}
			if err != nil {
				return err
			}
			if current != expectedVersion {
				return domain.ErrOptimisticLock
			}
		}

		// 1. Update identity scalars if present
		if u.Identity != nil {
			id := u.Identity
			_, err := tx.ExecContext(ctx, `
				UPDATE profiles
				SET preferred_name = $2, timezone = $3, nationality = $4,
				    headline = $5, bio = $6, photo_url = $7, cover_url = $8,
				    work_auth_status = $9, passport_nationality = $10,
				    preferred_contact_channel = $11, video_intro_url = $12,
				    location = $13
				WHERE user_id = $1`,
				userID, id.PreferredName, id.TimeZone, id.Nationality,
				id.Headline, id.Bio, id.PhotoURL, id.CoverURL,
				id.WorkAuthorization, id.Nationality,
				id.PreferredContactChannel, id.CoverURL, id.Location)
			if err != nil {
				return err
			}

			// Update languages
			if err := setLanguagesQ(ctx, tx, userID, id.Languages); err != nil {
				return err
			}
		}

		// 2. Update summary scalars if present
		if u.Summary != nil {
			sum := u.Summary
			_, err := tx.ExecContext(ctx, `
				UPDATE profiles
				SET executive_summary = $2, career_objectives = $3,
				    personal_brand_statement = $4, elevator_pitch = $5,
				    career_highlights = $6, functional_areas = $7
				WHERE user_id = $1`,
				userID, sum.ExecutiveSummary, sum.CareerObjectives,
				sum.PersonalBrandStatement, sum.ElevatorPitch,
				sum.CareerHighlights, sum.FunctionalAreas)
			if err != nil {
				return err
			}
		}

		// 3. Update preferences if present
		if u.Preferences != nil {
			pref := u.Preferences
			salMinEnc, _ := r.crypt.Encrypt(strconv.Itoa(pref.SalaryMin))
			salMaxEnc, _ := r.crypt.Encrypt(strconv.Itoa(pref.SalaryMax))
			salCurrEnc, _ := r.crypt.Encrypt(pref.SalaryCurrency)

			_, err := tx.ExecContext(ctx, `
				UPDATE profiles
				SET open_to_relocation = $2, notice_period = $3, work_mode = $4,
				    salary_min_enc = $5, salary_max_enc = $6, salary_currency_enc = $7,
				    travel_willingness = $8, company_size_preferences = $9,
				    preferred_countries = $10, preferred_cities = $11
				WHERE user_id = $1`,
				userID, pref.OpenToRelocation, pref.NoticePeriod, pref.RemotePreference,
				salMinEnc, salMaxEnc, salCurrEnc,
				pref.TravelWillingness, pref.CompanySizePreferences,
				pref.PreferredCountries, pref.PreferredCities)
			if err != nil {
				return err
			}

			if err := rebuildStringTable(ctx, tx, `profile_desired_roles`, `role`, userID, pref.DesiredRoles); err != nil {
				return err
			}
			if err := rebuildStringTable(ctx, tx, `profile_desired_industries`, `industry`, userID, pref.DesiredIndustries); err != nil {
				return err
			}
		}

		// 4. Update privacy settings
		if u.Privacy != nil {
			priv := u.Privacy
			visProfile := priv.FieldVisibility["profile"]
			visSalary := priv.FieldVisibility["salary"]
			visExp := priv.FieldVisibility["experience"]
			visEdu := priv.FieldVisibility["education"]
			visCert := priv.FieldVisibility["certifications"]
			visSkills := priv.FieldVisibility["skills"]
			visPortfolio := priv.FieldVisibility["portfolio"]

			_, err := tx.ExecContext(ctx, `
				UPDATE profiles
				SET visibility_profile = $2, visibility_salary = $3,
				    visibility_experience = $4, visibility_education = $5,
				    visibility_certifications = $6, visibility_skills = $7,
				    visibility_portfolio = $8
				WHERE user_id = $1`,
				userID, visProfile, visSalary, visExp, visEdu, visCert, visSkills, visPortfolio)
			if err != nil {
				return err
			}
		}

		// 5. Update draft flag
		if u.IsDraft != nil {
			_, err := tx.ExecContext(ctx, `UPDATE profiles SET is_draft = $2 WHERE user_id = $1`, userID, *u.IsDraft)
			if err != nil {
				return err
			}
		}

		// Reconcile child collections
		if u.Experiences != nil {
			if err := reconcileExperiences(ctx, tx, userID, *u.Experiences); err != nil {
				return err
			}
		}
		if u.Educations != nil {
			if err := reconcileEducations(ctx, tx, userID, *u.Educations); err != nil {
				return err
			}
		}
		if u.Skills != nil {
			if err := setSkillsQ(ctx, tx, userID, *u.Skills); err != nil {
				return err
			}
		}
		if u.Projects != nil {
			if err := reconcileProjects(ctx, tx, userID, *u.Projects); err != nil {
				return err
			}
		}
		if u.Certifications != nil {
			if err := reconcileCertifications(ctx, tx, userID, *u.Certifications); err != nil {
				return err
			}
		}
		if u.Achievements != nil {
			if err := reconcileAchievements(ctx, tx, userID, *u.Achievements); err != nil {
				return err
			}
		}

		// Bump version
		_, err := tx.ExecContext(ctx, `UPDATE profiles SET version = version + 1, updated_at = now() WHERE user_id = $1`, userID)
		return err
	})
}

// Helpers
func rebuildStringTable(ctx context.Context, q querier, table, col, userID string, values []string) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM `+table+` WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, err := q.ExecContext(ctx, `INSERT INTO `+table+` (user_id, `+col+`) VALUES ($1,$2)`, userID, v); err != nil {
			return err
		}
	}
	return nil
}

func reconcileExperiences(ctx context.Context, q querier, userID string, items []domain.WorkExperience) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM work_experiences WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		e := items[i]
		if e.ID != "" && isValidUUID(e.ID) {
			keep[e.ID] = true
			if err := updateExperienceQ(ctx, q, userID, e); err != nil {
				return err
			}
		} else {
			e.ID = ""
			if err := insertExperienceQ(ctx, q, userID, &e); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM work_experiences WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func reconcileEducations(ctx context.Context, q querier, userID string, items []domain.Education) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM educations WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		e := items[i]
		if e.ID != "" && isValidUUID(e.ID) {
			keep[e.ID] = true
			if err := updateEducationQ(ctx, q, userID, e); err != nil {
				return err
			}
		} else {
			e.ID = ""
			if err := insertEducationQ(ctx, q, userID, &e); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM educations WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func reconcileCertifications(ctx context.Context, q querier, userID string, items []domain.CertificationItem) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM certifications WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		c := items[i]
		if c.ID != "" && isValidUUID(c.ID) {
			keep[c.ID] = true
			if err := updateCertificationQ(ctx, q, userID, c); err != nil {
				return err
			}
		} else {
			c.ID = ""
			if err := insertCertificationQ(ctx, q, userID, &c); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM certifications WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func reconcileProjects(ctx context.Context, q querier, userID string, items []domain.ProjectItem) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM profile_projects WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		p := items[i]
		if p.ID != "" && isValidUUID(p.ID) {
			keep[p.ID] = true
			if err := updateProjectQ(ctx, q, userID, p); err != nil {
				return err
			}
		} else {
			p.ID = ""
			if err := insertProjectQ(ctx, q, userID, &p); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM profile_projects WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func reconcileAchievements(ctx context.Context, q querier, userID string, items []domain.AchievementItem) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM profile_achievements WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		a := items[i]
		if a.ID != "" && isValidUUID(a.ID) {
			keep[a.ID] = true
			if err := updateAchievementQ(ctx, q, userID, a); err != nil {
				return err
			}
		} else {
			a.ID = ""
			if err := insertAchievementQ(ctx, q, userID, &a); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM profile_achievements WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func existingIDs(ctx context.Context, q querier, query, userID string) (map[string]bool, error) {
	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

func isValidUUID(s string) bool {
	return len(s) == 36
}

func insertExperienceQ(ctx context.Context, q querier, userID string, e *domain.WorkExperience) error {
	// achievements/kpis/technologies/skills_used/attachments are Postgres text[]
	// columns; pgx encodes a Go []string to text[] directly (see loadExperiences).
	var start, end sql.NullTime
	if !e.StartDate.IsZero() {
		start = sql.NullTime{Time: e.StartDate, Valid: true}
	}
	if !e.EndDate.IsZero() {
		end = sql.NullTime{Time: e.EndDate, Valid: true}
	}

	return q.QueryRowContext(ctx, `
		INSERT INTO work_experiences (user_id, company, company_logo, position, employment_type, location, remote_type, start_date, end_date, is_current, responsibilities, achievements, kpis, technologies, skills_used, team_size, attachments)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING id`,
		userID, e.Company, e.CompanyLogo, e.Position, e.EmploymentType, e.Location, e.RemoteType, start, end, e.IsCurrent, e.Responsibilities, e.Achievements, e.KPIs, e.Technologies, e.SkillsUsed, e.TeamSize, e.Attachments).Scan(&e.ID)
}

func updateExperienceQ(ctx context.Context, q querier, userID string, e domain.WorkExperience) error {
	var start, end sql.NullTime
	if !e.StartDate.IsZero() {
		start = sql.NullTime{Time: e.StartDate, Valid: true}
	}
	if !e.EndDate.IsZero() {
		end = sql.NullTime{Time: e.EndDate, Valid: true}
	}

	res, err := q.ExecContext(ctx, `
		UPDATE work_experiences
		SET company=$3, company_logo=$4, position=$5, employment_type=$6, location=$7, remote_type=$8, start_date=$9, end_date=$10, is_current=$11, responsibilities=$12, achievements=$13, kpis=$14, technologies=$15, skills_used=$16, team_size=$17, attachments=$18, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.Company, e.CompanyLogo, e.Position, e.EmploymentType, e.Location, e.RemoteType, start, end, e.IsCurrent, e.Responsibilities, e.Achievements, e.KPIs, e.Technologies, e.SkillsUsed, e.TeamSize, e.Attachments)
	return owned(res, err)
}

func insertEducationQ(ctx context.Context, q querier, userID string, e *domain.Education) error {
	var grad sql.NullTime
	if !e.GraduationDate.IsZero() {
		grad = sql.NullTime{Time: e.GraduationDate, Valid: true}
	}

	return q.QueryRowContext(ctx, `
		INSERT INTO educations (user_id, school, degree, field_of_study, grade, description, major, minor, gpa, honors, activities, projects, research, thesis, graduation_date, verification_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id`,
		userID, e.Institution, e.Degree, e.FieldOfStudy, "", e.Thesis, e.Major, e.Minor, e.GPA, e.Honors, e.Activities, e.Projects, e.Research, e.Thesis, grad, e.VerificationStatus).Scan(&e.ID)
}

func updateEducationQ(ctx context.Context, q querier, userID string, e domain.Education) error {
	var grad sql.NullTime
	if !e.GraduationDate.IsZero() {
		grad = sql.NullTime{Time: e.GraduationDate, Valid: true}
	}

	res, err := q.ExecContext(ctx, `
		UPDATE educations
		SET school=$3, degree=$4, field_of_study=$5, grade=$6, description=$7, major=$8, minor=$9, gpa=$10, honors=$11, activities=$12, projects=$13, research=$14, thesis=$15, graduation_date=$16, verification_status=$17, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.Institution, e.Degree, e.FieldOfStudy, "", e.Thesis, e.Major, e.Minor, e.GPA, e.Honors, e.Activities, e.Projects, e.Research, e.Thesis, grad, e.VerificationStatus)
	return owned(res, err)
}

func insertCertificationQ(ctx context.Context, q querier, userID string, c *domain.CertificationItem) error {
	skills, _ := json.Marshal(c.SkillsCovered)
	var issue, expiry sql.NullTime
	if !c.IssueDate.IsZero() {
		issue = sql.NullTime{Time: c.IssueDate, Valid: true}
	}
	if !c.ExpirationDate.IsZero() {
		expiry = sql.NullTime{Time: c.ExpirationDate, Valid: true}
	}

	return q.QueryRowContext(ctx, `
		INSERT INTO certifications (user_id, name, issuer, credential_id, credential_url, skills_covered, issue_date, expiry_date, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id`,
		userID, c.Name, c.Issuer, c.CredentialID, c.VerificationURL, skills, issue, expiry, c.Status).Scan(&c.ID)
}

func updateCertificationQ(ctx context.Context, q querier, userID string, c domain.CertificationItem) error {
	skills, _ := json.Marshal(c.SkillsCovered)
	var issue, expiry sql.NullTime
	if !c.IssueDate.IsZero() {
		issue = sql.NullTime{Time: c.IssueDate, Valid: true}
	}
	if !c.ExpirationDate.IsZero() {
		expiry = sql.NullTime{Time: c.ExpirationDate, Valid: true}
	}

	res, err := q.ExecContext(ctx, `
		UPDATE certifications
		SET name=$3, issuer=$4, credential_id=$5, credential_url=$6, skills_covered=$7, issue_date=$8, expiry_date=$9, status=$10, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		c.ID, userID, c.Name, c.Issuer, c.CredentialID, c.VerificationURL, skills, issue, expiry, c.Status)
	return owned(res, err)
}

func insertProjectQ(ctx context.Context, q querier, userID string, p *domain.ProjectItem) error {
	screens, _ := json.Marshal(p.Images)
	techs, _ := json.Marshal(p.Technologies)

	return q.QueryRowContext(ctx, `
		INSERT INTO profile_projects (user_id, title, description, repository_url, live_demo_url, video_url, screenshots, technologies, timeline, team_size, metrics, awards, business_impact)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id`,
		userID, p.Title, p.Description, p.RepositoryURL, p.LiveDemoURL, p.VideoURL, screens, techs, p.Timeline, len(p.TeamMembers), p.Metrics, p.Awards, p.BusinessImpact).Scan(&p.ID)
}

func updateProjectQ(ctx context.Context, q querier, userID string, p domain.ProjectItem) error {
	screens, _ := json.Marshal(p.Images)
	techs, _ := json.Marshal(p.Technologies)

	res, err := q.ExecContext(ctx, `
		UPDATE profile_projects
		SET title=$3, description=$4, repository_url=$5, live_demo_url=$6, video_url=$7, screenshots=$8, technologies=$9, timeline=$10, team_size=$11, metrics=$12, awards=$13, business_impact=$14, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		p.ID, userID, p.Title, p.Description, p.RepositoryURL, p.LiveDemoURL, p.VideoURL, screens, techs, p.Timeline, len(p.TeamMembers), p.Metrics, p.Awards, p.BusinessImpact)
	return owned(res, err)
}

func insertAchievementQ(ctx context.Context, q querier, userID string, a *domain.AchievementItem) error {
	var d sql.NullTime
	if !a.Date.IsZero() {
		d = sql.NullTime{Time: a.Date, Valid: true}
	}
	return q.QueryRowContext(ctx, `
		INSERT INTO profile_achievements (user_id, title, issuer_or_org, date, category, description, evidence_url)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		userID, a.Title, a.IssuerOrOrg, d, a.Category, a.Description, a.EvidenceURL).Scan(&a.ID)
}

func updateAchievementQ(ctx context.Context, q querier, userID string, a domain.AchievementItem) error {
	var d sql.NullTime
	if !a.Date.IsZero() {
		d = sql.NullTime{Time: a.Date, Valid: true}
	}
	res, err := q.ExecContext(ctx, `
		UPDATE profile_achievements
		SET title=$3, issuer_or_org=$4, date=$5, category=$6, description=$7, evidence_url=$8
		WHERE id=$1 AND user_id=$2`,
		a.ID, userID, a.Title, a.IssuerOrOrg, d, a.Category, a.Description, a.EvidenceURL)
	return owned(res, err)
}

func setSkillsQ(ctx context.Context, q querier, userID string, skills []domain.SkillItem) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM profile_skills WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, sk := range skills {
		if sk.Name == "" {
			continue
		}
		var skillID string
		if err := q.QueryRowContext(ctx,
			`INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
			sk.Name).Scan(&skillID); err != nil {
			return err
		}
		if _, err := q.ExecContext(ctx,
			`INSERT INTO profile_skills (user_id, skill_id, category, proficiency_level, years_of_experience, last_used, verified, recruiter_demand_score, ai_recommendation_score)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT DO NOTHING`,
			userID, skillID, sk.Category, sk.Level, sk.YearsOfExperience, sk.LastUsed, sk.Verified, sk.RecruiterDemandScore, sk.AIRecommendationScore); err != nil {
			return err
		}
	}
	return nil
}

func setLanguagesQ(ctx context.Context, q querier, userID string, langs []domain.LanguageItem) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM profile_languages WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, l := range langs {
		if l.Name == "" {
			continue
		}
		var langID string
		if err := q.QueryRowContext(ctx,
			`INSERT INTO languages (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
			l.Name).Scan(&langID); err != nil {
			return err
		}
		if _, err := q.ExecContext(ctx,
			`INSERT INTO profile_languages (user_id, language_id, proficiency) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
			userID, langID, l.Proficiency); err != nil {
			return err
		}
	}
	return nil
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
