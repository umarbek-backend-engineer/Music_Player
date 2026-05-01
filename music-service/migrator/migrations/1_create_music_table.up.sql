create table music (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    title text not null,
    filepath text unique not null,
    is_public boolean default false,
    uploaded_at timestamp default now()
);

create index idx_music_user_id on music(user_id);