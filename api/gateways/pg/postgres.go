package pg

import (
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresGateway struct {
	settings settings.Settings
	logger   common.Logger
	queries  *bindings.Queries
}

func NewDatabaseGateway(settings settings.Settings, logger common.Logger, db *pgxpool.Pool) usecases.DatabaseGateway {
	queries := bindings.New(db)
	return &postgresGateway{
		settings,
		logger,
		queries,
	}
}
