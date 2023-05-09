package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
)

type s3Gateway struct {
	logger   common.Logger
	settings settings.Settings
	client   *s3.S3
}

func NewImageGateway(logger common.Logger, settings settings.Settings, client *s3.S3) usecases.ImageGateway {
	return &s3Gateway{
		logger,
		settings,
		client,
	}
}

func (g *s3Gateway) UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error) {
	bucket := g.settings.ImageBucket()
	_, err := g.client.PutObject(&s3.PutObjectInput{
		Bucket:       &bucket,
		Key:          &fileName,
		Body:         bytes.NewReader(*data),
		ContentType:  &contentType,
		CacheControl: aws.String("max-age=2592000"), // 30 days
	})

	if err != nil {
		return entities.Image{}, err
	}

	url := g.getImageUrl(fileName)

	return entities.NewImage(fileName, url, contentType), nil
}

// We get file header information to both verify that the file exists and to get the content type
func (g *s3Gateway) GetImageByFileName(ctx context.Context, fileName string) (entities.Image, error) {
	bucket := g.settings.ImageBucket()
	header, err := g.client.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
	})

	if err != nil {
		return entities.Image{}, err
	}

	url := g.getImageUrl(fileName)
	contentType := *header.ContentType

	return entities.NewImage(fileName, url, contentType), nil
}

func (g *s3Gateway) getImageUrl(fileName string) string {
	return fmt.Sprintf("%s/%s", g.settings.ImagePublicBaseURL(), fileName)
}
