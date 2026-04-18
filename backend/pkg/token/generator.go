package token

import (
	"crypto/rand"
	"encoding/hex"
)

type secureTokenGenerator struct{}

func NewSecureTokenGenerator() *secureTokenGenerator {
	return &secureTokenGenerator{}
}

func (g *secureTokenGenerator) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
