create table if not exists session(
    id char(64) not null primary key,
    account_id uuid not null references account (id),
    max_age int not null,
    user_agent text not null,
    ip inet not null,
    expires_at bigint not null,
    created_at timestamp with time zone default current_timestamp not null
);