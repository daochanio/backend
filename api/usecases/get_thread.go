package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
)

type GetThreadUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewGetThreadUseCase(dbGateway gateways.IDatabaseGateway) *GetThreadUseCase {
	return &GetThreadUseCase{
		dbGateway,
	}
}

// get thread input
type GetThreadInput struct {
	ThreadId int32
}

func (u *GetThreadUseCase) Execute(ctx context.Context, input GetThreadInput) (entities.Thread, error) {
	return u.dbGateway.GetThreadById(ctx, input.ThreadId)
}
