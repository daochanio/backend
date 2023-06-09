package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/controllers/distribute"
	"github.com/daochanio/backend/distributor/controllers/subscribe"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := newContainer(ctx)

	if err := container.Invoke(start); err != nil {
		panic(err)
	}
}

func start(ctx context.Context, logger common.Logger, distributor distribute.Distributor, subscriber subscribe.Subscriber) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		distributor.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		subscriber.Start(ctx)
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
