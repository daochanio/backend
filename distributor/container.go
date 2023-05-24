package main

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/controllers/distribute"
	"github.com/daochanio/backend/distributor/controllers/subscribe"
	"github.com/daochanio/backend/distributor/settings"
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
	if err := container.Provide(distribute.NewDistributor); err != nil {
		panic(err)
	}
	if err := container.Provide(subscribe.NewSubscriber); err != nil {
		panic(err)
	}

	return container
}
