CREATE TABLE IF NOT EXISTS account (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid () NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    nickname varchar(25) UNIQUE NOT NULL,
    password text NOT NULL,
    is_verified boolean DEFAULT FALSE NOT NULL,
    is_archived boolean DEFAULT FALSE NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS verification (
    id serial NOT NULL PRIMARY KEY,
    code char(64) UNIQUE NOT NULL,
    account_id uuid UNIQUE NOT NULL REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS password_token (
    id serial NOT NULL PRIMARY KEY,
    token char(64) UNIQUE NOT NULL,
    account_id uuid NOT NULL REFERENCES account (id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

