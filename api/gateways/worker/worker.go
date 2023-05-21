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

type nftMetadata struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Image       string                    `json:"image"`
	ImageURL    string                    `json:"image_url"`
	ImageData   string                    `json:"image_data"`
	Attributes  *[]map[string]interface{} `json:"attributes"`
	Properties  *map[string]interface{}   `json:"properties"`
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
	} else if metadata.ImageData != "" {
		return metadata.ImageData, nil
	}

	return "", errors.New("failed to detect nft image uri")
}

func (w *worker) safeProxy(ctx context.Context, uri string) (*http.Response, error) {
	url := w.parseURI(uri)
	proxyURL := fmt.Sprintf("%s/proxy?url=%s", w.settings.WokerURI(), url)

	resp, err := w.httpClient.Do(ctx, "GET", proxyURL, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("safe proxy error for uri %v %w", uri, err)
	}

	return resp, nil
}

func (w *worker) parseURI(uri string) string {
	if suffix, ok := strings.CutPrefix(uri, "ipfs://"); ok {
		if !strings.HasPrefix(suffix, "ipfs/") {
			suffix = fmt.Sprintf("ipfs/%s", suffix)
		}
		return fmt.Sprintf("%s/%s", w.settings.IPFSGatewayURI(), suffix)
	}
	if suffix, ok := strings.CutPrefix(uri, "ipns://"); ok {
		if !strings.HasPrefix(suffix, "ipns/") {
			suffix = fmt.Sprintf("ipns/%s", suffix)
		}
		return fmt.Sprintf("%s/%s", w.settings.IPFSGatewayURI(), suffix)
	}
	return uri
}
