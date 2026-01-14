package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"spotigo/internal/config"
	"spotigo/internal/ollama"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage AI models",
	Long: `View and manage AI models used by Spotigo.

Models are served via Ollama and run entirely locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	modelsCmd.AddCommand(modelsListCmd)
	modelsCmd.AddCommand(modelsStatusCmd)
	modelsCmd.AddCommand(modelsPullCmd)
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured models",
	Run: func(cmd *cobra.Command, args []string) {
		listModels()
	},
}

var modelsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check model availability in Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		checkModelStatus()
	},
}

var modelsPullCmd = &cobra.Command{
	Use:   "pull [model]",
	Short: "Pull a model from Ollama",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pullModel(args[0])
	},
}

func listModels() {
	// Load model configuration
	modelCfg, err := config.LoadModelConfig("./config")
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		fmt.Println("Using default configuration:")
		fmt.Println()
		fmt.Println("  Chat:")
		fmt.Println("    Primary:  granite4:1b")
		fmt.Println("    Fallback: qwen3:0.6b")
		fmt.Println()
		fmt.Println("  Fast Tasks:")
		fmt.Println("    Primary:  granite4:350m")
		fmt.Println("    Fallback: qwen3:0.6b")
		fmt.Println()
		fmt.Println("  Reasoning:")
		fmt.Println("    Primary:  qwen3:1.7b")
		fmt.Println("    Fallback: granite4:1b")
		fmt.Println()
		fmt.Println("  Embeddings:")
		fmt.Println("    Primary:  nomic-embed-text-v2-moe")
		fmt.Println("    Fallback: qwen3-embedding:0.6b")
		return
	}

	fmt.Println("Configured Models:")
	fmt.Println()

	// Chat models
	fmt.Println("  Chat:")
	fmt.Printf("    Primary:  %s\n", modelCfg.Models.Chat.Primary)
	fmt.Printf("    Fallback: %s\n", modelCfg.Models.Chat.Fallback)
	fmt.Printf("    %s\n", modelCfg.Models.Chat.Description)
	fmt.Println()

	// Fast models
	fmt.Println("  Fast Tasks:")
	fmt.Printf("    Primary:  %s\n", modelCfg.Models.Fast.Primary)
	fmt.Printf("    Fallback: %s\n", modelCfg.Models.Fast.Fallback)
	fmt.Printf("    %s\n", modelCfg.Models.Fast.Description)
	fmt.Println()

	// Reasoning models
	fmt.Println("  Reasoning:")
	fmt.Printf("    Primary:  %s\n", modelCfg.Models.Reasoning.Primary)
	fmt.Printf("    Fallback: %s\n", modelCfg.Models.Reasoning.Fallback)
	fmt.Printf("    %s\n", modelCfg.Models.Reasoning.Description)
	fmt.Println()

	// Tool models
	fmt.Println("  Tools:")
	fmt.Printf("    Primary:  %s\n", modelCfg.Models.Tools.Primary)
	fmt.Printf("    Fallback: %s\n", modelCfg.Models.Tools.Fallback)
	fmt.Printf("    %s\n", modelCfg.Models.Tools.Description)
	fmt.Println()

	// Embedding models
	fmt.Println("  Embeddings:")
	fmt.Printf("    Primary:  %s\n", modelCfg.Models.Embeddings.Primary)
	fmt.Printf("    Fallback: %s\n", modelCfg.Models.Embeddings.Fallback)
	fmt.Printf("    %s\n", modelCfg.Models.Embeddings.Description)
}

func checkModelStatus() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	fmt.Println("Checking Ollama connection...")
	fmt.Println()

	// Connect to Ollama
	client := ollama.NewClient(cfg.Ollama.Host, 10*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx); err != nil {
		fmt.Printf("❌ Ollama Status: Not connected\n")
		fmt.Printf("   Error: %v\n", err)
		fmt.Println()
		fmt.Printf("Configured Ollama Host: %s\n", cfg.Ollama.Host)
		fmt.Println()
		fmt.Println("Troubleshooting:")
		fmt.Println("1. Make sure Ollama is installed:")
		fmt.Println("   curl -fsSL https://ollama.ai/install.sh | sh")
		fmt.Println("2. Start Ollama service:")
		fmt.Println("   ollama serve")
		fmt.Println("3. Check if Ollama is accessible:")
		fmt.Printf("   curl %s/api/tags\n", cfg.Ollama.Host)
		return
	}

	fmt.Printf("✅ Ollama Status: Connected\n")
	fmt.Printf("   Host: %s\n", cfg.Ollama.Host)
	fmt.Println()

	// List available models
	models, err := client.ListModels(ctx)
	if err != nil {
		fmt.Printf("Error listing models: %v\n", err)
		return
	}

	if len(models) == 0 {
		fmt.Println("No models found in Ollama.")
		fmt.Println("Pull some models to get started:")
		fmt.Println("  ollama pull granite4:1b")
		fmt.Println("  ollama pull nomic-embed-text-v2-moe")
		return
	}

	fmt.Printf("Available Models (%d):\n", len(models))
	fmt.Println()

	// Get configured models for comparison
	modelCfg, _ := config.LoadModelConfig("./config")
	configuredModels := make(map[string]string)
	if modelCfg != nil {
		configuredModels["chat"], _ = modelCfg.GetModelForRole("chat")
		configuredModels["fast"], _ = modelCfg.GetModelForRole("fast")
		configuredModels["reasoning"], _ = modelCfg.GetModelForRole("reasoning")
		configuredModels["embeddings"], _ = modelCfg.GetModelForRole("embeddings")
	}

	for _, model := range models {
		status := "✅"
		// Check if this is a configured model
		for _, configuredName := range configuredModels {
			if model.Name == configuredName {
				status = "🎯"
				break
			}
		}

		// Get model size
		sizeMB := float64(model.Size) / 1024 / 1024
		fmt.Printf("  %s %-25s %8.1f MB  %s\n", status, model.Name, sizeMB, model.ModifiedAt.Format("2006-01-02"))
	}

	fmt.Println()
	fmt.Println("Status: 🎯 = Configured for use, ✅ = Available")
}

func pullModel(modelName string) {
	fmt.Printf("Pulling model: %s\n", modelName)
	fmt.Println("This may take a while depending on model size...")
	fmt.Println()

	// For now, we'll use the ollama CLI directly since our client doesn't support streaming pulls
	fmt.Println("Using 'ollama pull' command directly...")
	fmt.Printf("You can run: ollama pull %s\n", modelName)
	fmt.Println()
	fmt.Println("Note: Spotigo currently delegates model pulling to the ollama CLI for better progress reporting.")
	fmt.Println("Future versions will implement direct API pulling.")
}
