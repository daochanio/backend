package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type Buffer map[string]AggregateVotesInput

type AggregateVotes struct {
	logger   common.Logger
	database Database
}

func NewAggregateVotesUseCase(logger common.Logger, database Database) *AggregateVotes {
	return &AggregateVotes{
		logger,
		database,
	}
}

type AggregateVotesInput struct {
	Id   int64           `validate:"gt=0"`
	Type common.VoteType `validate:"oneof=thread comment"`
}

func (u *AggregateVotes) Execute(ctx context.Context, input AggregateVotesInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	if err := u.database.AggregateVotes(ctx, input.Id, input.Type); err != nil {
		u.logger.Warn(ctx).Err(err).Msgf("error aggregating thread votes for %v %v", input.Id, input.Type)
	}

	return nil
}
