package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
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

func NewStorageGateway(logger common.Logger, settings settings.Settings) usecases.Storage {
	sess, err := session.NewSession(settings.S3Config())
	if err != nil {
		panic(err)
	}
	client := s3.New(sess)
	return &s3Gateway{
		logger,
		settings,
		client,
	}
}

func (g *s3Gateway) UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error) {
	if data == nil {
		return entities.Image{}, fmt.Errorf("invalid image data")
	}

	bucket := g.settings.ImageBucket()
	_, err := g.client.PutObject(&s3.PutObjectInput{
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
	header, err := g.client.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
	})

	if awsError, ok := err.(awserr.Error); ok && awsError.Code() == "NotFound" {
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
