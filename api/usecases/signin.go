package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type Signin struct {
	logger   common.Logger
	database Database
	cache    Cache
	stream   Stream
}

type SigninInput struct {
	Address string `validate:"eth_addr"`
}

func NewSigninUseCase(logger common.Logger, database Database, cache Cache, stream Stream) *Signin {
	return &Signin{
		logger,
		database,
		cache,
		stream,
	}
}

func (u *Signin) Execute(ctx context.Context, input SigninInput) (entities.Challenge, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Challenge{}, err
	}

	challenge, err := u.getChallenge(ctx, input.Address)

	if err != nil {
		return entities.Challenge{}, err
	}

	err = u.upsertUser(ctx, input.Address)

	if err != nil {
		return entities.Challenge{}, err
	}

	return challenge, err
}

func (u *Signin) getChallenge(ctx context.Context, address string) (entities.Challenge, error) {
	challenge, err := u.cache.GetChallengeByAddress(ctx, address)

	if err == nil {
		return challenge, nil
	}

	newChallenge := entities.GenerateChallenge(address)

	err = u.cache.SaveChallenge(ctx, newChallenge)

	return newChallenge, err
}

func (u *Signin) upsertUser(ctx context.Context, address string) error {
	err := u.database.UpsertUser(ctx, address)

	if err != nil {
		return err
	}

	return u.stream.PublishSignin(ctx, address)
}
