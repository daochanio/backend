-- +goose Up
-- +goose StatementBegin

CREATE TABLE indexer_progress (
	version VARCHAR NOT NULL PRIMARY KEY,
	last_indexed_block NUMERIC NOT NULL,
	indexed_on TIMESTAMP NOT NULL DEFAULT NOW()
);

-- seed the starting block for the indexer
INSERT INTO indexer_progress (version, last_indexed_block)
VALUES ('1.0', '0');

CREATE TABLE transfers (
	block_number NUMERIC NOT NULL,
	transaction_id VARCHAR(66) NOT NULL,
	log_index BIGINT NOT NULL,
	from_address VARCHAR(42) NOT NULL REFERENCES users(address),
	to_address VARCHAR(42) NOT NULL REFERENCES users(address),
	amount NUMERIC NOT NULL,

	PRIMARY KEY (block_number, transaction_id, log_index)
);

-- +goose StatementEnd
