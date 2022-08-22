CREATE TABLE IF NOT EXISTS package
(
    id               serial                                             NOT NULL PRIMARY KEY,
    name             varchar(48)                                        NOT NULL,
    description      text,
    account_id       uuid                                               NOT NULL REFERENCES account (id),
    is_published     bool                     DEFAULT FALSE             NOT NULL,
    language_id      smallint                                           NOT NULL REFERENCES language (id),
    cover_filename   text,
    created_at       timestamptz DEFAULT current_timestamp NOT NULL,
    updated_at       timestamptz DEFAULT current_timestamp NOT NULL
);
