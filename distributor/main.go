package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

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
	go func() {
		distributor.Start(ctx)
	}()

	go func() {
		subscriber.Start(ctx)
	}()

	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	stopCtx := context.Background()

	// Allow distributor to finish if its currently in the middle of processing
	time.Sleep(time.Second * 10)

	distributor.Stop(stopCtx)

	subscriber.Stop(stopCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
