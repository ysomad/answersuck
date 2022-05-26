CREATE TYPE mime_type AS enum ('image/jpeg', 'image/png', 'audio/mp4', 'audio/aac', 'audio/mpeg');

CREATE TABLE IF NOT EXISTS media
(
    id           uuid          NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    url          varchar(2048) NOT NULL,
    mime_type    mime_type     NOT NULL,
    account_id   uuid          NOT NULL       REFERENCES account (id),
    created_at   timestamp     WITH TIME ZONE DEFAULT current_timestamp NOT NULL
);

