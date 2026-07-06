-- Admin suspend/activate: suspended users keep their data + tickets but
-- cannot purchase. Soft flag instead of hard delete.
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
