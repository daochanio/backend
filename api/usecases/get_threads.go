package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
)

type GetThreadsUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewGetThreadsUseCase(dbGateway gateways.IDatabaseGateway) *GetThreadsUseCase {
	return &GetThreadsUseCase{
		dbGateway,
	}
}

// get threads input
type GetThreadsInput struct {
	Offset uint32
	Limit  uint32
}

func (u *GetThreadsUseCase) Execute(ctx context.Context, input GetThreadsInput) ([]entities.Thread, error) {
	return u.dbGateway.GetThreads(ctx, input.Offset, input.Limit)
}