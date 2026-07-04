CREATE TABLE IF NOT EXISTS mentor_profiles (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  headline   text NOT NULL,
  bio        text,
  expertise  text,
  is_active  boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_mentor_profiles_active ON mentor_profiles(is_active);

CREATE TABLE IF NOT EXISTS mentorship_sessions (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  mentor_id    uuid NOT NULL REFERENCES mentor_profiles(id) ON DELETE CASCADE,
  mentee_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  topic        text,
  scheduled_at timestamptz NOT NULL,
  status       text NOT NULL DEFAULT 'requested'
                 CHECK (status IN ('requested','confirmed','completed','cancelled')),
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_mentorship_sessions_mentor ON mentorship_sessions(mentor_id);
CREATE INDEX IF NOT EXISTS idx_mentorship_sessions_mentee ON mentorship_sessions(mentee_id);

CREATE TABLE IF NOT EXISTS mentorship_reviews (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id uuid NOT NULL UNIQUE REFERENCES mentorship_sessions(id) ON DELETE CASCADE,
  rating     integer NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment    text,
  created_at timestamptz NOT NULL DEFAULT now()
);
