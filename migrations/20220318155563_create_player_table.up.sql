CREATE TABLE IF NOT EXISTS player (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    account_id uuid NOT NULL UNIQUE REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS player_avatar (
    id serial NOT NULL PRIMARY KEY,
    filename text NOT NULL,
    player_id uuid UNIQUE NOT NULL REFERENCES player (id)
);
