CREATE TABLE IF NOT EXISTS profiles (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  bio        text,
  location   text,
  website    text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
