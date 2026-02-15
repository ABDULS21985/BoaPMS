package service

import (
	"testing"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

// newTestEncryptionService creates an EncryptionService with a valid 32-byte
// test key (64 hex characters) for use in unit tests.
func newTestEncryptionService() EncryptionService {
	cfg := &config.Config{
		Encryption: config.EncryptionConfig{
			Key: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		},
	}
	return newEncryptionService(cfg, zerolog.Nop())
}

// newUnconfiguredEncryptionService creates an EncryptionService with an empty
// key, simulating a missing configuration.
func newUnconfiguredEncryptionService() EncryptionService {
	cfg := &config.Config{
		Encryption: config.EncryptionConfig{
			Key: "",
		},
	}
	return newEncryptionService(cfg, zerolog.Nop())
}

// ---------------------------------------------------------------------------
// Round-trip: encrypt then decrypt should return original
// ---------------------------------------------------------------------------

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	svc := newTestEncryptionService()

	plaintext := "Hello, Enterprise PMS!"
	encrypted, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt(%q) returned error: %v", plaintext, err)
	}
	if encrypted == "" {
		t.Fatal("Encrypt returned empty ciphertext")
	}
	if encrypted == plaintext {
		t.Fatal("Encrypt returned plaintext unchanged; expected ciphertext")
	}

	decrypted, err := svc.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("Decrypt = %q; want %q", decrypted, plaintext)
	}
}

// ---------------------------------------------------------------------------
// GCM nonce randomness: encrypting the same plaintext twice should produce
// different ciphertext.
// ---------------------------------------------------------------------------

func TestEncrypt_DifferentCiphertextEachTime(t *testing.T) {
	svc := newTestEncryptionService()

	plaintext := "same input"
	enc1, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("first Encrypt returned error: %v", err)
	}
	enc2, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("second Encrypt returned error: %v", err)
	}
	if enc1 == enc2 {
		t.Error("two encryptions of the same plaintext produced identical ciphertext; expected different nonces")
	}
}

// ---------------------------------------------------------------------------
// Decrypt with wrong / corrupted data should error
// ---------------------------------------------------------------------------

func TestDecrypt_InvalidData(t *testing.T) {
	svc := newTestEncryptionService()

	// "not-valid-base64!!" is not valid base64
	_, err := svc.Decrypt("not-valid-base64!!")
	if err == nil {
		t.Error("Decrypt with invalid base64 should return error")
	}

	// Valid base64 but wrong/corrupted ciphertext
	_, err = svc.Decrypt("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	if err == nil {
		t.Error("Decrypt with corrupted ciphertext should return error")
	}
}

// ---------------------------------------------------------------------------
// Empty plaintext round-trip
// ---------------------------------------------------------------------------

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	svc := newTestEncryptionService()

	encrypted, err := svc.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt('') returned error: %v", err)
	}
	if encrypted == "" {
		t.Fatal("Encrypt('') returned empty ciphertext; expected non-empty base64")
	}

	decrypted, err := svc.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}
	if decrypted != "" {
		t.Errorf("Decrypt = %q; want empty string", decrypted)
	}
}

// ---------------------------------------------------------------------------
// Nil cipher (unconfigured key) should return descriptive error
// ---------------------------------------------------------------------------

func TestEncrypt_UnconfiguredKey(t *testing.T) {
	svc := newUnconfiguredEncryptionService()

	_, err := svc.Encrypt("test")
	if err == nil {
		t.Fatal("Encrypt with unconfigured key should return error")
	}
	// The error message should mention "cipher not initialised"
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}

func TestDecrypt_UnconfiguredKey(t *testing.T) {
	svc := newUnconfiguredEncryptionService()

	_, err := svc.Decrypt("dGVzdA==")
	if err == nil {
		t.Fatal("Decrypt with unconfigured key should return error")
	}
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}
