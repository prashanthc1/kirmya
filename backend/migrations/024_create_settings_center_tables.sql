-- Migration 024: Kirmya Settings Center Schema Extensions
-- Adds support for accessibility preferences, AI personalization, learning preferences, unique usernames/custom urls, and cookie consent tracking.

-- 1. Extend user_settings table
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS font_size text NOT NULL DEFAULT 'medium' CHECK (font_size IN ('small','medium','large','extra-large'));
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS high_contrast boolean NOT NULL DEFAULT false;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS reduced_motion boolean NOT NULL DEFAULT false;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS compact_mode boolean NOT NULL DEFAULT false;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS default_landing_page text NOT NULL DEFAULT 'dashboard';
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS accessibility_keyboard_navigation boolean NOT NULL DEFAULT false;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS accessibility_screen_reader_improvements boolean NOT NULL DEFAULT false;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS accessibility_focus_indicators boolean NOT NULL DEFAULT false;

-- AI personalization preferences
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS enable_ai_assistant boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_job_recommendations boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_resume_suggestions boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_roadmap_suggestions boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_skill_gap_analysis boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_interview_prep boolean NOT NULL DEFAULT true;
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS ai_learning_recommendations boolean NOT NULL DEFAULT true;

-- Learning goals and preferences
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS learning_goals text[] NOT NULL DEFAULT '{}';
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS technologies_of_interest text[] NOT NULL DEFAULT '{}';
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS certification_goals text[] NOT NULL DEFAULT '{}';
ALTER TABLE user_settings ADD COLUMN IF NOT EXISTS learning_reminders boolean NOT NULL DEFAULT true;

-- 2. Create cookie_consents table
CREATE TABLE IF NOT EXISTS cookie_consents (
  user_id           uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  essential         boolean NOT NULL DEFAULT true,
  functional        boolean NOT NULL DEFAULT false,
  analytics         boolean NOT NULL DEFAULT false,
  ai_personalization boolean NOT NULL DEFAULT false,
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now()
);

-- 3. Add username unique identifier to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS username text;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username ON users (lower(username)) WHERE deleted_at IS NULL;

-- 4. Add custom_url unique slug to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS custom_url text;
CREATE UNIQUE INDEX IF NOT EXISTS uq_profiles_custom_url ON profiles (lower(custom_url)) WHERE deleted_at IS NULL;

-- 5. Expand oauth_accounts check constraint to support Apple, Microsoft, GitHub
ALTER TABLE oauth_accounts DROP CONSTRAINT IF EXISTS oauth_accounts_provider_check;
ALTER TABLE oauth_accounts ADD CONSTRAINT oauth_accounts_provider_check CHECK (provider IN ('google', 'linkedin', 'microsoft', 'apple', 'github'));
