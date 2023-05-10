package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type GetThreadsUseCase struct {
	dbGateway DatabaseGateway
	logger    common.Logger
}

func NewGetThreadsUseCase(dbGateway DatabaseGateway, logger common.Logger) *GetThreadsUseCase {
	return &GetThreadsUseCase{
		dbGateway,
		logger,
	}
}

type GetThreadsInput struct {
	Limit int64 `validate:"gte=0,lte=100"`
}

// Threads returned are random and thus the concept of pages/offset/count is not relevant
func (u *GetThreadsUseCase) Execute(ctx context.Context, input GetThreadsInput) ([]entities.Thread, error) {
	if err := common.ValidateStruct(input); err != nil {
		return nil, err
	}

	t1 := time.Now()
	threads, err := u.dbGateway.GetThreads(ctx, input.Limit)
	u.logger.Info(ctx).Msgf("get threads duration %v", time.Since(t1))

	return threads, err
}
