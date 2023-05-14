package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type HydrateUser struct {
	logger     common.Logger
	blockchain Blockchain
	database   Database
}

type HydrateUserInput struct {
	Address string `validate:"eth_addr"`
}

func NewHydrateUserUseCase(logger common.Logger, blockchain Blockchain, database Database) *HydrateUser {
	return &HydrateUser{
		logger,
		blockchain,
		database,
	}
}

func (u *HydrateUser) Execute(ctx context.Context, input HydrateUserInput) error {
	name, err := u.blockchain.GetNameFromAddress(ctx, input.Address)

	// If the user does not have an ENS name, we still hydrate the user with a nil value
	var ensName *string
	if err != nil {
		ensName = nil
	} else {
		ensName = &name
	}

	return u.database.UpdateUser(ctx, input.Address, ensName)
}
