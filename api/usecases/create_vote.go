package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateVoteUseCase struct {
	dbGateway      gateways.DatabaseGateway
	messageGateway gateways.MessageGateway
	logger         common.Logger
}

func NewCreateVoteUseCase(dbGateway gateways.DatabaseGateway, messageGateway gateways.MessageGateway, logger common.Logger) *CreateVoteUseCase {
	return &CreateVoteUseCase{
		dbGateway,
		messageGateway,
		logger,
	}
}

type CreateVoteInput struct {
	Id      int64            `validate:"gt=0"`
	Address string           `validate:"eth_addr"`
	Value   common.VoteValue `validate:"oneof=upvote downvote unvote"`
	Type    common.VoteType  `validate:"oneof=thread comment"`
}

func (u *CreateVoteUseCase) Execute(ctx context.Context, input CreateVoteInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	vote := entities.NewVote(input.Id, input.Address, input.Value, input.Type)

	switch input.Type {
	case common.ThreadVote:
		if err := u.dbGateway.CreateThreadVote(ctx, vote); err != nil {
			return err
		}
	case common.CommentVote:
		if err := u.dbGateway.CreateCommentVote(ctx, vote); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid vote type: %v", input.Type)
	}

	if err := u.messageGateway.PublishVote(ctx, vote); err != nil {
		// Current thinking, we don't want to fail the request if the message fails to publish
		// But internally its important to log the error at an elevated level
		u.logger.Error(ctx).Err(err).Msg("error publishing comment vote message")
	}
	return nil
}
