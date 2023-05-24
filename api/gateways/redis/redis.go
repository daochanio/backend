package redis

import (
	"context"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type redisCacheGateway struct {
	settings settings.Settings
	logger   common.Logger
	client   *redis.Client
	limiter  *redis_rate.Limiter
}

func NewCacheGateway(settings settings.Settings, logger common.Logger) usecases.Cache {
	return &redisCacheGateway{
		settings: settings,
		logger:   logger,
		client:   nil,
		limiter:  nil,
	}
}

func (r *redisCacheGateway) Start(ctx context.Context) {
	r.logger.Info(ctx).Msg("starting redis cache")
	r.client = redis.NewClient(r.settings.RegionalRedisOptions())
	r.limiter = redis_rate.NewLimiter(r.client)
}

func (r *redisCacheGateway) Shutdown(ctx context.Context) {
	r.logger.Info(ctx).Msg("shutting down redis cache")
	if err := r.client.Close(); err != nil {
		r.logger.Error(ctx).Err(err).Msg("error closing redis cache client")
	}
}

type redisStreamGateway struct {
	settings settings.Settings
	logger   common.Logger
	client   *redis.Client
}

func NewStreamGateway(settings settings.Settings, logger common.Logger) usecases.Stream {
	return &redisStreamGateway{
		settings: settings,
		logger:   logger,
		client:   nil,
	}
}

func (r *redisStreamGateway) Start(ctx context.Context) {
	r.logger.Info(ctx).Msg("starting redis stream")
	r.client = redis.NewClient(r.settings.GlobalRedisOptions())
}

func (r *redisStreamGateway) Shutdown(ctx context.Context) {
	r.logger.Info(ctx).Msg("shutting down redis stream")
	if err := r.client.Close(); err != nil {
		r.logger.Error(ctx).Err(err).Msg("error closing redis stream client")
	}
}
