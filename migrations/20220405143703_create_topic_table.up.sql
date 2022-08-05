CREATE TABLE IF NOT EXISTS topic
(
    id          serial                                             NOT NULL PRIMARY KEY,
    name        varchar(50)                                        NOT NULL,
    language_id smallint                                           NOT NULL REFERENCES language (id),
    created_at  timestamptz DEFAULT current_timestamp NOT NULL,
    updated_at  timestamptz DEFAULT current_timestamp NOT NULL
);
