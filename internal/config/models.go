// Package config handles model configuration
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ModelConfig holds AI model configuration
type ModelConfig struct {
	Models   ModelsSection   `yaml:"models"`
	Agents   AgentsSection   `yaml:"agents"`
	Strategy StrategySection `yaml:"strategy"`
	Ollama   OllamaSection   `yaml:"ollama"`
}

// ModelsSection defines available model roles
type ModelsSection struct {
	Chat       ModelRole `yaml:"chat"`
	Fast       ModelRole `yaml:"fast"`
	Reasoning  ModelRole `yaml:"reasoning"`
	Tools      ModelRole `yaml:"tools"`
	Embeddings ModelRole `yaml:"embeddings"`
}

// ModelRole defines a model and its fallback
type ModelRole struct {
	Primary     string  `yaml:"primary"`
	Fallback    string  `yaml:"fallback"`
	Description string  `yaml:"description"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
	Dimensions  int     `yaml:"dimensions,omitempty"` // For embeddings
}

// AgentsSection defines agent configurations
type AgentsSection map[string]AgentConfig

// AgentConfig defines an individual agent
type AgentConfig struct {
	ModelRole    string `yaml:"model_role"`
	SystemPrompt string `yaml:"system_prompt"`
}

// StrategySection defines model selection strategy
type StrategySection struct {
	Routing             string  `yaml:"routing"`
	EscalationThreshold float64 `yaml:"escalation_threshold"`
	MaxRetries          int     `yaml:"max_retries"`
	Timeout             int     `yaml:"timeout"`
}

// OllamaSection defines Ollama connection settings
type OllamaSection struct {
	Host string `yaml:"host"`
}

// LoadModelConfig loads model configuration from YAML
func LoadModelConfig(configDir string) (*ModelConfig, error) {
	// Clean path to prevent traversal attacks
	cleanConfigDir := filepath.Clean(configDir)
	configPath := filepath.Join(cleanConfigDir, "models.yaml")

	data, err := os.ReadFile(configPath) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return nil, fmt.Errorf("failed to read models.yaml: %w", err)
	}

	var cfg ModelConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse models.yaml: %w", err)
	}

	return &cfg, nil
}

// getRoleModel returns the ModelRole for a given role name
func (m *ModelConfig) getRoleModel(role string) (*ModelRole, error) {
	switch role {
	case "chat":
		return &m.Models.Chat, nil
	case "fast":
		return &m.Models.Fast, nil
	case "reasoning":
		return &m.Models.Reasoning, nil
	case "tools":
		return &m.Models.Tools, nil
	case "embeddings":
		return &m.Models.Embeddings, nil
	default:
		return nil, fmt.Errorf("unknown model role: %s", role)
	}
}

// GetModelForRole returns the primary model for a given role
func (m *ModelConfig) GetModelForRole(role string) (string, error) {
	mr, err := m.getRoleModel(role)
	if err != nil {
		return "", err
	}
	return mr.Primary, nil
}

// GetFallbackForRole returns the fallback model for a given role
func (m *ModelConfig) GetFallbackForRole(role string) (string, error) {
	mr, err := m.getRoleModel(role)
	if err != nil {
		return "", err
	}
	return mr.Fallback, nil
}
