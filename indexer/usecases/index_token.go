package usecases

import (
	"context"
	"math/big"

	cmn "github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/gateways"
	"github.com/daochanio/backend/indexer/settings"
)

type IndexTokenUseCase struct {
	logger            cmn.ILogger
	settings          settings.ISettings
	blockchainGateway gateways.IBlockchainGateway
}

func NewIndexTokenUseCase(logger cmn.ILogger, settings settings.ISettings, blockchainGateway gateways.IBlockchainGateway) *IndexTokenUseCase {
	return &IndexTokenUseCase{
		logger,
		settings,
		blockchainGateway,
	}
}

func (u *IndexTokenUseCase) Execute(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) error {
	u.logger.Info(ctx).Msgf("indexing token events from block %d to block %d", fromBlock, toBlock)

	events, err := u.blockchainGateway.GetTokenEvents(ctx, fromBlock, toBlock)

	if err != nil {
		u.logger.Error(ctx).Err(err).Msg("failed to read events")
		return err
	}

	u.logger.Info(ctx).Msgf("found %v events", len(events))

	return nil
}
