package usecases

import (
	"context"
	"errors"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
)

type IndexBlocksUseCase struct {
	logger            common.Logger
	settings          settings.Settings
	databaseGateway   DatabaseGateway
	blockchainGateway BlockchainGateway
	indexRepUseCase   *IndexTokenUseCase
}

func NewIndexBlocksUseCase(
	logger common.Logger,
	settings settings.Settings,
	indexRepUseCase *IndexTokenUseCase,
	databaseGateway DatabaseGateway,
	blockchainGateway BlockchainGateway) *IndexBlocksUseCase {
	return &IndexBlocksUseCase{
		logger,
		settings,
		databaseGateway,
		blockchainGateway,
		indexRepUseCase,
	}
}

// Execute checks the last block indexed (minus an offset) and the latest block produced and indexes all blocks in between.
// Executre returns an error if it failed to fully index the blocks.
// We want to make indexing idempotent and be resilient to re-orgs so we:
//   - keep track of last block indexed
//   - read events from last block indexed - n to lastest block
//   - always delete events that already exist for the same block being inserted (in single tx)
func (u *IndexBlocksUseCase) Execute(ctx context.Context) error {
	lastBlockNumber, err := u.databaseGateway.GetLastIndexedBlock(ctx)
	if err != nil {
		u.logger.Warn(ctx).Err(err).Msg("failed to get last indexed block")
		return err
	}
	fromBlock := big.NewInt(0).Sub(lastBlockNumber, big.NewInt(u.settings.ReorgOffset()))

	toBlock, err := u.blockchainGateway.GetLatestBlockNumber(ctx)
	if err != nil {
		u.logger.Warn(ctx).Err(err).Msg("failed to get latest block")
		return err
	}

	if lastBlockNumber.Cmp(toBlock) == 0 {
		return errors.New("no new blocks")
	}

	// TODO: call indexer usecases here

	err = u.databaseGateway.UpdateLastIndexedBlock(ctx, toBlock)

	if err != nil {
		u.logger.Warn(ctx).Err(err).Msg("failed to update last indexed block")
		return err
	}

	u.logger.Info(ctx).Msgf("indexed block %d to block %d", fromBlock, toBlock)

	return nil
}
