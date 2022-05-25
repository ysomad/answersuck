CREATE TABLE IF NOT EXISTS question
(
    id          serial                                             NOT NULL PRIMARY KEY,
    question    varchar(200)                                       NOT NULL,
    answer_id   int                                                NOT NULL REFERENCES answer (id),
    account_id  uuid                                               NOT NULL REFERENCES account (id),
    language_id int                                                NOT NULL REFERENCES language (id),
    media_id    uuid REFERENCES media (id),
    created_at  timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL,
    updated_at  timestamp WITH TIME ZONE DEFAULT current_timestamp NOT NULL
);