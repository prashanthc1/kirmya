-- Full name is edited in the profile workspace (Identity section) but had no
-- column; it was previously derived from the headline and never persisted.
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS full_name text;
