package service

import (
	"crypto/rand"
	"math/big"
)

const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

// passwordGenerator produces cryptographically random passwords.
type passwordGenerator struct{}

func newPasswordGenerator() PasswordGenerator {
	return &passwordGenerator{}
}

func (g *passwordGenerator) Generate(length int) string {
	if length < 8 {
		length = 8
	}
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordCharset))))
		if err != nil {
			b[i] = 'x' // fallback (should never happen)
			continue
		}
		b[i] = passwordCharset[n.Int64()]
	}
	return string(b)
}
