package images

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type images struct {
	logger common.Logger
	client *http.Client
	config *gateways.ImagesConfig
}

func NewImagesGateway(logger common.Logger) gateways.Images {
	return &images{
		logger,
		nil,
		nil,
	}
}

func (i *images) Start(ctx context.Context, config gateways.ImagesConfig) {
	i.logger.Info(ctx).Msg("starting images gateway")
	i.config = &config
	i.client = &http.Client{}
}

func (i *images) Shutdown(ctx context.Context) {
	i.logger.Info(ctx).Msg("shutting down images gateway")
	i.client.CloseIdleConnections()
}

func (i *images) UploadImage(ctx context.Context, body io.Reader) (*entities.Image, error) {
	body, err := i.do(ctx, "POST", "/images", body, "application/octet-stream")

	if err != nil {
		return nil, fmt.Errorf("post images response error %w", err)
	}

	return i.toImage(body)
}

func (i *images) UploadAvatar(ctx context.Context, uri string, isNFT bool) (*entities.Image, error) {
	body, err := json.Marshal(&avatarRequestJSON{
		URL:   uri,
		IsNFT: isNFT,
	})

	if err != nil {
		return nil, fmt.Errorf("marshal put avatar request error %w", err)
	}

	respBody, err := i.do(ctx, "PUT", "/avatars", bytes.NewReader(body), "application/json")

	if err != nil {
		return nil, fmt.Errorf("put avatar request error %w", err)
	}

	return i.toImage(respBody)
}

func (i *images) GetImageByFileName(ctx context.Context, fileName string) (*entities.Image, error) {
	body, err := i.do(ctx, "GET", fmt.Sprintf("/images/%s", fileName), nil, "application/json")

	if err != nil {
		return nil, fmt.Errorf("get images response error %w", err)
	}

	return i.toImage(body)
}

func (i *images) do(ctx context.Context, method string, path string, body io.Reader, contentType string) (io.Reader, error) {
	url := fmt.Sprintf("%s%s", i.config.BaseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, fmt.Errorf("http url %v", url)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", i.config.APIKey))

	req.Header.Add("Content-Type", contentType)

	resp, err := i.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http response %v err %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http invalid status code %v %v", resp.StatusCode, url)
	}

	return resp.Body, nil
}

func (i *images) toImage(reader io.Reader) (*entities.Image, error) {
	imageJSON, err := common.Decode[imageJSON](reader)

	if err != nil {
		return nil, fmt.Errorf("decoding get images response error %w", err)
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

type imageJSON struct {
	FileName  string     `json:"file_name"`
	Original  headerJSON `json:"original"`
	Formatted headerJSON `json:"formatted"`
}

type headerJSON struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

type avatarRequestJSON struct {
	URL   string `json:"url"`
	IsNFT bool   `json:"is_nft"`
}
