BEGIN
;

CREATE TABLE IF NOT EXISTS player (
    nickname varchar(25) PRIMARY KEY NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    display_name varchar(25),
    email_verified boolean DEFAULT FALSE NOT NULL,
    password text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS media (
    url varchar(2048) NOT NULL PRIMARY KEY,
    type smallint NOT NULL,
    uploaded_by varchar(25) NOT NULL REFERENCES player (nickname),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS package (
    id serial NOT NULL PRIMARY KEY,
    name varchar(64) NOT NULL,
    author varchar(25) NOT NULL REFERENCES player (nickname),
    is_published bool DEFAULT FALSE NOT NULL,
    cover_url varchar(2048) REFERENCES media (url),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,

    round_count smallint,
    topic_count smallint,
    question_count smallint,
    video_count smallint,
    audio_count smallint,
    image_count smallint
);

CREATE TABLE IF NOT EXISTS round (
    id serial NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    position smallint NOT NULL,
    package_id int NOT NULL REFERENCES package (id)
);

CREATE TABLE IF NOT EXISTS topic (
    id serial NOT NULL PRIMARY KEY,
    title varchar(50) NOT NULL,
    author varchar(25) NOT NULL REFERENCES player (nickname),
    round_id int NOT NULL REFERENCES round (id),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS answer (
    id serial NOT NULL PRIMARY KEY,
    text varchar(112) NOT NULL,
    author varchar(25) NOT NULL REFERENCES player (nickname),
    media_url varchar(2048) REFERENCES media (url),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS question (
    id serial NOT NULL PRIMARY KEY,
    text varchar(200) NOT NULL,
    answer_id int NOT NULL REFERENCES answer (id),
    author varchar(25) NOT NULL REFERENCES player (nickname),
    media_url varchar(2048) REFERENCES media (url),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS round_topic (
    round_id int NOT NULL REFERENCES round (id),
    topic_id int NOT NULL REFERENCES topic (id),
    PRIMARY KEY (round_id, topic_id)
);

CREATE TABLE IF NOT EXISTS round_question (
    id serial NOT NULL PRIMARY KEY,
    round_id int NOT NULL REFERENCES round (id),
    topic_id int NOT NULL REFERENCES topic (id),
    question_id int NOT NULL REFERENCES question (id),

    question_type smallint NOT NULL,
    cost smallint NOT NULL,
    answer_time smallint NOT NULL,
    host_comment text,

    secret_topic varchar(64),
    secret_cost smallint,
    transfer_type smallint,
    is_keepable boolean
);

CREATE TABLE IF NOT EXISTS tag (
    name varchar(16) NOT NULL PRIMARY KEY,
    author varchar(25) NOT NULL REFERENCES player (nickname),
    created_at timestamptz NOT NULL
);

ALTER TABLE tag ADD COLUMN ts tsvector
GENERATED ALWAYS AS
    (setweight(to_tsvector('russian', coalesce(name, '')), 'A')) STORED;

CREATE INDEX tag_gin_idx ON tag USING GIN (ts);

CREATE TABLE IF NOT EXISTS package_tag (
    package_id int NOT NULL REFERENCES package (id),
    tag varchar(16) NOT NULL REFERENCES tag (name),
    PRIMARY KEY (package_id, tag)
);

COMMIT;
