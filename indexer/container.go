package main

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/usecases"
	"github.com/daochanio/backend/ethereum"
	"github.com/daochanio/backend/indexer/index"
	postgres "github.com/daochanio/backend/postgres/gateway"
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
	if err := container.Provide(NewSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(postgres.NewDatabaseGateway); err != nil {
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
