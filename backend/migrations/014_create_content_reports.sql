-- Admin bounded context: the content-moderation report queue. Any authenticated
-- user can file a report against a post/comment/user/message; admins triage them
-- (resolve or dismiss) and may remove the offending content. Admin user- and
-- content-management actions reuse existing tables (users, posts, comments) and
-- are recorded in audit_logs.

CREATE TABLE IF NOT EXISTS content_reports (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_type  text NOT NULL CHECK (target_type IN ('post','comment','user','message')),
  target_id    uuid NOT NULL,
  reason       text NOT NULL,
  status       text NOT NULL DEFAULT 'open'
                 CHECK (status IN ('open','reviewing','resolved','dismissed')),
  action_taken text,
  resolved_by  uuid REFERENCES users(id) ON DELETE SET NULL,
  resolved_at  timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reports_status ON content_reports(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reports_target ON content_reports(target_type, target_id);
