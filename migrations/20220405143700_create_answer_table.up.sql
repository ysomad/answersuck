CREATE TABLE IF NOT EXISTS answer
(
    id        serial       NOT NULL PRIMARY KEY,
    text      varchar(100) NOT NULL,
    image     uuid         REFERENCES media (id)
);
