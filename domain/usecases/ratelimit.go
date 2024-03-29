package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
)

func NewVerifyRateLimitUseCase(
	logger common.Logger,
	validator common.Validator,
	cache gateways.Cache,
) *RateLimit {
	return &RateLimit{
		logger,
		validator,
		cache,
	}
}

type RateLimit struct {
	logger    common.Logger
	validator common.Validator
	cache     gateways.Cache
}

type RateLimitInput struct {
	Namespace string        `validate:"min=1,max=100"`
	IpAddress string        `validate:"ip"`
	Rate      int           `validate:"gt=0"`
	Period    time.Duration `validate:"gt=0"`
}

func (u *RateLimit) Execute(ctx context.Context, input *RateLimitInput) error {
	if err := u.validator.ValidateStruct(input); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", input.Namespace, input.IpAddress)

	if err := u.cache.VerifyRateLimit(ctx, key, input.Rate, input.Period); err != nil {
		return err
	}

	return nil
}
