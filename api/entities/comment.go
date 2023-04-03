package entities

import (
	"time"
)

type Comment struct {
	id              int64
	parentCommentId *int64
	threadId        int64
	address         string
	content         string
	isDeleted       bool
	createdAt       time.Time
	deletedAt       *time.Time
	votes           int64
}

func NewComment() Comment {
	return Comment{}
}

func (c Comment) SetId(id int64) Comment {
	c.id = id
	return c
}

func (c Comment) SetParentCommentId(parentCommentId int64) Comment {
	c.parentCommentId = &parentCommentId
	return c
}

func (c Comment) SetThreadId(threadId int64) Comment {
	c.threadId = threadId
	return c
}

func (c Comment) SetAddress(address string) Comment {
	c.address = address
	return c
}

func (c Comment) SetContent(content string) Comment {
	c.content = content
	return c
}

func (c Comment) SetIsDeleted(isDeleted bool) Comment {
	c.isDeleted = isDeleted
	return c
}

func (c Comment) SetCreatedAt(createdAt time.Time) Comment {
	c.createdAt = createdAt
	return c
}

func (c Comment) SetDeletedAt(deletedAt *time.Time) Comment {
	c.deletedAt = deletedAt
	return c
}

func (c Comment) SetVotes(votes int64) Comment {
	c.votes = votes
	return c
}

func (c Comment) GetId() int64 {
	return c.id
}

func (c Comment) GetParentCommentId() *int64 {
	return c.parentCommentId
}

func (c Comment) GetThreadId() int64 {
	return c.threadId
}

func (c Comment) GetAddress() string {
	return c.address
}

func (c Comment) GetContent() string {
	if c.isDeleted {
		return "This comment has been deleted."
	}

	return c.content
}

func (c Comment) GetIsDeleted() bool {
	return c.isDeleted
}

func (c Comment) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c Comment) GetDeletedAt() *time.Time {
	return c.deletedAt
}

func (c Comment) GetVotes() int64 {
	return c.votes
}
