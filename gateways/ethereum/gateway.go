package ethereum

import (
	"context"
	"errors"
	"fmt"

	com "github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
	"github.com/daochanio/backend/gateways/ethereum/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ethereumGateway struct {
	logger       com.Logger
	ethClient    *ethclient.Client
	eventSources []common.Address
	reputation   *bindings.Reputation
}

func NewEthereumGateway(logger com.Logger) gateways.Blockchain {
	return &ethereumGateway{
		logger,
		nil,
		[]common.Address{},
		nil,
	}
}

func (g *ethereumGateway) Start(ctx context.Context, config gateways.BlockchainConfig) {
	g.logger.Info(ctx).Msg("starting ethereum gateway")

	ethClient, err := ethclient.Dial(config.BlockchainURL)

	if err != nil {
		panic(err)
	}

	g.ethClient = ethClient

	reputationAddress := common.HexToAddress(config.ReputationAddress)
	reputation, err := bindings.NewReputation(reputationAddress, ethClient)

	if err != nil {
		panic(err)
	}

	g.eventSources = append(g.eventSources, reputationAddress)
	g.reputation = reputation
}

func (g *ethereumGateway) Shutdown(ctx context.Context) {
	g.logger.Info(ctx).Msg("shutting down ethereum gateway")

	g.ethClient.Close()
}

// Parse for go-ethereum http error to determine if its retryable.
// Wrap in common.ErrRetryable if status code is 429 and include msg
func (e *ethereumGateway) tryWrapRetryable(ctx context.Context, msg string, err error) error {
	if err == nil {
		return nil
	}

	var httpErr rpc.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == 429 {
		e.logger.Warn(ctx).Err(err).Msgf("eth rate limit: %v", msg)
		return fmt.Errorf("%v %v %w", msg, err, com.ErrRetryable)
	}

	return err
}
