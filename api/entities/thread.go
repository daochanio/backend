package entities

import (
	"time"
)

type Thread struct {
	id        int64
	address   string
	title     string
	content   string
	image     Image
	comments  *[]Comment
	isDeleted bool
	createdAt time.Time
	deletedAt *time.Time
	votes     int64
}

type ThreadParams struct {
	Id        int64
	Address   string
	Title     string
	Content   string
	Image     Image
	Comments  *[]Comment
	IsDeleted bool
	CreatedAt time.Time
	DeletedAt *time.Time
	Votes     int64
}

func NewThread(params ThreadParams) Thread {
	return Thread{
		id:        params.Id,
		address:   params.Address,
		title:     params.Title,
		content:   params.Content,
		image:     params.Image,
		comments:  params.Comments,
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

func (t *Thread) Title() string {
	return t.title
}

func (t *Thread) Content() string {
	if t.isDeleted {
		return "This thread has been deleted."
	}

	return t.content
}

func (t *Thread) Image() *Image {
	return &t.image
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

func (t *Thread) SetComments(comments *[]Comment) {
	t.comments = comments
}

// returned comments can be nil if not hydrated
func (t *Thread) Comments() *[]Comment {
	return t.comments
}
