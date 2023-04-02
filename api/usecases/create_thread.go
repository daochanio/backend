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
	Address string `json:"-" validate:"required,eth_addr"`
	Content string `json:"content" validate:"required,max=1000"`
}

func (u *CreateThreadUseCase) Execute(ctx context.Context, input CreateThreadInput) (int32, error) {
	if err := common.ValidateStruct(input); err != nil {
		return 0, err
	}

	return u.dbGateway.CreateThread(ctx, input.Address, input.Content)
}
