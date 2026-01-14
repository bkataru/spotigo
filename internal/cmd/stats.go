package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"spotigo/internal/config"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your listening statistics",
	Long: `View detailed statistics about your music library.

Available statistics:
  - Library overview (tracks, playlists, artists)
  - Top artists by track count
  - Genre distribution
  - Playlist analysis

Statistics are calculated from your backup data.
Run 'spotigo backup' first to generate statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStats()
	},
}

var (
	statsPeriod string
	statsTop    int
)

func init() {
	statsCmd.Flags().StringVar(&statsPeriod, "period", "all", "time period: all, year, month, week")
	statsCmd.Flags().IntVar(&statsTop, "top", 10, "number of top items to show")

	// Subcommands for specific stats
	statsCmd.AddCommand(statsTopCmd)
	statsCmd.AddCommand(statsGenresCmd)
	statsCmd.AddCommand(statsPlaylistsCmd)
}

var statsTopCmd = &cobra.Command{
	Use:   "top [artists|tracks|albums]",
	Short: "Show your top artists, tracks, or albums",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		itemType := "artists"
		if len(args) > 0 {
			itemType = args[0]
		}
		runStatsTop(itemType)
	},
}

var statsGenresCmd = &cobra.Command{
	Use:   "genres",
	Short: "Show genre distribution",
	Run: func(cmd *cobra.Command, args []string) {
		runStatsGenres()
	},
}

var statsPlaylistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "Show playlist statistics",
	Run: func(cmd *cobra.Command, args []string) {
		runStatsPlaylists()
	},
}

// LibraryStats holds computed statistics about the music library
type LibraryStats struct {
	TotalTracks    int
	TotalPlaylists int
	TotalArtists   int
	UniqueAlbums   int
	TopArtists     []ArtistCount
	TopGenres      []GenreCount
	PlaylistStats  []PlaylistStat
}

// ArtistCount tracks artist frequency
type ArtistCount struct {
	Name  string
	Count int
}

// GenreCount tracks genre frequency
type GenreCount struct {
	Name  string
	Count int
}

// PlaylistStat holds playlist statistics
type PlaylistStat struct {
	Name       string
	TrackCount int
	Owner      string
}

func runStats() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	stats, err := computeLibraryStats(cfg)
	if err != nil {
		fmt.Println("Error computing statistics:", err)
		fmt.Println()
		fmt.Println("Run 'spotigo backup' first to save your Spotify library.")
		return
	}

	fmt.Println("Library Statistics")
	fmt.Println("==================")
	fmt.Println()

	fmt.Printf("  Total Tracks:    %d\n", stats.TotalTracks)
	fmt.Printf("  Total Playlists: %d\n", stats.TotalPlaylists)
	fmt.Printf("  Unique Artists:  %d\n", stats.TotalArtists)
	fmt.Printf("  Unique Albums:   %d\n", stats.UniqueAlbums)
	fmt.Println()

	if len(stats.TopArtists) > 0 {
		fmt.Println("Top 5 Artists (by saved tracks):")
		limit := 5
		if len(stats.TopArtists) < limit {
			limit = len(stats.TopArtists)
		}
		for i, artist := range stats.TopArtists[:limit] {
			fmt.Printf("  %d. %s (%d tracks)\n", i+1, artist.Name, artist.Count)
		}
		fmt.Println()
	}

	if len(stats.TopGenres) > 0 {
		fmt.Println("Top 5 Genres:")
		limit := 5
		if len(stats.TopGenres) < limit {
			limit = len(stats.TopGenres)
		}
		for i, genre := range stats.TopGenres[:limit] {
			fmt.Printf("  %d. %s (%d artists)\n", i+1, genre.Name, genre.Count)
		}
		fmt.Println()
	}

	fmt.Println("Run 'spotigo stats top' for detailed rankings.")
	fmt.Println("Run 'spotigo stats genres' for full genre breakdown.")
	fmt.Println("Run 'spotigo stats playlists' for playlist analysis.")
}

func runStatsTop(itemType string) {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	stats, err := computeLibraryStats(cfg)
	if err != nil {
		fmt.Println("Error computing statistics:", err)
		return
	}

	switch itemType {
	case "artists":
		fmt.Printf("Top %d Artists (by saved tracks):\n", statsTop)
		fmt.Println()

		if len(stats.TopArtists) == 0 {
			fmt.Println("  No artist data available.")
			return
		}

		limit := statsTop
		if len(stats.TopArtists) < limit {
			limit = len(stats.TopArtists)
		}

		for i, artist := range stats.TopArtists[:limit] {
			bar := strings.Repeat("█", min(artist.Count, 20))
			fmt.Printf("%3d. %-30s %3d %s\n", i+1, truncate(artist.Name, 30), artist.Count, bar)
		}

	case "tracks":
		fmt.Println("Top tracks by play count is not available from backup data.")
		fmt.Println("Spotify's 'Top Tracks' feature requires real-time API access.")

	case "albums":
		fmt.Printf("Top %d Albums (by saved tracks):\n", statsTop)
		fmt.Println()

		albums, err := computeTopAlbums(cfg)
		if err != nil {
			fmt.Println("  No album data available.")
			return
		}

		limit := statsTop
		if len(albums) < limit {
			limit = len(albums)
		}

		for i, album := range albums[:limit] {
			bar := strings.Repeat("█", min(album.Count, 20))
			fmt.Printf("%3d. %-40s %3d %s\n", i+1, truncate(album.Name, 40), album.Count, bar)
		}

	default:
		fmt.Printf("Unknown type: %s\n", itemType)
		fmt.Println("Available types: artists, tracks, albums")
	}
}

func runStatsGenres() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	stats, err := computeLibraryStats(cfg)
	if err != nil {
		fmt.Println("Error computing statistics:", err)
		return
	}

	fmt.Println("Genre Distribution")
	fmt.Println("==================")
	fmt.Println()

	if len(stats.TopGenres) == 0 {
		fmt.Println("No genre data available.")
		fmt.Println("Genres are extracted from followed artists.")
		return
	}

	// Calculate max for scaling bars
	maxCount := 0
	for _, g := range stats.TopGenres {
		if g.Count > maxCount {
			maxCount = g.Count
		}
	}

	limit := statsTop
	if len(stats.TopGenres) < limit {
		limit = len(stats.TopGenres)
	}

	for i, genre := range stats.TopGenres[:limit] {
		barLen := (genre.Count * 30) / maxCount
		if barLen < 1 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)
		fmt.Printf("%3d. %-25s %3d %s\n", i+1, truncate(genre.Name, 25), genre.Count, bar)
	}

	fmt.Println()
	fmt.Printf("Total genres: %d\n", len(stats.TopGenres))
}

func runStatsPlaylists() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	stats, err := computeLibraryStats(cfg)
	if err != nil {
		fmt.Println("Error computing statistics:", err)
		return
	}

	fmt.Println("Playlist Statistics")
	fmt.Println("===================")
	fmt.Println()

	if len(stats.PlaylistStats) == 0 {
		fmt.Println("No playlist data available.")
		return
	}

	// Sort by track count
	sort.Slice(stats.PlaylistStats, func(i, j int) bool {
		return stats.PlaylistStats[i].TrackCount > stats.PlaylistStats[j].TrackCount
	})

	fmt.Printf("Total playlists: %d\n", len(stats.PlaylistStats))
	fmt.Println()

	// Calculate total tracks in playlists
	totalTracks := 0
	for _, p := range stats.PlaylistStats {
		totalTracks += p.TrackCount
	}
	fmt.Printf("Total tracks across playlists: %d\n", totalTracks)
	if len(stats.PlaylistStats) > 0 {
		fmt.Printf("Average tracks per playlist: %.1f\n", float64(totalTracks)/float64(len(stats.PlaylistStats)))
	}
	fmt.Println()

	fmt.Println("Largest playlists:")
	limit := statsTop
	if len(stats.PlaylistStats) < limit {
		limit = len(stats.PlaylistStats)
	}

	for i, p := range stats.PlaylistStats[:limit] {
		fmt.Printf("%3d. %-40s %4d tracks\n", i+1, truncate(p.Name, 40), p.TrackCount)
	}
}

func computeLibraryStats(cfg *config.Config) (*LibraryStats, error) {
	stats := &LibraryStats{}

	// Load saved tracks
	tracksPath := filepath.Join(cfg.Storage.DataDir, "saved_tracks.json")
	tracks, err := loadRawTracks(tracksPath)
	if err == nil {
		stats.TotalTracks = len(tracks)

		// Count artists and albums
		artistCounts := make(map[string]int)
		albumSet := make(map[string]bool)

		for _, track := range tracks {
			// Get artist names
			for _, artist := range getTrackArtists(track) {
				artistCounts[artist]++
			}

			// Get album name
			if album := getTrackAlbum(track); album != "" {
				albumSet[album] = true
			}
		}

		stats.UniqueAlbums = len(albumSet)
		stats.TotalArtists = len(artistCounts)

		// Convert to sorted slice
		for name, count := range artistCounts {
			stats.TopArtists = append(stats.TopArtists, ArtistCount{Name: name, Count: count})
		}
		sort.Slice(stats.TopArtists, func(i, j int) bool {
			return stats.TopArtists[i].Count > stats.TopArtists[j].Count
		})
	}

	// Load followed artists for genres
	artistsPath := filepath.Join(cfg.Storage.DataDir, "followed_artists.json")
	artists, err := loadRawArtists(artistsPath)
	if err == nil {
		genreCounts := make(map[string]int)

		for _, artist := range artists {
			genres := getArtistGenres(artist)
			for _, genre := range genres {
				genreCounts[genre]++
			}
		}

		for name, count := range genreCounts {
			stats.TopGenres = append(stats.TopGenres, GenreCount{Name: name, Count: count})
		}
		sort.Slice(stats.TopGenres, func(i, j int) bool {
			return stats.TopGenres[i].Count > stats.TopGenres[j].Count
		})
	}

	// Load playlists
	playlistsPath := filepath.Join(cfg.Storage.DataDir, "playlists.json")
	playlists, err := loadRawPlaylists(playlistsPath)
	if err == nil {
		stats.TotalPlaylists = len(playlists)

		for _, playlist := range playlists {
			ps := PlaylistStat{
				Name:       getPlaylistName(playlist),
				TrackCount: getPlaylistTrackCount(playlist),
				Owner:      getPlaylistOwner(playlist),
			}
			stats.PlaylistStats = append(stats.PlaylistStats, ps)
		}
	}

	// Check if we have any data
	if stats.TotalTracks == 0 && stats.TotalPlaylists == 0 {
		return nil, fmt.Errorf("no backup data found")
	}

	return stats, nil
}

func computeTopAlbums(cfg *config.Config) ([]ArtistCount, error) {
	tracksPath := filepath.Join(cfg.Storage.DataDir, "saved_tracks.json")
	tracks, err := loadRawTracks(tracksPath)
	if err != nil {
		return nil, err
	}

	albumCounts := make(map[string]int)
	for _, track := range tracks {
		if album := getTrackAlbum(track); album != "" {
			albumCounts[album]++
		}
	}

	var albums []ArtistCount
	for name, count := range albumCounts {
		albums = append(albums, ArtistCount{Name: name, Count: count})
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Count > albums[j].Count
	})

	return albums, nil
}

// Helper functions for loading and parsing JSON

func loadRawTracks(path string) ([]map[string]interface{}, error) {
	var tracks []map[string]interface{}
	if err := loadStatsJSONFile(path, &tracks); err != nil {
		return nil, err
	}
	return tracks, nil
}

func loadRawArtists(path string) ([]map[string]interface{}, error) {
	var artists []map[string]interface{}
	if err := loadStatsJSONFile(path, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func loadRawPlaylists(path string) ([]map[string]interface{}, error) {
	var playlists []map[string]interface{}
	if err := loadStatsJSONFile(path, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func loadStatsJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func getTrackArtists(track map[string]interface{}) []string {
	var names []string

	// Handle spotify.SavedTrack format
	trackData := track
	if t, ok := track["track"].(map[string]interface{}); ok {
		trackData = t
	}

	if artists, ok := trackData["artists"].([]interface{}); ok {
		for _, a := range artists {
			if artist, ok := a.(map[string]interface{}); ok {
				if name, ok := artist["name"].(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
	}
	return names
}

func getTrackAlbum(track map[string]interface{}) string {
	// Handle spotify.SavedTrack format
	trackData := track
	if t, ok := track["track"].(map[string]interface{}); ok {
		trackData = t
	}

	if album, ok := trackData["album"].(map[string]interface{}); ok {
		if name, ok := album["name"].(string); ok {
			return name
		}
	}
	return ""
}

func getArtistGenres(artist map[string]interface{}) []string {
	var genres []string
	if g, ok := artist["genres"].([]interface{}); ok {
		for _, genre := range g {
			if s, ok := genre.(string); ok {
				genres = append(genres, s)
			}
		}
	}
	return genres
}

func getPlaylistName(playlist map[string]interface{}) string {
	if name, ok := playlist["name"].(string); ok {
		return name
	}
	return ""
}

func getPlaylistTrackCount(playlist map[string]interface{}) int {
	if tracks, ok := playlist["tracks"].([]interface{}); ok {
		return len(tracks)
	}
	return 0
}

func getPlaylistOwner(playlist map[string]interface{}) string {
	if owner, ok := playlist["owner"].(string); ok {
		return owner
	}
	return ""
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
