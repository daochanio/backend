-- +goose Up
-- +goose StatementBegin

-- alter users table to make updated_at column nullable and default null
ALTER TABLE users ALTER COLUMN updated_at DROP NOT NULL;
ALTER TABLE users ALTER COLUMN updated_at DROP DEFAULT;

-- +goose StatementEnd