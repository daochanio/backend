package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
)

type worker struct {
	logger     common.Logger
	settings   settings.Settings
	httpClient common.HttpClient
}

func NewSafeProxyGateway(logger common.Logger, settings settings.Settings, httpClient common.HttpClient) usecases.SafeProxy {
	return &worker{
		logger,
		settings,
		httpClient,
	}
}

func (w *worker) DownloadImage(ctx context.Context, uri string) (*[]byte, string, error) {
	resp, err := w.safeProxy(ctx, uri)

	if err != nil {
		return nil, "", fmt.Errorf("failed to get data from uri: %w", err)
	}

	defer resp.Body.Close()
	limit := int64(1024 * 1024 * 5)
	limitedReader := io.LimitReader(resp.Body, limit)
	data, err := io.ReadAll(limitedReader)
	if err != nil && err != io.EOF {
		return nil, "", fmt.Errorf("failed to read uri response body: %w", err)
	}

	contentType := http.DetectContentType(data)

	if !strings.HasPrefix(contentType, "image") {
		return nil, "", fmt.Errorf("invalid content type: %s", contentType)
	}

	return &data, contentType, nil
}

func (w *worker) GetNFTImageURI(ctx context.Context, uri string) (string, error) {
	resp, err := w.safeProxy(ctx, uri)

	if err != nil {
		return "", fmt.Errorf("failed to get nft metadata: %w", err)
	}

	defer resp.Body.Close()
	metadata := &nftMetadata{}
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return "", fmt.Errorf("metadata decode response %w", err)
	}

	if metadata.Image != "" {
		return metadata.Image, nil
	} else if metadata.ImageURL != "" {
		return metadata.ImageURL, nil
	}

	return "", errors.New("failed to detect nft image uri")
}

func (w *worker) safeProxy(ctx context.Context, uri string) (*http.Response, error) {
	url := w.settings.IPFSGatewayURI(uri)
	proxyURL := fmt.Sprintf("%s/proxy?url=%s", w.settings.WokerURI(), url)

	resp, err := w.httpClient.Do(ctx, "GET", proxyURL, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("safe proxy error for uri %v %w", uri, err)
	}

	return resp, nil
}

type nftMetadata struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Image       string                    `json:"image"`
	ImageURL    string                    `json:"image_url"`
	Attributes  *[]map[string]interface{} `json:"attributes"`
	Properties  *map[string]interface{}   `json:"properties"`
}
