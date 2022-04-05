
create table if not exists package_cover(
    id bigserial not null primary key,
    url varchar(2048) not null
);

create table if not exists package(
    id bigserial not null primary key,
    name varchar(48) not null,
    account_id uuid not null references account (id),
    cover_id bigint not null references package_cover (id),
    is_published bool not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null
);

create table if not exists tag(
    id bigserial not null primary key,
    name varchar(32),
    created_at timestamp with time zone default current_timestamp not null
);

create table if not exists package_tag(
    package_id bigint not null references package (id),
    tag_id bigint not null references tag (id)
);