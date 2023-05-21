package ethereum

import (
	"context"
	"errors"
	"fmt"

	cmn "github.com/daochanio/backend/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

// The returned name is nil if no name can be resolved from the address
func (e *ethereumGateway) GetNameByAddress(ctx context.Context, address string) (*string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return nil, fmt.Errorf("failed to dial blockchain: %w", err)
	}

	name, err := cmn.FunctionRetrier(ctx, func() (string, error) {
		name, err := ens.ReverseResolve(client, common.HexToAddress(address))
		return name, e.tryWrapRetryable(ctx, "failed to get ens name", err)
	})

	if errors.Is(err, cmn.ErrRetryable) {
		return nil, err
	}

	// non transient errors are considered as no name
	if err != nil {
		return nil, nil
	}

	return &name, nil
}

// The returned avatar uri is nil if no avatar text record can be resolved from the name
func (e *ethereumGateway) GetAvatarURIByName(ctx context.Context, name string) (*string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return nil, fmt.Errorf("failed to dial blockchain: %w", err)
	}

	resolver, err := cmn.FunctionRetrier(ctx, func() (*ens.Resolver, error) {
		resolver, err := ens.NewResolver(client, name)
		return resolver, e.tryWrapRetryable(ctx, "failed to get resolver", err)
	})

	if err != nil {
		return nil, err
	}

	uri, err := cmn.FunctionRetrier(ctx, func() (string, error) {
		uri, err := resolver.Text("avatar")
		return uri, e.tryWrapRetryable(ctx, "failed to get avatar", err)
	})

	if err != nil {
		return nil, err
	}

	if uri == "" {
		return nil, nil
	}

	return &uri, nil
}
