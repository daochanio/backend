package ethereum

import (
	"context"
	"errors"
	"fmt"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ethereumGateway struct {
	settings settings.Settings
	logger   common.Logger
	client   *ethclient.Client
}

func NewBlockchainGateway(ctx context.Context, settings settings.Settings, logger common.Logger) gateways.BlockchainGateway {
	client, err := ethclient.DialContext(ctx, settings.BlockchainURI())

	if err != nil {
		panic(err)
	}

	return &ethereumGateway{
		settings,
		logger,
		client,
	}
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
		return fmt.Errorf("%v %v %w", msg, err, common.ErrRetryable)
	}

	return err
}
