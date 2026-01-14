package crypto

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTokenEncryptor(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}
	if enc == nil {
		t.Fatal("encryptor is nil")
	}
	if len(enc.key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(enc.key))
	}
}

func TestNewTokenEncryptorWithKey(t *testing.T) {
	// Test with short key (should fail)
	_, err := NewTokenEncryptorWithKey([]byte("short"))
	if err == nil {
		t.Error("expected error for short key")
	}

	// Test with valid key
	key := []byte("this-is-a-test-key-for-spotigo!")
	enc, err := NewTokenEncryptorWithKey(key)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}
	if enc == nil {
		t.Fatal("encryptor is nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"short", []byte("hello")},
		{"json-like", []byte(`{"access_token":"abc123","refresh_token":"xyz789"}`)},
		{"long", bytes.Repeat([]byte("a"), 10000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := enc.Encrypt(tc.data)
			if err != nil {
				t.Fatalf("encrypt failed: %v", err)
			}

			// Encrypted should be different from original (unless empty)
			if len(tc.data) > 0 && bytes.Equal(encrypted, tc.data) {
				t.Error("encrypted data equals plaintext")
			}

			decrypted, err := enc.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("decrypt failed: %v", err)
			}

			if !bytes.Equal(decrypted, tc.data) {
				t.Errorf("decrypted data doesn't match original")
			}
		})
	}
}

func TestEncryptDecryptBase64(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	original := []byte(`{"token":"test123"}`)

	encoded, err := enc.EncryptToBase64(original)
	if err != nil {
		t.Fatalf("encrypt to base64 failed: %v", err)
	}

	// Should be valid base64
	if encoded == "" {
		t.Error("encoded string is empty")
	}

	decoded, err := enc.DecryptFromBase64(encoded)
	if err != nil {
		t.Fatalf("decrypt from base64 failed: %v", err)
	}

	if !bytes.Equal(decoded, original) {
		t.Errorf("decoded data doesn't match original")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	// Test with too short data
	_, err = enc.Decrypt([]byte("short"))
	if err == nil {
		t.Error("expected error for short ciphertext")
	}

	// Test with invalid ciphertext
	_, err = enc.Decrypt(bytes.Repeat([]byte("x"), 50))
	if err == nil {
		t.Error("expected error for invalid ciphertext")
	}
}

func TestSaveLoadEncryptedFile(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "spotigo-crypto-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test_token.enc")
	original := []byte(`{"access_token":"abc123","refresh_token":"xyz789","expiry":"2025-01-01"}`)

	// Save
	if err := enc.SaveEncryptedFile(testFile, original); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Verify file exists and has correct permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	// Check file content is not plaintext
	rawContent, _ := os.ReadFile(testFile)
	if bytes.Contains(rawContent, []byte("access_token")) {
		t.Error("file contains plaintext - not encrypted")
	}

	// Load
	loaded, err := enc.LoadEncryptedFile(testFile)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !bytes.Equal(loaded, original) {
		t.Errorf("loaded data doesn't match original")
	}

	// Check file size is reasonable
	if info.Size() < int64(len(original)) {
		t.Errorf("encrypted file seems too small: %d bytes", info.Size())
	}
}

func TestIsEncryptedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-crypto-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test plaintext JSON file
	plaintextFile := filepath.Join(tmpDir, "plaintext.json")
	os.WriteFile(plaintextFile, []byte(`{"token":"test"}`), 0600)
	if IsEncryptedFile(plaintextFile) {
		t.Error("plaintext JSON incorrectly identified as encrypted")
	}

	// Test encrypted file
	enc, _ := NewTokenEncryptor()
	encryptedFile := filepath.Join(tmpDir, "encrypted.enc")
	enc.SaveEncryptedFile(encryptedFile, []byte(`{"token":"test"}`))
	if !IsEncryptedFile(encryptedFile) {
		t.Error("encrypted file incorrectly identified as plaintext")
	}

	// Test non-existent file
	if IsEncryptedFile(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("non-existent file incorrectly identified as encrypted")
	}
}

func TestDifferentEncryptorsWithSameKey(t *testing.T) {
	key := []byte("shared-encryption-key-for-test!!")

	enc1, err := NewTokenEncryptorWithKey(key)
	if err != nil {
		t.Fatalf("failed to create encryptor 1: %v", err)
	}

	enc2, err := NewTokenEncryptorWithKey(key)
	if err != nil {
		t.Fatalf("failed to create encryptor 2: %v", err)
	}

	original := []byte("test data for cross-encryptor test")

	// Encrypt with enc1
	encrypted, err := enc1.Encrypt(original)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Decrypt with enc2
	decrypted, err := enc2.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, original) {
		t.Error("cross-encryptor decryption failed")
	}
}

func TestEncryptionIsNonDeterministic(t *testing.T) {
	enc, err := NewTokenEncryptor()
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	original := []byte("same data")

	// Encrypt the same data twice
	encrypted1, _ := enc.Encrypt(original)
	encrypted2, _ := enc.Encrypt(original)

	// Due to random nonce, encrypted results should differ
	if bytes.Equal(encrypted1, encrypted2) {
		t.Error("encrypting same data twice produced identical ciphertext (nonce not random?)")
	}

	// But both should decrypt to same value
	decrypted1, _ := enc.Decrypt(encrypted1)
	decrypted2, _ := enc.Decrypt(encrypted2)

	if !bytes.Equal(decrypted1, decrypted2) {
		t.Error("different ciphertexts decrypted to different values")
	}
}
