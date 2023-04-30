package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
	"github.com/daochanio/backend/distributor/worker"
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
	if err := container.Provide(worker.NewWorker); err != nil {
		panic(err)
	}

	// start the app in a go routine
	if err := container.Invoke(startWorker); err != nil {
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

func startWorker(ctx context.Context, worker *worker.Worker, logger common.Logger) {
	go func() {
		// blocking call to start the worker
		if err := worker.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("worker stopped")
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
