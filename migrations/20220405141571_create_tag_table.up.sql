CREATE TABLE IF NOT EXISTS tag
(
    id          serial NOT NULL PRIMARY KEY,
    name        varchar(32),
    language_id int    NOT NULL REFERENCES language (id)
);
