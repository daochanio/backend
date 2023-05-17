package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/api/controllers/http"
	"github.com/daochanio/backend/api/controllers/subscribe"
	"github.com/daochanio/backend/api/gateways/ethereum"
	"github.com/daochanio/backend/api/gateways/postgres"
	"github.com/daochanio/backend/api/gateways/redis"
	"github.com/daochanio/backend/api/gateways/s3"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"go.uber.org/dig"
)

func main() {
	container := newContainer()

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

func newContainer() *dig.Container {
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
	if err := container.Provide(postgres.NewDatabaseGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(redis.NewCacheGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(redis.NewStreamGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(s3.NewImageGateway); err != nil {
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
	if err := container.Provide(usecases.NewHydrateUserUseCase); err != nil {
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

func appName() string {
	return "api"
}

func startHttpServer(ctx context.Context, httpServer http.HttpServer) {
	go func() {
		err := httpServer.Start(ctx)
		panic(err)
	}()
}

func startSubscriber(ctx context.Context, subscriber subscribe.Subscriber) {
	go func() {
		err := subscriber.Start(ctx)
		panic(err)
	}()
}

func awaitSigterm(ctx context.Context, logger common.Logger) {
	logger.Info(ctx).Msg("awaiting sigterm")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan

	logger.Info(ctx).Msgf("received sigterm %v", sig)
}
