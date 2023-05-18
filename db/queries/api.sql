-- upsert user every time they signin so we don't have to check if they exist
-- name: UpsertUser :exec
INSERT INTO users (address)
VALUES ($1)
ON CONFLICT (address) DO NOTHING;

-- name: UpdateUser :exec
UPDATE users
SET
	ens_name = $2,
	updated_at = NOW()
WHERE address = $1;

-- name: GetChallenge :one
SELECT *
FROM challenges
WHERE address = $1;

-- name: UpdateChallenge :exec
INSERT INTO challenges (address, message, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (address) DO UPDATE
SET message = $2, expires_at = $3;

-- name: CreateThread :one
INSERT INTO threads (address, title, content, image_file_name, image_url, image_content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateComment :one
INSERT INTO comments (address, thread_id, replied_to_comment_id, content, image_file_name, image_url, image_content_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- Order by random is not performant as we need to do a full table scan.
-- Move to TABLESAMPLE SYSTEM_ROWS(N) when performance becomes an issue.
-- Table sample is not random enough until the table gets big.
-- https://www.postgresql.org/docs/current/tsm-system-rows.html
-- name: GetThreads :many
SELECT threads.*
FROM threads
WHERE threads.is_deleted = FALSE
ORDER BY RANDOM()
LIMIT $1::bigint;

-- name: GetThread :one
SELECT threads.*
FROM threads
WHERE threads.id = $1
AND threads.is_deleted = FALSE;

-- name: GetComments :many
SELECT
	c.*,
	r.id as r_id,
	r.address as r_address,
	r.content as r_content,
	r.image_file_name as r_image_file_name,
	r.image_url as r_image_url,
	r.image_content_type as r_image_content_type,
	r.is_deleted as r_is_deleted,
	r.created_at as r_created_at,
	r.deleted_at as r_deleted_at,
	count(*) OVER() AS full_count
FROM comments c
LEFT JOIN comments r on c.replied_to_comment_id = r.id
WHERE c.thread_id = $1
ORDER BY c.created_at DESC
OFFSET $2::bigint
LIMIT $3::bigint;

-- name: GetComment :one
SELECT c.*
FROM comments c
WHERE c.id = $1;

-- name: DeleteThread :one
UPDATE threads
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id as thread_id;

-- name: DeleteComment :one
UPDATE comments
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id as comment_id;

-- name: CreateThreadUpVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 1, updated_at = NOW()
WHERE thread_votes.updated_at < TO_TIMESTAMP($3); -- only update the vote if the incoming vote is newer than the current vote

-- name: CreateThreadDownVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = -1, updated_at = NOW()
WHERE thread_votes.updated_at < TO_TIMESTAMP($3);

-- name: CreateThreadUnVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 0, updated_at = NOW()
WHERE thread_votes.updated_at < TO_TIMESTAMP($3);

-- name: CreateCommentUpVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 1, updated_at = NOW()
WHERE comment_votes.updated_at < TO_TIMESTAMP($3);

-- name: CreateCommentDownVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = -1, updated_at = NOW()
WHERE comment_votes.updated_at < TO_TIMESTAMP($3);

-- name: CreateCommentUnVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 0, updated_at = NOW()
WHERE comment_votes.updated_at < TO_TIMESTAMP($3);

-- name: AggregateThreadVotes :exec
UPDATE threads
SET votes = (
	SELECT COALESCE(SUM(vote), 0)
	FROM thread_votes
	WHERE thread_votes.thread_id = $1
)
WHERE threads.id = $1;

-- name: AggregateCommentVotes :exec
UPDATE comments
SET votes = (
	SELECT COALESCE(SUM(vote), 0)
	FROM comment_votes
	WHERE comment_votes.comment_id = $1
)
WHERE comments.id = $1;