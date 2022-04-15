CREATE TYPE question_type AS enum ('DEFAULT', 'BET', 'SECRET', 'SUPERSECRET', 'SAFE');

CREATE TABLE IF NOT EXISTS question_config
(
    id           serial        NOT NULL PRIMARY KEY,
    type         question_type NOT NULL DEFAULT 'DEFAULT',
    cost         int           NOT NULL,
    interval     int           NOT NULL,
    comment      text,

    secret_topic varchar(64),
    secret_cost  int,
    is_keepable  boolean,
    is_visible   boolean
);
