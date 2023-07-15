package settings

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Settings interface {
	Port() string
	JWTSecret() string
	PostgresConfig() *pgxpool.Config
	RedisCacheOptions() *redis.Options
	RedisStreamOptions() *redis.Options
	S3Config(context.Context) *aws.Config
	StaticPublicBaseURL() string
	ImageBucket() string
	BlockchainURI() string
	RealIPHeader() string
	IPFSGatewayURI() string
	WokerURI() string
}

type settings struct {
	port                        string
	pgConnectionString          string
	redisCacheConnectionString  string
	redisStreamConnectionString string
	jwtSecret                   string
	blockchainURI               string
	realIPHeader                string
	staticPublicBaseURL         string
	staticAccessKeyId           string
	staticSecretAccessKey       string
	staticURL                   string
	imageBucket                 string
	ipfsGatewayURI              string
	workerURI                   string
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
		staticPublicBaseURL:         os.Getenv("STATIC_PUBLIC_BASE_URL"),
		staticAccessKeyId:           os.Getenv("STATIC_ACCESS_KEY_ID"),
		staticSecretAccessKey:       os.Getenv("STATIC_SECRET_ACCESS_KEY"),
		staticURL:                   os.Getenv("STATIC_URL"),
		imageBucket:                 os.Getenv("IMAGE_BUCKET"),
		ipfsGatewayURI:              os.Getenv("IPFS_GATEWAY_URI"),
		workerURI:                   os.Getenv("WORKER_URI"),
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

func (s *settings) S3Config(ctx context.Context) *aws.Config {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s.staticURL,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.staticAccessKeyId, s.staticSecretAccessKey, "")))

	if err != nil {
		panic(err)
	}

	return &cfg
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

func (s *settings) StaticPublicBaseURL() string {
	return s.staticPublicBaseURL
}

func (s *settings) ImageBucket() string {
	return s.imageBucket
}

func (s *settings) IPFSGatewayURI() string {
	return s.ipfsGatewayURI
}

func (s *settings) WokerURI() string {
	return s.workerURI
}
