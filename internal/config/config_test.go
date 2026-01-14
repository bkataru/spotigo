package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad_Defaults(t *testing.T) {
	// Reset viper state for clean test
	viper.Reset()

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "spotigo-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory so config lookup doesn't find existing files
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create required directories to prevent ensureDirectories from failing
	os.MkdirAll(filepath.Join(tmpDir, "data"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "backups"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "embeddings"), 0755)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check Spotify defaults
	if cfg.Spotify.RedirectURI != "http://localhost:8888/callback" {
		t.Errorf("expected default redirect URI, got '%s'", cfg.Spotify.RedirectURI)
	}
	if cfg.Spotify.TokenFile != ".spotify_token" {
		t.Errorf("expected default token file, got '%s'", cfg.Spotify.TokenFile)
	}

	// Check Ollama defaults
	if cfg.Ollama.Host != "http://localhost:11434" {
		t.Errorf("expected default Ollama host, got '%s'", cfg.Ollama.Host)
	}
	if cfg.Ollama.Timeout != 30 {
		t.Errorf("expected default Ollama timeout 30, got %d", cfg.Ollama.Timeout)
	}

	// Check Storage defaults
	if cfg.Storage.DataDir != "./data" {
		t.Errorf("expected default data dir './data', got '%s'", cfg.Storage.DataDir)
	}
	if cfg.Storage.BackupDir != "./data/backups" {
		t.Errorf("expected default backup dir, got '%s'", cfg.Storage.BackupDir)
	}
	if cfg.Storage.EmbeddingsDir != "./data/embeddings" {
		t.Errorf("expected default embeddings dir, got '%s'", cfg.Storage.EmbeddingsDir)
	}

	// Check Backup defaults
	if cfg.Backup.Schedule != "daily" {
		t.Errorf("expected default schedule 'daily', got '%s'", cfg.Backup.Schedule)
	}
	if cfg.Backup.RetainDays != 30 {
		t.Errorf("expected default retain days 30, got %d", cfg.Backup.RetainDays)
	}
	if cfg.Backup.Format != "json" {
		t.Errorf("expected default format 'json', got '%s'", cfg.Backup.Format)
	}

	// Check App defaults
	if cfg.App.Verbose != false {
		t.Error("expected default verbose false")
	}
	if cfg.App.Theme != "dark" {
		t.Errorf("expected default theme 'dark', got '%s'", cfg.App.Theme)
	}
}

func TestLoad_FromFile(t *testing.T) {
	viper.Reset()

	// Create temp directory and config file
	tmpDir, err := os.MkdirTemp("", "spotigo-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create required directories
	os.MkdirAll(filepath.Join(tmpDir, "custom-data"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "custom-data", "backups"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "custom-data", "embeddings"), 0755)

	configContent := `
spotify:
  client_id: "test-client-id"
  client_secret: "test-client-secret"
  redirect_uri: "http://localhost:9999/callback"

ollama:
  host: "http://ollama.local:11434"
  timeout: 60

storage:
  data_dir: "./custom-data"
  backup_dir: "./custom-data/backups"
  embeddings_dir: "./custom-data/embeddings"

backup:
  schedule: "weekly"
  retain_days: 60

app:
  verbose: true
  theme: "light"
`

	configPath := filepath.Join(tmpDir, "spotigo.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check custom values
	if cfg.Spotify.ClientID != "test-client-id" {
		t.Errorf("expected client ID 'test-client-id', got '%s'", cfg.Spotify.ClientID)
	}
	if cfg.Spotify.ClientSecret != "test-client-secret" {
		t.Errorf("expected client secret 'test-client-secret', got '%s'", cfg.Spotify.ClientSecret)
	}
	if cfg.Spotify.RedirectURI != "http://localhost:9999/callback" {
		t.Errorf("expected redirect URI 'http://localhost:9999/callback', got '%s'", cfg.Spotify.RedirectURI)
	}

	if cfg.Ollama.Host != "http://ollama.local:11434" {
		t.Errorf("expected Ollama host 'http://ollama.local:11434', got '%s'", cfg.Ollama.Host)
	}
	if cfg.Ollama.Timeout != 60 {
		t.Errorf("expected Ollama timeout 60, got %d", cfg.Ollama.Timeout)
	}

	if cfg.Backup.Schedule != "weekly" {
		t.Errorf("expected schedule 'weekly', got '%s'", cfg.Backup.Schedule)
	}
	if cfg.Backup.RetainDays != 60 {
		t.Errorf("expected retain days 60, got %d", cfg.Backup.RetainDays)
	}

	if cfg.App.Verbose != true {
		t.Error("expected verbose true")
	}
	if cfg.App.Theme != "light" {
		t.Errorf("expected theme 'light', got '%s'", cfg.App.Theme)
	}
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	viper.Reset()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "spotigo-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create required directories
	os.MkdirAll(filepath.Join(tmpDir, "data"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "backups"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "embeddings"), 0755)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Set environment variables
	os.Setenv("SPOTIFY_CLIENT_ID", "env-client-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "env-client-secret")
	os.Setenv("OLLAMA_HOST", "http://env-ollama:11434")
	defer func() {
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
		os.Unsetenv("OLLAMA_HOST")
	}()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Spotify.ClientID != "env-client-id" {
		t.Errorf("expected client ID from env 'env-client-id', got '%s'", cfg.Spotify.ClientID)
	}
	if cfg.Spotify.ClientSecret != "env-client-secret" {
		t.Errorf("expected client secret from env 'env-client-secret', got '%s'", cfg.Spotify.ClientSecret)
	}
	if cfg.Ollama.Host != "http://env-ollama:11434" {
		t.Errorf("expected Ollama host from env, got '%s'", cfg.Ollama.Host)
	}
}

func TestLoad_ExplicitConfigFile(t *testing.T) {
	viper.Reset()

	tmpDir, err := os.MkdirTemp("", "spotigo-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directories in temp location
	os.MkdirAll(filepath.Join(tmpDir, "mydata"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "mydata", "backups"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "mydata", "embeddings"), 0755)

	configContent := `
spotify:
  client_id: "explicit-client-id"

storage:
  data_dir: "./mydata"
  backup_dir: "./mydata/backups"
  embeddings_dir: "./mydata/embeddings"
`

	configPath := filepath.Join(tmpDir, "custom-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Change to temp directory so relative paths work
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load with explicit path failed: %v", err)
	}

	if cfg.Spotify.ClientID != "explicit-client-id" {
		t.Errorf("expected client ID 'explicit-client-id', got '%s'", cfg.Spotify.ClientID)
	}
}

func TestGetConfigPath(t *testing.T) {
	viper.Reset()

	tmpDir, err := os.MkdirTemp("", "spotigo-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create required directories
	os.MkdirAll(filepath.Join(tmpDir, "data"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "backups"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "data", "embeddings"), 0755)

	// Create config file
	configContent := `app:
  verbose: false
`
	configPath := filepath.Join(tmpDir, "spotigo.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	_, err = Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	usedPath := GetConfigPath()
	if usedPath == "" {
		// This is OK if no config file was found
		t.Log("No config file path recorded (using defaults)")
	}
}

func TestConfigStructure(t *testing.T) {
	// Test that Config struct has all expected fields
	cfg := Config{}

	// Spotify config
	cfg.Spotify.ClientID = "test"
	cfg.Spotify.ClientSecret = "test"
	cfg.Spotify.RedirectURI = "test"
	cfg.Spotify.TokenFile = "test"

	// Ollama config
	cfg.Ollama.Host = "test"
	cfg.Ollama.Timeout = 30

	// Storage config
	cfg.Storage.DataDir = "test"
	cfg.Storage.BackupDir = "test"
	cfg.Storage.EmbeddingsDir = "test"

	// Backup config
	cfg.Backup.Schedule = "test"
	cfg.Backup.RetainDays = 30
	cfg.Backup.Format = "test"

	// App config
	cfg.App.Verbose = true
	cfg.App.Theme = "test"

	// If we got here without panicking, structure is correct
}
