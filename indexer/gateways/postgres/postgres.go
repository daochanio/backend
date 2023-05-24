package postgres

import (
	"context"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	logger   common.Logger
	settings settings.Settings
	db       *pgxpool.Pool
	queries  *bindings.Queries
}

func NewPostgresGateway(ctx context.Context, settings settings.Settings, logger common.Logger) usecases.Database {
	return &postgresGateway{
		logger:   logger,
		settings: settings,
		db:       nil,
		queries:  nil,
	}
}

func (g *postgresGateway) Start(ctx context.Context) {
	g.logger.Info(ctx).Msg("starting postgres database")
	db, err := pgxpool.NewWithConfig(ctx, g.settings.PostgresConfig())
	if err != nil {
		panic(err)
	}
	g.db = db
	g.queries = bindings.New(db)
}

func (g *postgresGateway) Shutdown(ctx context.Context) {
	g.logger.Info(ctx).Msg("shutting down postgres database")
	g.db.Close()
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
