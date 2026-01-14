// Package examples provides reference implementations for Spotigo's core functionality.
// These examples demonstrate how to use the various components of the Spotigo library.
package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/rag"
	"github.com/bkataru/spotigo/internal/spotify"
	"github.com/bkataru/spotigo/internal/storage"
)

// BackupExample demonstrates how to use the backup functionality of Spotigo.
// This example shows the conceptual structure of backup operations without
// requiring actual Spotify API calls.
func BackupExample() {
	fmt.Println("=== Spotigo Backup Example ===")

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		fmt.Println("Using default configuration...")

		// Create default config
		cfg = &config.Config{
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
	}

	// Create Spotify client
	spotifyCfg := spotify.Config{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		TokenFile:    cfg.Spotify.TokenFile,
	}

	client, err := spotify.NewClient(spotifyCfg)
	if err != nil {
		log.Printf("Failed to create Spotify client: %v", err)
		fmt.Println("Note: Actual backup requires valid Spotify credentials and authentication")
		// Continue with the example even if client creation fails
	}

	// Create storage
	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)

	// Demonstrate backup structure (without actual API calls)
	fmt.Println("\nBackup Structure:")
	fmt.Println("- Saved Tracks: Individual track metadata")
	fmt.Println("- Playlists: Playlist metadata and track lists")
	fmt.Println("- Followed Artists: Artist information")
	fmt.Println("- Albums: Saved album information")

	// Show how to create a backup (conceptually)
	fmt.Println("\n=== Conceptual Backup Process ===")
	fmt.Println("1. Authenticate with Spotify")
	fmt.Println("2. Fetch user library data")
	fmt.Println("3. Serialize data to JSON")
	fmt.Println("4. Store in local backup directory")
	fmt.Println("5. Optionally create embeddings for semantic search")

	// Example of what backup data looks like
	fmt.Println("\n=== Sample Backup Data Structure ===")
	sampleData := map[string]interface{}{
		"saved_tracks": []map[string]interface{}{
			{
				"id":     "track123",
				"name":   "Sample Track",
				"artist": "Sample Artist",
				"album":  "Sample Album",
				"added":  time.Now().Format(time.RFC3339),
			},
		},
		"playlists": []map[string]interface{}{
			{
				"id":     "playlist456",
				"name":   "My Playlist",
				"owner":  "user123",
				"tracks": []string{"track123", "track789"},
			},
		},
	}

	// In a real implementation, you would save this data:
	// err = store.SaveJSON("sample_backup.json", sampleData)
	// if err != nil {
	//     log.Printf("Failed to save backup: %v", err)
	// }

	fmt.Printf("Sample backup data structure created with %d tracks and %d playlists\n",
		len(sampleData["saved_tracks"].([]map[string]interface{})),
		len(sampleData["playlists"].([]map[string]interface{})))

	fmt.Println("\n=== Backup Management ===")
	fmt.Println("Use the CLI for actual backup operations:")
	fmt.Println("$ spotigo backup              # Create a new backup")
	fmt.Println("$ spotigo backup list         # List available backups")
	fmt.Println("$ spotigo backup restore      # Restore from backup")
	fmt.Println("$ spotigo backup status       # Show backup status")

	// These would be used in actual implementation
	_ = client
	_ = store
}

// ChatExample demonstrates how to use the chat functionality of Spotigo.
// This example shows how to interact with local LLMs for music-related conversations.
func ChatExample() {
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
	options := &ollama.Options{
		Temperature: 0.7,
		NumPredict:  200,
	}

	response, err := ollamaClient.Chat(ctx, ollama.ChatRequest{
		Model:    "granite4:3b",
		Messages: messages,
		Options:  options,
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
		Role:    "user",
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

	// These would be used in actual implementation
	_ = ollamaClient
	_ = messages
}

// RAGExample demonstrates how to use the RAG (Retrieval-Augmented Generation) functionality of Spotigo.
// This example shows how to create embeddings and perform semantic search on music data.
func RAGExample() {
	fmt.Println("=== Spotigo RAG Example ===")

	// Create Ollama client (requires Ollama running locally)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)

	// Check if Ollama is available
	ctx := context.Background()
	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Println("Warning: Ollama not available - RAG functionality will be limited")
		fmt.Println("To use RAG features, start Ollama: ollama serve")
		fmt.Println("Then pull required models: ollama pull nomic-embed-text")
	} else {
		fmt.Println("Ollama connection successful")
	}

	// Create RAG store
	store := rag.NewStore(ollamaClient, "nomic-embed-text", "./data/embeddings.json")
	fmt.Println("RAG store created successfully")

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
		{
			ID:      "track_2",
			Type:    "track",
			Content: "Stairway to Heaven by Led Zeppelin",
			Metadata: map[string]string{
				"artist": "Led Zeppelin",
				"album":  "Led Zeppelin IV",
				"genre":  "Rock",
				"year":   "1971",
			},
		},
		{
			ID:      "artist_1",
			Type:    "artist",
			Content: "Queen is a British rock band formed in London in 1970",
			Metadata: map[string]string{
				"name":    "Queen",
				"formed":  "1970",
				"origin":  "London, England",
				"genres":  "Rock, Glam rock, Hard rock",
				"members": "Freddie Mercury, Brian May, Roger Taylor, John Deacon",
			},
		},
	}

	fmt.Printf("\nCreated %d sample documents\n", len(documents))

	// Add documents to store (this would generate embeddings)
	fmt.Println("\n=== Adding Documents to Store ===")
	for _, doc := range documents {
		if err := store.Add(ctx, doc); err != nil {
			log.Printf("Failed to add document %s: %v", doc.ID, err)
		} else {
			fmt.Printf("Added document: %s (%s)\n", doc.ID, doc.Type)
		}
	}

	// Demonstrate search functionality
	fmt.Println("\n=== Semantic Search Example ===")
	query := "classic rock songs"
	fmt.Printf("Searching for: %q\n", query)

	results, err := store.Search(ctx, query, 5, "")
	if err != nil {
		log.Printf("Search failed: %v", err)
	} else {
		fmt.Printf("Found %d results:\n", len(results))
		for i, result := range results {
			fmt.Printf("%d. %s (similarity: %.2f%%)\n",
				i+1, result.Document.Content, result.Similarity*100)
			fmt.Printf("   Type: %s\n", result.Document.Type)
			if len(result.Document.Metadata) > 0 {
				fmt.Printf("   Metadata: %v\n", result.Document.Metadata)
			}
			fmt.Println()
		}
	}

	// Show store statistics
	fmt.Println("=== Store Statistics ===")
	counts := store.CountByType()
	total := store.Count()
	fmt.Printf("Total documents: %d\n", total)
	for typ, count := range counts {
		fmt.Printf("%s: %d\n", typ, count)
	}

	// Demonstrate persistence
	fmt.Println("\n=== Persistence ===")
	if err := store.Save(); err != nil {
		log.Printf("Failed to save store: %v", err)
	} else {
		fmt.Println("Store saved to disk successfully")
	}

	fmt.Println("\n=== RAG in Spotigo ===")
	fmt.Println("The RAG functionality enables:")
	fmt.Println("- Semantic search across your music library")
	fmt.Println("- AI-powered music recommendations")
	fmt.Println("- Context-aware chat about your music")
	fmt.Println("- Music analysis and insights")

	fmt.Println("\nUse the CLI for full RAG features:")
	fmt.Println("$ spotigo search \"upbeat songs for working out\"")
	fmt.Println("$ spotigo search index    # Build search index")
	fmt.Println("$ spotigo chat           # Chat with your music library")

	// These would be used in actual implementation
	_ = ollamaClient
	_ = store
}
