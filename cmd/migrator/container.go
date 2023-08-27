package main

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/usecases"
	"github.com/daochanio/backend/gateways/postgres"
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
	if err := container.Provide(usecases.NewMigrateDatabase); err != nil {
		panic(err)
	}

	return container
}
