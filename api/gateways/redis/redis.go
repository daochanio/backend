package redis

import (
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type redisGateway struct {
	settings settings.Settings
	logger   common.Logger
	client   *redis.Client
	limiter  *redis_rate.Limiter
}

func NewRedisGateway(settings settings.Settings, logger common.Logger) gateways.CacheGateway {
	opt, err := redis.ParseURL(settings.CacheConnectionString())

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	limiter := redis_rate.NewLimiter(client)

	return &redisGateway{
		settings: settings,
		logger:   logger,
		client:   client,
		limiter:  limiter,
	}
}
