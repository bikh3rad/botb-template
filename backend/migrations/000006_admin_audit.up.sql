-- Audit trail for admin mutations across all services (written via pkg/audit).
CREATE TABLE IF NOT EXISTS admin_audit_log (
    id          UUID PRIMARY KEY,
    actor_id    TEXT        NOT NULL DEFAULT '',
    actor_email TEXT        NOT NULL DEFAULT '',
    action      TEXT        NOT NULL,
    entity_type TEXT        NOT NULL,
    entity_id   TEXT        NOT NULL DEFAULT '',
    reason      TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_entity ON admin_audit_log (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_created ON admin_audit_log (created_at);
