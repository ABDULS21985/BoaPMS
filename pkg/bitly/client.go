package bitly

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

// Config holds Bitly API configuration.
type Config struct {
	URL       string `mapstructure:"url"`
	GroupGUID string `mapstructure:"group_guid"`
	Token     string `mapstructure:"token"`
	Domain    string `mapstructure:"domain"`
}

// Client is a Bitly API client for URL shortening.
// It is safe for concurrent use.
type Client struct {
	config     Config
	httpClient *http.Client
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client used by the Bitly client.
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

// ShortenRequest is the request body for the Bitly shorten API.
type ShortenRequest struct {
	GroupGUID string `json:"group_guid"`
	Domain    string `json:"domain"`
	LongURL   string `json:"long_url"`
}

// ShortenResponse is the response from the Bitly shorten API.
type ShortenResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

// NewClient creates a new Bitly API client.
//
//	cfg – configuration containing the API URL, group GUID, token, and domain
func NewClient(cfg Config, opts ...Option) *Client {
	c := &Client{
		config: Config{
			URL:       strings.TrimRight(cfg.URL, "/"),
			GroupGUID: cfg.GroupGUID,
			Token:     cfg.Token,
			Domain:    cfg.Domain,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return c
}

// Shorten creates a shortened URL for the given long URL.
func (c *Client) Shorten(ctx context.Context, longURL string) (*ShortenResponse, error) {
	reqBody := ShortenRequest{
		GroupGUID: c.config.GroupGUID,
		Domain:    c.config.Domain,
		LongURL:   longURL,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("bitly: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("bitly: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.Token)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bitly: shorten: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("bitly: read response body: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Body:       string(respBody),
		}
	}

	var resp ShortenResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("bitly: unmarshal response: %w", err)
	}

	return &resp, nil
}

// APIError is returned when the Bitly API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bitly: shorten returned HTTP %d: %s", e.StatusCode, e.Body)
}
