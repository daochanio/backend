package settings

import (
	"os"
)

type Settings interface {
	Port() string
	DbConnectionString() string
	CacheConnectionString() string
	BlockchainURI() string
	RealIPHeader() string
}

type settings struct {
	port                  string
	pgConnectionString    string
	redisConnectionString string
	blockchainURI         string
	realIPHeader          string
}

func NewSettings() Settings {
	return &settings{
		port:                  os.Getenv("PORT"),
		pgConnectionString:    os.Getenv("PG_CONNECTION_STRING"),
		redisConnectionString: os.Getenv("REDIS_CONNECTION_STRING"),
		blockchainURI:         os.Getenv("BLOCKCHAIN_URI"),
		realIPHeader:          os.Getenv("REAL_IP_HEADER"),
	}
}

func (s *settings) Port() string {
	return s.port
}

func (s *settings) DbConnectionString() string {
	return s.pgConnectionString
}

func (s *settings) CacheConnectionString() string {
	return s.redisConnectionString
}

func (s *settings) BlockchainURI() string {
	return s.blockchainURI
}

func (s *settings) RealIPHeader() string {
	return s.realIPHeader
}
