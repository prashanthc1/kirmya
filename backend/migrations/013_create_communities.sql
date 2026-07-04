CREATE TABLE IF NOT EXISTS communities (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug        text NOT NULL UNIQUE,
  name        text NOT NULL,
  description text,
  category    text,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS community_members (
  community_id uuid NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
  user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role         text NOT NULL DEFAULT 'member' CHECK (role IN ('member','moderator')),
  joined_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (community_id, user_id)
);

CREATE TABLE IF NOT EXISTS posts (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  community_id uuid NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
  author_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title        text NOT NULL,
  body         text,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_posts_community ON posts(community_id, created_at DESC);

CREATE TABLE IF NOT EXISTS comments (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id    uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  author_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  body       text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id, created_at);

CREATE TABLE IF NOT EXISTS reactions (
  post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind    text NOT NULL DEFAULT 'like',
  PRIMARY KEY (post_id, user_id, kind)
);

INSERT INTO communities (slug, name, description, category) VALUES
  ('facilities-management','Facilities Management','Operations, maintenance, and facilities professionals','Operations'),
  ('construction','Construction','Construction and built-environment careers','Construction'),
  ('logistics','Logistics','Supply chain, logistics, and transport','Logistics'),
  ('technology','Technology','Software, IT, and technology roles','Technology'),
  ('hr','Human Resources','People, talent, and HR professionals','HR'),
  ('operations','Operations','Operations and general management','Operations')
ON CONFLICT (slug) DO NOTHING;
