package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type SigninUseCase struct {
	logger          common.Logger
	databaseGateway DatabaseGateway
	cacheGateway    CacheGateway
	messageGateway  MessageGateway
}

type SigninUseCaseInput struct {
	Address string `validate:"eth_addr"`
}

func NewSigninUseCase(logger common.Logger, databaseGateway DatabaseGateway, cacheGateway CacheGateway, messageGateway MessageGateway) *SigninUseCase {
	return &SigninUseCase{
		logger,
		databaseGateway,
		cacheGateway,
		messageGateway,
	}
}

func (u *SigninUseCase) Execute(ctx context.Context, input SigninUseCaseInput) (entities.Challenge, error) {
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

func (u *SigninUseCase) getChallenge(ctx context.Context, address string) (entities.Challenge, error) {
	challenge, err := u.cacheGateway.GetChallengeByAddress(ctx, address)

	if err == nil {
		return challenge, nil
	}

	newChallenge := entities.GenerateChallenge(address)

	err = u.cacheGateway.SaveChallenge(ctx, newChallenge)

	return newChallenge, err
}

func (u *SigninUseCase) upsertUser(ctx context.Context, address string) error {
	err := u.databaseGateway.UpsertUser(ctx, address)

	if err != nil {
		return err
	}

	return u.messageGateway.PublishSignin(ctx, address)
}
