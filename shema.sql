-- users
create table if not exists users (
    id bigserial primary key,
    username text not null unique,
    password_hash text not null,
    created_at timestamptz not null default now()
);
-- messages
create table if not exists messages (
    id bigserial primary key,
    user_id bigint not null references users(id) on delete cascade,
    body text not null,
    created_at timestamptz not null default now()
);
create index if not exists idx_messages_created_at on messages(created_at desc);
create index if not exists idx_messages_user_id_created_at on messages(user_id, created_at desc);