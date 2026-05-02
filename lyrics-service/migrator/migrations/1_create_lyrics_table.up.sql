create table lyrics (
    id uuid primary key default gen_random_uuid(),
    music_id uuid not null references music(id) on delete cascade,
    title text not null,
    content JSONB not null,
    is_verified boolean default false,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create index idx_lyrics_music_id on lyrics(music_id);
create index idx_lyrics_user_id on lyrics(user_id);