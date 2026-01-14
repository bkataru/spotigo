//go:build integration
// +build integration

package integration

import (
	"testing"

	"spotigo/internal/spotify"
)

func TestAuth_ClientCreation(t *testing.T) {
	// Create Spotify client with test credentials
	cfg := spotify.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
		TokenFile:    "",
	}

	client, err := spotify.NewClient(cfg)
	AssertNoError(t, err, "Failed to create Spotify client")

	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	// Client should not be authenticated without token
	if client.IsAuthenticated() {
		t.Error("Expected client to not be authenticated without token")
	}
}

func TestAuth_GetAuthURL(t *testing.T) {
	cfg := spotify.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
		TokenFile:    "",
	}

	client, err := spotify.NewClient(cfg)
	AssertNoError(t, err, "Failed to create Spotify client")

	// Get auth URL
	state := "test-state-12345"
	authURL := client.GetAuthURL(state)

	if authURL == "" {
		t.Error("Expected non-empty auth URL")
	}

	// Verify URL contains required parameters
	if !contains(authURL, "client_id=test-client-id") {
		t.Error("Auth URL missing client_id parameter")
	}

	if !contains(authURL, "response_type=code") {
		t.Error("Auth URL missing response_type parameter")
	}

	if !contains(authURL, "redirect_uri=") {
		t.Error("Auth URL missing redirect_uri parameter")
	}

	if !contains(authURL, "state="+state) {
		t.Error("Auth URL missing state parameter")
	}

	t.Logf("Auth URL: %s", authURL)
}

func TestAuth_TokenEncryption(t *testing.T) {
	testCfg := SetupTestEnvironment(t)

	// Test data to encrypt
	testToken := []byte(`{
		"access_token": "test_access_token",
		"token_type": "Bearer",
		"refresh_token": "test_refresh_token",
		"expiry": "2024-12-31T23:59:59Z"
	}`)

	// Encrypt
	encrypted, err := testCfg.TokenEncryptor.Encrypt(testToken)
	AssertNoError(t, err, "Failed to encrypt token")

	if len(encrypted) == 0 {
		t.Fatal("Expected non-empty encrypted data")
	}

	// Encrypted data should be different from plaintext
	if string(encrypted) == string(testToken) {
		t.Error("Encrypted data should differ from plaintext")
	}

	// Decrypt
	decrypted, err := testCfg.TokenEncryptor.Decrypt(encrypted)
	AssertNoError(t, err, "Failed to decrypt token")

	// Decrypted should match original
	if string(decrypted) != string(testToken) {
		t.Error("Decrypted data doesn't match original")
	}
}

func TestAuth_TokenEncryptionBase64(t *testing.T) {
	testCfg := SetupTestEnvironment(t)

	testData := []byte("secret-spotify-token-data")

	// Encrypt to base64
	encoded, err := testCfg.TokenEncryptor.EncryptToBase64(testData)
	AssertNoError(t, err, "Failed to encrypt to base64")

	if encoded == "" {
		t.Fatal("Expected non-empty base64 string")
	}

	// Decrypt from base64
	decoded, err := testCfg.TokenEncryptor.DecryptFromBase64(encoded)
	AssertNoError(t, err, "Failed to decrypt from base64")

	if string(decoded) != string(testData) {
		t.Error("Decoded data doesn't match original")
	}
}

func TestAuth_TokenFileOperations(t *testing.T) {
	testCfg := SetupTestEnvironment(t)

	testData := []byte("test-token-data-for-file")
	tokenFile := testCfg.TempDir + "/test_token.enc"

	// Save encrypted file
	err := testCfg.TokenEncryptor.SaveEncryptedFile(tokenFile, testData)
	AssertNoError(t, err, "Failed to save encrypted file")

	// Verify file exists
	AssertFileExists(t, tokenFile)

	// Load encrypted file
	loaded, err := testCfg.TokenEncryptor.LoadEncryptedFile(tokenFile)
	AssertNoError(t, err, "Failed to load encrypted file")

	if string(loaded) != string(testData) {
		t.Error("Loaded data doesn't match original")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
