// Package postgres implements settings/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/settings/domain"

	"github.com/lib/pq"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const settingsCols = `
	user_id, language, timezone, theme, email_digest,
	profile_visibility, show_email, discoverable, allow_messages,
	notif_email_jobs, notif_email_mentorship, notif_email_messages, notif_email_referrals,
	notif_inapp_jobs, notif_inapp_mentorship, notif_inapp_messages, notif_inapp_referrals,
	login_alerts,
	font_size, high_contrast, reduced_motion, compact_mode, default_landing_page,
	accessibility_keyboard_navigation, accessibility_screen_reader, accessibility_focus_indicators,
	enable_ai_assistant, ai_job_recommendations, ai_resume_suggestions, ai_roadmap_suggestions,
	ai_skill_gap_analysis, ai_interview_prep, ai_learning_recommendations,
	learning_goals, technologies_of_interest, certification_goals, learning_reminders,
	version, created_at, updated_at`

func scanSettings(s interface{ Scan(...any) error }) (domain.UserSettings, error) {
	var x domain.UserSettings
	n := &x.Notifications
	var goals, techs, certs []string
	err := s.Scan(
		&x.UserID, &x.Language, &x.Timezone, &x.Theme, &x.EmailDigest,
		&x.ProfileVisibility, &x.ShowEmail, &x.Discoverable, &x.AllowMessages,
		&n.EmailJobs, &n.EmailMentorship, &n.EmailMessages, &n.EmailReferrals,
		&n.InAppJobs, &n.InAppMentorship, &n.InAppMessages, &n.InAppReferrals,
		&x.LoginAlerts,
		&x.FontSize, &x.HighContrast, &x.ReducedMotion, &x.CompactMode, &x.DefaultLandingPage,
		&x.AccessibilityKeyboardNavigation, &x.AccessibilityScreenReader, &x.AccessibilityFocusIndicators,
		&x.EnableAIAssistant, &x.AIJobRecommendations, &x.AIResumeSuggestions, &x.AIRoadmapSuggestions,
		&x.AISkillGapAnalysis, &x.AIInterviewPrep, &x.AILearningRecommendations,
		pq.Array(&goals), pq.Array(&techs), pq.Array(&certs), &x.LearningReminders,
		&x.Version, &x.CreatedAt, &x.UpdatedAt,
	)
	x.LearningGoals = goals
	x.TechnologiesOfInterest = techs
	x.CertificationGoals = certs
	return x, err
}

func (r *Repository) Get(ctx context.Context, userID string) (*domain.UserSettings, error) {
	x, err := scanSettings(r.db.QueryRowContext(ctx, `SELECT `+settingsCols+` FROM user_settings WHERE user_id = $1`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *Repository) EnsureDefaults(ctx context.Context, userID string) (*domain.UserSettings, error) {
	x, err := scanSettings(r.db.QueryRowContext(ctx, `
		INSERT INTO user_settings (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING `+settingsCols, userID))
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *Repository) Update(ctx context.Context, s *domain.UserSettings) error {
	n := s.Notifications
	err := r.db.QueryRowContext(ctx, `
		UPDATE user_settings SET
			language = $2, timezone = $3, theme = $4, email_digest = $5,
			profile_visibility = $6, show_email = $7, discoverable = $8, allow_messages = $9,
			notif_email_jobs = $10, notif_email_mentorship = $11, notif_email_messages = $12, notif_email_referrals = $13,
			notif_inapp_jobs = $14, notif_inapp_mentorship = $15, notif_inapp_messages = $16, notif_inapp_referrals = $17,
			login_alerts = $18,
			font_size = $19, high_contrast = $20, reduced_motion = $21, compact_mode = $22, default_landing_page = $23,
			accessibility_keyboard_navigation = $24, accessibility_screen_reader = $25, accessibility_focus_indicators = $26,
			enable_ai_assistant = $27, ai_job_recommendations = $28, ai_resume_suggestions = $29, ai_roadmap_suggestions = $30,
			ai_skill_gap_analysis = $31, ai_interview_prep = $32, ai_learning_recommendations = $33,
			learning_goals = $34, technologies_of_interest = $35, certification_goals = $36, learning_reminders = $37,
			version = version + 1, updated_at = now()
		WHERE user_id = $1 AND version = $38
		RETURNING version, updated_at`,
		s.UserID, s.Language, s.Timezone, s.Theme, s.EmailDigest,
		s.ProfileVisibility, s.ShowEmail, s.Discoverable, s.AllowMessages,
		n.EmailJobs, n.EmailMentorship, n.EmailMessages, n.EmailReferrals,
		n.InAppJobs, n.InAppMentorship, n.InAppMessages, n.InAppReferrals,
		s.LoginAlerts,
		s.FontSize, s.HighContrast, s.ReducedMotion, s.CompactMode, s.DefaultLandingPage,
		s.AccessibilityKeyboardNavigation, s.AccessibilityScreenReader, s.AccessibilityFocusIndicators,
		s.EnableAIAssistant, s.AIJobRecommendations, s.AIResumeSuggestions, s.AIRoadmapSuggestions,
		s.AISkillGapAnalysis, s.AIInterviewPrep, s.AILearningRecommendations,
		pq.Array(s.LearningGoals), pq.Array(s.TechnologiesOfInterest), pq.Array(s.CertificationGoals), s.LearningReminders,
		s.Version,
	).Scan(&s.Version, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrOptimisticLock
	}
	return err
}

// Connected Accounts
func (r *Repository) ListConnectedAccounts(ctx context.Context, userID string) ([]domain.ConnectedAccount, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, provider, provider_uid, created_at
		FROM oauth_accounts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ConnectedAccount
	for rows.Next() {
		var a domain.ConnectedAccount
		if err := rows.Scan(&a.ID, &a.Provider, &a.ProviderUID, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) DisconnectAccount(ctx context.Context, userID string, provider string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM oauth_accounts WHERE user_id = $1 AND provider = $2`, userID, provider)
	return err
}

// Cookie Consent
func (r *Repository) GetCookieConsent(ctx context.Context, userID string) (*domain.CookieConsent, error) {
	var cc domain.CookieConsent
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, essential, functional, analytics, ai_personalization, created_at, updated_at
		FROM cookie_consents WHERE user_id = $1`, userID).Scan(
		&cc.UserID, &cc.Essential, &cc.Functional, &cc.Analytics, &cc.AIPersonalization, &cc.CreatedAt, &cc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		// Return a default record
		return &domain.CookieConsent{
			UserID:    userID,
			Essential: true,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &cc, nil
}

func (r *Repository) SaveCookieConsent(ctx context.Context, cc *domain.CookieConsent) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cookie_consents (user_id, essential, functional, analytics, ai_personalization, updated_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (user_id) DO UPDATE SET
			essential = EXCLUDED.essential,
			functional = EXCLUDED.functional,
			analytics = EXCLUDED.analytics,
			ai_personalization = EXCLUDED.ai_personalization,
			updated_at = now()`,
		cc.UserID, cc.Essential, cc.Functional, cc.Analytics, cc.AIPersonalization,
	)
	return err
}

// Active Sessions
func (r *Repository) ListActiveSessions(ctx context.Context, userID string) ([]domain.ActiveSession, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, COALESCE(user_agent, 'Unknown Device'), COALESCE(ip, 'Unknown IP'), created_at, expires_at
		FROM refresh_tokens
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ActiveSession
	for rows.Next() {
		var s domain.ActiveSession
		if err := rows.Scan(&s.ID, &s.UserAgent, &s.IPAddress, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) RevokeSession(ctx context.Context, userID string, tokenID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE user_id = $1 AND id = $2`, userID, tokenID)
	return err
}

// Security History / Audit Logs
func (r *Repository) ListSecurityHistory(ctx context.Context, userID string) ([]domain.SecurityHistoryEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, action, COALESCE(ip, ''), created_at
		FROM audit_logs
		WHERE actor_id = $1 AND (target_type = 'security' OR action LIKE 'auth.%' OR action LIKE 'security.%')
		ORDER BY created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.SecurityHistoryEntry
	for rows.Next() {
		var h domain.SecurityHistoryEntry
		if err := rows.Scan(&h.ID, &h.Action, &h.IPAddress, &h.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) WriteSecurityLog(ctx context.Context, userID string, action string, ip string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, ip)
		VALUES ($1, $2, 'security', $3)`, userID, action, ip)
	return err
}

func (r *Repository) GetProfileSettings(ctx context.Context, userID string) (username, customURL, profileVisibility string, fieldVisibility map[string]string, openToWork, referralEligible, willingToMentor bool, err error) {
	// 1. Ensure profile exists
	_, err = r.db.ExecContext(ctx, `INSERT INTO profiles (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`, userID)
	if err != nil {
		return
	}

	// 2. Fetch username
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(username, '') FROM users WHERE id = $1`, userID).Scan(&username)
	if err != nil {
		return
	}

	// 3. Fetch profile settings
	var vs, vt, vex, ved, vc, vsks, vp, vr string
	var openRemote, openReloc bool
	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(custom_url, ''), visibility_profile, visibility_salary, visibility_transition_reason,
		       visibility_experience, visibility_education, visibility_certifications, visibility_skills,
		       visibility_portfolio, visibility_references, open_to_remote, open_to_relocation,
		       referral_eligible, willing_to_mentor
		FROM profiles WHERE user_id = $1`, userID).Scan(
		&customURL, &profileVisibility, &vs, &vt, &vex, &ved, &vc, &vsks, &vp, &vr, &openRemote, &openReloc, &referralEligible, &willingToMentor,
	)
	if err != nil {
		return
	}

	openToWork = openRemote || openReloc
	fieldVisibility = map[string]string{
		"salary":            vs,
		"transition_reason": vt,
		"experience":        vex,
		"education":         ved,
		"certifications":    vc,
		"skills":            vsks,
		"portfolio":         vp,
		"references":        vr,
	}
	return
}

func (r *Repository) UpdateProfileSettings(ctx context.Context, userID string, username, customURL, profileVisibility string, fieldVisibility map[string]string, openToWork, referralEligible, willingToMentor bool) error {
	// 1. Update username
	var dbUsername *string
	if username != "" {
		dbUsername = &username
	}
	_, err := r.db.ExecContext(ctx, `UPDATE users SET username = $2 WHERE id = $1`, userID, dbUsername)
	if err != nil {
		return err
	}

	// 2. Extract field visibility
	vs := fieldVisibility["salary"]
	vt := fieldVisibility["transition_reason"]
	vex := fieldVisibility["experience"]
	ved := fieldVisibility["education"]
	vc := fieldVisibility["certifications"]
	vsks := fieldVisibility["skills"]
	vp := fieldVisibility["portfolio"]
	vr := fieldVisibility["references"]

	// Set defaults if empty
	if vs == "" {
		vs = "private"
	}
	if vt == "" {
		vt = "private"
	}
	if vex == "" {
		vex = "public"
	}
	if ved == "" {
		ved = "public"
	}
	if vc == "" {
		vc = "public"
	}
	if vsks == "" {
		vsks = "public"
	}
	if vp == "" {
		vp = "public"
	}
	if vr == "" {
		vr = "private"
	}

	var dbCustomURL *string
	if customURL != "" {
		dbCustomURL = &customURL
	}

	// 3. Update profiles
	_, err = r.db.ExecContext(ctx, `
		UPDATE profiles SET
			custom_url = $2,
			visibility_profile = $3,
			visibility_salary = $4,
			visibility_transition_reason = $5,
			visibility_experience = $6,
			visibility_education = $7,
			visibility_certifications = $8,
			visibility_skills = $9,
			visibility_portfolio = $10,
			visibility_references = $11,
			open_to_remote = $12,
			open_to_relocation = $12,
			referral_eligible = $13,
			willing_to_mentor = $14,
			updated_at = now()
		WHERE user_id = $1`,
		userID, dbCustomURL, profileVisibility, vs, vt, vex, ved, vc, vsks, vp, vr, openToWork, referralEligible, willingToMentor,
	)
	return err
}
