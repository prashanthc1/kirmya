package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/cookies/domain"
)

// Repository implements domain.Repository on PostgreSQL.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func scanRow(row *sql.Row) (*domain.CookiePreferences, error) {
	var p domain.CookiePreferences
	var uID, aID sql.NullString
	var ip, cty, ua sql.NullString

	err := row.Scan(
		&p.ID, &uID, &aID, &p.Essential, &p.Functional, &p.Analytics,
		&p.Marketing, &p.Performance, &p.Personalization, &p.AIPreferences,
		&p.ConsentVersion, &p.AcceptedAt, &ip, &cty, &ua, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if uID.Valid {
		p.UserID = &uID.String
	}
	if aID.Valid {
		p.AnonymousID = &aID.String
	}
	p.IPAddress = ip.String
	p.Country = cty.String
	p.UserAgent = ua.String

	return &p, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) (*domain.CookiePreferences, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, anonymous_id, essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, accepted_at, ip_address, country, user_agent, created_at, updated_at
		FROM cookie_preferences
		WHERE user_id = $1`, userID)
	p, err := scanRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (r *Repository) GetByAnonymousID(ctx context.Context, anonymousID string) (*domain.CookiePreferences, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, anonymous_id, essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, accepted_at, ip_address, country, user_agent, created_at, updated_at
		FROM cookie_preferences
		WHERE anonymous_id = $1`, anonymousID)
	p, err := scanRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (r *Repository) Save(ctx context.Context, p *domain.CookiePreferences) error {
	if p.UserID != nil {
		// Logged in user: insert or update on conflict of user_id
		err := r.db.QueryRowContext(ctx, `
			INSERT INTO cookie_preferences (
				user_id, anonymous_id, essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, accepted_at, ip_address, country, user_agent, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now())
			ON CONFLICT (user_id) DO UPDATE SET
				anonymous_id = EXCLUDED.anonymous_id,
				essential = EXCLUDED.essential,
				functional = EXCLUDED.functional,
				analytics = EXCLUDED.analytics,
				marketing = EXCLUDED.marketing,
				performance = EXCLUDED.performance,
				personalization = EXCLUDED.personalization,
				ai_preferences = EXCLUDED.ai_preferences,
				consent_version = EXCLUDED.consent_version,
				accepted_at = EXCLUDED.accepted_at,
				ip_address = EXCLUDED.ip_address,
				country = EXCLUDED.country,
				user_agent = EXCLUDED.user_agent,
				updated_at = now()
			RETURNING id, created_at, updated_at`,
			p.UserID, p.AnonymousID, p.Essential, p.Functional, p.Analytics, p.Marketing, p.Performance, p.Personalization, p.AIPreferences, p.ConsentVersion, p.AcceptedAt, p.IPAddress, p.Country, p.UserAgent,
		).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
		return err
	}

	if p.AnonymousID == nil {
		return errors.New("cannot save preferences without user_id or anonymous_id")
	}

	// Anonymous user: check if exists
	var id string
	err := r.db.QueryRowContext(ctx, `SELECT id FROM cookie_preferences WHERE anonymous_id = $1`, p.AnonymousID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		// Insert new anonymous record
		err = r.db.QueryRowContext(ctx, `
			INSERT INTO cookie_preferences (
				anonymous_id, essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, accepted_at, ip_address, country, user_agent, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now())
			RETURNING id, created_at, updated_at`,
			p.AnonymousID, p.Essential, p.Functional, p.Analytics, p.Marketing, p.Performance, p.Personalization, p.AIPreferences, p.ConsentVersion, p.AcceptedAt, p.IPAddress, p.Country, p.UserAgent,
		).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
		return err
	} else if err != nil {
		return err
	}

	// Update existing anonymous record
	err = r.db.QueryRowContext(ctx, `
		UPDATE cookie_preferences SET
			essential = $2,
			functional = $3,
			analytics = $4,
			marketing = $5,
			performance = $6,
			personalization = $7,
			ai_preferences = $8,
			consent_version = $9,
			accepted_at = $10,
			ip_address = $11,
			country = $12,
			user_agent = $13,
			updated_at = now()
		WHERE anonymous_id = $1
		RETURNING id, created_at, updated_at`,
		p.AnonymousID, p.Essential, p.Functional, p.Analytics, p.Marketing, p.Performance, p.Personalization, p.AIPreferences, p.ConsentVersion, p.AcceptedAt, p.IPAddress, p.Country, p.UserAgent,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return err
}

func (r *Repository) Delete(ctx context.Context, userID string, anonymousID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM cookie_preferences
		WHERE (user_id = $1 AND $1 <> '') OR (anonymous_id = $2 AND $2 <> '')`, userID, anonymousID)
	return err
}

func (r *Repository) Merge(ctx context.Context, anonymousID string, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Fetch anonymous record
	var essential, functional, analytics, marketing, performance, personalization, aiPref bool
	var consentVersion, ip, country, ua string
	err = tx.QueryRowContext(ctx, `
		SELECT essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, COALESCE(ip_address, ''), COALESCE(country, ''), COALESCE(user_agent, '')
		FROM cookie_preferences
		WHERE anonymous_id = $1`, anonymousID).Scan(
		&essential, &functional, &analytics, &marketing, &performance, &personalization, &aiPref, &consentVersion, &ip, &country, &ua,
	)
	if errors.Is(err, sql.ErrNoRows) {
		// Nothing to merge
		return nil
	} else if err != nil {
		return err
	}

	// 2. Insert or update onto the user record
	_, err = tx.ExecContext(ctx, `
		INSERT INTO cookie_preferences (
			user_id, essential, functional, analytics, marketing, performance, personalization, ai_preferences, consent_version, accepted_at, ip_address, country, user_agent, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), $10, $11, $12, now())
		ON CONFLICT (user_id) DO UPDATE SET
			essential = EXCLUDED.essential,
			functional = EXCLUDED.functional,
			analytics = EXCLUDED.analytics,
			marketing = EXCLUDED.marketing,
			performance = EXCLUDED.performance,
			personalization = EXCLUDED.personalization,
			ai_preferences = EXCLUDED.ai_preferences,
			consent_version = EXCLUDED.consent_version,
			accepted_at = EXCLUDED.accepted_at,
			ip_address = EXCLUDED.ip_address,
			country = EXCLUDED.country,
			user_agent = EXCLUDED.user_agent,
			updated_at = now()`,
		userID, essential, functional, analytics, marketing, performance, personalization, aiPref, consentVersion, ip, country, ua,
	)
	if err != nil {
		return err
	}

	// 3. Delete the anonymous record
	_, err = tx.ExecContext(ctx, `DELETE FROM cookie_preferences WHERE anonymous_id = $1`, anonymousID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
