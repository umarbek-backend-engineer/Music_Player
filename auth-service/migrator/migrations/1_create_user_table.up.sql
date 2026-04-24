CREATE EXTENSION IF NOT EXISTS pgcrypto;
create table users (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    lastname text not null,
    email text not null unique,
    role text not null default 'user',
    password text not null,
    password_changed_at timestamptz,
    refreshtoken text,
    user_created_at timestamptz default now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    user_agent text,
    ip_address text,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    revoked BOOLEAN DEFAULT FALSE
);