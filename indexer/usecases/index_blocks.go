package usecases

import (
	"context"
	"math/big"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/gateways"
	"github.com/daochanio/backend/indexer/settings"
)

type IndexBlocksUseCase struct {
	logger            common.ILogger
	settings          settings.ISettings
	databaseGateway   gateways.IDatabaseGateway
	blockchainGateway gateways.IBlockchainGateway
	indexRepUseCase   *IndexTokenUseCase
}

func NewIndexBlocksUseCase(
	logger common.ILogger,
	settings settings.ISettings,
	indexRepUseCase *IndexTokenUseCase,
	databaseGateway gateways.IDatabaseGateway,
	blockchainGateway gateways.IBlockchainGateway) *IndexBlocksUseCase {
	return &IndexBlocksUseCase{
		logger,
		settings,
		databaseGateway,
		blockchainGateway,
		indexRepUseCase,
	}
}

// Execute is a blocking call that look for and indexes unindexed blocks
// Execute checks the last block indexed and gets the header for that block to see if it has been produced.
// If getting the header errors, we assume the block has not been produced yet and sleep before checking again.
// We do not return an error here because we want to keep checking for the next block forever.
// We want to make indexing idempotent and be resilient to re-orgs so we:
//   - keep track of last block indexed
//   - read events from last block indexed - n to last block indexed + 1
//   - always delete events that already exist for the same block being inserted (in single tx)
func (u *IndexBlocksUseCase) Execute(ctx context.Context) error {
	for {
		blockIndexed := u.indexNextBlock(ctx)

		// sleep a bit if a block was not indexed to avoid spamming / allow transient errors to recover
		if !blockIndexed {
			time.Sleep(5 * time.Second)
		}
	}
}

func (u *IndexBlocksUseCase) indexNextBlock(ctx context.Context) bool {
	blockNumber, err := u.databaseGateway.GetLastIndexedBlock(ctx)

	if err != nil {
		u.logger.Warn(ctx).Err(err).Msg("failed to get last indexed block")
		return false
	}

	nextBlock := blockNumber.Add(blockNumber, big.NewInt(1))

	exists := u.blockchainGateway.DoesBlockExist(ctx, nextBlock)

	if !exists {
		return false
	}

	u.logger.Info(ctx).Msgf("found block %d", nextBlock)

	fromBlock := big.NewInt(0).Sub(nextBlock, big.NewInt(u.settings.ReorgOffset()))
	toBlock := nextBlock

	// TODO: call indexer usecases here

	err = u.databaseGateway.UpdateLastIndexedBlock(ctx, nextBlock)

	if err != nil {
		u.logger.Warn(ctx).Err(err).Msg("failed to update last indexed block")
		return false
	}

	u.logger.Info(ctx).Msgf("indexed block %d to block %d", fromBlock, toBlock)

	return true
}
