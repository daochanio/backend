package ethereum

import (
	"context"
	"math/big"

	com "github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/entities"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethereumGateway struct {
	ethClient *ethclient.Client
	logger    com.Logger
	settings  settings.Settings
}

func NewEthereumGateway(logger com.Logger, settings settings.Settings) usecases.Blockchain {
	ethClient, err := ethclient.Dial(settings.BlockchainURI())
	if err != nil {
		panic(err)
	}

	return &ethereumGateway{
		ethClient,
		logger,
		settings,
	}
}

func (g *ethereumGateway) Start(ctx context.Context) {}

func (g *ethereumGateway) Shutdown(ctx context.Context) {}

func (g *ethereumGateway) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	header, err := g.ethClient.HeaderByNumber(ctx, nil)
	if err == ethereum.NotFound {
		return nil, err
	}

	return header.Number, nil
}

func (g *ethereumGateway) GetTokenEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]entities.TokenEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			common.HexToAddress(g.settings.TokenAddress()),
		},
	}

	logs, err := g.ethClient.FilterLogs(ctx, query)

	if err != nil {
		return nil, err
	}

	// TODO
	for _, log := range logs {
		g.logger.Info(ctx).Msgf("event log for address: %v at block: %d", log.Address.Hex(), log.BlockNumber)
	}

	return nil, nil
}
