package entities

import (
	"time"
)

type Thread struct {
	id        int64
	address   string
	content   string
	isDeleted bool
	createdAt time.Time
	deletedAt *time.Time
	votes     int64
}

func NewThread() Thread {
	return Thread{}
}

func (t Thread) SetId(id int64) Thread {
	t.id = id
	return t
}

func (t Thread) SetAddress(address string) Thread {
	t.address = address
	return t
}

func (t Thread) SetContent(content string) Thread {
	t.content = content
	return t
}

func (t Thread) SetIsDeleted(isDeleted bool) Thread {
	t.isDeleted = isDeleted
	return t
}

func (t Thread) SetCreatedAt(createdAt time.Time) Thread {
	t.createdAt = createdAt
	return t
}

func (t Thread) SetDeletedAt(deletedAt *time.Time) Thread {
	t.deletedAt = deletedAt
	return t
}

func (t Thread) SetVotes(votes int64) Thread {
	t.votes = votes
	return t
}

func (t Thread) GetId() int64 {
	return t.id
}

func (t Thread) GetAddress() string {
	return t.address
}

func (t Thread) GetContent() string {
	if t.isDeleted {
		return "This thread has been deleted."
	}

	return t.content
}

func (t Thread) GetIsDeleted() bool {
	return t.isDeleted
}

func (t Thread) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t Thread) GetDeletedAt() *time.Time {
	return t.deletedAt
}

func (t Thread) GetVotes() int64 {
	return t.votes
}
