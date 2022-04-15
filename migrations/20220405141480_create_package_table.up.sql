CREATE TABLE IF NOT EXISTS package_cover
(
    id  serial        NOT NULL PRIMARY KEY,
    url varchar(2048) NOT NULL
);

CREATE TABLE IF NOT EXISTS package
(
    id               serial                                             NOT NULL PRIMARY KEY,
    name             varchar(48)                                        NOT NULL,
    account_id       uuid                                               NOT NULL REFERENCES account (id),
    is_published     bool                     DEFAULT FALSE             NOT NULL,
    language_id      int                                                NOT NULL REFERENCES language (id),
    package_cover_id int REFERENCES package_cover (id),
    created_at       timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL,
    updated_at       timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL
);

