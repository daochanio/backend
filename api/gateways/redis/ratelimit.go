package redis

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

// this rate limit implementation is based on the leaky bucket algorithm
// it is more memory efficient than the fixed/sliding window algorithm
// but it is not as precise and may surprise people when first configuring the rate/period
// as it can allow up to 2x the rate limit in a single period
func (r *redisGateway) VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error {
	t1 := time.Now()
	res, err := r.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   rate,
		Period: period,
		Burst:  rate,
	})
	r.logger.Info(ctx).Msgf("rate limit duration %v", time.Since(t1))
	numCPUs := runtime.NumCPU()
	r.logger.Info(ctx).Msgf("num cpus %v", numCPUs)
	stats := r.client.PoolStats()
	r.logger.Info(ctx).Msgf("redis pool stats after: %+v", stats)

	if err != nil {
		return fmt.Errorf("rate limit error %w for key %v", err, key)
	}

	if res.Allowed == 0 {
		return fmt.Errorf("rate limit %v exceeded for %v", res.Limit, key)
	}

	return nil
}
