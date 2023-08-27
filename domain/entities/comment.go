package entities

import (
	"time"
)

type Comment struct {
	id               int64
	repliedToComment *Comment
	threadId         int64
	content          string
	image            Image
	user             User
	isDeleted        bool
	createdAt        time.Time
	deletedAt        *time.Time
	votes            int64
}

type CommentParams struct {
	Id               int64
	RepliedToComment *Comment
	ThreadId         int64
	Content          string
	Image            Image
	User             User
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
		content:          params.Content,
		image:            params.Image,
		user:             params.User,
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

func (c *Comment) SetRepliedToComment(comment *Comment) {
	c.repliedToComment = comment
}

func (c *Comment) ThreadId() int64 {
	return c.threadId
}

func (c *Comment) Content() string {
	if c.isDeleted {
		return "This comment has been deleted."
	}

	return c.content
}

func (c *Comment) Image() *Image {
	if c.isDeleted {
		return nil
	}

	return &c.image
}

func (c *Comment) User() User {
	return c.user
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
