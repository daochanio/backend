package pg

import (
	"database/sql"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/db/bindings"
	_ "github.com/lib/pq"
)

type postgresGateway struct {
	settings settings.Settings
	db       *sql.DB
	queries  *bindings.Queries
}

func NewPostgresGateway(settings settings.Settings) gateways.DatabaseGateway {
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
