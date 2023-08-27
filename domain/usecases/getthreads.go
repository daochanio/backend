package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type GetThreads struct {
	database gateways.Database
	logger   common.Logger
}

func NewGetThreadsUseCase(database gateways.Database, logger common.Logger) *GetThreads {
	return &GetThreads{
		database,
		logger,
	}
}

type GetThreadsInput struct {
	Limit int64 `validate:"gte=0,lte=100"`
}

// Threads returned are random and thus the concept of pages/offset/count is not relevant
func (u *GetThreads) Execute(ctx context.Context, input GetThreadsInput) ([]entities.Thread, error) {
	if err := common.ValidateStruct(input); err != nil {
		return nil, err
	}

	return u.database.GetThreads(ctx, input.Limit)
}
