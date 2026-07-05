-- Migration 019: comprehensive profile fields and normalized sub-tables
-- Extends the profiles table and adds supporting tables for recovery, compliance/consent logs, desired preferences, willing-to-mentor, and alerts.

ALTER TABLE profiles ADD COLUMN IF NOT EXISTS pronouns text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS career_status text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS transition_reason_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS target_comeback_timeline text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS open_to_remote boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS open_to_relocation boolean NOT NULL DEFAULT false;

-- Job Preferences
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS employment_type text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS salary_min_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS salary_max_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS salary_currency_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS salary_visible boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS work_mode text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS availability_date date;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS notice_period text;

-- Referral & Verification
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS referral_eligible boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS phone_verified boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS linkedin_verified boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS id_verified boolean NOT NULL DEFAULT false;

-- AI Coach Integration
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS career_narrative text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS coaching_metadata text;

-- Work Auth & Mobility
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS work_auth_status text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS passport_nationality text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS driving_license_bool boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS driving_license_type text;

-- Communication & Accessibility
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS preferred_contact_channel text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS accessibility_needs text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS video_intro_url text;

-- Two-sided mentorship
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS willing_to_mentor boolean NOT NULL DEFAULT false;

-- System-calculated fields
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS avg_response_time_hours double precision NOT NULL DEFAULT 0.0;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS profile_completeness_score integer NOT NULL DEFAULT 0;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS last_active_at timestamptz NOT NULL DEFAULT now();

-- Consent & Compliance
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS background_check_consent boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS background_check_consent_at timestamptz;

-- Job Alerts
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS job_alert_frequency text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS job_alert_channel text;

-- Privacy settings (default to public for general info, private for sensitive info)
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_profile text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_salary text NOT NULL DEFAULT 'private';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_transition_reason text NOT NULL DEFAULT 'private';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_experience text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_education text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_certifications text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_skills text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_portfolio text NOT NULL DEFAULT 'public';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS visibility_references text NOT NULL DEFAULT 'private';

-- Alter skill association
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS proficiency_level text;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS endorsed_count integer NOT NULL DEFAULT 0;

-- Achievements sub-table for Work Experiences
CREATE TABLE IF NOT EXISTS work_experience_achievements (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  experience_id uuid NOT NULL REFERENCES work_experiences(id) ON DELETE CASCADE,
  achievement   text NOT NULL,
  sort_order    integer NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_work_experience_achievements_exp ON work_experience_achievements(experience_id);

-- Supports Needed multi-select table
CREATE TABLE IF NOT EXISTS profile_supports (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  support text NOT NULL,
  PRIMARY KEY (user_id, support)
);

-- Relocation Cities/Countries mobility table
CREATE TABLE IF NOT EXISTS profile_relocation_locations (
  user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  location text NOT NULL,
  PRIMARY KEY (user_id, location)
);

-- Desired Roles matching inputs
CREATE TABLE IF NOT EXISTS profile_desired_roles (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role    text NOT NULL,
  PRIMARY KEY (user_id, role)
);

-- Desired Industries matching inputs
CREATE TABLE IF NOT EXISTS profile_desired_industries (
  user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  industry text NOT NULL,
  PRIMARY KEY (user_id, industry)
);

-- Endorsements/Recommendations
CREATE TABLE IF NOT EXISTS profile_endorsements (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  to_user_id   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  from_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  relationship text NOT NULL,
  text         text NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profile_endorsements_to ON profile_endorsements(to_user_id);

-- References (Private, exposed on request to recruiters)
CREATE TABLE IF NOT EXISTS profile_references (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id               uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name                  text NOT NULL,
  relationship          text NOT NULL,
  contact_info          text NOT NULL,
  permission_to_contact boolean NOT NULL DEFAULT false,
  created_at            timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profile_references_user ON profile_references(user_id);

-- GDPR / DPDP Consent Log
CREATE TABLE IF NOT EXISTS profile_consent_logs (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  consent_type  text NOT NULL,
  target_entity text NOT NULL,
  consented     boolean NOT NULL,
  ip_address    text,
  user_agent    text,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profile_consent_logs_user ON profile_consent_logs(user_id);
