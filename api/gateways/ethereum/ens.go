package ethereum

import (
	"context"

	cmn "github.com/daochanio/backend/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

func (e *ethereumGateway) GetENSNameFromAddress(ctx context.Context, address string) (string, error) {
	return cmn.FunctionRetrier(ctx, func() (string, error) {
		name, err := ens.ReverseResolve(e.client, common.HexToAddress(address))
		return name, e.tryWrapRetryable(ctx, "get ens name from address", err)
	})
}
