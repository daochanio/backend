package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type GetUser struct {
	logger   common.Logger
	database gateways.Database
}

func NewGetUserUseCase(logger common.Logger, database gateways.Database) *GetUser {
	return &GetUser{
		logger,
		database,
	}
}

type GetUserInput struct {
	Address string `validate:"eth_addr"`
}

func (g *GetUser) Execute(ctx context.Context, input GetUserInput) (entities.User, error) {
	user, err := g.database.GetUserByAddress(ctx, input.Address)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
