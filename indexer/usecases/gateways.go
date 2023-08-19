package usecases

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/indexer/entities"
)

type Database interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	GetLastIndexedBlock(ctx context.Context) (*big.Int, error)
	UpdateLastIndexedBlock(ctx context.Context, block *big.Int) error
	InsertTransferEvents(ctx context.Context, from *big.Int, to *big.Int, transfers []entities.Transfer) error
	UpdateReputation(ctx context.Context, addresses []string) error
}

type Blockchain interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	GetEvents(context.Context, *big.Int, *big.Int) (entities.Events, error)
}
