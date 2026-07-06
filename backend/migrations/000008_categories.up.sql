-- First-class competition categories. Fixed UUIDs so the seeder and tests can
-- reference them deterministically.
CREATE TABLE IF NOT EXISTS categories (
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL UNIQUE,
    slug       TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO categories (id, name, slug) VALUES
    ('c0000000-0000-4000-8000-000000000001', 'Cars',         'cars'),
    ('c0000000-0000-4000-8000-000000000002', 'Property',     'property'),
    ('c0000000-0000-4000-8000-000000000003', 'Instant Wins', 'instant-wins'),
    ('c0000000-0000-4000-8000-000000000004', 'Lifestyle',    'lifestyle'),
    ('c0000000-0000-4000-8000-000000000005', 'Tech',         'tech'),
    ('c0000000-0000-4000-8000-000000000006', 'Cash',         'cash')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE competitions ADD COLUMN IF NOT EXISTS category_id UUID NULL REFERENCES categories(id);

CREATE INDEX IF NOT EXISTS idx_competitions_category ON competitions (category_id);
