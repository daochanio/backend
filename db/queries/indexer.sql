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
