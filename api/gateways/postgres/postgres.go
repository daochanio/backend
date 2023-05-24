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
	return &postgresGateway{
		settings: settings,
		logger:   logger,
		db:       nil,
		queries:  nil,
	}
}

func (p *postgresGateway) Start(ctx context.Context) {
	p.logger.Info(ctx).Msg("starting postgres database")
	db, err := pgxpool.NewWithConfig(ctx, p.settings.PostgresConfig())
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
