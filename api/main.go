package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/daochanio/backend/api/controllers/http"
	"github.com/daochanio/backend/api/controllers/subscribe"
	"github.com/daochanio/backend/api/gateways/ethereum"
	"github.com/daochanio/backend/api/gateways/postgres"
	"github.com/daochanio/backend/api/gateways/redis"
	"github.com/daochanio/backend/api/gateways/s3"
	"github.com/daochanio/backend/api/gateways/worker"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"go.uber.org/dig"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container := newContainer(ctx)

	// start the http controller inside a go routine
	if err := container.Invoke(startHttpServer); err != nil {
		panic(err)
	}

	// start the message subscriber inside a go routine
	if err := container.Invoke(startSubscriber); err != nil {
		panic(err)
	}

	// blocking call in main go routine to await sigterm
	if err := container.Invoke(awaitSigterm); err != nil {
		panic(err)
	}
}

func newContainer(ctx context.Context) *dig.Container {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(func() string {
		return "api"
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewCommonSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewLogger); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewHttpClient); err != nil {
		panic(err)
	}
	if err := container.Provide(settings.NewSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(postgres.NewDatabaseGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(redis.NewCacheGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(redis.NewStreamGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(s3.NewStorageGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(worker.NewSafeProxyGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(ethereum.NewBlockchainGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewVerifyRateLimitUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetChallengeUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewAuthenticateUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewSigninUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetThreadUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetThreadsUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewCreateThreadUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewDeleteThreadUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewCreateVoteUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetCommentsUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewCreateCommentUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewDeleteCommentUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewUploadImageUsecase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewAggregateVotesUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewHydrateUsersUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetUserUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(http.NewHttpServer); err != nil {
		panic(err)
	}
	if err := container.Provide(subscribe.NewSubscriber); err != nil {
		panic(err)
	}
	return container
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
