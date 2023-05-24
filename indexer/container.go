package main

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/controllers/index"
	"github.com/daochanio/backend/indexer/gateways/ethereum"
	"github.com/daochanio/backend/indexer/gateways/postgres"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
	"go.uber.org/dig"
)

func newContainer(ctx context.Context) *dig.Container {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(common.NewSettings); err != nil {
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

	return container
}
