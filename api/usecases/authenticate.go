package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type Authenticate struct {
	cache Cache
}

func NewAuthenticateUseCase(cache Cache) *Authenticate {
	return &Authenticate{
		cache,
	}
}

type AuthenticateInput struct {
	Address string `validate:"eth_addr"`
	SigHex  string `validate:"hexadecimal,min=1"`
}

func (u *Authenticate) Execute(ctx context.Context, input *AuthenticateInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	challenge, err := u.cache.GetChallengeByAddress(ctx, input.Address)

	if err != nil {
		return err
	}

	return challenge.Verify(input.SigHex)
}
