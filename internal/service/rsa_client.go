package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

type rsaAuthService struct {
	client  *http.Client
	cfg     config.RSAConfig
	gs      GlobalSettingService
	log     zerolog.Logger
	baseURL string
	apiKey  string
	mu      sync.RWMutex
	inited  bool
}

func newRSAAuthService(cfg config.RSAConfig, gs GlobalSettingService, log zerolog.Logger) RSAAuthService {
	return &rsaAuthService{
		client: &http.Client{Timeout: 30 * time.Second},
		cfg:    cfg,
		gs:     gs,
		log:    log.With().Str("service", "rsa_auth").Logger(),
	}
}

// ensureInit lazily loads RSA_BASE_URL and RSA_API_KEY from config or
// GlobalSettingService (database). Config values take precedence.
func (s *rsaAuthService) ensureInit(ctx context.Context) error {
	s.mu.RLock()
	if s.inited {
		s.mu.RUnlock()
		return nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.inited {
		return nil
	}

	// Prefer config values; fall back to database settings.
	baseURL := s.cfg.BaseURL
	apiKey := s.cfg.APIKey

	if baseURL == "" && s.gs != nil {
		val, err := s.gs.GetStringValue(ctx, "RSA_BASE_URL")
		if err == nil && val != "" {
			baseURL = val
		}
	}
	if apiKey == "" && s.gs != nil {
		val, err := s.gs.GetStringValue(ctx, "RSA_API_KEY")
		if err == nil && val != "" {
			apiKey = val
		}
	}

	if baseURL == "" {
		return fmt.Errorf("rsa_auth: base URL is not configured")
	}

	s.baseURL = baseURL
	s.apiKey = apiKey
	s.inited = true
	s.log.Info().Str("baseURL", baseURL).Msg("RSA auth service initialised")
	return nil
}

func (s *rsaAuthService) Initialize(ctx context.Context, req *RSAInitializeRequest) (*RSAInitializeResponse, error) {
	if err := s.ensureInit(ctx); err != nil {
		return nil, err
	}

	var resp RSAInitializeResponse
	if err := s.doPost(ctx, s.baseURL+"/initialize", req, &resp); err != nil {
		return nil, fmt.Errorf("rsa_auth: initialize: %w", err)
	}
	return &resp, nil
}

func (s *rsaAuthService) Verify(ctx context.Context, req *RSAVerifyRequest) (*RSAVerifyResponse, error) {
	if err := s.ensureInit(ctx); err != nil {
		return nil, err
	}

	var resp RSAVerifyResponse
	if err := s.doPost(ctx, s.baseURL+"/verify", req, &resp); err != nil {
		return nil, fmt.Errorf("rsa_auth: verify: %w", err)
	}
	return &resp, nil
}

func (s *rsaAuthService) doPost(ctx context.Context, url string, body, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if s.apiKey != "" {
		httpReq.Header.Set("client-key", s.apiKey)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.log.Error().Err(err).Str("url", url).Msg("RSA API call failed")
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.log.Error().Int("statusCode", resp.StatusCode).Str("url", url).Msg("RSA API returned non-success status")
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

func init() {
	var _ RSAAuthService = (*rsaAuthService)(nil)
}
