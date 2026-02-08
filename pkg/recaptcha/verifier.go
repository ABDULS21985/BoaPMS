package recaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultVerifyURL is Google's reCAPTCHA server-side verification endpoint.
const DefaultVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Config holds Google reCAPTCHA configuration.
type Config struct {
	SiteKey   string `mapstructure:"site_key"`
	SecretKey string `mapstructure:"secret_key"`
	VerifyURL string `mapstructure:"verify_url"`
}

// VerifyResponse is the response from Google's reCAPTCHA verify API.
type VerifyResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

// Verifier handles server-side reCAPTCHA verification.
type Verifier struct {
	config     Config
	httpClient *http.Client
}

// NewVerifier creates a new reCAPTCHA verifier.
func NewVerifier(cfg Config) *Verifier {
	if cfg.VerifyURL == "" {
		cfg.VerifyURL = DefaultVerifyURL
	}

	return &Verifier{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Verify validates a reCAPTCHA token with Google's API.
// Returns the verification response or an error.
func (v *Verifier) Verify(ctx context.Context, token string) (*VerifyResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("recaptcha: token is empty")
	}

	// Build form-encoded POST body: secret=...&response=...
	form := url.Values{}
	form.Set("secret", v.config.SecretKey)
	form.Set("response", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.config.VerifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("recaptcha: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recaptcha: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recaptcha: unexpected status code %d from verify endpoint", resp.StatusCode)
	}

	var result VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("recaptcha: failed to decode response: %w", err)
	}

	return &result, nil
}

// VerifyWithThreshold validates a reCAPTCHA token and checks the score against a minimum threshold.
// Returns true if the token is valid and the score meets the threshold.
func (v *Verifier) VerifyWithThreshold(ctx context.Context, token string, threshold float64) (bool, error) {
	result, err := v.Verify(ctx, token)
	if err != nil {
		return false, err
	}

	if !result.Success {
		return false, fmt.Errorf("recaptcha: verification failed, error codes: %v", result.ErrorCodes)
	}

	return result.Score >= threshold, nil
}
