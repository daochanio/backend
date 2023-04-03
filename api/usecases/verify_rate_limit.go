package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

func NewVerifyRateLimitUseCase(
	logger common.ILogger,
	cacheGateway gateways.ICacheGateway,
) *VerifyRateLimitUseCase {
	return &VerifyRateLimitUseCase{
		logger,
		cacheGateway,
	}
}

type VerifyRateLimitUseCase struct {
	logger       common.ILogger
	cacheGateway gateways.ICacheGateway
}

type VerifyRateLimitInput struct {
	IpAddress string `validate:"required,ip"`
}

func (u *VerifyRateLimitUseCase) Execute(ctx context.Context, input *VerifyRateLimitInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	if err := u.cacheGateway.VerifyRateLimit(ctx, input.IpAddress); err != nil {
		return err
	}

	return nil
}
