package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/db/bindings"
	"github.com/daochanio/backend/indexer/gateways"
	"github.com/daochanio/backend/indexer/settings"
	_ "github.com/lib/pq"
)

type postgresGateway struct {
	settings settings.Settings
	db       *sql.DB
	queries  *bindings.Queries
}

func NewPostgresGateway(settings settings.Settings) gateways.IDatabaseGateway {
	db, err := sql.Open("postgres", settings.DbConnectionString())
	if err != nil {
		panic(err)
	}

	queries := bindings.New(db)

	return &postgresGateway{
		settings,
		db,
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
