package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/daochanio/backend/api/gateways/ethereum"
	"github.com/daochanio/backend/api/gateways/pg"
	"github.com/daochanio/backend/api/gateways/redis"
	"github.com/daochanio/backend/api/http"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(context.Background)
	container.Provide(appName)
	container.Provide(common.NewCommonSettings)
	container.Provide(common.NewLogger)
	container.Provide(settings.NewSettings)
	container.Provide(pg.NewPostgresGateway)
	container.Provide(redis.NewRedisGateway)
	container.Provide(ethereum.NewEthereumGateway)
	container.Provide(usecases.NewCreateUserUseCase)
	container.Provide(usecases.NewVerifyRateLimitUseCase)
	container.Provide(usecases.NewVerifyChallengeUseCase)
	container.Provide(usecases.NewGetChallengeUseCase)
	container.Provide(usecases.NewGetThreadUseCase)
	container.Provide(usecases.NewGetThreadsUseCase)
	container.Provide(usecases.NewCreateThreadUseCase)
	container.Provide(usecases.NewDeleteThreadUseCase)
	container.Provide(usecases.NewCreateThreadVoteUseCase)
	container.Provide(usecases.NewGetCommentsUseCase)
	container.Provide(usecases.NewCreateCommentUseCase)
	container.Provide(usecases.NewDeleteCommentUseCase)
	container.Provide(usecases.NewCreateCommentVoteUseCase)
	container.Provide(http.NewHttpServer)

	// start the http controller inside a go routine
	if err := container.Invoke(startHttpServer); err != nil {
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

func startHttpServer(ctx context.Context, httpServer http.IHttpServer) {
	go func() {
		err := httpServer.Start(ctx)
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
