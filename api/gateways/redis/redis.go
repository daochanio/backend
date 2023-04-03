package redis

import (
	"crypto/tls"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type redisGateway struct {
	settings settings.ISettings
	logger   common.ILogger
	client   *redis.Client
	limiter  *redis_rate.Limiter
}

func NewRedisGateway(settings settings.ISettings, logger common.ILogger) gateways.ICacheGateway {
	opt := &redis.Options{
		Addr:      settings.CacheAddress(),
		Password:  settings.CachePassword(),
		DB:        settings.CacheDb(),
		TLSConfig: nil,
	}

	if settings.CacheUseTLS() {
		opt.TLSConfig = &tls.Config{}
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
