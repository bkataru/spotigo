package jsonutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSONFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jsonutil-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write test JSON
	testPath := filepath.Join(tmpDir, "test.json")
	testData := `{"name": "test", "count": 42, "items": ["a", "b", "c"]}`
	if err := os.WriteFile(testPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var result map[string]interface{}
	err = LoadJSONFile(testPath, &result)
	if err != nil {
		t.Fatalf("LoadJSONFile failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("expected name 'test', got '%v'", result["name"])
	}
	if result["count"].(float64) != 42 {
		t.Errorf("expected count 42, got %v", result["count"])
	}
}

func TestLoadJSONFile_NotFound(t *testing.T) {
	var result map[string]interface{}
	err := LoadJSONFile("/nonexistent/path.json", &result)
	if err == nil {
		t.Error("LoadJSONFile should fail for nonexistent file")
	}
}

func TestGetString(t *testing.T) {
	m := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}

	if got := GetString(m, "name"); got != "John" {
		t.Errorf("expected 'John', got '%s'", got)
	}

	if got := GetString(m, "age"); got != "" {
		t.Errorf("expected empty string for non-string value, got '%s'", got)
	}

	if got := GetString(m, "missing"); got != "" {
		t.Errorf("expected empty string for missing key, got '%s'", got)
	}
}

func TestGetNestedString(t *testing.T) {
	m := map[string]interface{}{
		"album": map[string]interface{}{
			"name": "Abbey Road",
			"artist": map[string]interface{}{
				"name": "The Beatles",
			},
		},
	}

	if got := GetNestedString(m, "album", "name"); got != "Abbey Road" {
		t.Errorf("expected 'Abbey Road', got '%s'", got)
	}

	if got := GetNestedString(m, "album", "artist", "name"); got != "The Beatles" {
		t.Errorf("expected 'The Beatles', got '%s'", got)
	}

	if got := GetNestedString(m, "album", "missing"); got != "" {
		t.Errorf("expected empty string for missing key, got '%s'", got)
	}

	if got := GetNestedString(m, "missing", "name"); got != "" {
		t.Errorf("expected empty string for missing path, got '%s'", got)
	}
}

func TestGetStringSlice(t *testing.T) {
	m := map[string]interface{}{
		"genres": []interface{}{"rock", "pop", "indie"},
		"count":  42,
	}

	genres := GetStringSlice(m, "genres")
	if len(genres) != 3 {
		t.Fatalf("expected 3 genres, got %d", len(genres))
	}
	if genres[0] != "rock" {
		t.Errorf("expected 'rock', got '%s'", genres[0])
	}

	if got := GetStringSlice(m, "count"); got != nil {
		t.Error("expected nil for non-slice value")
	}

	if got := GetStringSlice(m, "missing"); got != nil {
		t.Error("expected nil for missing key")
	}
}

func TestGetArtistNames(t *testing.T) {
	m := map[string]interface{}{
		"artists": []interface{}{
			map[string]interface{}{"name": "Artist A", "id": "1"},
			map[string]interface{}{"name": "Artist B", "id": "2"},
		},
	}

	names := GetArtistNames(m)
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
	if names[0] != "Artist A" {
		t.Errorf("expected 'Artist A', got '%s'", names[0])
	}
	if names[1] != "Artist B" {
		t.Errorf("expected 'Artist B', got '%s'", names[1])
	}
}

func TestGetTrackArtists(t *testing.T) {
	// Test plain track format
	plainTrack := map[string]interface{}{
		"name": "Song",
		"artists": []interface{}{
			map[string]interface{}{"name": "Artist"},
		},
	}
	names := GetTrackArtists(plainTrack)
	if len(names) != 1 || names[0] != "Artist" {
		t.Errorf("plain track: expected ['Artist'], got %v", names)
	}

	// Test SavedTrack format (nested)
	savedTrack := map[string]interface{}{
		"added_at": "2024-01-01",
		"track": map[string]interface{}{
			"name": "Song",
			"artists": []interface{}{
				map[string]interface{}{"name": "Nested Artist"},
			},
		},
	}
	names = GetTrackArtists(savedTrack)
	if len(names) != 1 || names[0] != "Nested Artist" {
		t.Errorf("saved track: expected ['Nested Artist'], got %v", names)
	}
}

func TestGetTrackAlbum(t *testing.T) {
	// Test plain track format
	plainTrack := map[string]interface{}{
		"name": "Song",
		"album": map[string]interface{}{
			"name": "Album Name",
		},
	}
	if got := GetTrackAlbum(plainTrack); got != "Album Name" {
		t.Errorf("plain track: expected 'Album Name', got '%s'", got)
	}

	// Test SavedTrack format (nested)
	savedTrack := map[string]interface{}{
		"added_at": "2024-01-01",
		"track": map[string]interface{}{
			"name": "Song",
			"album": map[string]interface{}{
				"name": "Nested Album",
			},
		},
	}
	if got := GetTrackAlbum(savedTrack); got != "Nested Album" {
		t.Errorf("saved track: expected 'Nested Album', got '%s'", got)
	}

	// Test missing album
	noAlbum := map[string]interface{}{"name": "Song"}
	if got := GetTrackAlbum(noAlbum); got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}

func TestGetArtistGenres(t *testing.T) {
	artist := map[string]interface{}{
		"name":   "Artist",
		"genres": []interface{}{"rock", "alternative"},
	}

	genres := GetArtistGenres(artist)
	if len(genres) != 2 {
		t.Fatalf("expected 2 genres, got %d", len(genres))
	}
	if genres[0] != "rock" {
		t.Errorf("expected 'rock', got '%s'", genres[0])
	}

	// Test missing genres
	noGenres := map[string]interface{}{"name": "Artist"}
	if got := GetArtistGenres(noGenres); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestGetPlaylistName(t *testing.T) {
	playlist := map[string]interface{}{"name": "My Playlist"}
	if got := GetPlaylistName(playlist); got != "My Playlist" {
		t.Errorf("expected 'My Playlist', got '%s'", got)
	}

	empty := map[string]interface{}{}
	if got := GetPlaylistName(empty); got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}

func TestGetPlaylistTrackCount(t *testing.T) {
	playlist := map[string]interface{}{
		"name":   "Playlist",
		"tracks": []interface{}{"track1", "track2", "track3"},
	}
	if got := GetPlaylistTrackCount(playlist); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}

	noTracks := map[string]interface{}{"name": "Empty"}
	if got := GetPlaylistTrackCount(noTracks); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGetPlaylistOwner(t *testing.T) {
	playlist := map[string]interface{}{"owner": "user123"}
	if got := GetPlaylistOwner(playlist); got != "user123" {
		t.Errorf("expected 'user123', got '%s'", got)
	}

	noOwner := map[string]interface{}{"name": "Playlist"}
	if got := GetPlaylistOwner(noOwner); got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
		{"abc", 3, "abc"},
		{"abcdef", 5, "ab..."},
	}

	for _, tt := range tests {
		got := Truncate(tt.input, tt.maxLen)
		if got != tt.expected {
			t.Errorf("Truncate(%q, %d): expected %q, got %q",
				tt.input, tt.maxLen, tt.expected, got)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{0, 0, 0},
		{-1, 1, -1},
		{10, 10, 10},
	}

	for _, tt := range tests {
		got := Min(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("Min(%d, %d): expected %d, got %d", tt.a, tt.b, tt.expected, got)
		}
	}
}
