package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type GetThreadsUseCase struct {
	dbGateway DatabaseGateway
}

func NewGetThreadsUseCase(dbGateway DatabaseGateway) *GetThreadsUseCase {
	return &GetThreadsUseCase{
		dbGateway,
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

	return u.dbGateway.GetThreads(ctx, input.Limit)
}
