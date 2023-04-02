package entities

import (
	"time"
)

type Thread struct {
	ID        int32
	Address   string
	Content   string
	IsDeleted bool
	CreatedAt time.Time
	DeletedAt *time.Time
	Votes     int64
}

func NewThread() Thread {
	return Thread{}
}

func (t Thread) SetId(id int32) Thread {
	t.ID = id
	return t
}

func (t Thread) SetAddress(address string) Thread {
	t.Address = address
	return t
}

func (t Thread) SetContent(content string) Thread {
	t.Content = content
	return t
}

func (t Thread) SetIsDeleted(isDeleted bool) Thread {
	t.IsDeleted = isDeleted
	return t
}

func (t Thread) SetCreatedAt(createdAt time.Time) Thread {
	t.CreatedAt = createdAt
	return t
}

func (t Thread) SetDeletedAt(deletedAt *time.Time) Thread {
	t.DeletedAt = deletedAt
	return t
}

func (t Thread) SetVotes(votes int64) Thread {
	t.Votes = votes
	return t
}

func (t Thread) GetId() int32 {
	return t.ID
}

func (t Thread) GetAddress() string {
	return t.Address
}

func (t Thread) GetContent() string {
	if t.IsDeleted {
		return "This thread has been deleted."
	}

	return t.Content
}

func (t Thread) GetIsDeleted() bool {
	return t.IsDeleted
}

func (t Thread) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t Thread) GetDeletedAt() *time.Time {
	return t.DeletedAt
}

func (t Thread) GetVotes() int64 {
	return t.Votes
}
