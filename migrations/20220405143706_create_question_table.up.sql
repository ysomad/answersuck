create table if not exists answer(
    id bigserial not null primary key,
    answer text not null,
    image_url varchar(2048)
);

create type question_type as enum ('DEFAULT', 'BET', 'SECRET', 'SUPERSECRET', 'NORISK');

create table if not exists question(
    id bigserial not null primary key,
    question text not null,
    answer_id bigint not null references answer (id),
    account_id uuid not null references account (id),
    type question_type not null,
    cost int not null,
    time_interval int not null,
    host_comment text,

    secret_topic varchar(64),
    secret_cost int,
    is_keepable bool,
    is_known bool
);

create type media_type as enum ('AUDIO', 'IMG', 'VIDEO');

create table if not exists question_media(
    id bigserial not null primary key,
    url varchar(2048) not null,
    question_id bigint not null references question (id),
    media_type media_type not null
);

create table if not exists topic(
    id bigserial not null primary key,
    title varchar(64) not null
);

create table if not exists topic_question(
    topic_id bigint not null references topic (id),
    question_id bigint not null references question (id)
);
