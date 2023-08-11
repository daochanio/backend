package images

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
)

type images struct {
	logger     common.Logger
	settings   settings.Settings
	httpClient common.HttpClient
}

func NewImagesGateway(logger common.Logger, settings settings.Settings, httpClient common.HttpClient) usecases.Images {
	return &images{
		logger,
		settings,
		httpClient,
	}
}

func (i *images) Start(ctx context.Context) {}

func (i *images) Shutdown(ctx context.Context) {}

func (i *images) UploadImage(ctx context.Context, body io.Reader) (*entities.Image, error) {
	url := fmt.Sprintf("%s/images", i.settings.ImagesBaseURL())

	resp, err := i.httpClient.Do(ctx, "POST", url, body, &common.HttpOptions{
		Headers: []common.Header{
			i.getAuthorizationHeaders(),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("post images response error %w", err)
	}

	return i.toImage(resp.Body)
}

func (i *images) UploadAvatar(ctx context.Context, uri string, isNFT bool) (*entities.Image, error) {
	url := fmt.Sprintf("%s/avatars", i.settings.ImagesBaseURL())

	body, err := json.Marshal(&AvatarRequestJSON{
		URL:   uri,
		IsNFT: isNFT,
	})

	if err != nil {
		return nil, fmt.Errorf("marshal put avatar request error %w", err)
	}

	resp, err := i.httpClient.Do(ctx, "PUT", url, bytes.NewReader(body), &common.HttpOptions{
		Headers: []common.Header{
			i.getAuthorizationHeaders(),
			{
				Key:   "Content-Type",
				Value: "application/json",
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("put avatar request error %w", err)
	}

	return i.toImage(resp.Body)
}

func (i *images) GetImageByFileName(ctx context.Context, fileName string) (*entities.Image, error) {
	url := fmt.Sprintf("%s/images/%s", i.settings.ImagesBaseURL(), fileName)

	resp, err := i.httpClient.Do(ctx, "GET", url, nil, &common.HttpOptions{
		Headers: []common.Header{
			i.getAuthorizationHeaders(),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("get images response error %w", err)
	}

	return i.toImage(resp.Body)
}

func (i *images) getAuthorizationHeaders() common.Header {
	return common.Header{
		Key:   "Authorization",
		Value: fmt.Sprintf("Bearer %s", i.settings.ImagesAPIKey()),
	}
}

func (i *images) toImage(reader io.Reader) (*entities.Image, error) {
	imageJSON := &ImageJSON{}
	err := json.NewDecoder(reader).Decode(imageJSON)

	if err != nil {
		return nil, fmt.Errorf("unmarshal get images response error %w", err)
	}

	image := entities.NewImage(
		imageJSON.FileName,
		imageJSON.Original.URL,
		imageJSON.Original.ContentType,
		imageJSON.Formatted.URL,
		imageJSON.Formatted.ContentType,
	)

	return &image, nil
}

type ImageJSON struct {
	FileName  string     `json:"file_name"`
	Original  HeaderJSON `json:"original"`
	Formatted HeaderJSON `json:"formatted"`
}

type HeaderJSON struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

type AvatarRequestJSON struct {
	URL   string `json:"url"`
	IsNFT bool   `json:"is_nft"`
}
