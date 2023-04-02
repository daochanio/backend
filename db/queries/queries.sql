-- create/update user every time the sign-in
-- name: CreateOrUpdateUser :exec
INSERT INTO users (address, ens_name)
VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE SET ens_name = $2, updated_at = NOW();

-- name: CreateThread :one
INSERT INTO threads (address, content)
VALUES ($1, $2)
RETURNING id;

-- name: CreateComment :one
INSERT INTO comments (address, thread_id, content)
VALUES ($1, $2, $3)
RETURNING id;

-- name: CreateSelfClosure :exec
INSERT INTO comment_closures (parent_id, child_id, depth)
VALUES ($1, $1, 0);

-- name: CreateParentClosures :exec
-- only do if not root comment (i.e parent_id != child_id)
INSERT into comment_closures (parent_id, child_id, depth)
SELECT p.parent_id, c.child_id, p.depth + c.depth+1
FROM comment_closures p, comment_closures c
WHERE p.child_id = $1 and c.parent_id = $2;

-- name: GetThreads :many
-- TODO: We can order by random in the future
SELECT
  threads.*,
  SUM(COALESCE(thread_votes.vote, 0)) as votes
FROM threads
LEFT JOIN thread_votes ON thread_votes.thread_id = threads.id
WHERE threads.is_deleted = FALSE
GROUP BY threads.id
ORDER BY threads.created_at ASC
OFFSET $1
LIMIT $2;

-- name: GetThread :one
SELECT
	threads.*,
	SUM(COALESCE(thread_votes.vote, 0)) as votes
FROM threads
LEFT JOIN thread_votes ON thread_votes.thread_id = threads.id
WHERE threads.id = $1
AND threads.is_deleted = FALSE
GROUP BY threads.id;

-- name: GetRootAndFirstDepthComments :many
-- select root and first depth comments
-- left join incase there are comments with no votes
-- coalesce as well for no vote comments
-- TODO: This kind of works but we need to paginate this query 
-- But I think we need to paginate the root comments without the children, as including the children will throw of the pagination count
SELECT comments.*, SUM(COALESCE(comment_votes.vote, 0)) as votes
FROM comments
LEFT JOIN comment_votes on comments.id = comment_votes.comment_id
WHERE comments.id NOT IN (
	SELECT child_id
	FROM comment_closures
	where depth > 1
)
AND comments.thread_id = $1
GROUP BY comments.id
ORDER BY comments.created_at ASC;

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
