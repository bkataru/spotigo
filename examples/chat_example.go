// chat_example.go demonstrates how to use the chat functionality of Spotigo
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/ollama"
)

func main() {
	fmt.Println("=== Spotigo Chat Example ===")
	
	// Create Ollama client (requires Ollama running locally)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)
	
	// Check if Ollama is available
	ctx := context.Background()
	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Println("Warning: Ollama not available - chat functionality will be limited")
		fmt.Println("To use chat features, start Ollama: ollama serve")
		fmt.Println("Then pull required models: ollama pull granite4:3b")
		return
	}
	
	fmt.Println("Ollama connection successful")

	// Demonstrate chat functionality
	fmt.Println("\n=== Chat Example ===")
	fmt.Println("Spotigo chat uses local LLMs for privacy-first conversations about your music.")
	
	// Example system prompt (similar to what's used in the CLI)
	systemPrompt := `You are Spotigo, a friendly music intelligence assistant. You help users explore their Spotify library, discover insights about their listening habits, and find new music. Be conversational, knowledgeable about music, and helpful. Always be concise but informative.`

	// Example conversation
	messages := []ollama.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "What are some good workout songs?"},
	}

	fmt.Println("\nUser: What are some good workout songs?")
	
	// Simulate chat response
	response, err := ollamaClient.Chat(ctx, ollama.ChatRequest{
		Model:    "granite4:3b",
		Messages: messages,
		Options: ollama.Options{
			Temperature: 0.7,
			NumPredict:  200,
		},
	})
	
	if err != nil {
		log.Printf("Chat failed: %v", err)
		fmt.Println("In a real implementation, this would connect to your local Ollama instance.")
	} else {
		fmt.Printf("Spotigo: %s\n", response.Message.Content)
	}

	// Show how context is maintained
	fmt.Println("\n=== Context Management ===")
	fmt.Println("Spotigo maintains conversation context for coherent discussions.")
	fmt.Println("Each message is added to the conversation history.")
	fmt.Println("Context is automatically trimmed to stay within model limits.")

	// Example of adding to conversation
	messages = append(messages, response.Message)
	messages = append(messages, ollama.Message{
		Role: "user", 
		Content: "Can you recommend similar artists to those?",
	})

	fmt.Println("\nUser: Can you recommend similar artists to those?")
	
	// In a real implementation, this would generate another response
	fmt.Println("Spotigo: Based on energetic rock and pop artists, you might enjoy bands like Imagine Dragons, Foo Fighters, or Paramore. These artists create music with driving rhythms perfect for workouts.")

	fmt.Println("\n=== Model Configuration ===")
	fmt.Println("Spotigo supports different model tiers:")
	fmt.Println("- Small (fast): granite4:3b - Quick responses, lower-end hardware")
	fmt.Println("- Medium (balanced): granite4:7b - Default, good balance of speed and quality")
	fmt.Println("- Large (quality): granite4:70b or qwen2.5:72b - Best quality for complex analysis")

	fmt.Println("\n=== Using Chat in Spotigo ===")
	fmt.Println("Use the CLI for interactive chat:")
	fmt.Println("$ spotigo chat           # One-shot chat")
	fmt.Println("$ spotigo chat --tui     # Interactive TUI mode")
	fmt.Println("$ spotigo --tui          # Launch TUI with menu")
}