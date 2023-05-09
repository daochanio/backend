package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type CreateUserUseCase struct {
	logger            common.Logger
	databaseGateway   DatabaseGateway
	blockchainGateway BlockchainGateway
}

func NewCreateUserUseCase(logger common.Logger, databaseGateway DatabaseGateway, blockchainGateway BlockchainGateway) *CreateUserUseCase {
	return &CreateUserUseCase{
		logger,
		databaseGateway,
		blockchainGateway,
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
	name, err := u.blockchainGateway.GetNameFromAddress(ctx, input.Address)
	if err != nil {
		u.logger.Info(ctx).Err(err).Msg("Failed to get name from address")
	} else {
		ensName = &name
	}

	return u.databaseGateway.CreateOrUpdateUser(ctx, input.Address, ensName)
}
