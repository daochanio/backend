package usecases

import (
	"context"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

const ZeroAddress = "0x0000000000000000000000000000000000000000"

type IndexReputation struct {
	logger   common.Logger
	database gateways.Database
}

func NewIndexReputationUseCase(
	logger common.Logger,
	database gateways.Database,
) *IndexReputation {
	return &IndexReputation{
		logger,
		database,
	}
}

// Remove all pre-existing transfers events for the blocks being indexed as a re-org could create orphaned events that need to be cleaned up.
// Insert new transfers.
// Track dirty addresses and set new reputation values.
// We must zero all addresses first, as an address could have a negative transfer record but not positive and vise versa, throwing off the math.
func (u *IndexReputation) Execute(ctx context.Context, from *big.Int, to *big.Int, transfers []entities.Transfer) error {

	err := u.database.InsertTransferEvents(ctx, from, to, transfers)

	if err != nil {
		return fmt.Errorf("failed to insert transfer events: %w", err)
	}

	dirtyAddresses := map[string]bool{}
	for _, transfer := range transfers {
		if address := transfer.FromAddress(); address != ZeroAddress {
			dirtyAddresses[address] = true
		}

		if address := transfer.ToAddress(); address != ZeroAddress {
			dirtyAddresses[address] = true
		}
	}

	addresses := []string{}
	for address := range dirtyAddresses {
		addresses = append(addresses, address)
	}

	err = u.database.UpdateReputation(ctx, addresses)

	if err != nil {
		return fmt.Errorf("failed to update reputation: %w", err)
	}

	return nil
}
