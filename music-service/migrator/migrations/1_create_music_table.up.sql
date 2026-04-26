create table music (
    id uuid primary key default gen_random_uuid(),
    -- user_id uuid not null,
    filename text,
    filepath text unique not null,
    uploaded_at timestamp default now()
);