package pg

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	settings settings.Settings
	queries  *bindings.Queries
}

func NewPostgresGateway(ctx context.Context, settings settings.Settings) gateways.DatabaseGateway {
	config, err := pgxpool.ParseConfig(settings.DbConnectionString())

	if err != nil {
		panic(err)
	}

	config.MinConns = 10
	config.MaxConns = 1000

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	queries := bindings.New(db)

	return &postgresGateway{
		settings,
		queries,
	}
}
