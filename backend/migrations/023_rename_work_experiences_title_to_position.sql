-- Migration 023: Rename work_experiences title to position and add cover_url to profiles
-- Enforces alignment between the database schema and the v2 Profile Aggregate domain models.

ALTER TABLE work_experiences RENAME COLUMN title TO position;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS cover_url text;
