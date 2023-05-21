-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
ADD reputation BIGINT NOT NULL DEFAULT 0,
ADD ens_avatar_file_name VARCHAR(255) NULL DEFAULT NULL,
ADD ens_avatar_url VARCHAR(255) NULL DEFAULT NULL,
ADD ens_avatar_content_type VARCHAR(255) NULL DEFAULT NULL;

-- +goose StatementEnd