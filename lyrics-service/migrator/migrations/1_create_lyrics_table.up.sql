create table lyrics (
    id uuid primary key default gen_random_uuid(),
    music_id uuid not null,
    name text not null,
    content JSONB not null
);