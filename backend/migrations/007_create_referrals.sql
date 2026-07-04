-- Referral marketplace: a job seeker requests a referral (optionally directed
-- at a specific employee/referrer and/or tied to a job or company). The
-- referrer reviews and accepts/declines; after acceptance the outcome is
-- tracked through the hiring pipeline.

CREATE TABLE IF NOT EXISTS referral_requests (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  seeker_id   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  referrer_id uuid REFERENCES users(id) ON DELETE SET NULL,   -- null = open request, claimed on accept
  job_id      uuid REFERENCES jobs(id) ON DELETE SET NULL,
  company     text,
  message     text,
  status      text NOT NULL DEFAULT 'requested'
                CHECK (status IN ('requested','under_review','accepted','declined')),
  outcome     text
                CHECK (outcome IN ('application_submitted','interviewing','offer','hired','rejected','withdrawn')),
  decided_at  timestamptz,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz,
  version     integer NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_referrals_seeker   ON referral_requests(seeker_id)   WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_referrals_referrer ON referral_requests(referrer_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_referrals_status   ON referral_requests(status);
