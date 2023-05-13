package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type HydrateUserUseCase struct {
	logger            common.Logger
	blockchainGateway BlockchainGateway
	databaseGateway   DatabaseGateway
}

type HydrateUserUseCaseInput struct {
	Address string `validate:"eth_addr"`
}

func NewHydrateUserUseCase(logger common.Logger, blochainGateway BlockchainGateway, databaseGateway DatabaseGateway) *HydrateUserUseCase {
	return &HydrateUserUseCase{
		logger,
		blochainGateway,
		databaseGateway,
	}
}

func (u *HydrateUserUseCase) Execute(ctx context.Context, input HydrateUserUseCaseInput) error {
	name, err := u.blockchainGateway.GetNameFromAddress(ctx, input.Address)

	// If the user does not have an ENS name, we still hydrate the user with a nil value
	var ensName *string
	if err != nil {
		ensName = nil
	} else {
		ensName = &name
	}

	return u.databaseGateway.UpdateUser(ctx, input.Address, ensName)
}
