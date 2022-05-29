CREATE TABLE IF NOT EXISTS tag
(
    id          serial      NOT NULL PRIMARY KEY,
    name        varchar(32) UNIQUE NOT NULL,
    language_id int         NOT NULL REFERENCES language (id)
);
