package main

import (
	"os"
	"time"

	"github.com/daochanio/backend/api/http"
	"github.com/daochanio/backend/api/subscribe"
	"github.com/daochanio/backend/core/gateways"
)

type Settings interface {
	HttpConfig() http.HttpConfig
	SubscriberConfig() subscribe.SubscriberConfig
	DatabaseConfig() gateways.DatabaseConfig
	StreamConfig() gateways.StreamConfig
	CacheConfig() gateways.CacheConfig
	BlockchainConfig() gateways.BlockchainConfig
	ImagesConfig() gateways.ImagesConfig
}

type settings struct {
	port                        string
	pgConnectionString          string
	redisCacheConnectionString  string
	redisStreamConnectionString string
	jwtSecret                   string
	blockchainURL               string
	realIPHeader                string
	imagesBaseUrl               string
	imagesAPIKey                string
}

func NewSettings() Settings {
	return &settings{
		port:                        os.Getenv("PORT"),
		pgConnectionString:          os.Getenv("PG_CONNECTION_STRING"),
		redisCacheConnectionString:  os.Getenv("REDIS_CACHE_CONNECTION_STRING"),
		redisStreamConnectionString: os.Getenv("REDIS_STREAM_CONNECTION_STRING"),
		jwtSecret:                   os.Getenv("JWT_SECRET"),
		blockchainURL:               os.Getenv("BLOCKCHAIN_URI"),
		realIPHeader:                os.Getenv("REAL_IP_HEADER"),
		imagesBaseUrl:               os.Getenv("IMAGES_BASE_URL"),
		imagesAPIKey:                os.Getenv("IMAGES_API_KEY"),
	}
}

func (s *settings) HttpConfig() http.HttpConfig {
	return http.HttpConfig{
		Port:         s.port,
		JWTSecret:    s.jwtSecret,
		RealIPHeader: s.realIPHeader,
	}
}

func (s *settings) DatabaseConfig() gateways.DatabaseConfig {
	return gateways.DatabaseConfig{
		ConnectionString: s.pgConnectionString,
		MinConnections:   10,
		MaxConnections:   100,
	}
}

func (s *settings) CacheConfig() gateways.CacheConfig {
	return gateways.CacheConfig{
		ConnectionString: s.redisCacheConnectionString,
		DialTimeout:      10 * time.Second,
		MinIdleConns:     10,
		PoolSize:         100,
		ReadTimeout:      -1,
		WriteTimeout:     -1,
	}
}

func (s *settings) SubscriberConfig() subscribe.SubscriberConfig {
	return subscribe.SubscriberConfig{
		ConnectionString: s.redisStreamConnectionString,
		DialTimeout:      10 * time.Second,
		MinIdleConns:     10,
		PoolSize:         100,
		ReadTimeout:      -1,
		WriteTimeout:     -1,
	}
}

func (s *settings) StreamConfig() gateways.StreamConfig {
	return gateways.StreamConfig{
		ConnectionString: s.redisStreamConnectionString,
		DialTimeout:      10 * time.Second,
		MinIdleConns:     10,
		PoolSize:         100,
		ReadTimeout:      -1,
		WriteTimeout:     -1,
	}
}

func (s *settings) BlockchainConfig() gateways.BlockchainConfig {
	return gateways.BlockchainConfig{
		BlockchainURL:     s.blockchainURL,
		ReputationAddress: "", //TODO
	}
}

func (s *settings) ImagesConfig() gateways.ImagesConfig {
	return gateways.ImagesConfig{
		BaseURL: s.imagesBaseUrl,
		APIKey:  s.imagesAPIKey,
	}
}
