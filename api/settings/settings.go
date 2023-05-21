package settings

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Settings interface {
	Port() string
	JWTSecret() string
	PostgresConfig() *pgxpool.Config
	RegionalRedisOptions() *redis.Options
	GlobalRedisOptions() *redis.Options
	S3Config() *aws.Config
	StaticPublicBaseURL() string
	ImageBucket() string
	BlockchainURI() string
	RealIPHeader() string
	IPFSGatewayURI(uri string) string
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
	staticRegion                string
	imageBucket                 string
	ipfsGatewayURI              string
	workerURI                   string
}

func NewSettings() Settings {
	appRegion := strings.ToUpper(os.Getenv(os.Getenv("APP_REGION_KEY")))
	redisCacheConnectionString := fmt.Sprintf("REDIS_CACHE_CONNECTION_STRING_%s", appRegion)
	return &settings{
		port:                        os.Getenv("PORT"),
		pgConnectionString:          os.Getenv("PG_CONNECTION_STRING"),
		redisCacheConnectionString:  os.Getenv(redisCacheConnectionString),
		redisStreamConnectionString: os.Getenv("REDIS_STREAM_CONNECTION_STRING"),
		jwtSecret:                   os.Getenv("JWT_SECRET"),
		blockchainURI:               os.Getenv("BLOCKCHAIN_URI"),
		realIPHeader:                os.Getenv("REAL_IP_HEADER"),
		staticPublicBaseURL:         os.Getenv("STATIC_PUBLIC_BASE_URL"),
		staticAccessKeyId:           os.Getenv("STATIC_ACCESS_KEY_ID"),
		staticSecretAccessKey:       os.Getenv("STATIC_SECRET_ACCESS_KEY"),
		staticURL:                   os.Getenv("STATIC_URL"),
		staticRegion:                os.Getenv("STATIC_REGION"),
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

func (s *settings) RegionalRedisOptions() *redis.Options {
	return s.buildRedisConfig(s.redisCacheConnectionString)
}

func (s *settings) GlobalRedisOptions() *redis.Options {
	return s.buildRedisConfig(s.redisStreamConnectionString)
}

func (s *settings) buildRedisConfig(connStr string) *redis.Options {
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

func (s *settings) S3Config() *aws.Config {
	credentials := credentials.NewStaticCredentials(s.staticAccessKeyId, s.staticSecretAccessKey, "")
	config := aws.NewConfig().WithCredentials(credentials).WithEndpoint(s.staticURL).WithRegion(s.staticRegion)
	return config
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

func (s *settings) IPFSGatewayURI(uri string) string {
	if suffix, ok := strings.CutPrefix(uri, "ipfs://"); ok {
		if !strings.HasPrefix(suffix, "ipfs/") {
			suffix = fmt.Sprintf("ipfs/%s", suffix)
		}
		return fmt.Sprintf("%s/%s", s.ipfsGatewayURI, suffix)
	}
	return uri
}

func (s *settings) WokerURI() string {
	return s.workerURI
}
