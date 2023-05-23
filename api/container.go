package main

import (
	"context"

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

func newContainer(ctx context.Context) *dig.Container {
	container := dig.New()
	provideGeneral(ctx, container)
	provideGateways(container)
	provideUseCases(container)
	provideControllers(container)
	return container
}

func provideGeneral(ctx context.Context, container *dig.Container) {
	if err := container.Provide(func() context.Context {
		return ctx
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
}

func provideGateways(container *dig.Container) {
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
}

func provideUseCases(container *dig.Container) {
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
}

func provideControllers(container *dig.Container) {
	if err := container.Provide(http.NewHttpServer); err != nil {
		panic(err)
	}
	if err := container.Provide(subscribe.NewSubscriber); err != nil {
		panic(err)
	}
}
