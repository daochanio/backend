package main

import (
	"os"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
	"github.com/joho/godotenv"
)

type Settings interface {
	LoggerConfig() common.LoggerConfig
	DatabaseConfig() gateways.DatabaseConfig
}

type settings struct {
	env                string
	appname            string
	hostname           string
	pgConnectionString string
}

func NewSettings() Settings {
	_ = godotenv.Load(".env/.env.migrator.dev")

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	return &settings{
		env:                os.Getenv("ENV"),
		appname:            os.Getenv("APP_NAME"),
		hostname:           hostname,
		pgConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}

func (s *settings) LoggerConfig() common.LoggerConfig {
	return common.LoggerConfig{
		Env:      s.env,
		Appname:  s.appname,
		Hostname: s.hostname,
	}
}

func (s *settings) DatabaseConfig() gateways.DatabaseConfig {
	return gateways.DatabaseConfig{
		ConnectionString: s.pgConnectionString,
		MinConnections:   1,
		MaxConnections:   1,
	}
}
