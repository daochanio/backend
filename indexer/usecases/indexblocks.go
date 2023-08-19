package usecases

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
)

type IndexBlocks struct {
	logger          common.Logger
	settings        settings.Settings
	database        Database
	blockchain      Blockchain
	indexReputation *IndexReputation
}

func NewIndexBlocksUseCase(
	logger common.Logger,
	settings settings.Settings,
	database Database,
	blockchain Blockchain,
	indexReputation *IndexReputation) *IndexBlocks {
	return &IndexBlocks{
		logger,
		settings,
		database,
		blockchain,
		indexReputation,
	}
}

// Execute checks the last block indexed (minus an offset) and the latest block produced and indexes all blocks in between.
// Executre returns an error if it failed to fully index the blocks.
// We want to make indexing idempotent and be resilient to re-orgs so we:
//   - keep track of last block indexed
//   - read events from last block indexed minus offset to lastest block
//   - always delete existing events for the blocks being indexed
func (u *IndexBlocks) Execute(ctx context.Context) error {
	fromBlock, toBlock, err := u.getBlockRange(ctx)

	if err != nil {
		return fmt.Errorf("failed to get block range: %w", err)
	}

	u.logger.Info(ctx).Msgf("indexing from block %d to block %d", fromBlock, toBlock)

	events, err := u.blockchain.GetEvents(ctx, fromBlock, toBlock)

	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	var wg sync.WaitGroup
	var indexReputationErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		indexReputationErr = u.indexReputation.Execute(ctx, fromBlock, toBlock, events.Transfers())
	}()
	wg.Wait()

	if indexReputationErr != nil {
		return fmt.Errorf("failed to index reputation: %w", indexReputationErr)
	}

	err = u.database.UpdateLastIndexedBlock(ctx, toBlock)
	if err != nil {
		return fmt.Errorf("failed to update last indexed block: %w", err)
	}

	u.logger.Info(ctx).Msgf("indexed block %d to block %d", fromBlock, toBlock)

	return nil
}

func (u *IndexBlocks) getBlockRange(ctx context.Context) (*big.Int, *big.Int, error) {
	lastBlock, err := u.database.GetLastIndexedBlock(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get last indexed block: %w", err)
	}

	latestBlock, err := u.blockchain.GetLatestBlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	if lastBlock.Cmp(latestBlock) == 0 {
		return nil, nil, common.ErrNoNewBlocks
	}

	if lastBlock.Cmp(latestBlock) > 0 {
		u.logger.Warn(ctx).Msgf("last indexed block %d is greater than latest block %d", lastBlock, latestBlock)
		lastBlock = latestBlock
	}

	maxBlockRange := u.settings.MaxBlockRange()
	if big.NewInt(0).Sub(latestBlock, lastBlock).Cmp(big.NewInt(maxBlockRange)) > 0 {
		latestBlock = big.NewInt(0).Add(lastBlock, big.NewInt(maxBlockRange))
	}

	offsetBlock := big.NewInt(0).Sub(lastBlock, big.NewInt(u.settings.ReorgOffset()))

	if offsetBlock.Cmp(big.NewInt(0)) < 0 {
		offsetBlock = big.NewInt(0)
	}

	return offsetBlock, latestBlock, nil
}
