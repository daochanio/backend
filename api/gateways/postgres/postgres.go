package postgres

import (
	"context"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	settings settings.Settings
	logger   common.Logger
	db       *pgxpool.Pool
	queries  *bindings.Queries
}

func NewDatabaseGateway(ctx context.Context, settings settings.Settings, logger common.Logger) usecases.Database {
	db, err := pgxpool.NewWithConfig(ctx, settings.PostgresConfig())
	if err != nil {
		panic(err)
	}
	queries := bindings.New(db)
	return &postgresGateway{
		settings,
		logger,
		db,
		queries,
	}
}
