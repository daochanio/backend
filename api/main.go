package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daochanio/backend/api/controllers/http"
	"github.com/daochanio/backend/api/controllers/subscribe"
	"github.com/daochanio/backend/api/usecases"
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
	httpServer http.HttpServer,
	subscriber subscribe.Subscriber,
	database usecases.Database,
	cache usecases.Cache,
	stream usecases.Stream,
	blockchain usecases.Blockchain,
	storage usecases.Storage,
	safeProxy usecases.SafeProxy,
) {
	database.Start(ctx)
	cache.Start(ctx)
	stream.Start(ctx)
	blockchain.Start(ctx)
	storage.Start(ctx)
	safeProxy.Start(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		httpServer.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		subscriber.Start(ctx)
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
	storage.Shutdown(shutdownCtx)
	safeProxy.Shutdown(shutdownCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
