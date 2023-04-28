package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

// this rate limit implementation is based on the leaky bucket algorithm
// it is more memory efficient than the fixed/sliding window algorithm
// but it is not as precise and may surprise people when first configuring the rate/period
// as it can allow up to 2x the rate limit in a single period
func (r *redisGateway) VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error {
	res, err := r.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   rate,
		Period: period,
		Burst:  rate,
	})

	if err != nil {
		return fmt.Errorf("rate limit error %w for key %v", err, key)
	}

	if res.Allowed == 0 {
		return fmt.Errorf("rate limit %v exceeded for %v", res.Limit, key)
	}

	return nil
}

func getFullKey(namespace string, keys ...string) string {
	return fmt.Sprintf("%v:%v", namespace, strings.Join(keys, ":"))
}
