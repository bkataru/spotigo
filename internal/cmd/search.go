package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bkataru-workshop/spotigo/internal/config"
	"github.com/bkataru-workshop/spotigo/internal/jsonutil"
	"github.com/bkataru-workshop/spotigo/internal/ollama"
	"github.com/bkataru-workshop/spotigo/internal/rag"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Semantic search across your music library",
	Long: `Search your music library using natural language queries.

Examples:
  spotigo search "upbeat songs for working out"
  spotigo search "melancholic piano music"
  spotigo search "songs similar to Radiohead"
  spotigo search --type artists "indie rock bands"

The search uses AI embeddings to understand meaning, not just keywords.
Run 'spotigo search index' first to build the search index from your backups.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		runSearch(query)
	},
}

var searchIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Build or rebuild the search index from backup data",
	Run: func(cmd *cobra.Command, args []string) {
		runSearchIndex()
	},
}

var searchStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show search index status",
	Run: func(cmd *cobra.Command, args []string) {
		runSearchStatus()
	},
}

var (
	searchLimit  int
	searchType   string
	searchFormat string
	searchModel  string
)

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "maximum number of results")
	searchCmd.Flags().StringVar(&searchType, "type", "all", "search type: all, tracks, artists, playlists")
	searchCmd.Flags().StringVar(&searchFormat, "format", "table", "output format: table, json")
	searchCmd.Flags().StringVar(&searchModel, "model", "nomic-embed-text", "embedding model to use")

	searchCmd.AddCommand(searchIndexCmd)
	searchCmd.AddCommand(searchStatusCmd)
}

func runSearch(query string) {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	// Create Ollama client
	client := ollama.NewClient(cfg.Ollama.Host, time.Duration(cfg.Ollama.Timeout)*time.Second)

	// Check if Ollama is available
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		fmt.Println("Error: Ollama is not available")
		fmt.Printf("  %v\n", err)
		fmt.Println()
		fmt.Println("Make sure Ollama is running: ollama serve")
		return
	}

	// Load the vector store
	storePath := filepath.Join(cfg.Storage.EmbeddingsDir, "vectors.json")
	store := rag.NewStore(client, searchModel, storePath)

	if err := store.Load(); err != nil {
		fmt.Printf("Error loading search index: %v\n", err)
		return
	}

	if store.Count() == 0 {
		fmt.Println("Search index is empty.")
		fmt.Println()
		fmt.Println("To build the index:")
		fmt.Println("  1. Run 'spotigo backup' to save your Spotify library")
		fmt.Println("  2. Run 'spotigo search index' to build embeddings")
		return
	}

	fmt.Printf("Searching for: \"%s\"\n", query)
	fmt.Printf("  Type: %s, Limit: %d\n", searchType, searchLimit)
	fmt.Println()

	// Map type flag to document type
	docType := searchType
	if docType == "tracks" {
		docType = "track"
	} else if docType == "artists" {
		docType = "artist"
	} else if docType == "playlists" {
		docType = "playlist"
	}

	// Perform search
	results, err := store.Search(ctx, query, searchLimit, docType)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	// Display results
	fmt.Printf("Found %d results:\n\n", len(results))

	for i, result := range results {
		displaySearchResult(i+1, result)
	}
}

func displaySearchResult(rank int, result rag.SearchResult) {
	doc := result.Document
	similarity := result.Similarity * 100 // Convert to percentage

	switch doc.Type {
	case "track":
		fmt.Printf("%2d. [%.0f%%] 🎵 %s\n", rank, similarity, doc.Metadata["name"])
		fmt.Printf("         by %s\n", doc.Metadata["artists"])
		if doc.Metadata["album"] != "" {
			fmt.Printf("         📀 %s\n", doc.Metadata["album"])
		}
	case "artist":
		fmt.Printf("%2d. [%.0f%%] 🎤 %s\n", rank, similarity, doc.Metadata["name"])
		if doc.Metadata["genres"] != "" {
			fmt.Printf("         Genres: %s\n", doc.Metadata["genres"])
		}
	case "playlist":
		fmt.Printf("%2d. [%.0f%%] 📋 %s\n", rank, similarity, doc.Metadata["name"])
		if doc.Metadata["description"] != "" {
			desc := doc.Metadata["description"]
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			fmt.Printf("         %s\n", desc)
		}
		fmt.Printf("         %s tracks\n", doc.Metadata["track_count"])
	default:
		fmt.Printf("%2d. [%.0f%%] %s: %s\n", rank, similarity, doc.Type, doc.Content)
	}
	fmt.Println()
}

func runSearchIndex() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	// Create Ollama client
	client := ollama.NewClient(cfg.Ollama.Host, time.Duration(cfg.Ollama.Timeout)*time.Second)

	// Check if Ollama is available
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		fmt.Println("Error: Ollama is not available")
		fmt.Printf("  %v\n", err)
		fmt.Println()
		fmt.Println("Make sure Ollama is running and has an embedding model:")
		fmt.Println("  ollama pull nomic-embed-text")
		return
	}

	fmt.Println("Building search index...")
	fmt.Printf("  Using model: %s\n", searchModel)
	fmt.Println()

	// Load the vector store
	storePath := filepath.Join(cfg.Storage.EmbeddingsDir, "vectors.json")
	store := rag.NewStore(client, searchModel, storePath)

	// Load backup data and create documents
	docs := loadBackupDocuments(cfg)
	if len(docs) == 0 {
		fmt.Println("No backup data found.")
		fmt.Println("Run 'spotigo backup' first to save your Spotify library.")
		return
	}

	fmt.Printf("Found %d items to index\n", len(docs))
	fmt.Println()

	// Add documents with progress
	for i, doc := range docs {
		if err := store.Add(ctx, doc); err != nil {
			fmt.Printf("  Warning: failed to index %s: %v\n", doc.ID, err)
			continue
		}

		// Show progress every 10 items
		if (i+1)%10 == 0 || i+1 == len(docs) {
			fmt.Printf("  Indexed %d/%d items...\r", i+1, len(docs))
		}
	}
	fmt.Println()

	// Save the store
	if err := store.Save(); err != nil {
		fmt.Printf("Error saving index: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Printf("✅ Search index built successfully!\n")
	fmt.Printf("   Indexed %d items\n", store.Count())

	counts := store.CountByType()
	for typ, count := range counts {
		fmt.Printf("   - %s: %d\n", typ, count)
	}
}

func loadBackupDocuments(cfg *config.Config) []rag.Document {
	var docs []rag.Document

	// Load saved tracks
	tracksPath := filepath.Join(cfg.Storage.DataDir, "saved_tracks.json")
	tracks, err := loadTracksFromFile(tracksPath)
	if err == nil {
		for _, track := range tracks {
			docs = append(docs, rag.TrackToDocument(track))
		}
	}

	// Load followed artists
	artistsPath := filepath.Join(cfg.Storage.DataDir, "followed_artists.json")
	artists, err := loadArtistsFromFile(artistsPath)
	if err == nil {
		for _, artist := range artists {
			docs = append(docs, rag.ArtistToDocument(artist))
		}
	}

	// Load playlists
	playlistsPath := filepath.Join(cfg.Storage.DataDir, "playlists.json")
	playlists, err := loadPlaylistsFromFile(playlistsPath)
	if err == nil {
		for _, playlist := range playlists {
			docs = append(docs, rag.PlaylistToDocument(playlist))
		}
	}

	return docs
}

func loadTracksFromFile(path string) ([]rag.TrackData, error) {
	// Load raw JSON and convert to TrackData
	var rawTracks []map[string]interface{}
	if err := jsonutil.LoadJSONFile(path, &rawTracks); err != nil {
		return nil, err
	}

	var tracks []rag.TrackData
	for _, raw := range rawTracks {
		track := rag.TrackData{}

		// Extract track info - handle both spotify.SavedTrack and plain track formats
		if trackData, ok := raw["track"].(map[string]interface{}); ok {
			// spotify.SavedTrack format
			track.ID = jsonutil.GetString(trackData, "id")
			track.Name = jsonutil.GetString(trackData, "name")
			track.Album = jsonutil.GetNestedString(trackData, "album", "name")
			track.Artists = jsonutil.GetArtistNames(trackData)
		} else {
			// Plain track format
			track.ID = jsonutil.GetString(raw, "id")
			track.Name = jsonutil.GetString(raw, "name")
			track.Album = jsonutil.GetNestedString(raw, "album", "name")
			track.Artists = jsonutil.GetArtistNames(raw)
		}

		if track.ID != "" && track.Name != "" {
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

func loadArtistsFromFile(path string) ([]rag.ArtistData, error) {
	var rawArtists []map[string]interface{}
	if err := jsonutil.LoadJSONFile(path, &rawArtists); err != nil {
		return nil, err
	}

	var artists []rag.ArtistData
	for _, raw := range rawArtists {
		artist := rag.ArtistData{
			ID:     jsonutil.GetString(raw, "id"),
			Name:   jsonutil.GetString(raw, "name"),
			Genres: jsonutil.GetStringSlice(raw, "genres"),
		}

		if artist.ID != "" && artist.Name != "" {
			artists = append(artists, artist)
		}
	}

	return artists, nil
}

func loadPlaylistsFromFile(path string) ([]rag.PlaylistData, error) {
	var rawPlaylists []map[string]interface{}
	if err := jsonutil.LoadJSONFile(path, &rawPlaylists); err != nil {
		return nil, err
	}

	var playlists []rag.PlaylistData
	for _, raw := range rawPlaylists {
		playlist := rag.PlaylistData{
			ID:          jsonutil.GetString(raw, "id"),
			Name:        jsonutil.GetString(raw, "name"),
			Description: jsonutil.GetString(raw, "description"),
			Owner:       jsonutil.GetString(raw, "owner"),
		}

		// Get track names if available
		if tracks, ok := raw["tracks"].([]interface{}); ok {
			playlist.TrackCount = len(tracks)
			for _, t := range tracks {
				if trackMap, ok := t.(map[string]interface{}); ok {
					if trackData, ok := trackMap["track"].(map[string]interface{}); ok {
						if name := jsonutil.GetString(trackData, "name"); name != "" {
							playlist.TrackNames = append(playlist.TrackNames, name)
						}
					}
				}
			}
		}

		if playlist.ID != "" && playlist.Name != "" {
			playlists = append(playlists, playlist)
		}
	}

	return playlists, nil
}

func runSearchStatus() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	storePath := filepath.Join(cfg.Storage.EmbeddingsDir, "vectors.json")
	store := rag.NewStore(nil, "", storePath)

	if err := store.Load(); err != nil {
		fmt.Printf("Error loading search index: %v\n", err)
		return
	}

	fmt.Println("Search Index Status")
	fmt.Println("===================")
	fmt.Println()

	if store.Count() == 0 {
		fmt.Println("Index: Empty")
		fmt.Println()
		fmt.Println("Run 'spotigo search index' to build the index.")
		return
	}

	fmt.Printf("Total documents: %d\n", store.Count())
	fmt.Println()

	counts := store.CountByType()
	fmt.Println("By type:")
	for typ, count := range counts {
		fmt.Printf("  - %s: %d\n", typ, count)
	}
	fmt.Println()

	fmt.Printf("Index location: %s\n", storePath)
}
