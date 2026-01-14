// Package config handles application configuration
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Spotify API configuration
	Spotify SpotifyConfig `mapstructure:"spotify"`

	// Ollama configuration
	Ollama OllamaConfig `mapstructure:"ollama"`

	// Storage configuration
	Storage StorageConfig `mapstructure:"storage"`

	// Backup configuration
	Backup BackupConfig `mapstructure:"backup"`

	// App settings
	App AppConfig `mapstructure:"app"`
}

// SpotifyConfig holds Spotify API credentials
type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
	TokenFile    string `mapstructure:"token_file"`
}

// OllamaConfig holds Ollama connection settings
type OllamaConfig struct {
	Host    string `mapstructure:"host"`
	Timeout int    `mapstructure:"timeout"`
}

// StorageConfig holds data storage settings
type StorageConfig struct {
	DataDir       string `mapstructure:"data_dir"`
	BackupDir     string `mapstructure:"backup_dir"`
	EmbeddingsDir string `mapstructure:"embeddings_dir"`
}

// BackupConfig holds backup settings
type BackupConfig struct {
	Schedule   string `mapstructure:"schedule"`
	RetainDays int    `mapstructure:"retain_days"`
	Format     string `mapstructure:"format"`
}

// AppConfig holds general app settings
type AppConfig struct {
	Verbose bool   `mapstructure:"verbose"`
	Theme   string `mapstructure:"theme"`
}

// Load reads configuration from file and environment
func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for config in home directory and current directory
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetConfigName("spotigo")
		viper.SetConfigType("yaml")
	}

	// Set defaults
	setDefaults()

	// Read environment variables
	viper.SetEnvPrefix("SPOTIGO")
	viper.AutomaticEnv()

	// Map environment variables to config keys
	if err := viper.BindEnv("spotify.client_id", "SPOTIFY_CLIENT_ID", "SPOTIFY_ID"); err != nil {
		return nil, fmt.Errorf("failed to bind env spotify.client_id: %w", err)
	}
	if err := viper.BindEnv("spotify.client_secret", "SPOTIFY_CLIENT_SECRET", "SPOTIFY_SECRET"); err != nil {
		return nil, fmt.Errorf("failed to bind env spotify.client_secret: %w", err)
	}
	if err := viper.BindEnv("ollama.host", "OLLAMA_HOST"); err != nil {
		return nil, fmt.Errorf("failed to bind env ollama.host: %w", err)
	}
	if err := viper.BindEnv("storage.data_dir", "SPOTIGO_DATA_DIR"); err != nil {
		return nil, fmt.Errorf("failed to bind env storage.data_dir: %w", err)
	}

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundErr) {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// Config file not found is OK, use defaults + env
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	// Ensure directories exist
	if err := ensureDirectories(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults() {
	// Spotify defaults
	viper.SetDefault("spotify.redirect_uri", "http://127.0.0.1:8888/callback")
	viper.SetDefault("spotify.token_file", ".spotify_token")

	// Ollama defaults
	viper.SetDefault("ollama.host", "http://localhost:11434")
	viper.SetDefault("ollama.timeout", 30)

	// Storage defaults
	viper.SetDefault("storage.data_dir", "./data")
	viper.SetDefault("storage.backup_dir", "./data/backups")
	viper.SetDefault("storage.embeddings_dir", "./data/embeddings")

	// Backup defaults
	viper.SetDefault("backup.schedule", "daily")
	viper.SetDefault("backup.retain_days", 30)
	viper.SetDefault("backup.format", "json")

	// App defaults
	viper.SetDefault("app.verbose", false)
	viper.SetDefault("app.theme", "dark")
}

func ensureDirectories(cfg *Config) error {
	dirs := []string{
		cfg.Storage.DataDir,
		cfg.Storage.BackupDir,
		cfg.Storage.EmbeddingsDir,
	}

	for _, dir := range dirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("invalid directory path %s: %w", dir, err)
		}
		if err := os.MkdirAll(absDir, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absDir, err)
		}
	}

	return nil
}

// GetConfigPath returns the path of the loaded config file
func GetConfigPath() string {
	return viper.ConfigFileUsed()
}
