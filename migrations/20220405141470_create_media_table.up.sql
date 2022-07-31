CREATE TYPE media_type AS enum (
    'image/jpeg',
    'image/png',
    'audio/mp4',
    'audio/aac',
    'audio/mpeg'
);

CREATE TABLE IF NOT EXISTS media (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    url text NOT NULL,
    type media_type NOT NULL,
    account_id uuid NOT NULL REFERENCES account (id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

