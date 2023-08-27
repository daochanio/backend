package gateways

import (
	"context"
	"time"

	"github.com/daochanio/backend/core/entities"
)

type StreamConfig struct {
	ConnectionString string
	DialTimeout      time.Duration
	MinIdleConns     int
	PoolSize         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
}

type Stream interface {
	Start(ctx context.Context, config StreamConfig)
	Shutdown(ctx context.Context)
	PublishSignin(ctx context.Context, address string) error
	PublishVote(ctx context.Context, vote entities.Vote) error
}
