-- Migration script to add is_read column to existing databases
ALTER TABLE posts ADD COLUMN IF NOT EXISTS is_read boolean NOT NULL DEFAULT false;
