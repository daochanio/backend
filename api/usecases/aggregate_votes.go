package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type Buffer map[string]AggregateVotesInput

type AggregateVotesUseCase struct {
	logger          common.Logger
	databaseGateway gateways.DatabaseGateway
}

func NewAggregateVotesUseCase(logger common.Logger, databaseGateway gateways.DatabaseGateway) *AggregateVotesUseCase {
	return &AggregateVotesUseCase{
		logger,
		databaseGateway,
	}
}

type AggregateVotesInput struct {
	Id   int64           `validate:"gt=0"`
	Type common.VoteType `validate:"oneof=thread comment"`
}

func (u *AggregateVotesUseCase) Execute(ctx context.Context, input AggregateVotesInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	u.logger.Info(ctx).Msgf("aggregating votes for %v %v", input.Type, input.Id)
	switch input.Type {
	case common.ThreadVote:
		if err := u.databaseGateway.AggregateThreadVotes(ctx, input.Id); err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("error aggregating thread votes for %v %v", input.Id, input.Type)
		}
	case common.CommentVote:
		if err := u.databaseGateway.AggregateCommentVotes(ctx, input.Id); err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("error aggregating comment votes for %v %v", input.Id, input.Type)
		}
	default:
		u.logger.Warn(ctx).Msgf("invalid vote type: %v %v", input.Id, input.Type)
	}

	return nil
}
