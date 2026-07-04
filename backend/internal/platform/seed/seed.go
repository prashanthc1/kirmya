// Package seed populates a fresh database with realistic demo data so a brand-new
// `docker compose up` is instantly explorable. It runs only when SEED_DEMO_DATA
// is "true" and is idempotent: if the demo admin already exists it does nothing.
// All demo accounts share the password "Password123!".
package seed

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"workspace-app/internal/identity/infrastructure/crypto"
)

const (
	sentinelEmail = "admin@demo.kirmya.io"
	demoPassword  = "Password123!"
)

// Run seeds demo data when enabled. It is safe to call on every startup.
func Run(db *sql.DB) error {
	if os.Getenv("SEED_DEMO_DATA") != "true" {
		return nil
	}

	var exists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE lower(email) = $1)`, sentinelEmail).Scan(&exists); err != nil {
		return fmt.Errorf("seed: check sentinel: %w", err)
	}
	if exists {
		log.Printf("[seed] demo data already present; skipping")
		return nil
	}

	hash, err := crypto.NewArgon2Hasher().Hash(demoPassword)
	if err != nil {
		return fmt.Errorf("seed: hash password: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	s := &seeder{tx: tx, hash: hash}
	if err := s.run(); err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("[seed] demo data inserted (login with any demo email / %q)", demoPassword)
	return nil
}

type seeder struct {
	tx   *sql.Tx
	hash string
	err  error // first error encountered; helpers become no-ops after it is set
}

func (s *seeder) run() error {
	// --- Users (email_verified, active) with roles -----------------------
	s.user("admin@demo.kirmya.io", "Demo Admin", "admin")
	asha := s.user("asha.rao@demo.kirmya.io", "Asha Rao", "job_seeker")
	ben := s.user("ben.carter@demo.kirmya.io", "Ben Carter", "job_seeker")
	carla := s.user("carla.mendes@demo.kirmya.io", "Carla Mendes", "referrer")
	deepa := s.user("deepa.nair@demo.kirmya.io", "Deepa Nair", "mentor")
	omar := s.user("omar.farouk@demo.kirmya.io", "Omar Farouk", "mentor")
	rita := s.user("rita.shah@demo.kirmya.io", "Rita Shah", "recruiter")

	if s.err != nil {
		return s.err
	}

	// --- Profiles + skills ----------------------------------------------
	s.profile(asha, "Senior Operations Coordinator → aiming for Operations Manager", "10 years keeping facilities and teams running. Recently laid off; ready for the next step.", "Dubai, UAE")
	s.skills(asha, "Operations", "Facilities Management", "Vendor Management", "Budgeting")

	s.profile(ben, "Logistics Analyst | Supply Chain", "Data-driven logistics analyst focused on cost-to-serve and route optimization.", "Pune, India")
	s.skills(ben, "Logistics", "SQL", "Supply Chain", "Analytics")

	s.profile(carla, "Facilities Manager at BuildCo", "Happy to refer strong candidates into BuildCo's operations org.", "Lisbon, Portugal")
	s.skills(carla, "Facilities Management", "Leadership", "Health & Safety")

	s.profile(deepa, "Engineering Leader | Mentor", "Helping engineers grow into leadership. 15 years across startups and big tech.", "Bengaluru, India")
	s.skills(deepa, "Leadership", "System Design", "Go", "Career Coaching")

	s.profile(omar, "Career Coach & HR Director", "I coach mid-career professionals through transitions and interview prep.", "Cairo, Egypt")
	s.skills(omar, "Career Coaching", "Interviewing", "HR", "Negotiation")

	s.profile(rita, "Talent Partner at Acme Logistics", "Hiring across operations, logistics, and facilities roles.", "Remote")

	// --- Jobs (posted by the recruiter and the referrer) -----------------
	s.job(rita, "Operations Manager", "Acme Logistics", "Dubai, UAE",
		"Lead a 20-person operations team across two warehouses. Own SLAs, budgets, and continuous improvement. PMP a plus.", "$70k–$90k", "full_time")
	s.job(rita, "Logistics Analyst", "Acme Logistics", "Remote",
		"Analyze route efficiency and cost-to-serve. Build dashboards and partner with operations to act on insights. SQL required.", "$45k–$60k", "full_time")
	s.job(rita, "Facilities Coordinator", "Acme Logistics", "Pune, India",
		"Coordinate maintenance, vendors, and health & safety across our Pune site. Great stepping stone toward facilities management.", "$25k–$35k", "full_time")
	s.job(carla, "Assistant Facilities Manager", "BuildCo", "Lisbon, Portugal",
		"Support the facilities manager across building operations, contractor management, and compliance. Referral-friendly role.", "€30k–€40k", "full_time")
	s.job(rita, "Operations Intern", "Acme Logistics", "Remote",
		"Six-month internship rotating across planning, warehouse ops, and analytics. Ideal for career switchers.", "Stipend", "internship")

	// --- Mentor profiles -------------------------------------------------
	s.mentor(deepa, "Engineering Leader | Mentor", "I help engineers and analysts grow into leadership and navigate layoffs.", "Leadership, System Design, Career Growth")
	s.mentor(omar, "Career Coach & HR Director", "Interview prep, resume reviews, and salary negotiation for career changers.", "Interviewing, HR, Negotiation")

	// --- Communities: memberships, posts, comments -----------------------
	ops := s.community("operations")
	logi := s.community("logistics")
	tech := s.community("technology")

	s.join(ops, asha, "moderator")
	s.join(ops, carla, "member")
	s.join(logi, ben, "member")
	s.join(logi, rita, "member")
	s.join(tech, deepa, "moderator")

	p1 := s.post(ops, asha, "Bouncing back after a layoff — what worked for you?",
		"I was let go last month after 10 years. Sharing what's helping: updating my resume, asking for warm referrals, and one mentor call a week. What worked for you?")
	s.comment(p1, carla, "Warm referrals made the biggest difference for me. Happy to refer for BuildCo roles.")
	s.comment(p1, omar, "Block 30 minutes a day for outreach. Momentum compounds.")

	p2 := s.post(logi, ben, "Route optimization: build vs. buy?",
		"Curious how teams here approach route optimization — in-house models or off-the-shelf tools? Trade-offs?")
	s.comment(p2, rita, "We started off-the-shelf, then built once volume justified it.")

	s.post(tech, deepa, "From analyst to engineering lead — AMA",
		"I've made the jump from analyst to engineering leadership. Ask me anything about the transition, interviews, or skill gaps.")

	// --- One referral request to show the recovery loop ------------------
	s.referral(asha, carla, "BuildCo", "I'd love a referral for the Assistant Facilities Manager role — 10 years in facilities ops.")

	return s.err
}

// ----- Insert helpers (record the first error and become no-ops after) ---

func (s *seeder) user(email, name, role string) string {
	if s.err != nil {
		return ""
	}
	var id string
	s.err = s.tx.QueryRow(`
		INSERT INTO users (email, password_hash, full_name, email_verified, status)
		VALUES ($1, $2, $3, true, 'active') RETURNING id`, email, s.hash, name).Scan(&id)
	if s.err != nil {
		return ""
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2 ON CONFLICT DO NOTHING`, id, role)
	return id
}

func (s *seeder) profile(userID, headline, about, location string) {
	if s.err != nil || userID == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO profiles (user_id, headline, about, location)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET headline = EXCLUDED.headline, about = EXCLUDED.about, location = EXCLUDED.location`,
		userID, headline, about, location)
}

func (s *seeder) skills(userID string, names ...string) {
	for _, name := range names {
		if s.err != nil {
			return
		}
		var skillID string
		s.err = s.tx.QueryRow(`
			INSERT INTO skills (name) VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, name).Scan(&skillID)
		if s.err != nil {
			return
		}
		_, s.err = s.tx.Exec(`
			INSERT INTO profile_skills (user_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, skillID)
	}
}

func (s *seeder) job(postedBy, title, company, location, description, salary, jobType string) {
	if s.err != nil || postedBy == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO jobs (title, company, location, description, salary, job_type, posted_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, title, company, location, description, salary, jobType, postedBy)
}

func (s *seeder) mentor(userID, headline, bio, expertise string) {
	if s.err != nil || userID == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO mentor_profiles (user_id, headline, bio, expertise)
		VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO NOTHING`, userID, headline, bio, expertise)
}

func (s *seeder) community(slug string) string {
	if s.err != nil {
		return ""
	}
	var id string
	s.err = s.tx.QueryRow(`SELECT id FROM communities WHERE slug = $1`, slug).Scan(&id)
	return id
}

func (s *seeder) join(communityID, userID, role string) {
	if s.err != nil || communityID == "" || userID == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO community_members (community_id, user_id, role)
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, communityID, userID, role)
}

func (s *seeder) post(communityID, authorID, title, body string) string {
	if s.err != nil || communityID == "" || authorID == "" {
		return ""
	}
	var id string
	s.err = s.tx.QueryRow(`
		INSERT INTO posts (community_id, author_id, title, body)
		VALUES ($1, $2, $3, $4) RETURNING id`, communityID, authorID, title, body).Scan(&id)
	return id
}

func (s *seeder) comment(postID, authorID, body string) {
	if s.err != nil || postID == "" || authorID == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO comments (post_id, author_id, body) VALUES ($1, $2, $3)`, postID, authorID, body)
}

func (s *seeder) referral(seekerID, referrerID, company, message string) {
	if s.err != nil || seekerID == "" || referrerID == "" {
		return
	}
	_, s.err = s.tx.Exec(`
		INSERT INTO referral_requests (seeker_id, referrer_id, company, message, status)
		VALUES ($1, $2, $3, $4, 'requested')`, seekerID, referrerID, company, message)
}
