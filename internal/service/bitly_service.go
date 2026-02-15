package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

type bitlyService struct {
	client *http.Client
	cfg    config.BitlyConfig
	log    zerolog.Logger
}

func newBitlyService(cfg config.BitlyConfig, log zerolog.Logger) BitlyService {
	return &bitlyService{
		client: &http.Client{Timeout: 10 * time.Second},
		cfg:    cfg,
		log:    log.With().Str("service", "bitly").Logger(),
	}
}

// bitlyRequest is the JSON body for the Bitly v4 /shorten endpoint.
type bitlyRequest struct {
	GroupGUID string `json:"group_guid"`
	Domain    string `json:"domain"`
	LongURL   string `json:"long_url"`
}

// bitlyResponse is the JSON response from the Bitly v4 /shorten endpoint.
type bitlyResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

func (s *bitlyService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	if s.cfg.APIKey == "" {
		return "", fmt.Errorf("bitly: API key is not configured")
	}
	if s.cfg.BaseURL == "" {
		return "", fmt.Errorf("bitly: base URL is not configured")
	}

	reqBody := bitlyRequest{
		GroupGUID: s.cfg.GroupGUID,
		Domain:    s.cfg.Domain,
		LongURL:   longURL,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("bitly: marshaling request: %w", err)
	}

	url := s.cfg.BaseURL + "/shorten"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("bitly: creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.log.Error().Err(err).Str("url", longURL).Msg("bitly API call failed")
		return "", fmt.Errorf("bitly: sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.log.Error().Int("statusCode", resp.StatusCode).Str("url", longURL).Msg("bitly API returned non-success status")
		return "", fmt.Errorf("bitly: API returned status %d", resp.StatusCode)
	}

	var result bitlyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("bitly: decoding response: %w", err)
	}

	s.log.Info().Str("longURL", longURL).Str("shortLink", result.Link).Msg("URL shortened")
	return result.Link, nil
}

func init() {
	var _ BitlyService = (*bitlyService)(nil)
}
