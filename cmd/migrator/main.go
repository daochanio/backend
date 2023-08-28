package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/usecases"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := newContainer(ctx)

	if err := container.Invoke(start); err != nil {
		panic(err)
	}

}

func start(
	ctx context.Context,
	settings Settings,
	logger common.Logger,
	migrateDatabase *usecases.MigrateDatabase,
) {
	logger.Start(ctx, settings.LoggerConfig())

	logger.Info(ctx).Msg("starting migrator")

	databaseConfig := settings.DatabaseConfig()

	if err := migrateDatabase.Execute(ctx, databaseConfig); err != nil {
		logger.Error(ctx).Err(err).Msg("error migrating database")
		panic(err)
	}

	logger.Info(ctx).Msg("migrated database")
}
