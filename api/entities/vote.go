package entities

import "github.com/daochanio/backend/common"

type Vote struct {
	id       int64
	address  string
	value    common.VoteValue
	voteType common.VoteType
}

func NewVote(id int64, address string, value common.VoteValue, voteType common.VoteType) Vote {
	return Vote{
		id:       id,
		address:  address,
		value:    value,
		voteType: voteType,
	}
}

func (v *Vote) Id() int64 {
	return v.id
}

func (v *Vote) Address() string {
	return v.address
}

func (v *Vote) Value() common.VoteValue {
	return v.value
}

func (v *Vote) Type() common.VoteType {
	return v.voteType
}
