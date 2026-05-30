package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type TokenCipher struct {
	aead cipher.AEAD
}

func NewTokenCipher(key string) (*TokenCipher, error) {
	keyBytes, err := decodeKey(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("create token cipher failed: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create token cipher gcm failed: %w", err)
	}

	return &TokenCipher{aead: aead}, nil
}

func (c *TokenCipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("create token nonce failed: %w", err)
	}

	ciphertext := c.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *TokenCipher) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode encrypted token failed: %w", err)
	}

	nonceSize := c.aead.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("encrypted token is too short")
	}

	nonce := data[:nonceSize]
	encryptedToken := data[nonceSize:]

	plaintext, err := c.aead.Open(nil, nonce, encryptedToken, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt token failed: %w", err)
	}

	return string(plaintext), nil
}

func decodeKey(key string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN_ENCRYPTION_KEY not set")
	}

	if decoded, err := base64.StdEncoding.DecodeString(key); err == nil && len(decoded) == 32 {
		return decoded, nil
	}

	hashedKey := sha256.Sum256([]byte(key))
	return hashedKey[:], nil
}
