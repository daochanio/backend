package ethereum

import (
	"context"

	cmn "github.com/daochanio/backend/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

func (e *ethereumGateway) GetNameFromAddress(ctx context.Context, address string) (string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		panic(err)
	}

	return cmn.FunctionRetrier(ctx, func() (string, error) {
		name, err := ens.ReverseResolve(client, common.HexToAddress(address))
		return name, e.tryWrapRetryable(ctx, "get ens name from address", err)
	})
}
