package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
	"github.com/daochanio/backend/gateways/postgres/bindings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
	p.logger.Info(ctx).Msg("starting postgres gateway")

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
	p.logger.Info(ctx).Msg("shutting down postgres gateway")
	p.db.Close()
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (p *postgresGateway) Migrate(ctx context.Context, config gateways.DatabaseConfig) error {
	db, err := goose.OpenDBWithDriver("pgx", config.ConnectionString)

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	return nil
}

// Rollback returns an err but its idiomatic to call in a defer so we don't
// have the opportunity to check the error when defering Rollback directly.
func (p *postgresGateway) rollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); !errors.Is(err, pgx.ErrTxClosed) {
		p.logger.Error(ctx).Err(err).Msg("error rolling back transaction")
	}
}

func numericToBigInt(num pgtype.Numeric) *big.Int {
	return new(big.Int).Mul(num.Int, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(num.Exp)), nil))
}
