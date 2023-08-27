package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daochanio/backend/api/http"
	"github.com/daochanio/backend/api/subscribe"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/gateways"
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
	commonSettings common.Settings,
	settings Settings,
	logger common.Logger,
	httpServer http.HttpServer,
	subscriber subscribe.Subscriber,
	database gateways.Database,
	cache gateways.Cache,
	stream gateways.Stream,
	blockchain gateways.Blockchain,
	images gateways.Images,
) {
	database.Start(ctx, settings.DatabaseConfig())
	cache.Start(ctx, settings.CacheConfig())
	stream.Start(ctx, settings.StreamConfig())
	blockchain.Start(ctx, settings.BlockchainConfig())
	images.Start(ctx, settings.ImagesConfig())

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		httpServer.Start(ctx, settings.HttpConfig())
	}()

	go func() {
		defer wg.Done()
		subscriber.Start(ctx, settings.SubscriberConfig())
	}()

	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	shutdownCtx := context.Background()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx).Err(err).Msg("failed to shutdown http server")
	}

	// the ctx being marked as done should cause the subscriber to return from its blocking call
	wg.Wait()

	subscriber.Shutdown(shutdownCtx)

	database.Shutdown(shutdownCtx)
	cache.Shutdown(shutdownCtx)
	stream.Shutdown(shutdownCtx)
	blockchain.Shutdown(shutdownCtx)
	images.Shutdown(shutdownCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
