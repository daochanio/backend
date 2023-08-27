package main

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/usecases"
	"github.com/daochanio/backend/distributor/distribute"
	"github.com/daochanio/backend/distributor/subscribe"
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
	if err := container.Provide(usecases.NewDistribute); err != nil {
		panic(err)
	}
	if err := container.Provide(usecases.NewProcessVote); err != nil {
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
