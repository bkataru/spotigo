// backup_example.go demonstrates how to use the backup functionality of Spotigo
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/spotify"
	"github.com/bkataru/spotigo/internal/storage"
)

func main() {
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
		return
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
				"id":    "playlist456",
				"name":  "My Playlist",
				"owner": "user123",
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
}