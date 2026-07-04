-- Saved jobs (bookmarks) and a richer application status set for the Jobs module.

CREATE TABLE IF NOT EXISTS saved_jobs (
  user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  job_id   uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  saved_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, job_id)
);

CREATE INDEX IF NOT EXISTS idx_saved_jobs_user ON saved_jobs(user_id);
