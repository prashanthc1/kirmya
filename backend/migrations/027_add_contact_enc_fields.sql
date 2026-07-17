-- Migration 027: Add encrypted contact fields to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS email_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS phone_enc text;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS address_enc text;
