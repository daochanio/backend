package entities

import (
	"time"
)

type Comment struct {
	ID        int32
	ThreadID  int32
	Address   string
	Content   string
	IsDeleted bool
	CreatedAt time.Time
	DeletedAt *time.Time
	Votes     int64
}

func NewComment() Comment {
	return Comment{}
}

func (c Comment) SetId(id int32) Comment {
	c.ID = id
	return c
}

func (c Comment) SetThreadId(threadId int32) Comment {
	c.ThreadID = threadId
	return c
}

func (c Comment) SetAddress(address string) Comment {
	c.Address = address
	return c
}

func (c Comment) SetContent(content string) Comment {
	c.Content = content
	return c
}

func (c Comment) SetIsDeleted(isDeleted bool) Comment {
	c.IsDeleted = isDeleted
	return c
}

func (c Comment) SetCreatedAt(createdAt time.Time) Comment {
	c.CreatedAt = createdAt
	return c
}

func (c Comment) SetDeletedAt(deletedAt *time.Time) Comment {
	c.DeletedAt = deletedAt
	return c
}

func (c Comment) SetVotes(votes int64) Comment {
	c.Votes = votes
	return c
}

func (c Comment) GetId() int32 {
	return c.ID
}

func (c Comment) GetThreadId() int32 {
	return c.ThreadID
}

func (c Comment) GetAddress() string {
	return c.Address
}

func (c Comment) GetContent() string {
	if c.IsDeleted {
		return "This comment has been deleted."
	}

	return c.Content
}

func (c Comment) GetIsDeleted() bool {
	return c.IsDeleted
}

func (c Comment) GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c Comment) GetDeletedAt() *time.Time {
	return c.DeletedAt
}

func (c Comment) GetVotes() int64 {
	return c.Votes
}
