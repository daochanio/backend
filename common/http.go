package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(ctx context.Context, method string, url string, body io.Reader, options *HttpOptions) (*http.Response, error)
}

type Header struct {
	Key   string
	Value string
}
type HttpOptions struct {
	Headers []Header
}

type httpClient struct {
	logger Logger
	client *http.Client
}

func NewHttpClient(logger Logger) HttpClient {
	client := &http.Client{}
	return &httpClient{
		logger,
		client,
	}
}

func (c *httpClient) Do(ctx context.Context, method string, url string, body io.Reader, options *HttpOptions) (*http.Response, error) {
	return FunctionRetrier(ctx, func() (*http.Response, error) {
		return c.doInternal(ctx, method, url, body, options)
	})
}

func (c *httpClient) doInternal(ctx context.Context, method string, uri string, body io.Reader, options *HttpOptions) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, body)

	if err != nil {
		return nil, fmt.Errorf("http url %v", uri)
	}

	if options != nil {
		for _, header := range options.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http response %v err %w", uri, err)
	}

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("http rate limit %v %w", uri, ErrRetryable)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http invalid status code %v %v", resp.StatusCode, uri)
	}

	return resp, nil
}
