package gateways

import (
	"context"
	"time"
)

type CacheConfig struct {
	ConnectionString string
	DialTimeout      time.Duration
	MinIdleConns     int
	PoolSize         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
}

type Cache interface {
	Start(ctx context.Context, config CacheConfig)
	Shutdown(ctx context.Context)
	VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error
}
