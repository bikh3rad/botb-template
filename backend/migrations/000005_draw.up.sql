-- postgresql draws migration
-- A draw selects a winning ticket for a competition. Winner fields are NULL
-- until the draw is run. FKs reference the competition, user, and ticket tables
-- in the shared database.
CREATE TABLE draws (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    competition_id UUID NOT NULL REFERENCES competitions (id) ON DELETE CASCADE,
    winner_user_id UUID NULL REFERENCES users (id) ON DELETE SET NULL,
    winner_ticket_id UUID NULL REFERENCES tickets (id) ON DELETE SET NULL,
    prize TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    drawn_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_draws_competition ON draws (competition_id);
CREATE INDEX idx_draws_status ON draws (status);
