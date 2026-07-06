-- Draw voiding: void_reason records WHY a draw was voided (required by the
-- void endpoint; also mirrored into admin_audit_log).
ALTER TABLE draws ADD COLUMN IF NOT EXISTS void_reason TEXT NOT NULL DEFAULT '';
