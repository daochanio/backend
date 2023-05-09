package redis

import (
	"fmt"
	"strings"

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

func NewGateway(settings settings.Settings, logger common.Logger, client *redis.Client) *redisGateway {
	limiter := redis_rate.NewLimiter(client)
	return &redisGateway{
		settings: settings,
		logger:   logger,
		client:   client,
		limiter:  limiter,
	}
}

func getFullKey(namespace string, keys ...string) string {
	return fmt.Sprintf("%v:%v", namespace, strings.Join(keys, ":"))
}
