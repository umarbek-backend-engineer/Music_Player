create table lyrics (
    id uuid primary key default gen_random_uuid(),
    music_id uuid not null,
    title text not null,
    content JSONB not null,
    is_verified boolean default false,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
