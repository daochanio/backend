package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/gateways"
	"github.com/daochanio/backend/indexer/index"
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
	logger common.Logger,
	commonSettings common.Settings,
	indexer index.Indexer,
	settings Settings,
	database gateways.Database,
	blockchain gateways.Blockchain,
) {
	database.Start(ctx, settings.DatabaseConfig())
	blockchain.Start(ctx, settings.BlockchainConfig())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		indexer.Start(ctx, settings.IndexerConfig())
	}()

	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	wg.Wait()

	shutdownCtx := context.Background()

	indexer.Shutdown(shutdownCtx)

	database.Shutdown(shutdownCtx)
	blockchain.Shutdown(shutdownCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
