CREATE TABLE IF NOT EXISTS stage
(
    id         serial             NOT NULL PRIMARY KEY,
    name       varchar(32)        NOT NULL,
    is_final   bool DEFAULT FALSE NOT NULL,
    "order"    int                NOT NULL,
    package_id int                NOT NULL REFERENCES package (id)
);