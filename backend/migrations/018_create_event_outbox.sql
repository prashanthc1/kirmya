-- Migration: Create event_outbox table for transactional event outbox pattern.
-- This ensures existing deployments that already ran 001_create_users_table.sql
-- get the table created successfully.
CREATE TABLE IF NOT EXISTS event_outbox (
  id           bigserial PRIMARY KEY,
  aggregate_id uuid NOT NULL,
  event_type   text NOT NULL,
  payload      jsonb NOT NULL,
  occurred_at  timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_outbox_unpublished ON event_outbox(id) WHERE published_at IS NULL;
