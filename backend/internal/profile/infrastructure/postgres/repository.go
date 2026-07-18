// Package postgres implements profile/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"

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

func (r *Repository) ensureRow(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO profiles (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`, userID)
	return err
}

func (r *Repository) Get(ctx context.Context, userID string, includeDraft bool) (*domain.Profile, error) {
	if err := r.ensureRow(ctx, userID); err != nil {
		return nil, err
	}

	// 1. If public profile (includeDraft = false) requested, attempt to get latest version snapshot
	if !includeDraft {
		var snapshotJSON []byte
		err := r.db.QueryRowContext(ctx, `
			SELECT snapshot FROM profile_versions
			WHERE user_id = $1 ORDER BY version DESC LIMIT 1`, userID).Scan(&snapshotJSON)
		if err == nil {
			var p domain.Profile
			if err := json.Unmarshal(snapshotJSON, &p); err == nil {
				return &p, nil
			}
		}
	}

	// 2. Load draft/current state from normal tables
	p := &domain.Profile{UserID: userID}

	var availability sql.NullTime
	var lastActive time.Time
	var consentAt sql.NullTime
	var transReasonEnc, salMinEnc, salMaxEnc, salCurrEnc, emailEnc, phoneEnc, addressEnc string
	var bioOptimized string
	var visProfile, visSalary, visTransReason, visExp, visEdu, visCert, visSkills, visPortfolio, visRef string

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
		       p.is_draft, p.trust_score,
		       COALESCE(p.preferred_name, ''), COALESCE(p.timezone, ''), COALESCE(p.nationality, ''),
		       COALESCE(p.bio_optimized, ''), COALESCE(p.executive_summary, ''), COALESCE(p.career_objectives, ''),
		       COALESCE(p.personal_brand_statement, ''), COALESCE(p.elevator_pitch, ''),
		       p.email_verified, p.employment_verified, p.education_verified, p.certification_verified,
		       COALESCE(p.travel_willingness, ''),
		       COALESCE(p.email_enc, ''), COALESCE(p.phone_enc, ''), COALESCE(p.address_enc, ''),
		       COALESCE(p.full_name, '')
		FROM profiles p
		WHERE p.user_id = $1 AND p.deleted_at IS NULL`, userID).Scan(
		&p.Identity.Headline, &p.Identity.About, &p.Identity.PhotoURL,
		&p.Identity.Bio, &p.Identity.Location, &p.Identity.SocialLinks.Website, &p.Version,
		&p.Identity.Pronouns, &p.Identity.CareerStatus, &transReasonEnc,
		&p.Identity.Availability, &p.Preferences.OpenToRelocation, &p.Preferences.OpenToRelocation,
		&p.Preferences.NoticePeriod, &salMinEnc, &salMaxEnc,
		&salCurrEnc, &p.Preferences.OpenToRelocation, &p.Preferences.RemotePreference,
		&availability, &p.Preferences.NoticePeriod, &p.Verification.IdentityVerified,
		&p.AICareerAssistant.GapAnalysis, &p.AICareerAssistant.InterviewPrep, &p.Identity.WorkAuthorization,
		&p.Identity.Nationality, &p.Verification.IdentityVerified, &p.Identity.VisaStatus,
		&p.Identity.PreferredContactChannel, &p.Identity.VisaStatus, &p.Identity.CoverURL,
		&p.Identity.VisaStatus, &p.Analytics.ProfileViews, &p.ProfileCompletenessScore,
		&lastActive, &p.Verification.IdentityVerified, &consentAt,
		&p.Identity.VisaStatus, &p.Identity.VisaStatus,
		&visProfile, &visSalary, &visTransReason,
		&visExp, &visEdu, &visCert,
		&visSkills, &visPortfolio, &visRef,
		&p.Verification.PhoneVerified, &p.Verification.IdentityVerified, &p.Verification.IdentityVerified,
		&p.IsDraft, &p.TrustScore,
		&p.Identity.PreferredName, &p.Identity.TimeZone, &p.Identity.Nationality,
		&bioOptimized, &p.Summary.ExecutiveSummary, &p.Summary.CareerObjectives,
		&p.Summary.PersonalBrandStatement, &p.Summary.ElevatorPitch,
		&p.Verification.EmailVerified, &p.Verification.EmploymentVerified, &p.Verification.EducationVerified, &p.Verification.CertificationVerified,
		&p.Preferences.TravelWillingness,
		&emailEnc, &phoneEnc, &addressEnc, &p.Identity.FullName,
	)
	if err != nil {
		return nil, err
	}

	p.LastActiveAt = lastActive

	// Populate visibility map
	p.Privacy.FieldVisibility = map[string]string{
		"profile":           visProfile,
		"salary":            visSalary,
		"transition_reason": visTransReason,
		"experience":        visExp,
		"education":         visEdu,
		"certifications":    visCert,
		"skills":            visSkills,
		"portfolio":         visPortfolio,
		"references":        visRef,
	}

	// Decrypt sensitive fields
	_ = transReasonEnc // transition_reason is not surfaced on the identity aggregate; no longer clobbers Bio
	_ = bioOptimized
	if emailEnc != "" {
		p.Identity.Email, _ = r.crypt.Decrypt(emailEnc)
	}
	if phoneEnc != "" {
		p.Identity.Phone, _ = r.crypt.Decrypt(phoneEnc)
	}
	if addressEnc != "" {
		p.Identity.Address, _ = r.crypt.Decrypt(addressEnc)
	}
	if salCurrEnc != "" {
		p.Preferences.SalaryCurrency, _ = r.crypt.Decrypt(salCurrEnc)
	}
	if salMinEnc != "" {
		dec, _ := r.crypt.Decrypt(salMinEnc)
		p.Preferences.SalaryMin, _ = strconv.Atoi(dec)
	}
	if salMaxEnc != "" {
		dec, _ := r.crypt.Decrypt(salMaxEnc)
		p.Preferences.SalaryMax, _ = strconv.Atoi(dec)
	}

	// Load sub-resources
	var err2 error
	if p.Experiences, err2 = r.loadExperiences(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Educations, err2 = r.loadEducations(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Certifications, err2 = r.loadCertifications(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Skills, err2 = r.loadSkills(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Projects, err2 = r.loadProjects(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Achievements, err2 = r.loadAchievements(ctx, userID); err2 != nil {
		return nil, err2
	}
	if p.Resumes, err2 = r.loadResumes(ctx, userID); err2 != nil {
		return nil, err2
	}

	// Languages (written by setLanguagesQ; previously never read back)
	p.Identity.Languages, _ = r.loadLanguages(ctx, userID)

	// Load string slices
	p.Summary.CareerHighlights, _ = r.loadHighlights(ctx, userID)
	p.Summary.FunctionalAreas, _ = r.loadFunctionalAreas(ctx, userID)
	p.Summary.Industries, _ = r.loadIndustries(ctx, userID)
	p.Preferences.DesiredRoles, _ = r.loadDesiredRoles(ctx, userID)
	p.Preferences.DesiredIndustries, _ = r.loadDesiredIndustries(ctx, userID)
	p.Preferences.PreferredCountries, _ = r.loadPreferredCountries(ctx, userID)
	p.Preferences.PreferredCities, _ = r.loadPreferredCities(ctx, userID)
	p.Preferences.CompanySizePreferences, _ = r.loadCompanySizePreferences(ctx, userID)

	// Networking & Analytics summary loading
	p.Networking.ConnectionsCount, _ = r.countConnections(ctx, userID)
	p.Networking.FollowersCount, _ = r.countFollowers(ctx, userID)
	p.Networking.FollowingCount, _ = r.countFollowing(ctx, userID)
	p.Networking.Recommendations, _ = r.loadRecommendations(ctx, userID)

	p.Analytics.ProfileViews, _ = r.countAnalyticsEvents(ctx, userID, "view")
	p.Analytics.SearchAppearances, _ = r.countAnalyticsEvents(ctx, userID, "search_appearance")
	p.Analytics.RecruiterViews, _ = r.countAnalyticsEvents(ctx, userID, "recruiter_view")
	p.Analytics.ResumeDownloads, _ = r.countAnalyticsEvents(ctx, userID, "resume_download")

	return p, nil
}

// Sub-resource loaders
func (r *Repository) loadExperiences(ctx context.Context, userID string) ([]domain.WorkExperience, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company, COALESCE(company_logo,''), position, COALESCE(employment_type,''),
		       COALESCE(location,''), COALESCE(remote_type,''), start_date, end_date, is_current,
		       COALESCE(responsibilities,''), achievements, kpis, technologies, skills_used,
		       team_size, attachments
		FROM work_experiences WHERE user_id = $1 ORDER BY start_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.WorkExperience
	for rows.Next() {
		var e domain.WorkExperience
		var start, end sql.NullTime
		// achievements/kpis/technologies/skills_used/attachments are Postgres
		// text[] columns; the pgx stdlib driver surfaces them as the raw array
		// literal, so pq.Array parses them back into []string.
		if err := rows.Scan(
			&e.ID, &e.Company, &e.CompanyLogo, &e.Position, &e.EmploymentType,
			&e.Location, &e.RemoteType, &start, &end, &e.IsCurrent,
			&e.Responsibilities, pq.Array(&e.Achievements), pq.Array(&e.KPIs), pq.Array(&e.Technologies), pq.Array(&e.SkillsUsed),
			&e.TeamSize, pq.Array(&e.Attachments),
		); err != nil {
			return nil, err
		}
		if start.Valid {
			e.StartDate = start.Time
		}
		if end.Valid {
			e.EndDate = end.Time
		}

		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadEducations(ctx context.Context, userID string) ([]domain.Education, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, school, COALESCE(degree,''), COALESCE(field_of_study,''), start_date, end_date,
		       COALESCE(grade,''), COALESCE(description,''), COALESCE(major,''), COALESCE(minor,''),
		       gpa, COALESCE(honors,''), COALESCE(activities,''), COALESCE(projects,''),
		       COALESCE(research,''), COALESCE(thesis,''), graduation_date, COALESCE(verification_status,'unverified')
		FROM educations WHERE user_id = $1 ORDER BY start_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Education
	for rows.Next() {
		var e domain.Education
		var start, end, grad sql.NullTime
		var tempGrade, tempDesc string
		if err := rows.Scan(
			&e.ID, &e.Institution, &e.Degree, &e.FieldOfStudy, &start, &end,
			&tempGrade, &tempDesc, &e.Major, &e.Minor,
			&e.GPA, &e.Honors, &e.Activities, &e.Projects,
			&e.Research, &e.Thesis, &grad, &e.VerificationStatus,
		); err != nil {
			return nil, err
		}
		if grad.Valid {
			e.GraduationDate = grad.Time
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadCertifications(ctx context.Context, userID string) ([]domain.CertificationItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(issuer,''), COALESCE(credential_id,''), COALESCE(credential_url,''),
		       skills_covered, issue_date, expiry_date, COALESCE(status,'active')
		FROM certifications WHERE user_id = $1 ORDER BY issue_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.CertificationItem
	for rows.Next() {
		var c domain.CertificationItem
		var issue, expiry sql.NullTime
		var skills []byte
		if err := rows.Scan(&c.ID, &c.Name, &c.Issuer, &c.CredentialID, &c.VerificationURL, &skills, &issue, &expiry, &c.Status); err != nil {
			return nil, err
		}
		if issue.Valid {
			c.IssueDate = issue.Time
		}
		if expiry.Valid {
			c.ExpirationDate = expiry.Time
		}
		_ = json.Unmarshal(skills, &c.SkillsCovered)
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadSkills(ctx context.Context, userID string) ([]domain.SkillItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.name, COALESCE(ps.category,''), COALESCE(ps.proficiency_level,''), ps.years_of_experience,
		       COALESCE(ps.last_used,0), ps.verified, ps.recruiter_demand_score, ps.ai_recommendation_score
		FROM profile_skills ps JOIN skills s ON s.id = ps.skill_id
		WHERE ps.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.SkillItem
	for rows.Next() {
		var s domain.SkillItem
		if err := rows.Scan(&s.Name, &s.Category, &s.Level, &s.YearsOfExperience, &s.LastUsed, &s.Verified, &s.RecruiterDemandScore, &s.AIRecommendationScore); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadProjects(ctx context.Context, userID string) ([]domain.ProjectItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, COALESCE(description,''), COALESCE(repository_url,''), COALESCE(live_demo_url,''),
		       COALESCE(video_url,''), screenshots, technologies, COALESCE(timeline,''),
		       team_size, COALESCE(metrics,''), COALESCE(awards,''), COALESCE(business_impact,'')
		FROM profile_projects WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.ProjectItem
	for rows.Next() {
		var p domain.ProjectItem
		var teamSize int // team_size is an int column; TeamMembers is not persisted
		// screenshots/technologies are Postgres text[] columns (see insertProjectQ);
		// pq.Array parses them back into []string, matching loadExperiences.
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.RepositoryURL, &p.LiveDemoURL,
			&p.VideoURL, pq.Array(&p.Images), pq.Array(&p.Technologies), &p.Timeline,
			&teamSize, &p.Metrics, &p.Awards, &p.BusinessImpact,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadLanguages(ctx context.Context, userID string) ([]domain.LanguageItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.name, COALESCE(pl.proficiency, '')
		FROM profile_languages pl
		JOIN languages l ON l.id = pl.language_id
		WHERE pl.user_id = $1
		ORDER BY l.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.LanguageItem
	for rows.Next() {
		var li domain.LanguageItem
		if err := rows.Scan(&li.Name, &li.Proficiency); err != nil {
			return nil, err
		}
		out = append(out, li)
	}
	return out, rows.Err()
}

func (r *Repository) loadAchievements(ctx context.Context, userID string) ([]domain.AchievementItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, COALESCE(issuer_or_org,''), date, category, COALESCE(description,''), COALESCE(evidence_url,'')
		FROM profile_achievements WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.AchievementItem
	for rows.Next() {
		var a domain.AchievementItem
		var d sql.NullTime
		if err := rows.Scan(&a.ID, &a.Title, &a.IssuerOrOrg, &d, &a.Category, &a.Description, &a.EvidenceURL); err != nil {
			return nil, err
		}
		if d.Valid {
			a.Date = d.Time
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadResumes(ctx context.Context, userID string) ([]domain.ResumeVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rv.id, rv.filename, rv.size_bytes, COALESCE(rs.overall, 0), rs.suggestions, rv.created_at
		FROM resume_versions rv
		JOIN resumes r ON r.id = rv.resume_id
		LEFT JOIN resume_scores rs ON rs.version_id = rv.id
		WHERE r.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.ResumeVersion
	for rows.Next() {
		var rv domain.ResumeVersion
		var sug []byte
		if err := rows.Scan(&rv.ID, &rv.Name, &rv.FileSize, &rv.ATSScore, &sug, &rv.UploadedAt); err != nil {
			return nil, err
		}
		rv.KeywordAnalysis = sug
		out = append(out, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Slice loaders
func (r *Repository) loadHighlights(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT unnest(career_highlights) FROM profiles WHERE user_id = $1`, userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadFunctionalAreas(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT unnest(functional_areas) FROM profiles WHERE user_id = $1`, userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadIndustries(ctx context.Context, userID string) ([]string, error) {
	return r.loadStringTable(ctx, `profile_desired_industries`, `industry`, userID)
}
func (r *Repository) loadDesiredRoles(ctx context.Context, userID string) ([]string, error) {
	return r.loadStringTable(ctx, `profile_desired_roles`, `role`, userID)
}
func (r *Repository) loadDesiredIndustries(ctx context.Context, userID string) ([]string, error) {
	return r.loadStringTable(ctx, `profile_desired_industries`, `industry`, userID)
}
func (r *Repository) loadPreferredCountries(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT unnest(preferred_countries) FROM profiles WHERE user_id = $1`, userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
func (r *Repository) loadPreferredCities(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT unnest(preferred_cities) FROM profiles WHERE user_id = $1`, userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
func (r *Repository) loadCompanySizePreferences(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT unnest(company_size_preferences) FROM profiles WHERE user_id = $1`, userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) loadStringTable(ctx context.Context, table, col, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE user_id = $1 ORDER BY %s`, col, table, col), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err == nil {
			out = append(out, s)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Counters
func (r *Repository) countConnections(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM connections WHERE (user_a_id = $1 OR user_b_id = $1) AND status = 'accepted'`, userID).Scan(&count)
	return count, err
}

func (r *Repository) countFollowers(ctx context.Context, userID string) (int, error) {
	return 0, nil // Simulated
}

func (r *Repository) countFollowing(ctx context.Context, userID string) (int, error) {
	return 0, nil // Simulated
}

func (r *Repository) loadRecommendations(ctx context.Context, userID string) ([]domain.EndorsementSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.full_name, e.relationship, e.text, e.created_at
		FROM profile_endorsements e
		JOIN users u ON u.id = e.from_user_id
		WHERE e.to_user_id = $1 ORDER BY e.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.EndorsementSummary
	for rows.Next() {
		var sum domain.EndorsementSummary
		if err := rows.Scan(&sum.FromUserName, &sum.Relationship, &sum.Text, &sum.CreatedAt); err == nil {
			out = append(out, sum)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) countAnalyticsEvents(ctx context.Context, userID string, eventType string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM profile_analytics_events WHERE profile_id = $1 AND event_type = $2`, userID, eventType).Scan(&count)
	return count, err
}

// Versions Snapshot implementation
func (r *Repository) CreateVersionSnapshot(ctx context.Context, userID string, version int, p *domain.Profile) error {
	snapshotJSON, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO profile_versions (user_id, version, snapshot)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, version) DO UPDATE SET snapshot = EXCLUDED.snapshot`,
		userID, version, snapshotJSON)
	return err
}

func (r *Repository) GetVersionSnapshot(ctx context.Context, userID string, version int) (*domain.Profile, error) {
	var snapshotJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT snapshot FROM profile_versions
		WHERE user_id = $1 AND version = $2`, userID, version).Scan(&snapshotJSON)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var p domain.Profile
	if err := json.Unmarshal(snapshotJSON, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListVersions(ctx context.Context, userID string) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT version FROM profile_versions
		WHERE user_id = $1 ORDER BY version DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// Auditing
func (r *Repository) WriteAuditLog(ctx context.Context, log *domain.AuditLogEntry) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO profile_audit_logs (user_id, section, action, actor_id, old_value, new_value, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		log.UserID, log.Section, log.Action, log.ActorID, log.OldValue, log.NewValue, log.IPAddress, log.UserAgent)
	return err
}

// Analytics (simulated triggers/lookups)
func (r *Repository) RecordAnalyticsEvent(ctx context.Context, profileID string, eventType string, actorID *string, ip, ua string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO profile_analytics_events (profile_id, event_type, actor_id, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)`,
		profileID, eventType, actorID, ip, ua)
	return err
}

func (r *Repository) GetAnalytics(ctx context.Context, profileID string) (*domain.AnalyticsSummary, error) {
	var summary domain.AnalyticsSummary
	summary.ProfileViews, _ = r.countAnalyticsEvents(ctx, profileID, "view")
	summary.SearchAppearances, _ = r.countAnalyticsEvents(ctx, profileID, "search_appearance")
	summary.RecruiterViews, _ = r.countAnalyticsEvents(ctx, profileID, "recruiter_view")
	summary.ResumeDownloads, _ = r.countAnalyticsEvents(ctx, profileID, "resume_download")
	summary.WeeklyProfileViews = []int{0, 0, 0, 0, 0, 0, 0}
	return &summary, nil
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
	case "email_verified":
		query = `UPDATE profiles SET email_verified = $2 WHERE user_id = $1`
	case "employment_verified":
		query = `UPDATE profiles SET employment_verified = $2 WHERE user_id = $1`
	case "education_verified":
		query = `UPDATE profiles SET education_verified = $2 WHERE user_id = $1`
	case "certification_verified":
		query = `UPDATE profiles SET certification_verified = $2 WHERE user_id = $1`
	default:
		return fmt.Errorf("invalid verification field: %s", field)
	}
	_, err := r.db.ExecContext(ctx, query, userID, verified)
	return err
}

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
