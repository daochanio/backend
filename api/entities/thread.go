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

type ThreadParams struct {
	Id        int64
	Address   string
	Content   string
	IsDeleted bool
	CreatedAt time.Time
	DeletedAt *time.Time
	Votes     int64
}

func NewThread(params ThreadParams) Thread {
	return Thread{
		id:        params.Id,
		address:   params.Address,
		content:   params.Content,
		isDeleted: params.IsDeleted,
		createdAt: params.CreatedAt,
		deletedAt: params.DeletedAt,
		votes:     params.Votes,
	}
}

func (t *Thread) Id() int64 {
	return t.id
}

func (t *Thread) Address() string {
	return t.address
}

func (t *Thread) Content() string {
	if t.isDeleted {
		return "This thread has been deleted."
	}

	return t.content
}

func (t *Thread) IsDeleted() bool {
	return t.isDeleted
}

func (t *Thread) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Thread) DeletedAt() *time.Time {
	return t.deletedAt
}

func (t *Thread) Votes() int64 {
	return t.votes
}
