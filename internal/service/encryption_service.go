package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// encryptionService implements the EncryptionService interface.
// It provides AES-256-GCM encryption and decryption. The 32-byte key is read
// from configuration (Config.Encryption.Key) as a hex-encoded string.
// Ciphertext is returned as a base64-encoded string that contains the GCM
// nonce prepended to the sealed ciphertext.
// ---------------------------------------------------------------------------

type encryptionService struct {
	gcm cipher.AEAD
	log zerolog.Logger
}

// newEncryptionService creates a new EncryptionService backed by AES-256-GCM.
// The encryption key is read from cfg.Encryption.Key as a 64-character
// hex-encoded string (32 bytes). If the key is empty or invalid the service
// still starts, but Encrypt/Decrypt will return descriptive errors.
func newEncryptionService(cfg *config.Config, log zerolog.Logger) EncryptionService {
	l := log.With().Str("service", "encryption").Logger()

	keyHex := cfg.Encryption.Key
	if keyHex == "" {
		l.Warn().Msg("encryption key is not configured; encrypt/decrypt operations will fail until a key is set")
		return &encryptionService{gcm: nil, log: l}
	}

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		l.Error().Err(err).Msg("failed to hex-decode encryption key; encrypt/decrypt operations will fail")
		return &encryptionService{gcm: nil, log: l}
	}

	if len(keyBytes) != 32 {
		l.Error().Int("keyLength", len(keyBytes)).Msg("encryption key must be exactly 32 bytes (64 hex chars); encrypt/decrypt operations will fail")
		return &encryptionService{gcm: nil, log: l}
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		l.Error().Err(err).Msg("failed to create AES cipher; encrypt/decrypt operations will fail")
		return &encryptionService{gcm: nil, log: l}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		l.Error().Err(err).Msg("failed to create GCM cipher; encrypt/decrypt operations will fail")
		return &encryptionService{gcm: nil, log: l}
	}

	l.Info().Msg("encryption service initialised with AES-256-GCM")
	return &encryptionService{gcm: gcm, log: l}
}

// Encrypt encrypts the plaintext using AES-256-GCM and returns a
// base64-encoded string. The output format is base64(nonce + ciphertext).
func (s *encryptionService) Encrypt(plaintext string) (string, error) {
	if s.gcm == nil {
		return "", fmt.Errorf("encryption: cipher not initialised (check encryption key configuration)")
	}

	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		s.log.Error().Err(err).Msg("failed to generate nonce")
		return "", fmt.Errorf("encryption: generating nonce: %w", err)
	}

	// Seal appends the ciphertext to the nonce slice so the result is
	// nonce + ciphertext in a single byte slice.
	sealed := s.gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	encoded := base64.StdEncoding.EncodeToString(sealed)
	return encoded, nil
}

// Decrypt decodes the base64 ciphertext and decrypts it using AES-256-GCM.
// It expects the input format produced by Encrypt: base64(nonce + ciphertext).
func (s *encryptionService) Decrypt(ciphertext string) (string, error) {
	if s.gcm == nil {
		return "", fmt.Errorf("decryption: cipher not initialised (check encryption key configuration)")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to base64-decode ciphertext")
		return "", fmt.Errorf("decryption: base64 decoding: %w", err)
	}

	nonceSize := s.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("decryption: ciphertext too short (expected at least %d bytes for nonce, got %d)", nonceSize, len(data))
	}

	nonce, sealed := data[:nonceSize], data[nonceSize:]

	plaintext, err := s.gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		s.log.Error().Err(err).Msg("AES-GCM decryption failed")
		return "", fmt.Errorf("decryption: AES-GCM open: %w", err)
	}

	return string(plaintext), nil
}

func init() {
	// Compile-time interface compliance check.
	var _ EncryptionService = (*encryptionService)(nil)
}
