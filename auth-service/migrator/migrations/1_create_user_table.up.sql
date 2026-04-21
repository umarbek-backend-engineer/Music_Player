create table users (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    lastname text not null,
    email text not null unique,
    password text not null,
    password_changed_at timestamp,
    refreshtoken text,
    user_created_at timestamp default now()
);