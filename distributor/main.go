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
	container.Provide(context.Background)
	container.Provide(appName)
	container.Provide(common.NewCommonSettings)
	container.Provide(common.NewLogger)
	container.Provide(settings.NewSettings)
	container.Provide(worker.NewWorker)

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

func startWorker(ctx context.Context, worker *worker.Worker, logger common.ILogger) {
	go func() {
		// blocking call to start the worker
		if err := worker.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("worker stopped")
			panic(err)
		}
	}()
}

func awaitSigterm(ctx context.Context, logger common.ILogger) {
	logger.Info(ctx).Msg("awaiting sigterm")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan

	logger.Info(ctx).Msgf("received signal %v", sig)
}
