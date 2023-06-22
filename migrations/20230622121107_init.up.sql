BEGIN
;

CREATE TABLE IF NOT EXISTS player (
    id uuid PRIMARY KEY NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    nickname varchar(25) UNIQUE NOT NULL,
    display_name varchar(25) NOT NULL,
    email_verified boolean DEFAULT FALSE NOT NULL,
    password text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE TYPE media_type AS enum (
    'IMAGE',
    'AUDIO',
    'VIDEO'
);

CREATE TABLE IF NOT EXISTS media (
    url varchar(2048) NOT NULL PRIMARY KEY,
    type media_type NOT NULL,
    player_id uuid NOT NULL REFERENCES player (id),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS language (
    id smallserial NOT NULL PRIMARY KEY,
    name varchar(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS package (
    id serial NOT NULL PRIMARY KEY,
    name varchar(64) NOT NULL,
    author_id uuid NOT NULL REFERENCES player (id),
    is_published bool DEFAULT FALSE NOT NULL,
    language_id smallint NOT NULL REFERENCES language (id),
    cover_url varchar(2048) REFERENCES media (url),
    created_at timestamptz NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS tag (
    id serial NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    language_id smallint NOT NULL REFERENCES language (id),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS package_tag (
    package_id int NOT NULL REFERENCES package (id),
    tag_id int NOT NULL REFERENCES tag (id),
    PRIMARY KEY (package_id, tag_id)
);

CREATE TABLE IF NOT EXISTS answer (
    id serial NOT NULL PRIMARY KEY,
    text varchar(112) NOT NULL,
    media_url varchar(2048) REFERENCES media (url),
    language_id smallint NOT NULL REFERENCES language (id),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS topic (
    id serial NOT NULL PRIMARY KEY,
    name varchar(50) NOT NULL,
    language_id smallint NOT NULL REFERENCES language (id),
    created_at timestamptz NOT NULL
);

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
    author_id uuid NOT NULL REFERENCES player (id),
    language_id smallint NOT NULL REFERENCES language (id),
    media_url varchar(2048) REFERENCES media (url),
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS stage (
    id serial NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    "order" smallint NOT NULL,
    package_id int NOT NULL REFERENCES package (id)
);

COMMIT;