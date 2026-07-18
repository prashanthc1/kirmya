-- Migration 028: Add origin column to connections table
ALTER TABLE connections ADD COLUMN IF NOT EXISTS origin VARCHAR(50) DEFAULT 'manual_request';
