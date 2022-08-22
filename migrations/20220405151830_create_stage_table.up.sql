CREATE TABLE IF NOT EXISTS stage
(
    id         serial                NOT NULL PRIMARY KEY,
    name       varchar(32)           NOT NULL,
    is_final   boolean DEFAULT FALSE NOT NULL,
    "order"    smallint                   NOT NULL,
    package_id int                   NOT NULL REFERENCES package (id)
);
