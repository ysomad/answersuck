BEGIN
;

CREATE TABLE IF NOT EXISTS player (
    nickname varchar(25) PRIMARY KEY NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    display_name varchar(25) NOT NULL,
    email_verified boolean DEFAULT FALSE NOT NULL,
    password text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
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
    updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS stage (
    id serial NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    "order" smallint NOT NULL,
    package_id int NOT NULL REFERENCES package (id)
);

CREATE TABLE IF NOT EXISTS topic (
    id serial NOT NULL PRIMARY KEY,
    name varchar(50) NOT NULL,
    author varchar(25) NOT NULL REFERENCES player (nickname),
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

CREATE TABLE IF NOT EXISTS package_stage_question (
    stage_id int NOT NULL REFERENCES stage (id),
    topic_id int NOT NULL REFERENCES topic (id),
    question_id int NOT NULL REFERENCES question (id),

    question_type smallint NOT NULL,
    cost smallint NOT NULL,
    interval smallint NOT NULL,
    host_comment text,

    secret_topic varchar(64),
    secret_cost smallint,
    is_keepable boolean,
    is_visible boolean
);

CREATE TABLE IF NOT EXISTS tag (
    id serial NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS package_tag (
    package_id int NOT NULL REFERENCES package (id),
    tag_id int NOT NULL REFERENCES tag (id),
    PRIMARY KEY (package_id, tag_id)
);

COMMIT;