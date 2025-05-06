-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.user (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE,
    hash_password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_email ON users.user(email);

CREATE TABLE IF NOT EXISTS users.address (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    city TEXT NOT NULL,
    street TEXT NOT NULL,
    building TEXT NOT NULL,
    apartment TEXT,
    notes TEXT,

    FOREIGN KEY (user_id) REFERENCES users.user(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users.address;
DROP TABLE IF EXISTS users.user;
DROP SCHEMA IF EXISTS users;
-- +goose StatementEnd
