CREATE TABLE IF NOT EXISTS account
(
    id          uuid                                               NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email       varchar(255) UNIQUE                                NOT NULL,
    username    varchar(16) UNIQUE                                 NOT NULL,
    password    varchar(255)                                       NOT NULL,
    is_verified boolean                  DEFAULT FALSE             NOT NULL,
    is_archived boolean                  DEFAULT FALSE             NOT NULL,
    created_at  timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL,
    updated_at  timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS account_verification_code
(
    id         serial          NOT NULL PRIMARY KEY,
    code       char(64) UNIQUE NOT NULL,
    account_id uuid UNIQUE     NOT NULL REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS account_avatar
(
    id         serial        NOT NULL PRIMARY KEY,
    url        varchar(2048) NOT NULL,
    account_id uuid          NOT NULL REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS account_password_reset_token
(
    id         serial                                             NOT NULL PRIMARY KEY,
    token      char(64) UNIQUE                                    NOT NULL,
    account_id uuid                                               NOT NULL REFERENCES account (id),
    created_at timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL
);
