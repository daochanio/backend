package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daochanio/backend/cmd/distributor/distribute"
	"github.com/daochanio/backend/cmd/distributor/subscribe"
	"github.com/daochanio/backend/common"
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
	settings Settings,
	distributor distribute.Distributor,
	subscriber subscribe.Subscriber,
) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		distributor.Start(ctx, settings.DistributorConfig())
	}()

	go func() {
		defer wg.Done()
		subscriber.Start(ctx, settings.SubscribeConfig())
	}()

	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	wg.Wait()

	shutdownCtx := context.Background()

	distributor.Shutdown(shutdownCtx)

	subscriber.Shutdown(shutdownCtx)

	logger.Info(shutdownCtx).Msgf("shutdown complete")
}
