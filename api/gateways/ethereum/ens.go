package ethereum

import (
	"context"
	"fmt"

	cmn "github.com/daochanio/backend/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

func (e *ethereumGateway) GetNameFromAddress(ctx context.Context, address string) (string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return "", fmt.Errorf("failed to dial blockchain: %w", err)
	}

	return cmn.FunctionRetrier(ctx, func() (string, error) {
		name, err := ens.ReverseResolve(client, common.HexToAddress(address))
		return name, e.tryWrapRetryable(ctx, "failed to get ens name", err)
	})
}

// TODO:
// - ens gateway: name -> avatar url
//   - https
//   - ipfs: convert to https ipfs gateway
//   - data: data uris can go directly into an image tag
//   - erc721 & erc1155: need to be owned by address
//
// - http gateway: resolve url to bytes/content-type
// - s3 gateway: upload file
func (e *ethereumGateway) GetAvatarFromAddress(ctx context.Context, name string) (string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return "", fmt.Errorf("failed to dial blockchain: %w", err)
	}

	resolver, err := cmn.FunctionRetrier(ctx, func() (*ens.Resolver, error) {
		resolver, err := ens.NewResolver(client, name)
		return resolver, e.tryWrapRetryable(ctx, "failed to get resolver", err)
	})

	if err != nil {
		return "", err
	}

	avatar, err := cmn.FunctionRetrier(ctx, func() (string, error) {
		avatar, err := resolver.Text("avatar")
		return avatar, e.tryWrapRetryable(ctx, "failed to get avatar", err)
	})

	if err != nil {
		return "", err
	}

	// TODO

	return avatar, nil
}
