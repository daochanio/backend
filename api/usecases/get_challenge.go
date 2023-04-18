package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

func NewGetChallengeUseCase(cacheGateway gateways.CacheGateway) *GetChallengeUseCase {
	return &GetChallengeUseCase{
		cacheGateway,
	}
}

type GetChallengeUseCase struct {
	cacheGateway gateways.CacheGateway
}

type GetChallengeInput struct {
	Address string `validate:"eth_addr"`
}

func (u *GetChallengeUseCase) Execute(ctx context.Context, input *GetChallengeInput) (entities.Challenge, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Challenge{}, err
	}

	challenge, err := u.cacheGateway.GetChallengeByAddress(ctx, input.Address)

	if err == nil {
		return challenge, nil
	}

	newChallenge := entities.GenerateChallenge(input.Address)

	err = u.cacheGateway.SaveChallenge(ctx, newChallenge)

	return newChallenge, err
}
