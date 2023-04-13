-- create/update user every time the sign-in
-- name: CreateOrUpdateUser :one
INSERT INTO users (address, ens_name)
VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE SET ens_name = $2, updated_at = NOW()
RETURNING *;

-- name: CreateThread :one
INSERT INTO threads (address, content)
VALUES ($1, $2)
RETURNING id;

-- name: CreateComment :one
INSERT INTO comments (address, thread_id, replied_to_comment_id, content)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- Order by random is not performant as we need to do a full table scan.
-- Move to TABLESAMPLE SYSTEM_ROWS(N) when performance becomes an issue.
-- Table sample is not random enough until the table gets big.
-- https://www.postgresql.org/docs/current/tsm-system-rows.html
-- name: GetThreads :many
SELECT
  threads.*,
  SUM(COALESCE(thread_votes.vote, 0)) as votes
FROM threads
LEFT JOIN thread_votes ON thread_votes.thread_id = threads.id
WHERE threads.is_deleted = FALSE
GROUP BY threads.id
ORDER BY RANDOM()
LIMIT $1;

-- name: GetThread :one
SELECT
	threads.*,
	SUM(COALESCE(thread_votes.vote, 0)) as votes
FROM threads
LEFT JOIN thread_votes ON thread_votes.thread_id = threads.id
WHERE threads.id = $1
AND threads.is_deleted = FALSE
GROUP BY threads.id;

-- name: GetComments :many
SELECT
	c.*,
	SUM(COALESCE(cv.vote, 0)) as votes,
	r.id as r_id,
	r.address as r_address,
	r.content as r_content,
	r.is_deleted as r_is_deleted,
	r.created_at as r_created_at,
	r.deleted_at as r_deleted_at
FROM comments c
LEFT JOIN comment_votes cv on c.id = cv.comment_id
LEFT JOIN comments r on c.replied_to_comment_id = r.id
WHERE c.thread_id = $1
GROUP BY c.id, r.id
ORDER BY c.created_at DESC
OFFSET $2
LIMIT $3;

-- name: GetComment :one
SELECT
	c.*,
	SUM(COALESCE(cv.vote, 0)) as votes
FROM comments c
LEFT JOIN comment_votes cv on c.id = cv.comment_id
WHERE c.id = $1;

-- name: DeleteThread :one
UPDATE threads
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id;

-- name: DeleteComment :one
UPDATE comments
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id;

-- name: CreateThreadUpVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 1, updated_at = NOW();

-- name: CreateThreadDownVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = -1, updated_at = NOW();

-- name: CreateThreadUnVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 0, updated_at = NOW();

-- name: CreateCommentUpVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 1, updated_at = NOW();

-- name: CreateCommentDownVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = -1, updated_at = NOW();

-- name: CreateCommentUnVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 0, updated_at = NOW();
