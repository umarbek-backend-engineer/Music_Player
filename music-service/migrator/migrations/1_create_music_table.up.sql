create table music (
    id uuid primary key default gen_random_uuid(),
    filename text,
    filepath text not null,
    uploaded_at timestamp default now()
);