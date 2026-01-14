//go:build integration
// +build integration

// Package integration provides integration test utilities for Spotigo.
// Run with: go test -tags=integration ./test/integration/...
package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"spotigo/internal/crypto"
	"spotigo/internal/ollama"
	"spotigo/internal/rag"
	"spotigo/internal/spotify"
	"spotigo/internal/storage"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	TempDir        string
	SpotifyClient  *spotify.Client
	OllamaClient   *ollama.Client
	Storage        *storage.Store
	RAGStore       *rag.Store
	TokenEncryptor *crypto.TokenEncryptor
}

// SetupTestEnvironment creates a complete test environment
func SetupTestEnvironment(t *testing.T) *TestConfig {
	t.Helper()

	// Create temporary directory for test data
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")

	// Initialize crypto
	encryptor, err := crypto.NewTokenEncryptor()
	if err != nil {
		t.Fatalf("Failed to create token encryptor: %v", err)
	}

	// Initialize storage
	storageStore := storage.NewStore(tempDir, backupDir)

	// Initialize Ollama client (localhost default)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)

	// Initialize RAG store
	ragStorePath := filepath.Join(tempDir, "embeddings.json")
	ragStore := rag.NewStore(ollamaClient, "nomic-embed-text", ragStorePath)

	return &TestConfig{
		TempDir:        tempDir,
		OllamaClient:   ollamaClient,
		Storage:        storageStore,
		RAGStore:       ragStore,
		TokenEncryptor: encryptor,
	}
}

// SetupMockSpotifyClient creates a mock Spotify client for testing
func SetupMockSpotifyClient(t *testing.T) *spotify.Client {
	t.Helper()

	// Note: This would require modifying the Spotify client to accept a custom HTTP client
	// For now, we'll skip actual Spotify client creation in tests
	// Real integration tests would need actual Spotify API credentials

	return nil // Placeholder - tests should mock Spotify API calls
}

// MockSpotifyServer creates a mock HTTP server for Spotify API
func MockSpotifyServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	// Mock OAuth callback endpoint
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Authorization successful")
	})

	// Mock saved tracks endpoint
	mux.HandleFunc("/v1/me/tracks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"items": [
				{
					"track": {
						"id": "track1",
						"name": "Test Track 1",
						"artists": [{"id": "artist1", "name": "Test Artist"}],
						"album": {"id": "album1", "name": "Test Album"}
					}
				}
			],
			"total": 1,
			"limit": 50,
			"offset": 0
		}`
		fmt.Fprint(w, response)
	})

	// Mock playlists endpoint
	mux.HandleFunc("/v1/me/playlists", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"items": [
				{
					"id": "playlist1",
					"name": "Test Playlist",
					"tracks": {"total": 10}
				}
			],
			"total": 1
		}`
		fmt.Fprint(w, response)
	})

	// Mock user profile endpoint
	mux.HandleFunc("/v1/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"id": "testuser",
			"display_name": "Test User",
			"email": "test@example.com"
		}`
		fmt.Fprint(w, response)
	})

	return httptest.NewServer(mux)
}

// CheckOllamaAvailable checks if Ollama is running and accessible
func CheckOllamaAvailable(t *testing.T) bool {
	t.Helper()

	client := ollama.NewClient("http://localhost:11434", 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.ListModels(ctx)
	return err == nil
}

// SkipIfNoOllama skips the test if Ollama is not available
func SkipIfNoOllama(t *testing.T) {
	t.Helper()

	if !CheckOllamaAvailable(t) {
		t.Skip("Ollama not available, skipping integration test")
	}
}

// CheckModelAvailable checks if a specific Ollama model is available
func CheckModelAvailable(t *testing.T, modelName string) bool {
	t.Helper()

	client := ollama.NewClient("http://localhost:11434", 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	models, err := client.ListModels(ctx)
	if err != nil {
		return false
	}

	// ListModels returns []ModelInfo
	for _, model := range models {
		if model.Name == modelName {
			return true
		}
	}

	return false
}

// SkipIfNoModel skips the test if the specified model is not available
func SkipIfNoModel(t *testing.T, modelName string) {
	t.Helper()

	SkipIfNoOllama(t)

	if !CheckModelAvailable(t, modelName) {
		t.Skipf("Model %s not available, skipping integration test", modelName)
	}
}

// CreateTestDocuments generates test documents for RAG testing
func CreateTestDocuments(count int) []rag.Document {
	docs := make([]rag.Document, count)
	for i := 0; i < count; i++ {
		docs[i] = rag.Document{
			ID:      fmt.Sprintf("track-%d", i),
			Type:    "track",
			Content: fmt.Sprintf("Test Track %d by Artist %d from Album %d", i, i%10, i%5),
			Metadata: map[string]string{
				"name":   fmt.Sprintf("Track %d", i),
				"artist": fmt.Sprintf("Artist %d", i%10),
				"album":  fmt.Sprintf("Album %d", i%5),
			},
		}
	}
	return docs
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("Expected file not to exist: %s", path)
	}
}

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()

	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// CleanupTestData removes test data files
func CleanupTestData(t *testing.T, paths ...string) {
	t.Helper()

	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			t.Logf("Warning: failed to cleanup %s: %v", path, err)
		}
	}
}
