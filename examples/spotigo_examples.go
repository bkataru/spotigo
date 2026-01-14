// Package examples provides example code demonstrating how to use the Spotigo library.
package examples

import (
	"context"
	"time"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/rag"
	"github.com/bkataru/spotigo/internal/spotify"
	"github.com/bkataru/spotigo/internal/storage"
)

// BackupExample demonstrates how to use the backup functionality of Spotigo
func BackupExample() {
	// This is just a placeholder to show the structure
	// In practice, you would use the actual backup functionality
	
	// Load configuration
	cfg := &config.Config{
		Spotify: config.SpotifyConfig{
			ClientID:     "your_client_id_here",
			ClientSecret: "your_client_secret_here",
			RedirectURI:  "http://localhost:8888/callback",
			TokenFile:    ".spotify_token",
		},
		Storage: config.StorageConfig{
			DataDir:       "./data",
			BackupDir:     "./data/backups",
			EmbeddingsDir: "./data/embeddings",
		},
	}

	// Create Spotify client
	spotifyCfg := spotify.Config{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		TokenFile:    cfg.Spotify.TokenFile,
	}

	// Create storage
	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)
	
	// These would be used in actual implementation
	_ = spotifyCfg
	_ = store
}

// RAGExample demonstrates how to use the RAG (Retrieval-Augmented Generation) functionality of Spotigo
func RAGExample() {
	// Create Ollama client (requires Ollama running locally)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)
	
	// Create RAG store
	store := rag.NewStore(ollamaClient, "nomic-embed-text", "./data/embeddings.json")
	
	// Create sample documents
	documents := []rag.Document{
		{
			ID:      "track_1",
			Type:    "track",
			Content: "Bohemian Rhapsody by Queen",
			Metadata: map[string]string{
				"artist": "Queen",
				"album":  "A Night at the Opera",
				"genre":  "Rock",
				"year":   "1975",
			},
		},
	}

	// Add documents to store (this would generate embeddings)
	ctx := context.Background()
	for _, doc := range documents {
		_ = store.Add(ctx, doc)
	}
	
	// These would be used in actual implementation
	_ = ollamaClient
	_ = store
}

// ChatExample demonstrates how to use the chat functionality of Spotigo
func ChatExample() {
	// Create Ollama client (requires Ollama running locally)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)
	
	// Example system prompt (similar to what's used in the CLI)
	systemPrompt := `You are Spotigo, a friendly music intelligence assistant. You help users explore their Spotify library, discover insights about their listening habits, and find new music. Be conversational, knowledgeable about music, and helpful. Always be concise but informative.`

	// Example conversation
	messages := []ollama.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "What are some good workout songs?"},
	}

	// These would be used in actual implementation
	_ = ollamaClient
	_ = messages
}