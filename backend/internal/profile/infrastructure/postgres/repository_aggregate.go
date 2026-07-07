package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"workspace-app/internal/profile/domain"
)

// querier is the subset of *sql.DB / *sql.Tx used by the write helpers, so the
// same SQL can run either standalone (r.db) or inside a shared transaction (tx).
type querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// UpdateAggregate applies the scalar fields plus any provided child collections
// atomically. When expectedVersion > 0 it takes a row lock, compares the stored
// version, and returns domain.ErrOptimisticLock on a stale write. Nil collection
// pointers are left untouched; non-nil (possibly empty) slices fully reconcile.
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

		// Scalars (this also bumps version = version + 1).
		if err := r.applyScalars(ctx, tx, userID, u.Scalars); err != nil {
			return err
		}

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
		if u.Certifications != nil {
			if err := reconcileCertifications(ctx, tx, userID, *u.Certifications); err != nil {
				return err
			}
		}
		if u.Skills != nil {
			if err := setSkillsQ(ctx, tx, userID, *u.Skills); err != nil {
				return err
			}
		}
		if u.Languages != nil {
			if err := setLanguagesQ(ctx, tx, userID, *u.Languages); err != nil {
				return err
			}
		}
		if u.Portfolio != nil {
			if err := setPortfolioQ(ctx, tx, userID, *u.Portfolio); err != nil {
				return err
			}
		}
		if u.References != nil {
			if err := reconcileReferences(ctx, tx, userID, *u.References); err != nil {
				return err
			}
		}
		return nil
	})
}

// applyScalars encrypts sensitive fields and writes the scalar columns plus the
// normalized string tables (supports / relocations / roles / industries),
// bumping the row version. Shared by UpdateScalars and UpdateAggregate.
func (r *Repository) applyScalars(ctx context.Context, q querier, userID string, s domain.Scalars) error {
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

	if _, err = q.ExecContext(ctx, `
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
	); err != nil {
		return err
	}

	if err = rebuildStringTable(ctx, q, `profile_supports`, `support`, userID, s.SupportsNeeded); err != nil {
		return err
	}
	if err = rebuildStringTable(ctx, q, `profile_relocation_locations`, `location`, userID, s.RelocationLocations); err != nil {
		return err
	}
	if err = rebuildStringTable(ctx, q, `profile_desired_roles`, `role`, userID, s.DesiredRoles); err != nil {
		return err
	}
	if err = rebuildStringTable(ctx, q, `profile_desired_industries`, `industry`, userID, s.DesiredIndustries); err != nil {
		return err
	}
	return nil
}

// rebuildStringTable replaces the rows of a (user_id, <col>) join table.
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

// --- collection reconcilers (upsert incoming, delete the rest) ---

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

func reconcileCertifications(ctx context.Context, q querier, userID string, items []domain.Certification) error {
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

func reconcileReferences(ctx context.Context, q querier, userID string, items []domain.Reference) error {
	existing, err := existingIDs(ctx, q, `SELECT id FROM profile_references WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	keep := make(map[string]bool)
	for i := range items {
		rf := items[i]
		if rf.ID != "" && isValidUUID(rf.ID) {
			keep[rf.ID] = true
			if err := updateReferenceQ(ctx, q, userID, rf); err != nil {
				return err
			}
		} else {
			rf.ID = ""
			if err := insertReferenceQ(ctx, q, userID, &rf); err != nil {
				return err
			}
		}
	}
	for id := range existing {
		if !keep[id] {
			if _, err := q.ExecContext(ctx, `DELETE FROM profile_references WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
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

// isValidUUID reports whether s has canonical UUID hyphen placement (used to
// distinguish server-assigned IDs from temporary/mock client IDs).
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		}
	}
	return true
}

// --- querier-bound single-item writers (shared by public methods + reconcilers) ---

func insertExperienceQ(ctx context.Context, q querier, userID string, e *domain.WorkExperience) error {
	if err := q.QueryRowContext(ctx, `
		INSERT INTO work_experiences (user_id, title, company, location, employment_type, start_date, end_date, is_current, description)
		VALUES ($1,$2,$3,$4,$5,NULLIF($6,'')::date,NULLIF($7,'')::date,$8,$9)
		RETURNING id`,
		userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description).
		Scan(&e.ID); err != nil {
		return err
	}
	for i, ach := range e.Achievements {
		if ach == "" {
			continue
		}
		if _, err := q.ExecContext(ctx, `
			INSERT INTO work_experience_achievements (experience_id, achievement, sort_order)
			VALUES ($1,$2,$3)`, e.ID, ach, i); err != nil {
			return err
		}
	}
	return nil
}

func updateExperienceQ(ctx context.Context, q querier, userID string, e domain.WorkExperience) error {
	res, err := q.ExecContext(ctx, `
		UPDATE work_experiences
		SET title=$3, company=$4, location=$5, employment_type=$6,
		    start_date=NULLIF($7,'')::date, end_date=NULLIF($8,'')::date, is_current=$9, description=$10, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.Title, e.Company, e.Location, e.EmploymentType, e.StartDate, e.EndDate, e.IsCurrent, e.Description)
	if err = owned(res, err); err != nil {
		return err
	}
	if _, err = q.ExecContext(ctx, `DELETE FROM work_experience_achievements WHERE experience_id = $1`, e.ID); err != nil {
		return err
	}
	for i, ach := range e.Achievements {
		if ach == "" {
			continue
		}
		if _, err = q.ExecContext(ctx, `
			INSERT INTO work_experience_achievements (experience_id, achievement, sort_order)
			VALUES ($1,$2,$3)`, e.ID, ach, i); err != nil {
			return err
		}
	}
	return nil
}

func insertEducationQ(ctx context.Context, q querier, userID string, e *domain.Education) error {
	return q.QueryRowContext(ctx, `
		INSERT INTO educations (user_id, school, degree, field_of_study, start_date, end_date, grade, description)
		VALUES ($1,$2,$3,$4,NULLIF($5,'')::date,NULLIF($6,'')::date,$7,$8)
		RETURNING id`,
		userID, e.School, e.Degree, e.FieldOfStudy, e.StartDate, e.EndDate, e.Grade, e.Description).Scan(&e.ID)
}

func updateEducationQ(ctx context.Context, q querier, userID string, e domain.Education) error {
	res, err := q.ExecContext(ctx, `
		UPDATE educations
		SET school=$3, degree=$4, field_of_study=$5, start_date=NULLIF($6,'')::date,
		    end_date=NULLIF($7,'')::date, grade=$8, description=$9, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		e.ID, userID, e.School, e.Degree, e.FieldOfStudy, e.StartDate, e.EndDate, e.Grade, e.Description)
	return owned(res, err)
}

func insertCertificationQ(ctx context.Context, q querier, userID string, c *domain.Certification) error {
	return q.QueryRowContext(ctx, `
		INSERT INTO certifications (user_id, name, issuer, issue_date, expiry_date, credential_id, credential_url)
		VALUES ($1,$2,$3,NULLIF($4,'')::date,NULLIF($5,'')::date,$6,$7)
		RETURNING id`,
		userID, c.Name, c.Issuer, c.IssueDate, c.ExpiryDate, c.CredentialID, c.CredentialURL).Scan(&c.ID)
}

func updateCertificationQ(ctx context.Context, q querier, userID string, c domain.Certification) error {
	res, err := q.ExecContext(ctx, `
		UPDATE certifications
		SET name=$3, issuer=$4, issue_date=NULLIF($5,'')::date, expiry_date=NULLIF($6,'')::date,
		    credential_id=$7, credential_url=$8, updated_at=now()
		WHERE id=$1 AND user_id=$2`,
		c.ID, userID, c.Name, c.Issuer, c.IssueDate, c.ExpiryDate, c.CredentialID, c.CredentialURL)
	return owned(res, err)
}

func insertReferenceQ(ctx context.Context, q querier, userID string, rf *domain.Reference) error {
	return q.QueryRowContext(ctx, `
		INSERT INTO profile_references (user_id, name, relationship, contact_info, permission_to_contact)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, rf.Name, rf.Relationship, rf.ContactInfo, rf.PermissionToContact).Scan(&rf.ID)
}

func updateReferenceQ(ctx context.Context, q querier, userID string, rf domain.Reference) error {
	res, err := q.ExecContext(ctx, `
		UPDATE profile_references
		SET name = $3, relationship = $4, contact_info = $5, permission_to_contact = $6
		WHERE id = $1 AND user_id = $2`,
		rf.ID, userID, rf.Name, rf.Relationship, rf.ContactInfo, rf.PermissionToContact)
	return owned(res, err)
}

func setSkillsQ(ctx context.Context, q querier, userID string, skills []domain.ProfileSkill) error {
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
			`INSERT INTO profile_skills (user_id, skill_id, proficiency_level, endorsed_count) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
			userID, skillID, sk.ProficiencyLevel, sk.EndorsedCount); err != nil {
			return err
		}
	}
	return nil
}

func setLanguagesQ(ctx context.Context, q querier, userID string, langs []domain.Language) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM profile_languages WHERE user_id = $1`, userID); err != nil {
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
		if err := q.QueryRowContext(ctx,
			`INSERT INTO languages (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
			l.Name).Scan(&langID); err != nil {
			return err
		}
		if _, err := q.ExecContext(ctx,
			`INSERT INTO profile_languages (user_id, language_id, proficiency) VALUES ($1,$2,$3)
			 ON CONFLICT (user_id, language_id) DO UPDATE SET proficiency = EXCLUDED.proficiency`,
			userID, langID, prof); err != nil {
			return err
		}
	}
	return nil
}

func setPortfolioQ(ctx context.Context, q querier, userID string, links []domain.PortfolioLink) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM portfolio_links WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, l := range links {
		if l.URL == "" {
			continue
		}
		if _, err := q.ExecContext(ctx,
			`INSERT INTO portfolio_links (user_id, label, url) VALUES ($1,$2,$3)`, userID, l.Platform, l.URL); err != nil {
			return err
		}
	}
	return nil
}
