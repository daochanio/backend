package postgres

import (
	"context"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/db/bindings"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	queries *bindings.Queries
}

func NewPostgresGateway(ctx context.Context, settings settings.Settings) usecases.Database {
	db, err := pgxpool.NewWithConfig(ctx, settings.PostgresConfig())
	if err != nil {
		panic(err)
	}
	queries := bindings.New(db)
	return &postgresGateway{
		queries,
	}
}

func (g *postgresGateway) GetLastIndexedBlock(ctx context.Context) (*big.Int, error) {
	blockNumberStr, err := g.queries.GetLastIndexedBlock(ctx)

	if err != nil {
		return nil, err
	}

	blockNumber := new(big.Int)
	blockNumber, ok := blockNumber.SetString(blockNumberStr, 10)

	if !ok {
		return nil, fmt.Errorf("failed to parse block number %s", blockNumberStr)
	}

	return blockNumber, nil
}

func (g *postgresGateway) UpdateLastIndexedBlock(ctx context.Context, blockNumber *big.Int) error {
	return g.queries.UpdateLastIndexedBlock(ctx, blockNumber.String())
}
