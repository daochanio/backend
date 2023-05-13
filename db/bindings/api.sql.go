// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: api.sql

package bindings

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const aggregateCommentVotes = `-- name: AggregateCommentVotes :exec
UPDATE comments
SET votes = (
	SELECT COALESCE(SUM(vote), 0)
	FROM comment_votes
	WHERE comment_votes.comment_id = $1
)
WHERE comments.id = $1
`

func (q *Queries) AggregateCommentVotes(ctx context.Context, commentID int64) error {
	_, err := q.db.Exec(ctx, aggregateCommentVotes, commentID)
	return err
}

const aggregateThreadVotes = `-- name: AggregateThreadVotes :exec
UPDATE threads
SET votes = (
	SELECT COALESCE(SUM(vote), 0)
	FROM thread_votes
	WHERE thread_votes.thread_id = $1
)
WHERE threads.id = $1
`

func (q *Queries) AggregateThreadVotes(ctx context.Context, threadID int64) error {
	_, err := q.db.Exec(ctx, aggregateThreadVotes, threadID)
	return err
}

const createComment = `-- name: CreateComment :one
INSERT INTO comments (address, thread_id, replied_to_comment_id, content, image_file_name, image_url, image_content_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, thread_id, replied_to_comment_id, address, content, image_file_name, image_url, image_content_type, votes, is_deleted, created_at, deleted_at
`

type CreateCommentParams struct {
	Address            string
	ThreadID           int64
	RepliedToCommentID pgtype.Int8
	Content            string
	ImageFileName      string
	ImageUrl           string
	ImageContentType   string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRow(ctx, createComment,
		arg.Address,
		arg.ThreadID,
		arg.RepliedToCommentID,
		arg.Content,
		arg.ImageFileName,
		arg.ImageUrl,
		arg.ImageContentType,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.ThreadID,
		&i.RepliedToCommentID,
		&i.Address,
		&i.Content,
		&i.ImageFileName,
		&i.ImageUrl,
		&i.ImageContentType,
		&i.Votes,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createCommentDownVote = `-- name: CreateCommentDownVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = -1, updated_at = NOW()
`

type CreateCommentDownVoteParams struct {
	Address   string
	CommentID int64
}

func (q *Queries) CreateCommentDownVote(ctx context.Context, arg CreateCommentDownVoteParams) error {
	_, err := q.db.Exec(ctx, createCommentDownVote, arg.Address, arg.CommentID)
	return err
}

const createCommentUnVote = `-- name: CreateCommentUnVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 0, updated_at = NOW()
`

type CreateCommentUnVoteParams struct {
	Address   string
	CommentID int64
}

func (q *Queries) CreateCommentUnVote(ctx context.Context, arg CreateCommentUnVoteParams) error {
	_, err := q.db.Exec(ctx, createCommentUnVote, arg.Address, arg.CommentID)
	return err
}

const createCommentUpVote = `-- name: CreateCommentUpVote :exec
INSERT INTO comment_votes (address, comment_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, comment_id) DO UPDATE SET vote = 1, updated_at = NOW()
`

type CreateCommentUpVoteParams struct {
	Address   string
	CommentID int64
}

func (q *Queries) CreateCommentUpVote(ctx context.Context, arg CreateCommentUpVoteParams) error {
	_, err := q.db.Exec(ctx, createCommentUpVote, arg.Address, arg.CommentID)
	return err
}

const createThread = `-- name: CreateThread :one
INSERT INTO threads (address, title, content, image_file_name, image_url, image_content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, address, title, content, image_file_name, image_url, image_content_type, votes, is_deleted, created_at, deleted_at
`

type CreateThreadParams struct {
	Address          string
	Title            string
	Content          string
	ImageFileName    string
	ImageUrl         string
	ImageContentType string
}

func (q *Queries) CreateThread(ctx context.Context, arg CreateThreadParams) (Thread, error) {
	row := q.db.QueryRow(ctx, createThread,
		arg.Address,
		arg.Title,
		arg.Content,
		arg.ImageFileName,
		arg.ImageUrl,
		arg.ImageContentType,
	)
	var i Thread
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Title,
		&i.Content,
		&i.ImageFileName,
		&i.ImageUrl,
		&i.ImageContentType,
		&i.Votes,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createThreadDownVote = `-- name: CreateThreadDownVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, -1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = -1, updated_at = NOW()
`

type CreateThreadDownVoteParams struct {
	Address  string
	ThreadID int64
}

func (q *Queries) CreateThreadDownVote(ctx context.Context, arg CreateThreadDownVoteParams) error {
	_, err := q.db.Exec(ctx, createThreadDownVote, arg.Address, arg.ThreadID)
	return err
}

const createThreadUnVote = `-- name: CreateThreadUnVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 0)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 0, updated_at = NOW()
`

type CreateThreadUnVoteParams struct {
	Address  string
	ThreadID int64
}

func (q *Queries) CreateThreadUnVote(ctx context.Context, arg CreateThreadUnVoteParams) error {
	_, err := q.db.Exec(ctx, createThreadUnVote, arg.Address, arg.ThreadID)
	return err
}

const createThreadUpVote = `-- name: CreateThreadUpVote :exec
INSERT INTO thread_votes (address, thread_id, vote)
VALUES ($1, $2, 1)
ON CONFLICT (address, thread_id) DO UPDATE SET vote = 1, updated_at = NOW()
`

type CreateThreadUpVoteParams struct {
	Address  string
	ThreadID int64
}

func (q *Queries) CreateThreadUpVote(ctx context.Context, arg CreateThreadUpVoteParams) error {
	_, err := q.db.Exec(ctx, createThreadUpVote, arg.Address, arg.ThreadID)
	return err
}

const deleteComment = `-- name: DeleteComment :one
UPDATE comments
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id as comment_id
`

func (q *Queries) DeleteComment(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteComment, id)
	var comment_id int64
	err := row.Scan(&comment_id)
	return comment_id, err
}

const deleteThread = `-- name: DeleteThread :one
UPDATE threads
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1
RETURNING id as thread_id
`

func (q *Queries) DeleteThread(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteThread, id)
	var thread_id int64
	err := row.Scan(&thread_id)
	return thread_id, err
}

const getComment = `-- name: GetComment :one
SELECT c.id, c.thread_id, c.replied_to_comment_id, c.address, c.content, c.image_file_name, c.image_url, c.image_content_type, c.votes, c.is_deleted, c.created_at, c.deleted_at
FROM comments c
WHERE c.id = $1
`

func (q *Queries) GetComment(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRow(ctx, getComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.ThreadID,
		&i.RepliedToCommentID,
		&i.Address,
		&i.Content,
		&i.ImageFileName,
		&i.ImageUrl,
		&i.ImageContentType,
		&i.Votes,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getComments = `-- name: GetComments :many
SELECT
	c.id, c.thread_id, c.replied_to_comment_id, c.address, c.content, c.image_file_name, c.image_url, c.image_content_type, c.votes, c.is_deleted, c.created_at, c.deleted_at,
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
LIMIT $3::bigint
`

type GetCommentsParams struct {
	ThreadID int64
	Column2  int64
	Column3  int64
}

type GetCommentsRow struct {
	ID                 int64
	ThreadID           int64
	RepliedToCommentID pgtype.Int8
	Address            string
	Content            string
	ImageFileName      string
	ImageUrl           string
	ImageContentType   string
	Votes              int64
	IsDeleted          bool
	CreatedAt          pgtype.Timestamp
	DeletedAt          pgtype.Timestamp
	RID                pgtype.Int8
	RAddress           pgtype.Text
	RContent           pgtype.Text
	RImageFileName     pgtype.Text
	RImageUrl          pgtype.Text
	RImageContentType  pgtype.Text
	RIsDeleted         pgtype.Bool
	RCreatedAt         pgtype.Timestamp
	RDeletedAt         pgtype.Timestamp
	FullCount          int64
}

func (q *Queries) GetComments(ctx context.Context, arg GetCommentsParams) ([]GetCommentsRow, error) {
	rows, err := q.db.Query(ctx, getComments, arg.ThreadID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsRow
	for rows.Next() {
		var i GetCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.ThreadID,
			&i.RepliedToCommentID,
			&i.Address,
			&i.Content,
			&i.ImageFileName,
			&i.ImageUrl,
			&i.ImageContentType,
			&i.Votes,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.RID,
			&i.RAddress,
			&i.RContent,
			&i.RImageFileName,
			&i.RImageUrl,
			&i.RImageContentType,
			&i.RIsDeleted,
			&i.RCreatedAt,
			&i.RDeletedAt,
			&i.FullCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThread = `-- name: GetThread :one
SELECT threads.id, threads.address, threads.title, threads.content, threads.image_file_name, threads.image_url, threads.image_content_type, threads.votes, threads.is_deleted, threads.created_at, threads.deleted_at
FROM threads
WHERE threads.id = $1
AND threads.is_deleted = FALSE
`

func (q *Queries) GetThread(ctx context.Context, id int64) (Thread, error) {
	row := q.db.QueryRow(ctx, getThread, id)
	var i Thread
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Title,
		&i.Content,
		&i.ImageFileName,
		&i.ImageUrl,
		&i.ImageContentType,
		&i.Votes,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getThreads = `-- name: GetThreads :many
SELECT threads.id, threads.address, threads.title, threads.content, threads.image_file_name, threads.image_url, threads.image_content_type, threads.votes, threads.is_deleted, threads.created_at, threads.deleted_at
FROM threads
WHERE threads.is_deleted = FALSE
ORDER BY RANDOM()
LIMIT $1::bigint
`

// Order by random is not performant as we need to do a full table scan.
// Move to TABLESAMPLE SYSTEM_ROWS(N) when performance becomes an issue.
// Table sample is not random enough until the table gets big.
// https://www.postgresql.org/docs/current/tsm-system-rows.html
func (q *Queries) GetThreads(ctx context.Context, dollar_1 int64) ([]Thread, error) {
	rows, err := q.db.Query(ctx, getThreads, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Thread
	for rows.Next() {
		var i Thread
		if err := rows.Scan(
			&i.ID,
			&i.Address,
			&i.Title,
			&i.Content,
			&i.ImageFileName,
			&i.ImageUrl,
			&i.ImageContentType,
			&i.Votes,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
	ens_name = $2,
	updated_at = NOW()
WHERE address = $1
`

type UpdateUserParams struct {
	Address string
	EnsName pgtype.Text
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.Address, arg.EnsName)
	return err
}

const upsertUser = `-- name: UpsertUser :exec
INSERT INTO users (address)
VALUES ($1)
ON CONFLICT (address) DO NOTHING
`

// upsert user every time they signin so we don't have to check if they exist
func (q *Queries) UpsertUser(ctx context.Context, address string) error {
	_, err := q.db.Exec(ctx, upsertUser, address)
	return err
}
