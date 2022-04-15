CREATE TABLE IF NOT EXISTS topic
(
    id          serial      NOT NULL PRIMARY KEY,
    name        varchar(64) NOT NULL,
    language_id int         NOT NULL REFERENCES language (id)
);
