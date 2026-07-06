-- adminauth: a completely separate admin identity store. Lives in its own
-- Postgres schema with ZERO foreign keys to the site's public schema — admin
-- accounts are NOT site users and must never share tables or models with them.
CREATE SCHEMA IF NOT EXISTS adminauth;

CREATE TABLE IF NOT EXISTS adminauth.admin_accounts (
    id            UUID PRIMARY KEY,
    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL,
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS adminauth.admin_refresh_tokens (
    id         UUID PRIMARY KEY,
    admin_id   UUID        NOT NULL REFERENCES adminauth.admin_accounts(id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    rotated_at TIMESTAMPTZ NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_admin_refresh_admin ON adminauth.admin_refresh_tokens (admin_id);
