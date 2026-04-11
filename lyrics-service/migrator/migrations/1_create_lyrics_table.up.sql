create table lyrics (
    id uuid primary key default gen_random_uuid(),
    music_id uuid not null,
    music_content text not null
);