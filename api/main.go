package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/api/controllers/http"
	"github.com/daochanio/backend/api/controllers/subscribe"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/gateways/cloudfront"
	"github.com/daochanio/backend/api/gateways/ethereum"
	"github.com/daochanio/backend/api/gateways/pg"
	"github.com/daochanio/backend/api/gateways/redis"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
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
	if err := container.Provide(pg.NewDatabaseGateway); err != nil {
		panic(err)
	}
	// redis provides multiple interface implementations
	if err := container.Provide(redis.NewGateway, dig.As(new(gateways.MessageGateway), new(gateways.CacheGateway))); err != nil {
		panic(err)
	}
	if err := container.Provide(cloudfront.NewImageGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(ethereum.NewBlockchainGateway); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewCreateUserUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewVerifyRateLimitUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewVerifyChallengeUseCase); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewGetChallengeUseCase); err != nil {
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
	if err := container.Provide(http.NewHttpServer); err != nil {
		panic(err)
	}
	if err := container.Provide(subscribe.NewSubscriber); err != nil {
		panic(err)
	}

	// start the http controller inside a go routine
	if err := container.Invoke(startHttpServer); err != nil {
		panic(err)
	}

	// start the message subscriber inside a go routine
	if err := container.Invoke(startMessageSubscriber); err != nil {
		panic(err)
	}

	// blocking call in main go routine to await sigterm
	if err := container.Invoke(awaitSigterm); err != nil {
		panic(err)
	}
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

func startMessageSubscriber(ctx context.Context, messageSubscriber subscribe.Subscriber) {
	go func() {
		err := messageSubscriber.Start(ctx)
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
