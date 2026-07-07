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
	return r.tx(ctx, func(tx *sql.Tx) error {
		return r.applyScalars(ctx, tx, userID, s)
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
		return insertExperienceQ(ctx, tx, userID, e)
	})
}

func (r *Repository) UpdateExperience(ctx context.Context, userID string, e domain.WorkExperience) error {
	return r.tx(ctx, func(tx *sql.Tx) error {
		return updateExperienceQ(ctx, tx, userID, e)
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
	return insertEducationQ(ctx, r.db, userID, e)
}

func (r *Repository) UpdateEducation(ctx context.Context, userID string, e domain.Education) error {
	return updateEducationQ(ctx, r.db, userID, e)
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
	return insertCertificationQ(ctx, r.db, userID, c)
}

func (r *Repository) UpdateCertification(ctx context.Context, userID string, c domain.Certification) error {
	return updateCertificationQ(ctx, r.db, userID, c)
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
		return setSkillsQ(ctx, tx, userID, skills)
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
		return setLanguagesQ(ctx, tx, userID, langs)
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
		return setPortfolioQ(ctx, tx, userID, links)
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
	return insertReferenceQ(ctx, r.db, userID, rf)
}

func (r *Repository) UpdateReference(ctx context.Context, userID string, rf domain.Reference) error {
	return updateReferenceQ(ctx, r.db, userID, rf)
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
