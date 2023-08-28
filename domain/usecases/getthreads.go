package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type GetThreads struct {
	logger    common.Logger
	validator common.Validator
	database  gateways.Database
}

func NewGetThreadsUseCase(logger common.Logger, validator common.Validator, database gateways.Database) *GetThreads {
	return &GetThreads{
		logger,
		validator,
		database,
	}
}

type GetThreadsInput struct {
	Limit int64 `validate:"gte=0,lte=100"`
}

// Threads returned are random and thus the concept of pages/offset/count is not relevant
func (u *GetThreads) Execute(ctx context.Context, input GetThreadsInput) ([]entities.Thread, error) {
	if err := u.validator.ValidateStruct(input); err != nil {
		return nil, err
	}

	return u.database.GetThreads(ctx, input.Limit)
}
