CREATE TABLE IF NOT EXISTS player (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    account_id uuid NOT NULL UNIQUE REFERENCES account (id),
    avatar_filename text,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);
