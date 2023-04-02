-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
	address VARCHAR(42) PRIMARY KEY,
	ens_name VARCHAR(255) NULL DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE threads (
	id SERIAL PRIMARY KEY,
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	content TEXT NOT NULL,
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE thread_votes (
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	thread_id INTEGER NOT NULL REFERENCES threads(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	vote INTEGER NOT NULL,
	PRIMARY KEY (address, thread_id)
);

CREATE TABLE comments (
	id SERIAL PRIMARY KEY,
	thread_id INTEGER NOT NULL REFERENCES threads(id),
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	content TEXT NOT NULL,
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE comment_closures (
	parent_id INTEGER NOT NULL REFERENCES comments(id),
	child_id INTEGER NOT NULL REFERENCES comments(id),
	depth INTEGER NOT NULL,
  PRIMARY KEY (parent_id, child_id)
);

CREATE TABLE comment_votes (
	address VARCHAR(42) NOT NULL REFERENCES users(address),
	comment_id INTEGER NOT NULL REFERENCES comments(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	vote INTEGER NOT NULL,
	PRIMARY KEY (address, comment_id)
);

CREATE INDEX threads_created_at_idx ON threads(created_at);

CREATE INDEX comments_thread_id_idx ON comments(thread_id);

CREATE INDEX comment_closures_parent_id_idx ON comment_closures(parent_id);

CREATE INDEX comment_closures_depth_idx ON comment_closures(depth);

CREATE INDEX comment_votes_comment_id_idx ON comment_votes(comment_id);

-- +goose StatementEnd
