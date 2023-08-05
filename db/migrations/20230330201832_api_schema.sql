-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
	address VARCHAR(42) PRIMARY KEY,
	ens_name VARCHAR(255) NULL DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NULL DEFAULT NULL,
	reputation BIGINT NOT NULL DEFAULT 0,
	ens_avatar_file_name VARCHAR(255) NULL DEFAULT NULL,
	ens_avatar_url VARCHAR(255) NULL DEFAULT NULL
);

CREATE TABLE challenges (
	address VARCHAR(42) PRIMARY KEY,
	message VARCHAR(255) NOT NULL,
	expires_at BIGINT NOT NULL
);

CREATE TABLE threads (
	id BIGSERIAL PRIMARY KEY,
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	image_file_name TEXT NOT NULL,
	image_original_url TEXT NOT NULL,
	image_thumbnail_url TEXT NOT NULL,
	votes BIGINT NOT NULL DEFAULT 0,
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE thread_votes (
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	thread_id BIGINT NOT NULL REFERENCES threads(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	vote SMALLINT NOT NULL,
	PRIMARY KEY (address, thread_id)
);

CREATE TABLE comments (
	id BIGSERIAL PRIMARY KEY,
	thread_id BIGINT NOT NULL REFERENCES threads(id),
	replied_to_comment_id BIGINT NULL REFERENCES comments(id),
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	content TEXT NOT NULL,
	image_file_name TEXT NOT NULL,
	image_original_url TEXT NOT NULL,
	image_thumbnail_url TEXT NOT NULL,
	votes BIGINT NOT NULL DEFAULT 0,
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE comment_votes (
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	comment_id BIGINT NOT NULL REFERENCES comments(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	vote SMALLINT NOT NULL,
	PRIMARY KEY (address, comment_id)
);

CREATE INDEX threads_created_at_idx ON threads(created_at);

CREATE INDEX thread_votes_thread_id_idx ON thread_votes(thread_id);

CREATE INDEX comments_thread_id_idx ON comments(thread_id);

CREATE INDEX comments_created_at_idx ON comments(created_at);

CREATE INDEX comment_votes_comment_id_idx ON comment_votes(comment_id);

-- +goose StatementEnd
