package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type AuthenticateUseCase struct {
	cacheGateway CacheGateway
}

func NewAuthenticateUseCase(cacheGateway CacheGateway) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		cacheGateway,
	}
}

type AuthenticateInput struct {
	Address string `validate:"eth_addr"`
	SigHex  string `validate:"hexadecimal,min=1"`
}

func (u *AuthenticateUseCase) Execute(ctx context.Context, input *AuthenticateInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	challenge, err := u.cacheGateway.GetChallengeByAddress(ctx, input.Address)

	if err != nil {
		return err
	}

	return challenge.Verify(input.SigHex)
}
