-- Resume module: a resume has an ordered history of versions; each version is
-- parsed to text and scored (formatting / keywords / ATS) with suggestions.

CREATE TABLE IF NOT EXISTS resumes (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title      text NOT NULL DEFAULT 'My Resume',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);
CREATE INDEX IF NOT EXISTS idx_resumes_user ON resumes(user_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS resume_versions (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id      uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  version_no     integer NOT NULL,
  filename       text NOT NULL,
  content_type   text NOT NULL,
  size_bytes     bigint NOT NULL DEFAULT 0,
  storage_key    text NOT NULL,
  extracted_text text NOT NULL DEFAULT '',
  created_at     timestamptz NOT NULL DEFAULT now(),
  UNIQUE (resume_id, version_no)
);
CREATE INDEX IF NOT EXISTS idx_resume_versions_resume ON resume_versions(resume_id);

CREATE TABLE IF NOT EXISTS resume_scores (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  version_id  uuid NOT NULL UNIQUE REFERENCES resume_versions(id) ON DELETE CASCADE,
  overall     integer NOT NULL,
  formatting  integer NOT NULL,
  keywords    integer NOT NULL,
  ats         integer NOT NULL,
  suggestions jsonb NOT NULL DEFAULT '[]',
  created_at  timestamptz NOT NULL DEFAULT now()
);
