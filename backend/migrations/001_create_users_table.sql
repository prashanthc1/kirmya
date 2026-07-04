-- Identity bounded context (PostgreSQL). Replaces the former MySQL users schema.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email          text NOT NULL,
  password_hash  text,
  full_name      text NOT NULL DEFAULT '',
  email_verified boolean NOT NULL DEFAULT false,
  status         text NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','deactivated')),
  mfa_enabled    boolean NOT NULL DEFAULT false,
  last_login_at  timestamptz,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  created_by     uuid,
  updated_by     uuid,
  deleted_at     timestamptz,
  version        integer NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email_active ON users (lower(email)) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS roles (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id    uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  granted_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS oauth_accounts (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider     text NOT NULL CHECK (provider IN ('google','linkedin')),
  provider_uid text NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_uid)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  text NOT NULL,
  family_id   uuid NOT NULL,
  expires_at  timestamptz NOT NULL,
  revoked_at  timestamptz,
  replaced_by uuid,
  created_at  timestamptz NOT NULL DEFAULT now(),
  user_agent  text,
  ip          text
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id) WHERE revoked_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_refresh_token_hash ON refresh_tokens(token_hash);

CREATE TABLE IF NOT EXISTS mfa_credentials (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type         text NOT NULL DEFAULT 'totp' CHECK (type IN ('totp')),
  secret_enc   text NOT NULL,
  confirmed_at timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS email_verification_tokens (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL,
  expires_at timestamptz NOT NULL,
  used_at    timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL,
  expires_at timestamptz NOT NULL,
  used_at    timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_id    uuid,
  action      text NOT NULL,
  target_type text,
  target_id   uuid,
  metadata    jsonb NOT NULL DEFAULT '{}',
  ip          text,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_actor_time ON audit_logs(actor_id, created_at DESC);

CREATE TABLE IF NOT EXISTS event_outbox (
  id           bigserial PRIMARY KEY,
  aggregate_id uuid NOT NULL,
  event_type   text NOT NULL,
  payload      jsonb NOT NULL,
  occurred_at  timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_outbox_unpublished ON event_outbox(id) WHERE published_at IS NULL;

INSERT INTO roles (name) VALUES ('job_seeker'), ('referrer'), ('mentor'), ('recruiter'), ('admin')
  ON CONFLICT (name) DO NOTHING;
