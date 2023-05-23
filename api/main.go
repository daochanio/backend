package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/daochanio/backend/api/controllers/http"
	"github.com/daochanio/backend/api/controllers/subscribe"
	"github.com/daochanio/backend/common"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := newContainer(ctx, "api")

	if err := container.Invoke(startHttpServer); err != nil {
		panic(err)
	}

	if err := container.Invoke(startSubscriber); err != nil {
		panic(err)
	}

	if err := container.Invoke(awaitSigterm); err != nil {
		panic(err)
	}
}

func startHttpServer(ctx context.Context, httpServer http.HttpServer) {
	go func() {
		httpServer.Start(ctx)
	}()
}

func startSubscriber(ctx context.Context, subscriber subscribe.Subscriber) {
	go func() {
		subscriber.Start(ctx)
	}()
}

func awaitSigterm(ctx context.Context, logger common.Logger, httpServer http.HttpServer, subscriber subscribe.Subscriber) {
	logger.Info(ctx).Msg("awaiting kill signal")

	<-ctx.Done()

	logger.Info(ctx).Msgf("received kill signal")

	shutdownCtx := context.Background()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		logger.Error(ctx).Err(err).Msg("failed to shutdown http server")
	}

	// See https://github.com/redis/go-redis/issues/2276 and https://github.com/redis/go-redis/pull/2455
	// Blocking calls to redis client methods will not be interrupted by the shutdown context.
	// We need to wait before calling Stop() to ensure that the subscriber has finished processing its latest loop.
	// This way we ensure that no new messages will be written to the buffer after flushing inside the Stop() method.
	time.Sleep(10 * time.Second)

	subscriber.Stop(shutdownCtx)

	logger.Info(ctx).Msgf("shutdown complete")
}
