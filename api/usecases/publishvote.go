package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type CreateVote struct {
	stream Stream
	logger common.Logger
}

func NewCreateVoteUseCase(stream Stream, logger common.Logger) *CreateVote {
	return &CreateVote{
		stream,
		logger,
	}
}

type CreateVoteInput struct {
	Id      int64            `validate:"gt=0"`
	Address string           `validate:"eth_addr"`
	Value   common.VoteValue `validate:"oneof=upvote downvote unvote"`
	Type    common.VoteType  `validate:"oneof=thread comment"`
}

func (u *CreateVote) Execute(ctx context.Context, input CreateVoteInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	vote := entities.NewVote(input.Id, input.Address, input.Value, input.Type, time.Now().UnixMilli())

	return u.stream.PublishVote(ctx, vote)
}
