package usecases

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/indexer/entities"
)

type Database interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	GetLastIndexedBlock(context.Context) (*big.Int, error)
	UpdateLastIndexedBlock(context.Context, *big.Int) error
}

type Blockchain interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	GetTokenEvents(context.Context, *big.Int, *big.Int) ([]entities.TokenEvent, error)
}
