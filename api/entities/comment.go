package entities

import (
	"time"
)

type Comment struct {
	id               int64
	repliedToComment *Comment
	threadId         int64
	address          string
	content          string
	isDeleted        bool
	createdAt        time.Time
	deletedAt        *time.Time
	votes            int64
}

type CommentParams struct {
	Id               int64
	RepliedToComment *Comment
	ThreadId         int64
	Address          string
	Content          string
	IsDeleted        bool
	CreatedAt        time.Time
	DeletedAt        *time.Time
	Votes            int64
}

func NewComment(params CommentParams) Comment {
	return Comment{
		id:               params.Id,
		repliedToComment: params.RepliedToComment,
		threadId:         params.ThreadId,
		address:          params.Address,
		content:          params.Content,
		isDeleted:        params.IsDeleted,
		createdAt:        params.CreatedAt,
		deletedAt:        params.DeletedAt,
		votes:            params.Votes,
	}
}

func (c *Comment) Id() int64 {
	return c.id
}

func (c *Comment) RepliedToComment() *Comment {
	return c.repliedToComment
}

func (c *Comment) ThreadId() int64 {
	return c.threadId
}

func (c *Comment) Address() string {
	return c.address
}

func (c *Comment) Content() string {
	if c.isDeleted {
		return "This comment has been deleted."
	}
	return c.content
}

func (c *Comment) IsDeleted() bool {
	return c.isDeleted
}

func (c *Comment) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Comment) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *Comment) Votes() int64 {
	return c.votes
}
