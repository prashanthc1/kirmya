-- Profile bounded context: extend the base `profiles` table and add the
-- experience / education / certification / skills / languages / portfolio
-- child tables. Child rows reference users(id) directly and cascade on delete.

ALTER TABLE profiles ADD COLUMN IF NOT EXISTS headline   text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS about      text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS photo_url  text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS version    integer NOT NULL DEFAULT 1;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

CREATE TABLE IF NOT EXISTS work_experiences (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title           text NOT NULL,
  company         text NOT NULL,
  location        text,
  employment_type text,
  start_date      date,
  end_date        date,
  is_current      boolean NOT NULL DEFAULT false,
  description     text,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_work_experiences_user ON work_experiences(user_id);

CREATE TABLE IF NOT EXISTS educations (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school         text NOT NULL,
  degree         text,
  field_of_study text,
  start_date     date,
  end_date       date,
  grade          text,
  description    text,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_educations_user ON educations(user_id);

CREATE TABLE IF NOT EXISTS certifications (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name           text NOT NULL,
  issuer         text,
  issue_date     date,
  expiry_date    date,
  credential_id  text,
  credential_url text,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_certifications_user ON certifications(user_id);

CREATE TABLE IF NOT EXISTS skills (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS profile_skills (
  user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  skill_id uuid NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, skill_id)
);

CREATE TABLE IF NOT EXISTS languages (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS profile_languages (
  user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  language_id uuid NOT NULL REFERENCES languages(id) ON DELETE CASCADE,
  proficiency text NOT NULL DEFAULT 'professional'
                CHECK (proficiency IN ('basic','conversational','professional','fluent','native')),
  PRIMARY KEY (user_id, language_id)
);

CREATE TABLE IF NOT EXISTS portfolio_links (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  label      text,
  url        text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_portfolio_links_user ON portfolio_links(user_id);
