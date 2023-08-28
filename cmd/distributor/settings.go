package main

import (
	"os"
	"strconv"
	"time"

	"github.com/daochanio/backend/cmd/distributor/distribute"
	"github.com/daochanio/backend/cmd/distributor/subscribe"
	"github.com/daochanio/backend/common"
	"github.com/joho/godotenv"
)

type Settings interface {
	LoggerConfig() common.LoggerConfig
	DistributorConfig() distribute.DistributorConfig
	SubscribeConfig() subscribe.SubscriberConfig
}

type settings struct {
	env                   string
	appname               string
	hostname              string
	interval              time.Duration
	redisConnectionString string
}

func NewSettings() Settings {
	_ = godotenv.Load(".env/.env.distributor.dev")

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	intervalMinutes, err := strconv.Atoi(os.Getenv("INTERVAL_MINUTES"))
	if err != nil {
		panic(err)
	}
	interval := time.Duration(intervalMinutes) * time.Minute

	return &settings{
		env:                   os.Getenv("ENV"),
		appname:               os.Getenv("APP_NAME"),
		hostname:              hostname,
		interval:              interval,
		redisConnectionString: os.Getenv("REDIS_CONNECTION_STRING"),
	}
}

func (s *settings) LoggerConfig() common.LoggerConfig {
	return common.LoggerConfig{
		Env:      s.env,
		Appname:  s.appname,
		Hostname: s.hostname,
	}
}

func (s *settings) DistributorConfig() distribute.DistributorConfig {
	return distribute.DistributorConfig{
		Interval: s.interval,
	}
}

func (s *settings) SubscribeConfig() subscribe.SubscriberConfig {
	return subscribe.SubscriberConfig{
		Group:            s.appname,
		Consumer:         s.hostname,
		ConnectionString: s.redisConnectionString,
		DialTimeout:      10 * time.Second,
		MinIdleConns:     10,
		PoolSize:         100,
		ReadTimeout:      -1,
		WriteTimeout:     -1,
	}
}
