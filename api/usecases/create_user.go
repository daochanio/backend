package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateUserUseCase struct {
	logger common.Logger
	db     gateways.DatabaseGateway
	bc     gateways.BlockchainGateway
}

func NewCreateUserUseCase(logger common.Logger, db gateways.DatabaseGateway, bc gateways.BlockchainGateway) *CreateUserUseCase {
	return &CreateUserUseCase{
		logger,
		db,
		bc,
	}
}

type CreateUserInput struct {
	Address string `validate:"eth_addr"`
}

func (u *CreateUserUseCase) Execute(ctx context.Context, input *CreateUserInput) (entities.User, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.User{}, err
	}

	var ensName *string
	name, err := u.bc.GetENSNameFromAddress(ctx, input.Address)
	if err != nil {
		u.logger.Info(ctx).Err(err).Msg("Failed to get ENS name from address")
	} else {
		ensName = &name
	}

	return u.db.CreateOrUpdateUser(ctx, input.Address, ensName)
}
