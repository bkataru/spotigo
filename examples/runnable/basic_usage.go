// Basic usage example for the Spotigo library
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/rag"
	"github.com/bkataru/spotigo/internal/spotify"
)

func main() {
	// Example 1: Initialize Spotify client
	fmt.Println("=== Spotify Client Example ===")
	spotifyClient, err := spotify.NewClient(spotify.Config{
		ClientID:     "your_client_id_here",
		ClientSecret: "your_client_secret_here",
		RedirectURI:  "http://127.0.0.1:8888/callback",
		TokenFile:    ".spotify_token",
	})
	if err != nil {
		log.Printf("Failed to create Spotify client: %v", err)
	} else {
		fmt.Println("Spotify client created successfully")
		if spotifyClient.IsAuthenticated() {
			fmt.Println("Client is authenticated")
		} else {
			fmt.Println("Client needs authentication - use GetAuthURL() and HandleCallback()")
		}
	}

	// Example 2: Initialize Ollama client
	fmt.Println("\n=== Ollama Client Example ===")
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)
	fmt.Println("Ollama client created successfully")

	// Example 3: Initialize RAG store
	fmt.Println("\n=== RAG Store Example ===")
	store := rag.NewStore(ollamaClient, "nomic-embed-text-v2-moe", "./data/store.json")
	fmt.Println("RAG store created successfully")

	// Example 4: Create sample documents
	fmt.Println("\n=== Document Creation Example ===")
	documents := []rag.Document{
		{
			ID:      "track_1",
			Type:    "track",
			Content: "The Beatles - Hey Jude",
			Metadata: map[string]string{
				"artist": "The Beatles",
				"genre":  "rock",
				"year":   "1968",
			},
		},
		{
			ID:      "track_2",
			Type:    "track",
			Content: "Radiohead - Paranoid Android",
			Metadata: map[string]string{
				"artist": "Radiohead",
				"genre":  "alternative",
				"year":   "1997",
			},
		},
	}

	fmt.Printf("Created %d sample documents\n", len(documents))
	_ = store // Store is initialized but not used in this example

	// Example 5: Show how to use chat (requires Ollama running)
	fmt.Println("\n=== Ollama Chat Example (conceptual) ===")
	fmt.Println("To use chat functionality:")
	fmt.Println("1. Ensure Ollama is running: ollama serve")
	fmt.Println("2. Pull a model: ollama pull granite4:1b")
	fmt.Println("3. Use client.Chat() method")

	fmt.Println("\n=== Library Ready for Use ===")
	fmt.Println("The Spotigo library provides:")
	fmt.Println("- spotify.Client: Spotify API integration with OAuth")
	fmt.Println("- ollama.Client: Local LLM inference")
	fmt.Println("- rag.Store: Vector-based semantic search")
	fmt.Println("- internal/storage: Local file persistence")
	fmt.Println("- internal/config: Configuration management")
	fmt.Println("- internal/crypto: Secure token encryption")
}
