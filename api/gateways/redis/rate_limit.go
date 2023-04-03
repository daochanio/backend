package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis_rate/v10"
)

func (r *redisGateway) VerifyRateLimit(ctx context.Context, ipAddress string) error {
	res, err := r.limiter.Allow(ctx, ipAddress, redis_rate.PerSecond(10))

	if err != nil {
		return fmt.Errorf("rate limit error %w", err)
	}

	if res.Remaining == 0 {
		return fmt.Errorf("rate limit exceeded %v", res.Limit)
	}

	return nil
}

func getFullKey(namespace string, keys ...string) string {
	return fmt.Sprintf("%v-%v", namespace, strings.Join(keys, "-"))
}
