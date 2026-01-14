package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/rag"
	spotifyclient "github.com/bkataru/spotigo/internal/spotify"
	"github.com/bkataru/spotigo/internal/storage"
)

const (
	// maxConcurrentPlaylistFetches limits concurrent playlist track fetches
	maxConcurrentPlaylistFetches = 5
	// writeBufferSize is the buffer size for file writes (64KB)
	writeBufferSize = 64 * 1024
)

var (
	backupFull  bool
	backupType  string
	backupIndex bool
	backupCmd   = &cobra.Command{
		Use:   "backup",
		Short: "Backup your Spotify library",
		Long: `Backup your complete Spotify library to local JSON/CSV files.

This includes:
  - All saved tracks
  - All playlists (with tracks)
  - Followed artists
  - Saved albums
  - Recently played tracks
  - Top tracks and artists

Data is stored in the configured data directory (default: ./data/backups).

Use --index to automatically build the search index after backup.`,
		Run: func(cmd *cobra.Command, args []string) {
			runBackup()
		},
	}
)

func init() {
	backupCmd.Flags().BoolVar(&backupFull, "full", false, "perform full backup including all data types")
	backupCmd.Flags().StringVar(&backupType, "type", "all", "backup type: all, tracks, playlists, artists")
	backupCmd.Flags().BoolVar(&backupIndex, "index", false, "build search index after backup (requires Ollama)")

	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupStatusCmd)
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available backups",
	Run: func(cmd *cobra.Command, args []string) {
		listBackups()
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore [backup-id]",
	Short: "Restore from a backup",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		restoreBackup(args)
	},
}

var backupStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show backup status and schedule",
	Run: func(cmd *cobra.Command, args []string) {
		showBackupStatus()
	},
}

// backupResult holds the result of a concurrent backup operation
type backupResult struct {
	name string
	data interface{}
	err  error
}

// playlistData holds playlist information with tracks
type playlistData struct {
	ID     spotify.ID  `json:"id"`
	Name   string      `json:"name"`
	Owner  string      `json:"owner"`
	Public bool        `json:"public"`
	Tracks interface{} `json:"tracks"`
}

func runBackup() {
	startTime := time.Now()
	fmt.Println("Starting Spotify library backup...")
	fmt.Printf("  Type: %s\n", backupType)
	fmt.Printf("  Full: %v\n", backupFull)
	if backupIndex {
		fmt.Printf("  Index: enabled\n")
	}
	fmt.Println()

	// Check authentication
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	// Create Spotify client
	spotifyCfg := spotifyclient.Config{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		TokenFile:    cfg.Spotify.TokenFile,
	}

	client, err := spotifyclient.NewClient(spotifyCfg)
	if err != nil {
		fmt.Printf("Error creating Spotify client: %v\n", err)
		return
	}

	if !client.IsAuthenticated() {
		fmt.Println("Not authenticated with Spotify")
		fmt.Println("Run 'spotigo auth' to authenticate first.")
		return
	}

	// Create storage
	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)

	// Perform backup
	if err := performBackupConcurrent(client, store, backupType); err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		return
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nBackup completed successfully in %s!\n", elapsed.Round(time.Millisecond))

	// Build search index if requested
	if backupIndex {
		fmt.Println()
		fmt.Println("Building search index...")
		if err := buildSearchIndex(cfg); err != nil {
			fmt.Printf("Warning: Failed to build search index: %v\n", err)
			fmt.Println("You can run 'spotigo search index' manually later.")
		} else {
			fmt.Println("Search index built successfully!")
		}
	}
}

// performBackupConcurrent performs backup operations concurrently
func performBackupConcurrent(client *spotifyclient.Client, store *storage.Store, backupType string) error {
	ctx := context.Background()

	// Determine which backups to run
	runTracks := backupType == "all" || backupType == "tracks"
	runPlaylists := backupType == "all" || backupType == "playlists"
	runArtists := backupType == "all" || backupType == "artists"

	// Channel for collecting results
	results := make(chan backupResult, 3)
	var wg sync.WaitGroup

	// Launch concurrent fetches
	if runTracks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  🎵 Fetching saved tracks...")
			tracks, err := client.GetSavedTracks(ctx)
			if err == nil {
				fmt.Printf("    Found %d saved tracks\n", len(tracks))
			}
			results <- backupResult{name: "saved_tracks", data: tracks, err: err}
		}()
	}

	if runArtists {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  🎤 Fetching followed artists...")
			artists, err := client.GetFollowedArtists(ctx)
			if err == nil {
				fmt.Printf("    Found %d followed artists\n", len(artists))
			}
			results <- backupResult{name: "followed_artists", data: artists, err: err}
		}()
	}

	if runPlaylists {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  📋 Fetching playlists...")
			playlistData, err := fetchPlaylistsConcurrent(ctx, client)
			results <- backupResult{name: "playlists", data: playlistData, err: err}
		}()
	}

	// Close results channel when all fetches complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	backupData := make(map[string]interface{})
	var errors []error

	for result := range results {
		if result.err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", result.name, result.err))
			continue
		}
		backupData[result.name] = result.data
	}

	// Check for critical errors
	if len(errors) > 0 && len(backupData) == 0 {
		return fmt.Errorf("all backup operations failed: %v", errors)
	}

	// Report any partial failures
	for _, err := range errors {
		fmt.Printf("  Warning: %v\n", err)
	}

	// Write files concurrently
	fmt.Println("\n  💾 Saving backup files...")
	if err := saveBackupFilesConcurrent(store, backupData, backupType); err != nil {
		return fmt.Errorf("failed to save backup files: %w", err)
	}

	return nil
}

// fetchPlaylistsConcurrent fetches all playlists and their tracks concurrently
func fetchPlaylistsConcurrent(ctx context.Context, client *spotifyclient.Client) ([]playlistData, error) {
	// First, get all playlists
	playlists, err := client.GetPlaylists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	fmt.Printf("    Found %d playlists\n", len(playlists))

	if len(playlists) == 0 {
		return []playlistData{}, nil
	}

	// Create worker pool for fetching playlist tracks
	type playlistJob struct {
		index    int
		playlist spotify.SimplePlaylist
	}

	type playlistResult struct {
		index int
		data  playlistData
		err   error
	}

	jobs := make(chan playlistJob, len(playlists))
	results := make(chan playlistResult, len(playlists))

	// Start workers
	numWorkers := maxConcurrentPlaylistFetches
	if len(playlists) < numWorkers {
		numWorkers = len(playlists)
	}

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				items, err := client.GetPlaylistTracks(ctx, job.playlist.ID)
				if err != nil {
					results <- playlistResult{
						index: job.index,
						err:   fmt.Errorf("playlist %s: %w", job.playlist.Name, err),
					}
					continue
				}

				results <- playlistResult{
					index: job.index,
					data: playlistData{
						ID:     job.playlist.ID,
						Name:   job.playlist.Name,
						Owner:  job.playlist.Owner.DisplayName,
						Public: job.playlist.IsPublic,
						Tracks: items,
					},
				}
			}
		}()
	}

	// Send jobs
	for i, playlist := range playlists {
		jobs <- playlistJob{index: i, playlist: playlist}
	}
	close(jobs)

	// Wait for workers and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results maintaining order
	playlistDataSlice := make([]playlistData, len(playlists))
	var errors []error
	completed := 0

	for result := range results {
		completed++
		if result.err != nil {
			errors = append(errors, result.err)
			continue
		}
		playlistDataSlice[result.index] = result.data
		fmt.Printf("      %s: %d tracks\n", result.data.Name, countPlaylistTracks(result.data.Tracks))
	}

	// Filter out empty entries (from errors)
	var validPlaylists []playlistData
	for _, pd := range playlistDataSlice {
		if pd.Name != "" {
			validPlaylists = append(validPlaylists, pd)
		}
	}

	if len(errors) > 0 {
		fmt.Printf("    Warning: %d playlist(s) failed to fetch\n", len(errors))
	}

	return validPlaylists, nil
}

// countPlaylistTracks counts tracks in a playlist
func countPlaylistTracks(tracks interface{}) int {
	if items, ok := tracks.([]spotify.PlaylistItem); ok {
		return len(items)
	}
	return 0
}

// saveBackupFilesConcurrent saves all backup files concurrently with buffered I/O
func saveBackupFilesConcurrent(store *storage.Store, backupData map[string]interface{}, backupType string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(backupData)+1)

	// Save individual data files concurrently
	for name, data := range backupData {
		wg.Add(1)
		go func(filename string, content interface{}) {
			defer wg.Done()
			if err := saveJSONBuffered(filepath.Join(store.GetDataDir(), filename+".json"), content); err != nil {
				errChan <- fmt.Errorf("failed to save %s: %w", filename, err)
			}
		}(name, data)
	}

	// Save combined backup file
	wg.Add(1)
	go func() {
		defer wg.Done()
		metadata, err := createBackupBuffered(store, backupType, backupData)
		if err != nil {
			errChan <- fmt.Errorf("failed to create backup: %w", err)
			return
		}
		fmt.Printf("    Saved backup: %s (%d bytes)\n", metadata.ID, metadata.Size)
	}()

	// Wait for all writes to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("save errors: %v", errors)
	}

	return nil
}

// saveJSONBuffered saves JSON data with buffered I/O for better performance
func saveJSONBuffered(path string, data interface{}) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Use buffered writer for better I/O performance
	writer := bufio.NewWriterSize(file, writeBufferSize)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Flush the buffer
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}

// createBackupBuffered creates a backup file with buffered I/O
func createBackupBuffered(store *storage.Store, backupType string, data interface{}) (*storage.BackupMetadata, error) {
	timestamp := time.Now()
	timestampStr := timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", backupType, timestampStr)
	path := filepath.Join(store.GetBackupDir(), filename)

	// Ensure backup directory exists
	if err := os.MkdirAll(store.GetBackupDir(), 0750); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Use buffered writer
	writer := bufio.NewWriterSize(file, writeBufferSize)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode backup: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %w", err)
	}

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	return &storage.BackupMetadata{
		ID:        filename,
		Timestamp: timestamp,
		Type:      backupType,
		Size:      info.Size(),
	}, nil
}

func listBackups() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)
	backups, err := store.ListBackups()
	if err != nil {
		fmt.Printf("Error listing backups: %v\n", err)
		return
	}

	fmt.Println("Available backups:")
	fmt.Println()

	if len(backups) == 0 {
		fmt.Println("  No backups found.")
		fmt.Println("  Run 'spotigo backup' to create one.")
		return
	}

	for _, backup := range backups {
		fmt.Printf("  %s", backup.ID)
		fmt.Printf("  %s", backup.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  %d bytes", backup.Size)
		if backup.Type != "" {
			fmt.Printf("  [%s]", backup.Type)
		}
		fmt.Println()
	}
}

func restoreBackup(args []string) {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)

	// Get list of available backups
	backups, err := store.ListBackups()
	if err != nil {
		fmt.Printf("Error listing backups: %v\n", err)
		return
	}

	if len(backups) == 0 {
		fmt.Println("No backups available to restore.")
		fmt.Println("Run 'spotigo backup' to create one first.")
		return
	}

	// Sort backups by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	// Determine which backup to restore
	var selectedBackup *storage.BackupMetadata
	backupID := "latest"
	if len(args) > 0 {
		backupID = args[0]
	}

	if backupID == "latest" {
		selectedBackup = &backups[0]
		fmt.Printf("Selected latest backup: %s\n", selectedBackup.ID)
	} else {
		// Find matching backup
		for i := range backups {
			if backups[i].ID == backupID || strings.Contains(backups[i].ID, backupID) {
				selectedBackup = &backups[i]
				break
			}
		}

		if selectedBackup == nil {
			fmt.Printf("Backup not found: %s\n", backupID)
			fmt.Println("\nAvailable backups:")
			for _, b := range backups {
				fmt.Printf("  %s (%s)\n", b.ID, b.Timestamp.Format("2006-01-02 15:04:05"))
			}
			return
		}
	}

	fmt.Printf("Restoring from backup: %s\n", selectedBackup.ID)
	fmt.Printf("  Created: %s\n", selectedBackup.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Size: %d bytes\n", selectedBackup.Size)
	fmt.Println()

	// Load the backup data
	var backupData map[string]interface{}
	if err := store.LoadBackupJSON(selectedBackup.ID, &backupData); err != nil {
		fmt.Printf("Error loading backup: %v\n", err)
		return
	}

	// Restore files concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	restored := 0
	restoreErrors := make([]string, 0)

	restoreFile := func(key, filename string, data interface{}) {
		defer wg.Done()
		if err := store.SaveJSON(filename, data); err != nil {
			mu.Lock()
			restoreErrors = append(restoreErrors, fmt.Sprintf("%s: %v", key, err))
			mu.Unlock()
			return
		}
		count := countItems(data)
		mu.Lock()
		restored++
		fmt.Printf("  Restored %s (%d items)\n", filename, count)
		mu.Unlock()
	}

	// Restore each data type concurrently
	if tracks, ok := backupData["saved_tracks"]; ok {
		wg.Add(1)
		go restoreFile("saved_tracks", "saved_tracks.json", tracks)
	}

	if playlists, ok := backupData["playlists"]; ok {
		wg.Add(1)
		go restoreFile("playlists", "playlists.json", playlists)
	}

	if artists, ok := backupData["followed_artists"]; ok {
		wg.Add(1)
		go restoreFile("followed_artists", "followed_artists.json", artists)
	}

	wg.Wait()

	// Report errors
	for _, errMsg := range restoreErrors {
		fmt.Printf("  Error: %s\n", errMsg)
	}

	fmt.Println()
	if restored > 0 {
		fmt.Printf("Restore completed successfully! (%d data files restored)\n", restored)
		fmt.Println("\nNote: You may want to rebuild the search index:")
		fmt.Println("  spotigo search index")
	} else {
		fmt.Println("No data was restored. The backup may be empty or corrupted.")
	}
}

// countItems returns the count of items in a slice interface
func countItems(data interface{}) int {
	if slice, ok := data.([]interface{}); ok {
		return len(slice)
	}
	return 0
}

func showBackupStatus() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)
	backups, err := store.ListBackups()
	if err != nil {
		fmt.Printf("Error checking backup status: %v\n", err)
		return
	}

	fmt.Println("Backup Status:")
	fmt.Println()

	if len(backups) == 0 {
		fmt.Println("  Last backup: Never")
	} else {
		latest := backups[0]
		fmt.Printf("  Last backup: %s\n", latest.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Latest ID: %s\n", latest.ID)

		// Calculate total size
		var totalSize int64
		for _, backup := range backups {
			totalSize += backup.Size
		}
		fmt.Printf("  Total storage: %.2f MB\n", float64(totalSize)/1024/1024)
	}

	fmt.Printf("  Schedule: %s\n", cfg.Backup.Schedule)
	fmt.Printf("  Retention: %d days\n", cfg.Backup.RetainDays)
}

// buildSearchIndex creates vector embeddings for semantic search
func buildSearchIndex(cfg *config.Config) error {
	// Create Ollama client
	ollamaClient := ollama.NewClient(cfg.Ollama.Host, time.Duration(cfg.Ollama.Timeout)*time.Second)

	// Check if Ollama is available
	ctx := context.Background()
	if err := ollamaClient.Ping(ctx); err != nil {
		return fmt.Errorf("ollama not available: %w", err)
	}

	// Load the vector store
	embeddingModel := "nomic-embed-text"
	storePath := filepath.Join(cfg.Storage.EmbeddingsDir, "vectors.json")
	store := rag.NewStore(ollamaClient, embeddingModel, storePath)

	// Load backup data and create documents
	docs := loadBackupDocuments(cfg)

	if len(docs) == 0 {
		return fmt.Errorf("no backup data to index")
	}

	fmt.Printf("  Indexing %d items...\n", len(docs))

	// Add documents with progress
	for i, doc := range docs {
		if err := store.Add(ctx, doc); err != nil {
			// Log but continue on individual failures
			continue
		}

		// Show progress every 25 items
		if (i+1)%25 == 0 || i+1 == len(docs) {
			fmt.Printf("  Indexed %d/%d items...\r", i+1, len(docs))
		}
	}
	fmt.Println()

	// Save the store
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	counts := store.CountByType()
	for typ, count := range counts {
		fmt.Printf("  - %s: %d\n", typ, count)
	}

	return nil
}
