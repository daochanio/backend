package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/controllers/distribute"
	"github.com/daochanio/backend/distributor/controllers/subscribe"
	"github.com/daochanio/backend/distributor/settings"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	if err := container.Provide(context.Background); err != nil {
		panic(err)
	}
	if err := container.Provide(appName); err != nil {
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

func appName() string {
	return "distributor"
}

func startDistributor(ctx context.Context, distributor distribute.Distributor, logger common.Logger) {
	go func() {
		if err := distributor.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("distributor stopped")
			panic(err)
		}
	}()
}

func startMessageSubscriber(ctx context.Context, subscriber subscribe.Subscriber, logger common.Logger) {
	go func() {
		if err := subscriber.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("subscriber stopped")
			panic(err)
		}
	}()
}

func awaitSigterm(ctx context.Context, logger common.Logger) {
	logger.Info(ctx).Msg("awaiting sigterm")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan

	logger.Info(ctx).Msgf("received signal %v", sig)
}
