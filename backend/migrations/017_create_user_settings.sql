-- Settings bounded context: one preferences row per user covering general,
-- privacy, notification, and security-preference settings. Created lazily on
-- first read/write (UPSERT) so existing users need no backfill. The `version`
-- column powers optimistic locking, matching the platform convention.

CREATE TABLE IF NOT EXISTS user_settings (
  user_id              uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

  -- General
  language             text    NOT NULL DEFAULT 'en',
  timezone             text    NOT NULL DEFAULT 'UTC',
  theme                text    NOT NULL DEFAULT 'system'   CHECK (theme IN ('light','dark','system')),
  email_digest         text    NOT NULL DEFAULT 'weekly'   CHECK (email_digest IN ('off','daily','weekly')),

  -- Privacy
  profile_visibility   text    NOT NULL DEFAULT 'public'   CHECK (profile_visibility IN ('public','network','private')),
  show_email           boolean NOT NULL DEFAULT false,
  discoverable         boolean NOT NULL DEFAULT true,
  allow_messages       text    NOT NULL DEFAULT 'everyone' CHECK (allow_messages IN ('everyone','network','none')),

  -- Notifications (per channel x category)
  notif_email_jobs       boolean NOT NULL DEFAULT true,
  notif_email_mentorship boolean NOT NULL DEFAULT true,
  notif_email_messages   boolean NOT NULL DEFAULT true,
  notif_email_referrals  boolean NOT NULL DEFAULT true,
  notif_inapp_jobs       boolean NOT NULL DEFAULT true,
  notif_inapp_mentorship boolean NOT NULL DEFAULT true,
  notif_inapp_messages   boolean NOT NULL DEFAULT true,
  notif_inapp_referrals  boolean NOT NULL DEFAULT true,

  -- Security preferences
  login_alerts         boolean NOT NULL DEFAULT true,

  version              integer     NOT NULL DEFAULT 1,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz NOT NULL DEFAULT now()
);
