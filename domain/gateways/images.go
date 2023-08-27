package gateways

import (
	"context"
	"io"

	"github.com/daochanio/backend/domain/entities"
)

type ImagesConfig struct {
	BaseURL string
	APIKey  string
}

type Images interface {
	Start(ctx context.Context, config ImagesConfig)
	Shutdown(ctx context.Context)
	UploadImage(ctx context.Context, reader io.Reader) (*entities.Image, error)
	GetImageByFileName(ctx context.Context, fileName string) (*entities.Image, error)
	UploadAvatar(ctx context.Context, uri string, isNFT bool) (*entities.Image, error)
}
