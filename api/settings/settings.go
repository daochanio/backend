package settings

import (
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Settings interface {
	Port() string
	JWTSecret() string
	PostgresConfig() *pgxpool.Config
	RedisCacheOptions() *redis.Options
	RedisStreamOptions() *redis.Options
	BlockchainURI() string
	RealIPHeader() string
	ImagesBaseURL() string
	ImagesAPIKey() string
}

type settings struct {
	port                        string
	pgConnectionString          string
	redisCacheConnectionString  string
	redisStreamConnectionString string
	jwtSecret                   string
	blockchainURI               string
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
		blockchainURI:               os.Getenv("BLOCKCHAIN_URI"),
		realIPHeader:                os.Getenv("REAL_IP_HEADER"),
		imagesBaseUrl:               os.Getenv("IMAGES_BASE_URL"),
		imagesAPIKey:                os.Getenv("IMAGES_API_KEY"),
	}
}

func (s *settings) Port() string {
	return s.port
}

func (s *settings) PostgresConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(s.pgConnectionString)

	if err != nil {
		panic(err)
	}

	config.MinConns = 10
	config.MaxConns = 100
	return config
}

func (s *settings) RedisCacheOptions() *redis.Options {
	return s.buildRedisOptions(s.redisCacheConnectionString)
}

func (s *settings) RedisStreamOptions() *redis.Options {
	return s.buildRedisOptions(s.redisStreamConnectionString)
}

func (s *settings) buildRedisOptions(connStr string) *redis.Options {
	opt, err := redis.ParseURL(connStr)

	if err != nil {
		panic(err)
	}

	opt.DialTimeout = 10 * time.Second
	opt.MinIdleConns = 10
	opt.PoolSize = 100
	// timeouts are handled through request context
	opt.ReadTimeout = -1
	opt.WriteTimeout = -1
	return opt
}

func (s *settings) JWTSecret() string {
	return s.jwtSecret
}

func (s *settings) BlockchainURI() string {
	return s.blockchainURI
}

func (s *settings) RealIPHeader() string {
	return s.realIPHeader
}

func (s *settings) ImagesBaseURL() string {
	return s.imagesBaseUrl
}

func (s *settings) ImagesAPIKey() string {
	return s.imagesAPIKey
}
