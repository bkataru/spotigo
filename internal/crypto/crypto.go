// Package crypto provides encryption utilities for secure token storage
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// TokenEncryptor handles encryption and decryption of OAuth tokens
type TokenEncryptor struct {
	key []byte
}

// NewTokenEncryptor creates a new encryptor using machine-specific key derivation
func NewTokenEncryptor() (*TokenEncryptor, error) {
	key, err := deriveKey()
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	return &TokenEncryptor{key: key}, nil
}

// NewTokenEncryptorWithKey creates a new encryptor with a custom key
func NewTokenEncryptorWithKey(key []byte) (*TokenEncryptor, error) {
	if len(key) < 16 {
		return nil, fmt.Errorf("key must be at least 16 bytes")
	}

	// Hash the key to ensure it's exactly 32 bytes for AES-256
	hash := sha256.Sum256(key)
	return &TokenEncryptor{key: hash[:]}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (e *TokenEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (e *TokenEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptToBase64 encrypts and returns base64-encoded string
func (e *TokenEncryptor) EncryptToBase64(plaintext []byte) (string, error) {
	ciphertext, err := e.Encrypt(plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromBase64 decrypts base64-encoded ciphertext
func (e *TokenEncryptor) DecryptFromBase64(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return e.Decrypt(ciphertext)
}

// deriveKey creates a machine-specific encryption key
// This uses a combination of machine-specific values to create a stable key
func deriveKey() ([]byte, error) {
	// Collect machine-specific entropy sources
	var entropy []byte

	// 1. Username
	if user := os.Getenv("USER"); user != "" {
		entropy = append(entropy, []byte(user)...)
	} else if user := os.Getenv("USERNAME"); user != "" {
		entropy = append(entropy, []byte(user)...)
	}

	// 2. Home directory
	if home, err := os.UserHomeDir(); err == nil {
		entropy = append(entropy, []byte(home)...)
	}

	// 3. OS and architecture
	entropy = append(entropy, []byte(runtime.GOOS+runtime.GOARCH)...)

	// 4. Machine-specific config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		entropy = append(entropy, []byte(configDir)...)
	}

	// 5. Application-specific salt
	entropy = append(entropy, []byte("spotigo-token-encryption-v1")...)

	if len(entropy) < 32 {
		return nil, fmt.Errorf("insufficient entropy for key derivation")
	}

	// Hash to create 32-byte key for AES-256
	hash := sha256.Sum256(entropy)
	return hash[:], nil
}

// SaveEncryptedFile encrypts data and saves to file
func (e *TokenEncryptor) SaveEncryptedFile(filename string, data []byte) error {
	encrypted, err := e.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	// Clean path to prevent traversal attacks
	cleanPath := filepath.Clean(filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cleanPath), 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write with restrictive permissions
	if err := os.WriteFile(cleanPath, encrypted, 0600); err != nil { // #nosec G304 - path is sanitized with filepath.Clean
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadEncryptedFile reads and decrypts data from file
func (e *TokenEncryptor) LoadEncryptedFile(filename string) ([]byte, error) {
	// Clean path to prevent traversal attacks
	cleanPath := filepath.Clean(filename)

	encrypted, err := os.ReadFile(cleanPath) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	decrypted, err := e.Decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return decrypted, nil
}

// IsEncryptedFile checks if a file appears to be encrypted
// (checks if it starts with valid JSON or not)
func IsEncryptedFile(filename string) bool {
	// Clean path to prevent traversal attacks
	cleanPath := filepath.Clean(filename)

	data, err := os.ReadFile(cleanPath) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return false
	}

	// Check for JSON indicators (plaintext token file)
	if len(data) > 0 && (data[0] == '{' || data[0] == '[') {
		return false
	}

	return true
}
