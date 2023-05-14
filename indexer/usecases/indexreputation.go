package usecases

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
)

type IndexReputation struct {
	logger            common.Logger
	settings          settings.Settings
	blockchainGateway Blockchain
}

func NewIndexReputationUseCase(logger common.Logger, settings settings.Settings, blockchainGateway Blockchain) *IndexReputation {
	return &IndexReputation{
		logger,
		settings,
		blockchainGateway,
	}
}

func (u *IndexReputation) Execute(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) error {
	u.logger.Info(ctx).Msgf("indexing reputation events from block %d to block %d", fromBlock, toBlock)

	events, err := u.blockchainGateway.GetTokenEvents(ctx, fromBlock, toBlock)

	if err != nil {
		u.logger.Error(ctx).Err(err).Msg("failed to read events")
		return err
	}

	u.logger.Info(ctx).Msgf("found %v events", len(events))

	return nil
}
