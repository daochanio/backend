package settings

import (
	"fmt"
	"os"
	"strconv"
)

type ISettings interface {
	Port() string
	DbConnectionString() string
	CacheAddress() string
	CachePassword() string
	CacheDb() int
	CacheUseTLS() bool
	BlockchainURI() string
}

type settings struct {
	port               string
	pgConnectionString string
	redisAddr          string
	redisDB            string
	redisPassword      string
	redisUseTLS        string
	blockchainURI      string
}

func NewSettings() ISettings {
	return &settings{
		port:               os.Getenv("PORT"),
		pgConnectionString: os.Getenv("PG_CONNECTION_STRING"),
		redisAddr:          os.Getenv("REDIS_ADDRESS"),
		redisPassword:      os.Getenv("REDIS_PASSWORD"),
		redisDB:            os.Getenv("REDIS_DB"),
		redisUseTLS:        os.Getenv("REDIS_USE_TLS"),
		blockchainURI:      os.Getenv("BLOCKCHAIN_URI"),
	}
}

func (s *settings) Port() string {
	return s.port
}

func (s *settings) DbConnectionString() string {
	return s.pgConnectionString
}

func (s *settings) CacheAddress() string {
	return s.redisAddr
}

func (s *settings) CacheDb() int {
	redisDB, err := strconv.Atoi(s.redisDB)

	if err != nil {
		panic(fmt.Errorf("invalid redis db value %v %w", s.redisDB, err))
	}

	return redisDB
}

func (s *settings) CachePassword() string {
	return s.redisPassword
}

func (s *settings) CacheUseTLS() bool {
	useTLS, err := strconv.ParseBool(s.redisUseTLS)

	if err != nil {
		panic(fmt.Errorf("invalid redis use tls value %v %w", s.redisUseTLS, err))
	}

	return useTLS
}

func (s *settings) BlockchainURI() string {
	return s.blockchainURI
}
