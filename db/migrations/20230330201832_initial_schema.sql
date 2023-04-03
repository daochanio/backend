-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
	address VARCHAR(42) PRIMARY KEY,
	ens_name VARCHAR(255) NULL DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE threads (
	id BIGSERIAL PRIMARY KEY,
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	content TEXT NOT NULL,
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
