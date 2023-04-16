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
	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	GetTokenEvents(context.Context, *big.Int, *big.Int) ([]entities.TokenEvent, error)
}
