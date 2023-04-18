package settings

import (
	"os"
	"strconv"
	"time"
)

type Settings interface {
	Interval() time.Duration
}

type settings struct {
	interval time.Duration
}

func NewSettings() Settings {
	intervalMinutes, err := strconv.Atoi(os.Getenv("INTERVAL_MINUTES"))
	if err != nil {
		panic(err)
	}
	interval := time.Duration(intervalMinutes) * time.Minute

	return &settings{
		interval,
	}
}

func (s *settings) Interval() time.Duration {
	return s.interval
}
