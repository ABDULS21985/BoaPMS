package rsa

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

// Client is an HTTP client for the RSA SecurID authentication API.
// It is safe for concurrent use.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client used by the RSA client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the timeout on the default http.Client.
// Ignored when WithHTTPClient is also supplied.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{Timeout: d}
		}
	}
}

// NewClient creates an RSA SecurID API client.
//
//	baseURL – the RSA SecurID REST API base URL (e.g. "https://rsa.example.com/api/v1")
//	apiKey  – the value sent in the "client-key" header on every request
func NewClient(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return c
}

// Initialize starts an RSA SecurID authentication attempt by calling the
// "initialize" endpoint.
func (c *Client) Initialize(ctx context.Context, req *InitializeRequest) (*Response, error) {
	return c.do(ctx, "initialize", req)
}

// Verify submits a token code against an in-progress authentication attempt
// by calling the "verify" endpoint.
func (c *Client) Verify(ctx context.Context, req *VerifyRequest) (*Response, error) {
	return c.do(ctx, "verify", req)
}

// do is the shared HTTP round-trip helper.
func (c *Client) do(ctx context.Context, endpoint string, payload any) (*Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("rsa: marshal request: %w", err)
	}

	url := c.baseURL + "/" + endpoint

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("rsa: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("client-key", c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("rsa: %s: %w", endpoint, err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("rsa: read response body: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Endpoint:   endpoint,
			Body:       string(respBody),
		}
	}

	var resp Response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("rsa: unmarshal response: %w", err)
	}

	return &resp, nil
}

// APIError is returned when the RSA API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Endpoint   string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("rsa: %s returned HTTP %d: %s", e.Endpoint, e.StatusCode, e.Body)
}
