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
	ImagePublicBaseURL() string
	ImageAccessKeyId() string
	ImageSecretAccessKey() string
	ImageURL() string
	ImageRegion() string
	ImageBucket() string
}

type settings struct {
	port                  string
	pgConnectionString    string
	redisConnectionString string
	blockchainURI         string
	realIPHeader          string
	imagePublicBaseURL    string
	imageAccessKeyId      string
	imageSecretAccessKey  string
	imageURL              string
	imageRegion           string
	imageBucket           string
}

func NewSettings() Settings {
	return &settings{
		port:                  os.Getenv("PORT"),
		pgConnectionString:    os.Getenv("PG_CONNECTION_STRING"),
		redisConnectionString: os.Getenv("REDIS_CONNECTION_STRING"),
		blockchainURI:         os.Getenv("BLOCKCHAIN_URI"),
		realIPHeader:          os.Getenv("REAL_IP_HEADER"),
		imagePublicBaseURL:    os.Getenv("IMAGE_PUBLIC_BASE_URL"),
		imageAccessKeyId:      os.Getenv("IMAGE_ACCESS_KEY_ID"),
		imageSecretAccessKey:  os.Getenv("IMAGE_SECRET_ACCESS_KEY"),
		imageURL:              os.Getenv("IMAGE_URL"),
		imageRegion:           os.Getenv("IMAGE_REGION"),
		imageBucket:           os.Getenv("IMAGE_BUCKET"),
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

func (s *settings) ImagePublicBaseURL() string {
	return s.imagePublicBaseURL
}

func (s *settings) ImageAccessKeyId() string {
	return s.imageAccessKeyId
}

func (s *settings) ImageSecretAccessKey() string {
	return s.imageSecretAccessKey
}

func (s *settings) ImageURL() string {
	return s.imageURL
}

func (s *settings) ImageRegion() string {
	return s.imageRegion
}

func (s *settings) ImageBucket() string {
	return s.imageBucket
}
