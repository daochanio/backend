package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
)

type VerifyChallengeUseCase struct {
	cacheGateway CacheGateway
}

func NewVerifyChallengeUseCase(cacheGateway CacheGateway) *VerifyChallengeUseCase {
	return &VerifyChallengeUseCase{
		cacheGateway,
	}
}

type VerifyChallengeInput struct {
	Address string `validate:"eth_addr"`
	SigHex  string `validate:"hexadecimal,min=1"`
}

func (u *VerifyChallengeUseCase) Execute(ctx context.Context, input *VerifyChallengeInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	challenge, err := u.cacheGateway.GetChallengeByAddress(ctx, input.Address)

	if err != nil {
		return err
	}

	return challenge.Verify(input.SigHex)
}
