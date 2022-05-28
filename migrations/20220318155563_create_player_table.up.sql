CREATE TABLE IF NOT EXISTS player
(
    id       uuid        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname varchar(25) NOT NULL UNIQUE REFERENCES account (nickname)
);

CREATE TABLE IF NOT EXISTS player_avatar
(
    id         serial        NOT NULL PRIMARY KEY,
    url        varchar(2048) NOT NULL,
    player_id  uuid          NOT NULL REFERENCES player (id),
    created_at timestamptz   DEFAULT current_timestamp NOT NULL
);