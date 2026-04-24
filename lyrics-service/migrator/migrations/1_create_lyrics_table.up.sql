create table lyrics (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    music_id uuid not null unique,
    name text not null,
    content JSONB not null
);