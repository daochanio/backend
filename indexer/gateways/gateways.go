package gateways

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/indexer/entities"
)

type IDatabaseGateway interface {
	GetLastIndexedBlock(context.Context) (*big.Int, error)
	UpdateLastIndexedBlock(context.Context, *big.Int) error
}

type IBlockchainGateway interface {
	DoesBlockExist(ctx context.Context, blockNumber *big.Int) bool
	GetTokenEvents(context.Context, *big.Int, *big.Int) ([]entities.TokenEvent, error)
}
