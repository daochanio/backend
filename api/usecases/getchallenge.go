package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

func NewGetChallengeUseCase(cacheGateway Cache) *GetChallenge {
	return &GetChallenge{
		cacheGateway,
	}
}

type GetChallenge struct {
	cache Cache
}

type GetChallengeInput struct {
	Address string `validate:"eth_addr"`
}

func (u *GetChallenge) Execute(ctx context.Context, input *GetChallengeInput) (entities.Challenge, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Challenge{}, err
	}

	challenge, err := u.cache.GetChallengeByAddress(ctx, input.Address)

	if err == nil {
		return challenge, nil
	}

	newChallenge := entities.GenerateChallenge(input.Address)

	err = u.cache.SaveChallenge(ctx, newChallenge)

	return newChallenge, err
}
