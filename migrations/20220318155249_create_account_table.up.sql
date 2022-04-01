create table if not exists account(
    id uuid not null primary key default gen_random_uuid(),
    email varchar(255) unique not null,
    username varchar(16) unique not null,
    password varchar(255) not null,
    is_verified boolean default false not null,
    is_archived boolean default false not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null
);

create table if not exists account_verification_code(
    id bigserial not null primary key,
    code char(64) unique not null,
    account_id uuid unique not null references account (id)
);

create table if not exists account_avatar(
    id bigserial not null primary key,
    url varchar(2048) not null,
    account_id uuid unique not null references account (id)
);

create table if not exists account_password_reset_token(
   id bigserial not null primary key,
   token char(64) unique not null,
   account_id uuid not null references account (id),
   created_at timestamp with time zone default current_timestamp not null
);
