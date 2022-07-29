CREATE TYPE question_type AS enum (
    'DEFAULT',
    'BET',
    'SECRET',
    'SUPERSECRET',
    'SAFE'
);

CREATE TABLE IF NOT EXISTS question (
    id serial NOT NULL PRIMARY KEY,
    text varchar(200) NOT NULL,
    answer_id int NOT NULL REFERENCES answer (id),
    account_id uuid NOT NULL REFERENCES account (id),
    language_id int NOT NULL REFERENCES
    LANGUAGE (id),
    media_id uuid REFERENCES media (id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
);

