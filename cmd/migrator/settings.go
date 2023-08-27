package main

import (
	"os"

	"github.com/daochanio/backend/domain/gateways"
)

type Settings interface {
	DatabaseConfig() gateways.DatabaseConfig
}

type settings struct {
	pgConnectionString string
}

func NewSettings() Settings {
	return &settings{
		pgConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}
func (s *settings) DatabaseConfig() gateways.DatabaseConfig {
	return gateways.DatabaseConfig{
		ConnectionString: s.pgConnectionString,
		MinConnections:   1,
		MaxConnections:   1,
	}
}
