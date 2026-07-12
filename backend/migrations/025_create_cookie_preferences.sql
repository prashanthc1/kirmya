-- Migration 025: Create cookie_preferences table for the CMP system
CREATE TABLE IF NOT EXISTS cookie_preferences (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id            uuid REFERENCES users(id) ON DELETE CASCADE,
  anonymous_id       text,
  essential          boolean NOT NULL DEFAULT true,
  functional         boolean NOT NULL DEFAULT false,
  analytics          boolean NOT NULL DEFAULT false,
  marketing          boolean NOT NULL DEFAULT false,
  performance        boolean NOT NULL DEFAULT false,
  personalization    boolean NOT NULL DEFAULT false,
  ai_preferences     boolean NOT NULL DEFAULT false,
  consent_version    text NOT NULL DEFAULT '1.0',
  accepted_at        timestamptz NOT NULL DEFAULT now(),
  ip_address         text,
  country            text,
  user_agent         text,
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now()
);

-- Ensure a user has only one cookie preference record.
CREATE UNIQUE INDEX IF NOT EXISTS uq_cookie_pref_user ON cookie_preferences (user_id) WHERE user_id IS NOT NULL;
-- Optimize lookups for anonymous users.
CREATE INDEX IF NOT EXISTS idx_cookie_pref_anon ON cookie_preferences (anonymous_id) WHERE anonymous_id IS NOT NULL;
