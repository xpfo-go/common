package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Options struct {
	Timeout          time.Duration
	Retry            int
	RetryWait        time.Duration
	RetryStatusCodes []int
}

type Client struct {
	httpClient       *http.Client
	retry            int
	retryWait        time.Duration
	retryStatusCodes map[int]struct{}
}

func NewClient(options Options) *Client {
	if options.Timeout <= 0 {
		options.Timeout = 10 * time.Second
	}
	if options.Retry < 0 {
		options.Retry = 0
	}
	if options.RetryWait <= 0 {
		options.RetryWait = 100 * time.Millisecond
	}

	codes := options.RetryStatusCodes
	if len(codes) == 0 {
		codes = []int{
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		}
	}
	retryStatusCodes := make(map[int]struct{}, len(codes))
	for _, code := range codes {
		retryStatusCodes[code] = struct{}{}
	}

	return &Client{
		httpClient: &http.Client{Timeout: options.Timeout},
		retry:      options.Retry,
		retryWait:  options.RetryWait,

		retryStatusCodes: retryStatusCodes,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	var bodyBytes []byte
	if req.Body != nil {
		bs, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		_ = req.Body.Close()
		bodyBytes = bs
	}

	attempts := c.retry + 1
	for attempt := 0; attempt < attempts; attempt++ {
		attemptReq := req.Clone(ctx)
		if bodyBytes != nil {
			attemptReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			attemptReq.ContentLength = int64(len(bodyBytes))
		}

		resp, err := c.httpClient.Do(attemptReq)
		if err != nil {
			if attempt < attempts-1 && sleepWithContext(ctx, c.retryWait) == nil {
				continue
			}
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, err
		}

		if _, ok := c.retryStatusCodes[resp.StatusCode]; ok && attempt < attempts-1 {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if err := sleepWithContext(ctx, c.retryWait); err != nil {
				return nil, err
			}
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request retry exhausted")
}

func (c *Client) JSON(ctx context.Context, method, url string, in any, headers map[string]string, out any) error {
	var body io.Reader
	if in != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(in); err != nil {
			return fmt.Errorf("encode request json: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if in != nil && strings.TrimSpace(req.Header.Get("Content-Type")) == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		bs, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(bs)))
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil && err != io.EOF {
		return fmt.Errorf("decode response json: %w", err)
	}
	return nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
