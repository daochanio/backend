package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateThreadUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewCreateThreadUseCase(dbGateway gateways.IDatabaseGateway) *CreateThreadUseCase {
	return &CreateThreadUseCase{
		dbGateway,
	}
}

type CreateThreadInput struct {
	Address string `validate:"eth_addr"`
	Content string `validate:"max=1000"`
}

func (u *CreateThreadUseCase) Execute(ctx context.Context, input CreateThreadInput) (int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return 0, err
	}

	return u.dbGateway.CreateThread(ctx, input.Address, input.Content)
}
