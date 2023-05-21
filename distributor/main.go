package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/controllers/distribute"
	"github.com/daochanio/backend/distributor/controllers/subscribe"
	"github.com/daochanio/backend/distributor/settings"
	"go.uber.org/dig"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(func() string {
		return "distributor"
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewCommonSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewLogger); err != nil {
		panic(err)
	}
	if err := container.Provide(settings.NewSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(distribute.NewDistributor); err != nil {
		panic(err)
	}
	if err := container.Provide(subscribe.NewSubscriber); err != nil {
		panic(err)
	}

	// start the app in a go routine
	if err := container.Invoke(startDistributor); err != nil {
		panic(err)
	}

	if err := container.Invoke(startMessageSubscriber); err != nil {
		panic(err)
	}

	// blocking call in main go routine to await sigterm
	if err := container.Invoke(awaitSigterm); err != nil {
		panic(err)
	}
}

func startDistributor(ctx context.Context, distributor distribute.Distributor, logger common.Logger) {
	go func() {
		distributor.Start(ctx)
	}()
}

func startMessageSubscriber(ctx context.Context, subscriber subscribe.Subscriber, logger common.Logger) {
	go func() {
		subscriber.Start(ctx)
	}()
}

func awaitSigterm(ctx context.Context, logger common.Logger, distributor distribute.Distributor, subscriber subscribe.Subscriber) {
	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal, beginning shutdown")

	shutdownCtx := context.Background()

	// Allow distributor to finish if its currently in the middle of processing
	time.Sleep(time.Second * 10)

	distributor.Stop(shutdownCtx)

	subscriber.Stop(shutdownCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
