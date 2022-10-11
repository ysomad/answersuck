CREATE TABLE IF NOT EXISTS answer
(
    id        serial       NOT NULL PRIMARY KEY,
    text      varchar(100) NOT NULL,
    media_id  uuid         REFERENCES media (id),
    language_id  smallint  NOT NULL REFERENCES language (id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);
