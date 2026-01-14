package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"spotigo/internal/config"
	"spotigo/internal/ollama"
)

var (
	chatModel   string
	chatContext int
)

func init() {
	chatCmd.Flags().StringVar(&chatModel, "model", "", "override the default chat model")
	chatCmd.Flags().IntVar(&chatContext, "context", 4096, "context window size")
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an AI chat session about your music",
	Long: `Start an interactive AI chat session to discuss your music library.

You can ask questions like:
  - "What are my most played genres?"
  - "Find songs similar to [track name]"
  - "When did I start listening to [artist]?"
  - "What's my musical taste evolution over time?"
  - "Recommend something based on my recent listening"

The AI runs entirely locally using Ollama. No data leaves your machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		runChat()
	},
}

func runChat() {
	fmt.Println("Spotigo AI Chat")
	fmt.Println("===============")
	fmt.Println("Ask me anything about your music library!")
	fmt.Println("Type 'exit' or 'quit' to end the session.")
	fmt.Println()

	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Warning: Configuration not loaded. Using defaults.")
	}

	// Initialize Ollama client
	ollamaClient := ollama.NewClient(cfg.Ollama.Host, time.Duration(cfg.Ollama.Timeout)*time.Second)

	// Test Ollama connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Printf("❌ Cannot connect to Ollama: %v\n", err)
		fmt.Println("Make sure Ollama is running:")
		fmt.Println("  ollama serve")
		fmt.Println()
		fmt.Println("Or check your Ollama configuration:")
		fmt.Printf("  Host: %s\n", cfg.Ollama.Host)
		return
	}

	// Load model configuration
	modelCfg, err := config.LoadModelConfig("./config")
	if err != nil {
		fmt.Printf("Warning: Could not load model config: %v\n", err)
		fmt.Println("Using default model: granite4:1b")
	}

	// Determine which model to use
	modelName := chatModel
	if modelName == "" {
		if modelCfg != nil {
			modelName, _ = modelCfg.GetModelForRole("chat")
		}
		if modelName == "" {
			modelName = "granite4:1b"
		}
	}

	fmt.Printf("✅ Connected to Ollama using model: %s\n", modelName)
	fmt.Println()

	// Load system prompt
	var systemPrompt string
	if modelCfg != nil {
		if agentCfg, exists := modelCfg.Agents["chat_agent"]; exists {
			systemPrompt = agentCfg.SystemPrompt
		}
	}
	if systemPrompt == "" {
		systemPrompt = `You are Spotigo, a friendly music intelligence assistant. You help users explore their Spotify library, discover insights about their listening habits, and find new music. Be conversational, knowledgeable about music, and helpful. Always be concise but informative.`
	}

	// Start chat loop
	reader := bufio.NewReader(os.Stdin)

	// Initialize conversation with system message
	messages := []ollama.Message{
		{Role: "system", Content: systemPrompt},
	}

	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Add user message to conversation
		messages = append(messages, ollama.Message{
			Role:    "user",
			Content: input,
		})

		fmt.Print("Spotigo: ")

		// Send to Ollama
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		req := ollama.ChatRequest{
			Model:    modelName,
			Messages: messages,
			Options: &ollama.Options{
				Temperature: 0.7,
				NumPredict:  chatContext,
			},
		}

		resp, err := ollamaClient.Chat(ctx, req)
		cancel()

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Println("Retrying with fallback model...")

			// Try fallback model
			fallbackModel := "qwen3:0.6b"
			if modelCfg != nil {
				fallbackModel, _ = modelCfg.GetFallbackForRole("chat")
			}

			req.Model = fallbackModel
			ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
			resp, err = ollamaClient.Chat(ctx, req)
			cancel()

			if err != nil {
				fmt.Printf("❌ Fallback also failed: %v\n", err)
				// Remove user message from history if chat failed
				messages = messages[:len(messages)-1]
				continue
			}
		}

		// Add assistant response to conversation
		messages = append(messages, resp.Message)
		fmt.Println(resp.Message.Content)
		fmt.Println()
	}
}
