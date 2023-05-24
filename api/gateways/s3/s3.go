package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
)

type s3Gateway struct {
	logger   common.Logger
	settings settings.Settings
	client   *s3.Client
}

func NewStorageGateway(ctx context.Context, logger common.Logger, settings settings.Settings) usecases.Storage {
	return &s3Gateway{
		logger:   logger,
		settings: settings,
		client:   nil,
	}
}

func (g *s3Gateway) Start(ctx context.Context) {
	g.client = s3.NewFromConfig(*g.settings.S3Config(ctx))
}

func (g *s3Gateway) Shutdown(ctx context.Context) {}

func (g *s3Gateway) UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error) {
	if data == nil {
		return entities.Image{}, fmt.Errorf("invalid image data")
	}

	bucket := g.settings.ImageBucket()
	_, err := g.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       &bucket,
		Key:          &fileName,
		Body:         bytes.NewReader(*data),
		ContentType:  &contentType,
		CacheControl: aws.String("max-age=31536000"), // 1yr
	})

	if err != nil {
		return entities.Image{}, err
	}

	url := g.getExternalURL(fileName)

	return entities.NewImage(fileName, url, contentType), nil
}

func (g *s3Gateway) GetImageByFileName(ctx context.Context, fileName string) (*entities.Image, error) {
	bucket := g.settings.ImageBucket()
	header, err := g.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
	})

	var responseError *awshttp.ResponseError
	if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	url := g.getExternalURL(fileName)
	contentType := *header.ContentType

	image := entities.NewImage(fileName, url, contentType)

	return &image, nil
}

func (g *s3Gateway) getExternalURL(fileName string) string {
	return fmt.Sprintf("%s/%s", g.settings.StaticPublicBaseURL(), fileName)
}
