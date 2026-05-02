create table music (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    title text not null,
    filepath text not null,
    is_public boolean default false,
    uploaded_at timestamptz default now()
);
