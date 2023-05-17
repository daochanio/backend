package redis

import (
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
	client := redis.NewClient(settings.RegionalRedisOptions())
	limiter := redis_rate.NewLimiter(client)
	return &redisCacheGateway{
		settings: settings,
		logger:   logger,
		client:   client,
		limiter:  limiter,
	}
}

type redisStreamGateway struct {
	settings settings.Settings
	logger   common.Logger
	client   *redis.Client
}

func NewStreamGateway(settings settings.Settings, logger common.Logger) usecases.Stream {
	client := redis.NewClient(settings.GlobalRedisOptions())
	return &redisStreamGateway{
		settings: settings,
		logger:   logger,
		client:   client,
	}
}
