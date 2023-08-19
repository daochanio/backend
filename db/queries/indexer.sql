-- name: GetLastIndexedBlock :one
SELECT last_indexed_block
FROM indexer_progress
WHERE version = '1.0'
LIMIT 1;

-- name: UpdateLastIndexedBlock :exec
UPDATE indexer_progress
SET 
  last_indexed_block = $1,
  indexed_on = NOW()
WHERE version = '1.0';

-- name: DeleteTransfers :exec
DELETE FROM transfers
WHERE block_number >= $1
AND block_number <= $2;

-- name: InsertTransfers :copyfrom
INSERT INTO transfers (
  block_number,
  transaction_id,
  log_index,
  from_address,
  to_address,
  amount
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);

-- name: ZeroReputation :exec
-- set the reputation of all users to 0
UPDATE users
SET reputation = 0
WHERE address = ANY($1::varchar(42)[]);

-- name: AddReputation :exec
-- if the user address is the to_address, then add the amount
UPDATE users
SET reputation = reputation + sub.amount
FROM (
  SELECT t.to_address as address, SUM(t.amount) as amount
  FROM transfers t
  WHERE t.to_address = ANY($1::varchar(42)[])
  GROUP BY to_address
) as sub
WHERE users.address = sub.address;

-- name: DeductReputation :exec
-- if the user address is the from_address, then deduct the amount
UPDATE users
SET reputation = reputation - sub.amount
FROM (
  SELECT t.from_address as address, SUM(t.amount) as amount
  FROM transfers t
  WHERE t.from_address = ANY($1::varchar(42)[])
  GROUP BY from_address
) as sub
WHERE users.address = sub.address;