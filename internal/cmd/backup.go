package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"spotigo/internal/spotify"
	"spotigo/internal/storage"
)

var (
	backupFull     bool
	backupType     string
	backupSchedule string
	backupCmd      = &cobra.Command{
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

Data is stored in the configured data directory (default: ./data/backups).`,
		Run: func(cmd *cobra.Command, args []string) {
			runBackup(cmd)
		},
	}
)

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

func runBackup(cmd *cobra.Command) {
	fmt.Println("Starting Spotify library backup...")
	fmt.Printf("  Type: %s\n", backupType)
	fmt.Printf("  Full: %v\n", backupFull)
	fmt.Println()

	// Check authentication
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
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
		fmt.Printf("Error creating Spotify client: %v\n", err)
		return
	}

	if !client.IsAuthenticated() {
		fmt.Println("❌ Not authenticated with Spotify")
		fmt.Println("Run 'spotigo auth' to authenticate first.")
		return
	}

	// Create storage
	store := storage.NewStore(cfg.Storage.DataDir, cfg.Storage.BackupDir)

	// Perform backup
	if err := performBackup(client, store, backupType); err != nil {
		fmt.Printf("❌ Backup failed: %v\n", err)
		return
	}

	fmt.Println("✅ Backup completed successfully!")
}

func performBackup(client *spotify.Client, store *storage.Store, backupType string) error {
	ctx := context.Background()

	// Create backup data structure
	backupData := make(map[string]interface{})

	switch backupType {
	case "all", "tracks":
		if err := backupTracks(ctx, client, store, backupData); err != nil {
			return fmt.Errorf("failed to backup tracks: %w", err)
		}
	}

	switch backupType {
	case "all", "playlists":
		if err := backupPlaylists(ctx, client, store, backupData); err != nil {
			return fmt.Errorf("failed to backup playlists: %w", err)
		}
	}

	switch backupType {
	case "all", "artists":
		if err := backupArtists(ctx, client, store, backupData); err != nil {
			return fmt.Errorf("failed to backup artists: %w", err)
		}
	}

	// Save backup
	metadata, err := store.CreateBackup(backupType, backupData)
	if err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	fmt.Printf("  Saved backup: %s\n", metadata.ID)
	fmt.Printf("  Size: %d bytes\n", metadata.Size)

	return nil
}

func backupTracks(ctx context.Context, client *spotify.Client, store *storage.Store, backupData map[string]interface{}) error {
	fmt.Println("  🎵 Backing up saved tracks...")

	tracks, err := client.GetSavedTracks(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("    Found %d saved tracks\n", len(tracks))
	backupData["saved_tracks"] = tracks

	// Also save as separate file for easy access
	if err := store.SaveJSON("saved_tracks.json", tracks); err != nil {
		return fmt.Errorf("failed to save tracks: %w", err)
	}

	return nil
}

func backupPlaylists(ctx context.Context, client *spotify.Client, store *storage.Store, backupData map[string]interface{}) error {
	fmt.Println("  📋 Backing up playlists...")

	playlists, err := client.GetPlaylists(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("    Found %d playlists\n", len(playlists))

	// Get tracks for each playlist
	playlistData := make([]map[string]interface{}, len(playlists))
	for i, playlist := range playlists {
		items, err := client.GetPlaylistTracks(ctx, playlist.ID)
		if err != nil {
			return fmt.Errorf("failed to get tracks for playlist %s: %w", playlist.Name, err)
		}

		playlistData[i] = map[string]interface{}{
			"id":     playlist.ID,
			"name":   playlist.Name,
			"owner":  playlist.Owner.DisplayName,
			"public": playlist.IsPublic,
			"tracks": items,
		}

		fmt.Printf("      %s: %d tracks\n", playlist.Name, len(items))
	}

	backupData["playlists"] = playlistData

	// Save as separate file
	if err := store.SaveJSON("playlists.json", playlistData); err != nil {
		return fmt.Errorf("failed to save playlists: %w", err)
	}

	return nil
}

func backupArtists(ctx context.Context, client *spotify.Client, store *storage.Store, backupData map[string]interface{}) error {
	fmt.Println("  🎤 Backing up followed artists...")

	artists, err := client.GetFollowedArtists(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("    Found %d followed artists\n", len(artists))
	backupData["followed_artists"] = artists

	// Also save as separate file
	if err := store.SaveJSON("followed_artists.json", artists); err != nil {
		return fmt.Errorf("failed to save artists: %w", err)
	}

	return nil
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
	backupID := "latest"
	if len(args) > 0 {
		backupID = args[0]
	}

	fmt.Printf("Restoring from backup: %s\n", backupID)
	fmt.Println("🚧 Restore functionality coming soon...")
	fmt.Println("For now, you can manually restore from:")
	fmt.Println("  ./data/backups/")
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
