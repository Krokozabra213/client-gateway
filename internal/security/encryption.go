package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type Service struct {
	gcm cipher.AEAD
}

func NewEncryptionService(secret string) (*Service, error) {
	hash := sha256.Sum256([]byte(secret))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	return &Service{
		gcm: aesGCM,
	}, nil
}

func (s *Service) Encrypt(plainText string) (string, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := s.gcm.Seal(nil, nonce, []byte(plainText), nil)

	result := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func (s *Service) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	nonceSize := s.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid ciphertext: too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := s.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
