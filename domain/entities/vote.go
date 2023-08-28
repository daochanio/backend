package entities

type Vote struct {
	id        int64
	address   string
	value     VoteValue
	voteType  VoteType
	updatedAt int64
}

type VoteValue string

const (
	Upvote   VoteValue = "upvote"
	Downvote VoteValue = "downvote"
	Unvote   VoteValue = "unvote"
)

type VoteType string

const (
	ThreadVote  VoteType = "thread"
	CommentVote VoteType = "comment"
)

func NewVote(id int64, address string, value VoteValue, voteType VoteType, updatedAt int64) Vote {
	return Vote{
		id:        id,
		address:   address,
		value:     value,
		voteType:  voteType,
		updatedAt: updatedAt,
	}
}

func (v *Vote) Id() int64 {
	return v.id
}

func (v *Vote) Address() string {
	return v.address
}

func (v *Vote) Value() VoteValue {
	return v.value
}

func (v *Vote) Type() VoteType {
	return v.voteType
}

func (v *Vote) UpdatedAt() int64 {
	return v.updatedAt
}
