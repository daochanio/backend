-- +goose Up
-- +goose StatementBegin

CREATE TABLE indexer_progress (
	version VARCHAR NOT NULL PRIMARY KEY,
	last_indexed_block VARCHAR NOT NULL,
	indexed_on TIMESTAMP NOT NULL DEFAULT NOW()
);

-- the starting block for the indexer
INSERT INTO indexer_progress (version, last_indexed_block)
VALUES ('1.0', '588316');

-- +goose StatementEnd
