-- postgresql competition migration
-- A competition can have many associated media rows (see migration 000002),
-- resolved by owner_type='competition', owner_id=competitions.id.
CREATE TABLE competitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    prize TEXT NOT NULL DEFAULT '',
    ticket_price_pence BIGINT NOT NULL DEFAULT 0,
    tickets_total BIGINT NOT NULL DEFAULT 0,
    tickets_sold BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_competitions_status ON competitions (status);
