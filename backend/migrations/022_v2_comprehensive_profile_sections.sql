-- Migration 022: Kirmya v2 15-Section Profile Schema
-- Extends the profiles table and child tables to completely represent the 15-Section Career OS Profile Aggregate.

-- 1. Identity & Summary extensions (Sections 1 & 2)
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS preferred_name text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS timezone text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS nationality text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS bio_optimized text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS executive_summary text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS career_objectives text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS career_highlights text[];
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS functional_areas text[];
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS personal_brand_statement text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS elevator_pitch text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS is_draft boolean NOT NULL DEFAULT true;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS trust_score integer NOT NULL DEFAULT 0;

-- 2. Work Experience extensions (Section 3)
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS company_logo text;
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS remote_type text;
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS responsibilities text;
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS achievements text[];
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS kpis text[];
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS technologies text[];
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS skills_used text[];
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS team_size integer NOT NULL DEFAULT 1;
ALTER TABLE work_experiences ADD COLUMN IF NOT EXISTS attachments text[];

-- 3. Education extensions (Section 4)
ALTER TABLE educations ADD COLUMN IF NOT EXISTS major text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS minor text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS gpa double precision;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS honors text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS activities text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS projects text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS research text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS thesis text;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS graduation_date date;
ALTER TABLE educations ADD COLUMN IF NOT EXISTS verification_status text NOT NULL DEFAULT 'unverified';

-- 4. Skills extensions (Section 5)
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS category text;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS years_of_experience double precision NOT NULL DEFAULT 0.0;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS last_used integer;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS verified boolean NOT NULL DEFAULT false;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS recruiter_demand_score double precision NOT NULL DEFAULT 0.0;
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS ai_recommendation_score double precision NOT NULL DEFAULT 0.0;

-- 5. Profile Projects (Section 6)
CREATE TABLE IF NOT EXISTS profile_projects (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title          text NOT NULL,
  description    text,
  cover_image    text,
  repository_url text,
  live_demo_url  text,
  video_url      text,
  screenshots    text[],
  technologies   text[],
  timeline       text,
  team_size      integer NOT NULL DEFAULT 1,
  metrics        text,
  awards         text,
  business_impact text,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_projects_user ON profile_projects(user_id);

-- 6. Certifications extensions (Section 7)
ALTER TABLE certifications ADD COLUMN IF NOT EXISTS skills_covered text[];
ALTER TABLE certifications ADD COLUMN IF NOT EXISTS status text NOT NULL DEFAULT 'active';

-- 7. Profile Achievements (Section 8)
CREATE TABLE IF NOT EXISTS profile_achievements (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title          text NOT NULL,
  issuer_or_org  text,
  date           date,
  category       text NOT NULL CHECK (category IN ('award', 'patent', 'publication', 'conference', 'hackathon', 'volunteer', 'open_source', 'speaking')),
  description    text,
  evidence_url   text,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_achievements_user ON profile_achievements(user_id);

-- 8. Career Preferences extensions (Section 10)
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS company_size_preferences text[];
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS preferred_countries text[];
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS preferred_cities text[];
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS travel_willingness text;

-- 9. Verification & Trust extensions (Section 11)
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS email_verified boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS employment_verified boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS education_verified boolean NOT NULL DEFAULT false;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS certification_verified boolean NOT NULL DEFAULT false;

-- 10. Profile Versions & History Snapshots
CREATE TABLE IF NOT EXISTS profile_versions (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  version        integer NOT NULL,
  snapshot       jsonb NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT now(),
  created_by     uuid REFERENCES users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_versions_user_ver ON profile_versions(user_id, version);

-- 11. Profile Audit Logs
CREATE TABLE IF NOT EXISTS profile_audit_logs (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  section        text NOT NULL,
  action         text NOT NULL,
  actor_id       uuid NOT NULL,
  old_value      jsonb,
  new_value      jsonb,
  ip_address     text,
  user_agent     text,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profile_audit_logs_user ON profile_audit_logs(user_id);

-- 12. Profile Analytics Events (Section 13)
CREATE TABLE IF NOT EXISTS profile_analytics_events (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  profile_id     uuid NOT NULL REFERENCES profiles(user_id) ON DELETE CASCADE,
  event_type     text NOT NULL,
  actor_id       uuid REFERENCES users(id) ON DELETE SET NULL,
  ip_address     text,
  user_agent     text,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profile_analytics_events_prof ON profile_analytics_events(profile_id);
