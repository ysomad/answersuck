CREATE TABLE IF NOT EXISTS player
(
    id       uuid        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id uuid NOT NULL UNIQUE REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS player_avatar
(
    id         serial        NOT NULL PRIMARY KEY,
    url        varchar(2048) NOT NULL,
    player_id  uuid          NOT NULL REFERENCES player (id),
    created_at timestamptz   DEFAULT current_timestamp NOT NULL
);