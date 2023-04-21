package cloudfront

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
)

type cloudfrontGateway struct {
	logger   common.Logger
	settings settings.Settings
	client   *s3.S3
}

func NewCloudfrontGateway(logger common.Logger, settings settings.Settings) gateways.ImageGateway {
	credentials := credentials.NewStaticCredentials(settings.ImageAccessKeyId(), settings.ImageSecretAccessKey(), "")
	config := aws.NewConfig().WithCredentials(credentials).WithEndpoint(settings.ImageURL()).WithRegion(settings.ImageRegion())
	sess, err := session.NewSession(config)

	if err != nil {
		panic(err)
	}

	client := s3.New(sess)

	return &cloudfrontGateway{
		logger,
		settings,
		client,
	}
}

func (g *cloudfrontGateway) UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error) {
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
func (g *cloudfrontGateway) GetImageById(ctx context.Context, fileName string) (entities.Image, error) {
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

func (g *cloudfrontGateway) getImageUrl(fileName string) string {
	return fmt.Sprintf("%s/%s", g.settings.ImagePublicBaseURL(), fileName)
}
