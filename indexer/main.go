package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/controllers/index"
	"github.com/daochanio/backend/indexer/gateways/ethereum"
	"github.com/daochanio/backend/indexer/gateways/postgres"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
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
		return "indexer"
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
	if err := container.Provide(postgres.NewPostgresGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(ethereum.NewEthereumGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewIndexBlocksUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewIndexReputationUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(index.NewIndexer); err != nil {
		panic(err)
	}

	// start the app in a go routine
	if err := container.Invoke(startIndexer); err != nil {
		panic(err)
	}

	// blocking call in main go routine to await sigterm
	if err := container.Invoke(awaitSigterm); err != nil {
		panic(err)
	}
}

func startIndexer(ctx context.Context, indexer index.Indexer, logger common.Logger) {
	go func() {
		indexer.Start(ctx)
	}()
}

func awaitSigterm(ctx context.Context, logger common.Logger, indexer index.Indexer, settings settings.Settings) {
	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	// allow indexer to finish if its currently in the middle of indexing
	time.Sleep(time.Second * 10)

	stopCtx := context.Background()

	indexer.Stop(stopCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
