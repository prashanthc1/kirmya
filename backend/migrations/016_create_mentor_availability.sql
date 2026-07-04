-- Mentorship bounded context: mentor availability windows. A mentor opens slots
-- that mentees can book against; booking a session marks the slot is_booked=true.

CREATE TABLE IF NOT EXISTS mentor_availability (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  mentor_id  uuid NOT NULL REFERENCES mentor_profiles(id) ON DELETE CASCADE,
  starts_at  timestamptz NOT NULL,
  ends_at    timestamptz NOT NULL,
  is_booked  boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (ends_at > starts_at)
);

CREATE INDEX IF NOT EXISTS idx_mentor_availability_open
  ON mentor_availability(mentor_id, starts_at) WHERE is_booked = false;
