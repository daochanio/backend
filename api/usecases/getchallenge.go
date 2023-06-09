package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

func NewGetChallengeUseCase(database Database) *GetChallenge {
	return &GetChallenge{
		database,
	}
}

type GetChallenge struct {
	database Database
}

type GetChallengeInput struct {
	Address string `validate:"eth_addr"`
}

func (u *GetChallenge) Execute(ctx context.Context, input *GetChallengeInput) (entities.Challenge, error) {
	if err := common.ValidateStruct(input); err != nil {
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
