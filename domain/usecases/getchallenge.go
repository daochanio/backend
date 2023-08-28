package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

func NewGetChallengeUseCase(validator common.Validator, database gateways.Database) *GetChallenge {
	return &GetChallenge{
		validator,
		database,
	}
}

type GetChallenge struct {
	validator common.Validator
	database  gateways.Database
}

type GetChallengeInput struct {
	Address string `validate:"eth_addr"`
}

func (u *GetChallenge) Execute(ctx context.Context, input *GetChallengeInput) (entities.Challenge, error) {
	if err := u.validator.ValidateStruct(input); err != nil {
		return entities.Challenge{}, err
	}

	challenge, err := u.database.GetChallengeByAddress(ctx, input.Address)

	if err == nil {
		return challenge, nil
	}

	newChallenge := entities.GenerateChallenge(input.Address)

	err = u.database.SaveChallenge(ctx, newChallenge)

	return newChallenge, err
}
