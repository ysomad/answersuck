CREATE TABLE IF NOT EXISTS tag
(
    id          serial      NOT NULL PRIMARY KEY,
    name        varchar(32) NOT NULL,
    language_id smallint    NOT NULL REFERENCES language (id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);
