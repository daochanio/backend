package settings

import (
	"os"
	"strconv"
	"time"
)

type Settings interface {
	Interval() time.Duration
	StreamConnectionString() string
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

func (s *settings) Interval() time.Duration {
	return s.interval
}

func (s *settings) StreamConnectionString() string {
	return s.redisConnectionString
}
