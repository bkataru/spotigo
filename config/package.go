package config

import "os"

type SpotifyConfig struct {
	ID string
	SECRET string
}

type Config struct {
	SpotifyConfig
}

func getEnv(key string) (string, bool) {
	if value, exists := os.LookupEnv(key); exists {
		return value, true
	}

	return defaultVal, false
}

func New() *Config {
	return &Config{
		ID:
	}
}
