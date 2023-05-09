package usecases

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/indexer/entities"
)

type DatabaseGateway interface {
	GetLastIndexedBlock(context.Context) (*big.Int, error)
	UpdateLastIndexedBlock(context.Context, *big.Int) error
}

type BlockchainGateway interface {
	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	GetTokenEvents(context.Context, *big.Int, *big.Int) ([]entities.TokenEvent, error)
}
