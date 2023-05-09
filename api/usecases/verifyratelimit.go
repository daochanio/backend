package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/common"
)

func NewVerifyRateLimitUseCase(
	logger common.Logger,
	cacheGateway CacheGateway,
) *VerifyRateLimitUseCase {
	return &VerifyRateLimitUseCase{
		logger,
		cacheGateway,
	}
}

type VerifyRateLimitUseCase struct {
	logger       common.Logger
	cacheGateway CacheGateway
}

type VerifyRateLimitInput struct {
	Namespace string        `validate:"min=1,max=100"`
	IpAddress string        `validate:"ip"`
	Rate      int           `validate:"gt=0"`
	Period    time.Duration `validate:"gt=0"`
}

func (u *VerifyRateLimitUseCase) Execute(ctx context.Context, input *VerifyRateLimitInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", input.Namespace, input.IpAddress)

	if err := u.cacheGateway.VerifyRateLimit(ctx, key, input.Rate, input.Period); err != nil {
		return err
	}

	return nil
}
