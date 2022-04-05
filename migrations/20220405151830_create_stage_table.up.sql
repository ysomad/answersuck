create table if not exists stage(
    id bigserial not null primary key,
    title varchar(24),
    is_final bool not null,
    "order" int not null,
    package_id bigint not null references package (id)
);

create table if not exists stage_topic(
    stage_id bigint not null references stage (id),
    topic_id bigint not null references topic (id)
);
