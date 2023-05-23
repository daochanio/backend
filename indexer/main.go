package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/controllers/index"
	"github.com/daochanio/backend/indexer/settings"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := newContainer(ctx)

	if err := container.Invoke(start); err != nil {
		panic(err)
	}
}

func start(ctx context.Context, logger common.Logger, indexer index.Indexer, settings settings.Settings) {
	go func() {
		indexer.Start(ctx)
	}()

	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	// allow indexer to finish if its currently in the middle of indexing
	time.Sleep(time.Second * 10)

	stopCtx := context.Background()

	indexer.Stop(stopCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
