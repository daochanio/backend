package gateway

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/gateways"
	"github.com/daochanio/backend/postgres/bindings"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	logger  common.Logger
	db      *pgxpool.Pool
	queries *bindings.Queries
}

func NewDatabaseGateway(ctx context.Context, logger common.Logger) gateways.Database {
	return &postgresGateway{
		logger:  logger,
		db:      nil,
		queries: nil,
	}
}

func (p *postgresGateway) Start(ctx context.Context, config gateways.DatabaseConfig) {
	p.logger.Info(ctx).Msg("starting postgres database")

	poolConfig, err := pgxpool.ParseConfig(config.ConnectionString)

	if err != nil {
		panic(err)
	}

	poolConfig.MinConns = config.MinConnections
	poolConfig.MaxConns = config.MaxConnections

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		panic(err)
	}

	p.db = db
	p.queries = bindings.New(db)
}

func (p *postgresGateway) Shutdown(ctx context.Context) {
	p.logger.Info(ctx).Msg("shutting down postgres database")
	p.db.Close()
}
