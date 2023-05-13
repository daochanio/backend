package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type Buffer map[string]AggregateVotesInput

type AggregateVotesUseCase struct {
	logger          common.Logger
	databaseGateway DatabaseGateway
}

func NewAggregateVotesUseCase(logger common.Logger, databaseGateway DatabaseGateway) *AggregateVotesUseCase {
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

	if err := u.databaseGateway.AggregateVotes(ctx, input.Id, input.Type); err != nil {
		u.logger.Warn(ctx).Err(err).Msgf("error aggregating thread votes for %v %v", input.Id, input.Type)
	}

	return nil
}
