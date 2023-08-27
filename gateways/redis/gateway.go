package redis

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type redisCacheGateway struct {
	logger  common.Logger
	client  *redis.Client
	limiter *redis_rate.Limiter
}

func NewCacheGateway(logger common.Logger) gateways.Cache {
	return &redisCacheGateway{
		logger:  logger,
		client:  nil,
		limiter: nil,
	}
}

func (r *redisCacheGateway) Start(ctx context.Context, config gateways.CacheConfig) {
	r.logger.Info(ctx).Msg("starting redis cache gateway")

	opt, err := redis.ParseURL(config.ConnectionString)

	if err != nil {
		panic(err)
	}

	opt.DialTimeout = config.DialTimeout
	opt.MinIdleConns = config.MinIdleConns
	opt.PoolSize = config.PoolSize
	opt.ReadTimeout = config.ReadTimeout
	opt.WriteTimeout = config.WriteTimeout

	r.client = redis.NewClient(opt)

	r.limiter = redis_rate.NewLimiter(r.client)
}

func (r *redisCacheGateway) Shutdown(ctx context.Context) {
	r.logger.Info(ctx).Msg("shutting down redis cache gateway")

	if err := r.client.Close(); err != nil {
		r.logger.Error(ctx).Err(err).Msg("error closing redis cache client")
	}
}

type redisStreamGateway struct {
	logger common.Logger
	client *redis.Client
}

func NewStreamGateway(logger common.Logger) gateways.Stream {
	return &redisStreamGateway{
		logger: logger,
		client: nil,
	}
}

func (r *redisStreamGateway) Start(ctx context.Context, config gateways.StreamConfig) {
	r.logger.Info(ctx).Msg("starting redis stream gateway")

	opt, err := redis.ParseURL(config.ConnectionString)

	if err != nil {
		panic(err)
	}

	opt.DialTimeout = config.DialTimeout
	opt.MinIdleConns = config.MinIdleConns
	opt.PoolSize = config.PoolSize
	opt.ReadTimeout = config.ReadTimeout
	opt.WriteTimeout = config.WriteTimeout

	r.client = redis.NewClient(opt)
}

func (r *redisStreamGateway) Shutdown(ctx context.Context) {
	r.logger.Info(ctx).Msg("shutting down redis stream gateway")

	if err := r.client.Close(); err != nil {
		r.logger.Error(ctx).Err(err).Msg("error closing redis stream client")
	}
}
