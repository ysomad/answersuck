CREATE TYPE media_type AS enum ('AUDIO', 'IMG', 'VIDEO');

CREATE TABLE IF NOT EXISTS question_media
(
    id   serial        NOT NULL PRIMARY KEY,
    url  varchar(2048) NOT NULL,
    type media_type    NOT NULL
);

CREATE TABLE IF NOT EXISTS question
(
    id          serial       NOT NULL PRIMARY KEY,
    question    varchar(255) NOT NULL,
    answer_id   int          NOT NULL REFERENCES answer (id),
    account_id  uuid         NOT NULL REFERENCES account (id),
    language_id int          NOT NULL REFERENCES language (id),
    media_id    int REFERENCES question_media (id)
);