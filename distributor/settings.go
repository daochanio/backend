package main

import (
	"os"
	"strconv"
	"time"

	"github.com/daochanio/backend/distributor/distribute"
	"github.com/daochanio/backend/distributor/subscribe"
)

type Settings interface {
	DistributorConfig() distribute.DistributorConfig
	SubscribeConfig() subscribe.SubscriberConfig
}

type settings struct {
	interval              time.Duration
	redisConnectionString string
}

func NewSettings() Settings {
	intervalMinutes, err := strconv.Atoi(os.Getenv("INTERVAL_MINUTES"))
	if err != nil {
		panic(err)
	}
	interval := time.Duration(intervalMinutes) * time.Minute

	return &settings{
		interval:              interval,
		redisConnectionString: os.Getenv("REDIS_CONNECTION_STRING"),
	}
}

func (s *settings) DistributorConfig() distribute.DistributorConfig {
	return distribute.DistributorConfig{
		Interval: s.interval,
	}
}

func (s *settings) SubscribeConfig() subscribe.SubscriberConfig {
	return subscribe.SubscriberConfig{
		ConnectionString: s.redisConnectionString,
		DialTimeout:      10 * time.Second,
		MinIdleConns:     10,
		PoolSize:         100,
		ReadTimeout:      -1,
		WriteTimeout:     -1,
	}
}
