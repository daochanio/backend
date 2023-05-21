package s3

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(*data)

	if err != nil {
		return entities.Image{}, fmt.Errorf("failed to compress image data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return entities.Image{}, err
	}

	bucket := g.settings.ImageBucket()
	contentEncoding := "gzip"
	_, err = g.client.PutObject(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(buf.Bytes()),
		// Body: bytes.NewReader(*data),
		ContentEncoding: &contentEncoding,
		ContentType:     &contentType,
		CacheControl:    aws.String("max-age=31536000"), // 1yr
	})

	if err != nil {
		return entities.Image{}, err
	}

	url := g.getExternalURL(fileName)

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

	url := g.getExternalURL(fileName)
	contentType := *header.ContentType

	return entities.NewImage(fileName, url, contentType), nil
}

func (g *s3Gateway) getExternalURL(fileName string) string {
	return fmt.Sprintf("%s/%s", g.settings.StaticPublicBaseURL(), fileName)
}
