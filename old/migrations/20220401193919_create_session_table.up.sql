CREATE TABLE IF NOT EXISTS session (
    id char(64) NOT NULL PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES account (id),
    max_age int NOT NULL,
    user_agent text NOT NULL,
    ip varchar(15) NOT NULL,
    expires_at bigint NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

