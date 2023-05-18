package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type HydrateUsers struct {
	logger     common.Logger
	blockchain Blockchain
	database   Database
}

type HydrateUsersInput struct {
	Addresses []string
}

func NewHydrateUsersUseCase(logger common.Logger, blockchain Blockchain, database Database) *HydrateUsers {
	return &HydrateUsers{
		logger,
		blockchain,
		database,
	}
}

// We dedupe addresses to ensure we only processes each address once regardless of multiple updates
func (u *HydrateUsers) Execute(ctx context.Context, input HydrateUsersInput) {
	addresses := map[string]bool{}
	for _, address := range input.Addresses {
		addresses[address] = true
	}

	u.logger.Info(ctx).Msgf("hydrating %v users", len(addresses))

	for address := range addresses {
		name, err := u.blockchain.GetNameFromAddress(ctx, address)

		// If the user does not have an ENS name, we still hydrate the user with a nil value
		var ensName *string
		if err != nil {
			ensName = nil
		} else {
			ensName = &name
		}

		err = u.database.UpdateUser(ctx, address, ensName)
		if err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error hydrating user %v", address)
		}
	}
}
