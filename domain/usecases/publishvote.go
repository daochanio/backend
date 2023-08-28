package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type CreateVote struct {
	logger    common.Logger
	validator common.Validator
	stream    gateways.Stream
}

func NewCreateVoteUseCase(logger common.Logger, validator common.Validator, stream gateways.Stream) *CreateVote {
	return &CreateVote{
		logger,
		validator,
		stream,
	}
}

type CreateVoteInput struct {
	Id      int64              `validate:"gt=0"`
	Address string             `validate:"eth_addr"`
	Value   entities.VoteValue `validate:"oneof=upvote downvote unvote"`
	Type    entities.VoteType  `validate:"oneof=thread comment"`
}

func (u *CreateVote) Execute(ctx context.Context, input CreateVoteInput) error {
	if err := u.validator.ValidateStruct(input); err != nil {
		return err
	}

	vote := entities.NewVote(input.Id, input.Address, input.Value, input.Type, time.Now().UnixMilli())

	return u.stream.PublishVote(ctx, vote)
}
