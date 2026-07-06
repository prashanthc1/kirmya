// Package postgres implements profile/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"workspace-app/internal/profile/domain"
	"workspace-app/internal/profile/infrastructure/crypto"
)

type Repository struct {
	db    *sql.DB
	crypt *crypto.Encryptor
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:    db,
		crypt: crypto.NewEncryptor(),
	}
}

const dateLayout = "2006-01-02"

func dateStr(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format(dateLayout)
}

func timeStr(t time.Time) string {
	return t.Format(time.RFC3339)
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

	var availability sql.NullTime
	var lastActive time.Time
	var consentAt sql.NullTime
	var transReasonEnc, salMinEnc, salMaxEnc, salCurrEnc string

	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(p.headline,''), COALESCE(p.about,''), COALESCE(p.photo_url,''),
		       COALESCE(p.bio,''), COALESCE(p.location,''), COALESCE(p.website,''), p.version,
		       COALESCE(p.pronouns,''), COALESCE(p.career_status,''), COALESCE(p.transition_reason_enc,''),
		       COALESCE(p.target_comeback_timeline,''), p.open_to_remote, p.open_to_relocation,
		       COALESCE(p.employment_type,''), COALESCE(p.salary_min_enc,''), COALESCE(p.salary_max_enc,''),
		       COALESCE(p.salary_currency_enc,''), p.salary_visible, COALESCE(p.work_mode,''),
		       p.availability_date, COALESCE(p.notice_period,''), p.referral_eligible,
		       COALESCE(p.career_narrative,''), COALESCE(p.coaching_metadata,''), COALESCE(p.work_auth_status,''),
		       COALESCE(p.passport_nationality,''), p.driving_license_bool, COALESCE(p.driving_license_type,''),
		       COALESCE(p.preferred_contact_channel,''), COALESCE(p.accessibility_needs,''), COALESCE(p.video_intro_url,''),
		       p.willing_to_mentor, p.avg_response_time_hours, p.profile_completeness_score,
		       p.last_active_at, p.background_check_consent, p.background_check_consent_at,
		       COALESCE(p.job_alert_frequency,''), COALESCE(p.job_alert_channel,''),
		       COALESCE(p.visibility_profile,'public'), COALESCE(p.visibility_salary,'private'), COALESCE(p.visibility_transition_reason,'private'),
		       COALESCE(p.visibility_experience,'public'), COALESCE(p.visibility_education,'public'), COALESCE(p.visibility_certifications,'public'),
		       COALESCE(p.visibility_skills,'public'), COALESCE(p.visibility_portfolio,'public'), COALESCE(p.visibility_references,'private'),
		       p.phone_verified, p.linkedin_verified, p.id_verified,
		       u.email_verified
		FROM profiles p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $1 AND p.deleted_at IS NULL`, userID).
		Scan(
			&p.Headline, &p.About, &p.PhotoURL, &p.Bio, &p.Location, &p.Website, &p.Version,
			&p.Pronouns, &p.CareerStatus, &transReasonEnc,
			&p.TargetComebackTimeline, &p.OpenToRemote, &p.OpenToRelocation,
			&p.EmploymentType, &salMinEnc, &salMaxEnc,
			&salCurrEnc, &p.SalaryVisible, &p.WorkMode,
			&availability, &p.NoticePeriod, &p.ReferralEligible,
			&p.CareerNarrative, &p.CoachingMetadata, &p.WorkAuthStatus,
			&p.PassportNationality, &p.DrivingLicenseBool, &p.DrivingLicenseType,
			&p.PreferredContactChannel, &p.AccessibilityNeeds, &p.VideoIntroURL,
			&p.WillingToMentor, &p.AvgResponseTimeHours, &p.ProfileCompletenessScore,
			&lastActive, &p.BackgroundCheckConsent, &consentAt,
			&p.JobAlertFrequency, &p.JobAlertChannel,
			&p.VisibilityProfile, &p.VisibilitySalary, &p.VisibilityTransitionReason,
			&p.VisibilityExperience, &p.VisibilityEducation, &p.VisibilityCertifications,
			&p.VisibilitySkills, &p.VisibilityPortfolio, &p.VisibilityReferences,
			&p.PhoneVerified, &p.LinkedinVerified, &p.IdVerified,
			&p.EmailVerified,
		)
	if err != nil {
		return nil, err
	}

	p.AvailabilityDate = dateStr(availability)
	p.BackgroundCheckConsentAt = dateStr(consentAt)
	p.LastActiveAt = timeStr(lastActive)

	// Decrypt sensitive fields
	if transReasonEnc != "" {
		p.TransitionReason, _ = r.crypt.Decrypt(transReasonEnc)
	}
	if salCurrEnc != "" {
		p.SalaryCurrency, _ = r.crypt.Decrypt(salCurrEnc)
	}
	if salMinEnc != "" {
		dec, _ := r.crypt.Decrypt(salMinEnc)
		p.SalaryMin, _ = strconv.Atoi(dec)
	}
	if salMaxEnc != "" {
		dec, _ := r.crypt.Decrypt(salMaxEnc)
		p.SalaryMax, _ = strconv.Atoi(dec)
	}

	// Load sub-resources
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
	if p.SupportsNeeded, err = r.loadSupports(ctx, userID); err != nil {
		return nil, err
	}
	if p.RelocationLocations, err = r.loadRelocationLocations(ctx, userID); err != nil {
		return nil, err
	}
	if p.DesiredRoles, err = r.loadDesiredRoles(ctx, userID); err != nil {
		return nil, err
	}
	if p.DesiredIndustries, err = r.loadDesiredIndustries(ctx, userID); err != nil {
		return nil, err
	}
	if p.Endorsements, err = r.loadEndorsements(ctx, userID); err != nil {
		return nil, err
	}
	if p.References, err = r.loadReferences(ctx, userID); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) UpdateScalars(ctx context.Context, userID string, s domain.Scalars) error {
	if err := r.ensureRow(ctx, userID); err != nil {
		return err
	}

	// Encrypt sensitive fields
	transReasonEnc, err := r.crypt.Encrypt(s.TransitionReason)
	if err != nil {
		return err
	}
	salMinEnc, err := r.crypt.Encrypt(strconv.Itoa(s.SalaryMin))
	if err != nil {
		return err
	}
	salMaxEnc, err := r.crypt.Encrypt(strconv.Itoa(s.SalaryMax))
	if err != nil {
		return err
	}
	salCurrEnc, err := r.crypt.Encrypt(s.SalaryCurrency)
	if err != nil {
		return err
	}

	return r.tx(ctx, func(tx *sql.Tx) error {
		// Update scalar fields
		_, err = tx.ExecContext(ctx, `
			UPDATE profiles
			SET headline = $2, about = $3, photo_url = $4, bio = $5, location = $6, website = $7,
			    pronouns = $8, career_status = $9, transition_reason_enc = $10, target_comeback_timeline = $11,
			    open_to_remote = $12, open_to_relocation = $13, employment_type = $14,
			    salary_min_enc = $15, salary_max_enc = $16, salary_currency_enc = $17, salary_visible = $18,
			    work_mode = $19, availability_date = NULLIF($20,'')::date, notice_period = $21,
			    referral_eligible = $22, career_narrative = $23, coaching_metadata = $24,
			    work_auth_status = $25, passport_nationality = $26, driving_license_bool = $27, driving_license_type = $28,
			    preferred_contact_channel = $29, accessibility_needs = $30, video_intro_url = $31,
			    willing_to_mentor = $32, background_check_consent = $33,
			    background_check_consent_at = NULLIF($34,'')::timestamptz, job_alert_frequency = $35, job_alert_channel = $36,
			    visibility_profile = $37, visibility_salary = $38, visibility_transition_reason = $39,
			    visibility_experience = $40, visibility_education = $41, visibility_certifications = $42,
			    visibility_skills = $43, visibility_portfolio = $44, visibility_references = $45,
			    updated_at = now(), version = version + 1
			WHERE user_id = $1 AND deleted_at IS NULL`,
			userID, s.Headline, s.About, s.PhotoURL, s.Bio, s.Location, s.Website,
			s.Pronouns, s.CareerStatus, transReasonEnc, s.TargetComebackTimeline,
			s.OpenToRemote, s.OpenToRelocation, s.EmploymentType,
			salMinEnc, salMaxEnc, salCurrEnc, s.SalaryVisible,
			s.WorkMode, s.AvailabilityDate, s.NoticePeriod,
			s.ReferralEligible, s.CareerNarrative, s.CoachingMetadata,
			s.WorkAuthStatus, s.PassportNationality, s.DrivingLicenseBool, s.DrivingLicenseType,
			s.PreferredContactChannel, s.AccessibilityNeeds, s.VideoIntroURL,
			s.WillingToMentor, s.BackgroundCheckConsent, s.BackgroundCheckConsentAt,
			s.JobAlertFrequency, s.JobAlertChannel,
			s.VisibilityProfile, s.VisibilitySalary, s.VisibilityTransitionReason,
			s.VisibilityExperience, s.VisibilityEducation, s.VisibilityCertifications,
			s.VisibilitySkills, s.VisibilityPortfolio, s.VisibilityReferences,
		)
		if err != nil {
			return err
		}

		// Rebuild support table
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_supports WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, sup := range s.SupportsNeeded {
			if sup == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO profile_supports (user_id, support) VALUES ($1,$2)`, userID, sup); err != nil {
				return err
			}
		}

		// Rebuild relocation locations
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_relocation_locations WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, loc := range s.RelocationLocations {
			if loc == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO profile_relocation_locations (user_id, location) VALUES ($1,$2)`, userID, loc); err != nil {
				return err
			}
		}

		// Rebuild desired roles
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_desired_roles WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, role := range s.DesiredRoles {
			if role == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO profile_desired_roles (user_id, role) VALUES ($1,$2)`, userID, role); err != nil {
				return err
			}
		}

		// Rebuild desired industries
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_desired_industries WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, ind := range s.DesiredIndustries {
			if ind == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO profile_desired_industries (user_id, industry) VALUES ($1,$2)`, userID, ind); err != nil {
				return err
			}
		}

		return nil
	})
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

	var out []domain.WorkExperience
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
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Load Achievements for all experiences
	achRows, err := r.db.QueryContext(ctx, `
		SELECT a.experience_id, a.achievement
		FROM work_experience_achievements a
		JOIN work_experiences e ON e.id = a.experience_id
		WHERE e.user_id = $1
		ORDER BY a.sort_order ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer achRows.Close()

	achMap := make(map[string][]string)
	for achRows.Next() {
		var expID, ach string
		if err := achRows.Scan(&expID, &ach); err != nil {
			return nil, err
		}
		achMap[expID] = append(achMap[expID], ach)
	}
	if err = achRows.Err(); err != nil {
		return nil, err
	}

	for i := range out {
		if achs, ok := achMap[out[i].ID]; ok {
			out[i].Achievements = achs
		} else {
			out[i].Achievements = []string{}
		}
	}

	return out, nil
}

func (r *Repository) AddExperience(ctx context.Context, userID string, e *domain.WorkExperience) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO work_experiences (user_id, title, company, location, employment_type, start_date, end_date, is_current, description)
			VALUES ($1,$2,$3,$4,$5,NULLIF($6,'')::date,NULLIF($7,'')::date,$8,$9)
			RETURNING id`,
			userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description).
			Scan(&e.ID)
		if err != nil {
			return err
		}
		for i, ach := range e.Achievements {
			if ach == "" {
				continue
			}
			_, err = tx.ExecContext(ctx, `
				INSERT INTO work_experience_achievements (experience_id, achievement, sort_order)
				VALUES ($1,$2,$3)`, e.ID, ach, i)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) UpdateExperience(ctx context.Context, userID string, e domain.WorkExperience) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `
			UPDATE work_experiences
			SET title=$3, company=$4, location=$5, employment_type=$6,
			    start_date=NULLIF($7,'')::date, end_date=NULLIF($8,'')::date, is_current=$9, description=$10, updated_at=now()
			WHERE id=$1 AND user_id=$2`,
			e.ID, userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description)
		if err = owned(res, err); err != nil {
			return err
		}

		if _, err = tx.ExecContext(ctx, `DELETE FROM work_experience_achievements WHERE experience_id = $1`, e.ID); err != nil {
			return err
		}

		for i, ach := range e.Achievements {
			if ach == "" {
				continue
			}
			_, err = tx.ExecContext(ctx, `
				INSERT INTO work_experience_achievements (experience_id, achievement, sort_order)
				VALUES ($1,$2,$3)`, e.ID, ach, i)
			if err != nil {
				return err
			}
		}
		return nil
	})
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

func (r *Repository) loadSkills(ctx context.Context, userID string) ([]domain.ProfileSkill, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.name, COALESCE(ps.proficiency_level, ''), ps.endorsed_count
		FROM profile_skills ps JOIN skills s ON s.id = ps.skill_id
		WHERE ps.user_id = $1 ORDER BY s.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.ProfileSkill{}
	for rows.Next() {
		var s domain.ProfileSkill
		if err := rows.Scan(&s.Name, &s.ProficiencyLevel, &s.EndorsedCount); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) SetSkills(ctx context.Context, userID string, skills []domain.ProfileSkill) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile_skills WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, sk := range skills {
			if sk.Name == "" {
				continue
			}
			var skillID string
			if err := tx.QueryRowContext(ctx,
				`INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				sk.Name).Scan(&skillID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO profile_skills (user_id, skill_id, proficiency_level, endorsed_count) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
				userID, skillID, sk.ProficiencyLevel, sk.EndorsedCount); err != nil {
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
		if err := rows.Scan(&l.ID, &l.Platform, &l.URL); err != nil {
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
				`INSERT INTO portfolio_links (user_id, label, url) VALUES ($1,$2,$3)`, userID, l.Platform, l.URL); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- newly normalized tables (supports, relocations, roles, industries) ---

func (r *Repository) loadSupports(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT support FROM profile_supports WHERE user_id = $1 ORDER BY support`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) loadRelocationLocations(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT location FROM profile_relocation_locations WHERE user_id = $1 ORDER BY location`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) loadDesiredRoles(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT role FROM profile_desired_roles WHERE user_id = $1 ORDER BY role`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) loadDesiredIndustries(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT industry FROM profile_desired_industries WHERE user_id = $1 ORDER BY industry`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) loadEndorsements(ctx context.Context, userID string) ([]domain.Endorsement, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, to_user_id, from_user_id, relationship, text, created_at
		FROM profile_endorsements WHERE to_user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Endorsement
	for rows.Next() {
		var e domain.Endorsement
		var t time.Time
		if err := rows.Scan(&e.ID, &e.ToUserID, &e.FromUserID, &e.Relationship, &e.Text, &t); err != nil {
			return nil, err
		}
		e.CreatedAt = timeStr(t)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *Repository) loadReferences(ctx context.Context, userID string) ([]domain.Reference, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, relationship, contact_info, permission_to_contact
		FROM profile_references WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Reference
	for rows.Next() {
		var rf domain.Reference
		if err := rows.Scan(&rf.ID, &rf.Name, &rf.Relationship, &rf.ContactInfo, &rf.PermissionToContact); err != nil {
			return nil, err
		}
		out = append(out, rf)
	}
	return out, rows.Err()
}

// --- resource mutation methods ---

func (r *Repository) AddEndorsement(ctx context.Context, toUserID string, e *domain.Endorsement) error {
	var t time.Time
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO profile_endorsements (to_user_id, from_user_id, relationship, text)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		toUserID, e.FromUserID, e.Relationship, e.Text).Scan(&e.ID, &t)
	if err != nil {
		return err
	}
	e.CreatedAt = timeStr(t)
	return nil
}

func (r *Repository) AddReference(ctx context.Context, userID string, rf *domain.Reference) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO profile_references (user_id, name, relationship, contact_info, permission_to_contact)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, rf.Name, rf.Relationship, rf.ContactInfo, rf.PermissionToContact).Scan(&rf.ID)
}

func (r *Repository) UpdateReference(ctx context.Context, userID string, rf domain.Reference) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE profile_references
		SET name = $3, relationship = $4, contact_info = $5, permission_to_contact = $6
		WHERE id = $1 AND user_id = $2`,
		rf.ID, userID, rf.Name, rf.Relationship, rf.ContactInfo, rf.PermissionToContact)
	return owned(res, err)
}

func (r *Repository) DeleteReference(ctx context.Context, userID, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM profile_references WHERE id = $1 AND user_id = $2`, id, userID)
	return owned(res, err)
}

func (r *Repository) AddConsentLog(ctx context.Context, cl *domain.ConsentLog) error {
	var t time.Time
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO profile_consent_logs (user_id, consent_type, target_entity, consented, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		cl.UserID, cl.ConsentType, cl.TargetEntity, cl.Consented, cl.IPAddress, cl.UserAgent).Scan(&cl.ID, &t)
	if err != nil {
		return err
	}
	cl.CreatedAt = timeStr(t)
	return nil
}

func (r *Repository) SetVerificationStatus(ctx context.Context, userID string, field string, verified bool) error {
	var query string
	switch field {
	case "phone_verified":
		query = `UPDATE profiles SET phone_verified = $2 WHERE user_id = $1`
	case "linkedin_verified":
		query = `UPDATE profiles SET linkedin_verified = $2 WHERE user_id = $1`
	case "id_verified":
		query = `UPDATE profiles SET id_verified = $2 WHERE user_id = $1`
	default:
		return fmt.Errorf("invalid verification field: %s", field)
	}
	_, err := r.db.ExecContext(ctx, query, userID, verified)
	return err
}

func (r *Repository) UpdateCalculatedFields(ctx context.Context, userID string, completeness int, avgResponse float64, lastActive string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE profiles
		SET profile_completeness_score = $2, avg_response_time_hours = $3, last_active_at = COALESCE(NULLIF($4,'')::timestamptz, now())
		WHERE user_id = $1`,
		userID, completeness, avgResponse, lastActive)
	return err
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
