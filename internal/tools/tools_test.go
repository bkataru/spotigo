package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bkataru/spotigo/internal/ollama"
)

// setupTestData creates temporary test music data files
func setupTestData(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Create saved_tracks.json
	savedTracks := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"track": map[string]interface{}{
					"id":   "track1",
					"name": "Bohemian Rhapsody",
					"artists": []map[string]interface{}{
						{"name": "Queen"},
					},
					"album": map[string]interface{}{
						"name":         "A Night at the Opera",
						"release_date": "1975-11-21",
					},
					"duration_ms": 354000,
					"popularity":  95,
				},
				"added_at": "2024-01-15T10:30:00Z",
			},
			{
				"track": map[string]interface{}{
					"id":   "track2",
					"name": "Stairway to Heaven",
					"artists": []map[string]interface{}{
						{"name": "Led Zeppelin"},
					},
					"album": map[string]interface{}{
						"name":         "Led Zeppelin IV",
						"release_date": "1971-11-08",
					},
					"duration_ms": 482000,
					"popularity":  92,
				},
				"added_at": "2024-01-14T09:00:00Z",
			},
			{
				"track": map[string]interface{}{
					"id":   "track3",
					"name": "We Will Rock You",
					"artists": []map[string]interface{}{
						{"name": "Queen"},
					},
					"album": map[string]interface{}{
						"name":         "News of the World",
						"release_date": "1977-10-28",
					},
					"duration_ms": 122000,
					"popularity":  90,
				},
				"added_at": "2024-01-16T14:20:00Z",
			},
		},
	}

	savedTracksData, err := json.MarshalIndent(savedTracks, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal saved_tracks: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "saved_tracks.json"), savedTracksData, 0644); err != nil {
		t.Fatalf("Failed to write saved_tracks.json: %v", err)
	}

	// Create playlists.json
	playlists := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"id":   "playlist1",
				"name": "Rock Classics",
				"owner": map[string]interface{}{
					"display_name": "testuser",
				},
				"tracks": map[string]interface{}{
					"total": 25,
				},
				"description": "Classic rock hits",
			},
			{
				"id":   "playlist2",
				"name": "Workout Mix",
				"owner": map[string]interface{}{
					"display_name": "testuser",
				},
				"tracks": map[string]interface{}{
					"total": 30,
				},
				"description": "High energy tracks",
			},
		},
	}

	playlistsData, err := json.MarshalIndent(playlists, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal playlists: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "playlists.json"), playlistsData, 0644); err != nil {
		t.Fatalf("Failed to write playlists.json: %v", err)
	}

	// Create followed_artists.json
	followedArtists := map[string]interface{}{
		"artists": map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":   "artist1",
					"name": "Queen",
					"genres": []string{
						"classic rock",
						"rock",
					},
					"followers": map[string]interface{}{
						"total": 35000000,
					},
					"popularity": 88,
				},
				{
					"id":   "artist2",
					"name": "Led Zeppelin",
					"genres": []string{
						"classic rock",
						"hard rock",
					},
					"followers": map[string]interface{}{
						"total": 25000000,
					},
					"popularity": 82,
				},
			},
		},
	}

	followedArtistsData, err := json.MarshalIndent(followedArtists, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal followed_artists: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "followed_artists.json"), followedArtistsData, 0644); err != nil {
		t.Fatalf("Failed to write followed_artists.json: %v", err)
	}

	return tmpDir
}

func TestNewMusicTools(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	if tools == nil {
		t.Fatal("NewMusicTools returned nil")
	}

	if tools.queryHelper == nil {
		t.Fatal("queryHelper is nil")
	}
}

func TestGetToolDefinitions(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	defs := tools.GetToolDefinitions()

	if len(defs) == 0 {
		t.Fatal("GetToolDefinitions returned empty slice")
	}

	// Expected tool names
	expectedTools := map[string]bool{
		"get_library_stats":         true,
		"search_tracks":             true,
		"get_tracks_by_artist":      true,
		"get_recently_added_tracks": true,
		"get_all_artists":           true,
		"get_playlist_by_name":      true,
		"query_music_data":          true,
	}

	foundTools := make(map[string]bool)
	for _, def := range defs {
		if def.Type != "function" {
			t.Errorf("Expected type 'function', got '%s'", def.Type)
		}

		foundTools[def.Function.Name] = true

		// Verify each tool has required fields
		if def.Function.Name == "" {
			t.Error("Tool has empty name")
		}
		if def.Function.Description == "" {
			t.Error("Tool has empty description")
		}
		if def.Function.Parameters == nil {
			t.Errorf("Tool %s has nil parameters", def.Function.Name)
		}
	}

	// Verify all expected tools are present
	for expectedTool := range expectedTools {
		if !foundTools[expectedTool] {
			t.Errorf("Missing expected tool: %s", expectedTool)
		}
	}
}

func TestExecuteGetLibraryStats(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	result, err := tools.executeGetLibraryStats()
	if err != nil {
		t.Fatalf("executeGetLibraryStats failed: %v", err)
	}

	var statsResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &statsResult); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify result structure
	if statsResult["error"] != nil && statsResult["error"] != "" {
		t.Errorf("Expected no error, got: %v", statsResult["error"])
	}

	data, ok := statsResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Result data is not a map")
	}

	// Check expected stats fields
	if data["saved_tracks"] == nil {
		t.Error("Missing saved_tracks in stats")
	}
	if data["playlists"] == nil {
		t.Error("Missing playlists in stats")
	}
	if data["followed_artists"] == nil {
		t.Error("Missing followed_artists in stats")
	}
}

func TestExecuteSearchTracks(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
		checkResult func(t *testing.T, result string)
	}{
		{
			name: "search for Queen",
			args: map[string]interface{}{
				"query": "Queen",
				"limit": float64(10),
			},
			expectError: false,
			checkResult: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Result is empty")
				}
				// Result should contain tracks by Queen
				var resultData map[string]interface{}
				if err := json.Unmarshal([]byte(result), &resultData); err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
				}
			},
		},
		{
			name: "search with default limit",
			args: map[string]interface{}{
				"query": "Bohemian",
			},
			expectError: false,
			checkResult: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Result is empty")
				}
			},
		},
		{
			name:        "missing query parameter",
			args:        map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.executeSearchTracks(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestExecuteGetTracksByArtist(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
		checkResult func(t *testing.T, result string)
	}{
		{
			name: "get tracks by Queen",
			args: map[string]interface{}{
				"artist_name": "Queen",
			},
			expectError: false,
			checkResult: func(t *testing.T, result string) {
				var resultData map[string]interface{}
				if err := json.Unmarshal([]byte(result), &resultData); err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
				}
			},
		},
		{
			name:        "missing artist_name parameter",
			args:        map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.executeGetTracksByArtist(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestExecuteGetRecentlyAddedTracks(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
	}{
		{
			name: "with limit",
			args: map[string]interface{}{
				"limit": float64(5),
			},
			expectError: false,
		},
		{
			name:        "with default limit",
			args:        map[string]interface{}{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.executeGetRecentlyAddedTracks(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == "" {
				t.Error("Result is empty")
			}
		})
	}
}

func TestExecuteGetAllArtists(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	result, err := tools.executeGetAllArtists()
	if err != nil {
		t.Fatalf("executeGetAllArtists failed: %v", err)
	}

	if result == "" {
		t.Error("Result is empty")
	}

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}
}

func TestExecuteGetPlaylistByName(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
	}{
		{
			name: "find Rock Classics",
			args: map[string]interface{}{
				"playlist_name": "Rock Classics",
			},
			expectError: false,
		},
		{
			name:        "missing playlist_name parameter",
			args:        map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.executeGetPlaylistByName(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == "" {
				t.Error("Result is empty")
			}
		})
	}
}

func TestExecuteQueryMusicData(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
	}{
		{
			name: "select operation",
			args: map[string]interface{}{
				"source":    "saved_tracks.json",
				"operation": "select",
				"limit":     float64(5),
			},
			expectError: false,
		},
		{
			name: "count operation",
			args: map[string]interface{}{
				"source":    "saved_tracks.json",
				"operation": "count",
			},
			expectError: false,
		},
		{
			name: "filter operation",
			args: map[string]interface{}{
				"source":    "saved_tracks.json",
				"operation": "filter",
				"filters": []interface{}{
					map[string]interface{}{
						"field":    "track.popularity",
						"operator": "gte",
						"value":    float64(90),
					},
				},
			},
			expectError: false,
		},
		{
			name: "sort operation",
			args: map[string]interface{}{
				"source":     "saved_tracks.json",
				"operation":  "select",
				"sort_by":    "track.popularity",
				"sort_order": "desc",
				"limit":      float64(3),
			},
			expectError: false,
		},
		{
			name:        "missing source",
			args:        map[string]interface{}{},
			expectError: true,
		},
		{
			name: "missing operation",
			args: map[string]interface{}{
				"source": "saved_tracks.json",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.executeQueryMusicData(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == "" {
				t.Error("Result is empty")
			}
		})
	}
}

func TestExecuteToolCall(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	tests := []struct {
		name        string
		toolCall    ollama.ToolCall
		expectError bool
		checkResult func(t *testing.T, result string)
	}{
		{
			name: "get_library_stats",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "get_library_stats",
					Arguments: "{}",
				},
			},
			expectError: false,
			checkResult: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Result is empty")
				}
			},
		},
		{
			name: "search_tracks",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "search_tracks",
					Arguments: `{"query": "Queen", "limit": 10}`,
				},
			},
			expectError: false,
			checkResult: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Result is empty")
				}
			},
		},
		{
			name: "unknown tool",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "unknown_tool",
					Arguments: "{}",
				},
			},
			expectError: true,
		},
		{
			name: "invalid arguments JSON",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "get_library_stats",
					Arguments: "invalid json",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tools.ExecuteToolCall(tt.toolCall)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestExecuteToolCallWithEmptyDataDir(t *testing.T) {
	// Create empty data directory
	tmpDir := t.TempDir()

	tools := NewMusicTools(tmpDir)

	// Try to execute a tool that requires data
	toolCall := ollama.ToolCall{
		Function: ollama.FunctionCall{
			Name:      "get_library_stats",
			Arguments: "{}",
		},
	}

	result, err := tools.ExecuteToolCall(toolCall)

	// Should not panic, but may return empty results or error in the result JSON
	if err != nil {
		t.Logf("Got expected error with empty data dir: %v", err)
	}

	if result == "" {
		t.Logf("Got empty result with empty data dir (expected)")
	}
}

func TestToolParameterSchemas(t *testing.T) {
	dataDir := setupTestData(t)
	tools := NewMusicTools(dataDir)

	defs := tools.GetToolDefinitions()

	for _, def := range defs {
		t.Run(def.Function.Name, func(t *testing.T) {
			params := def.Function.Parameters

			// Verify type is object
			if params["type"] != "object" {
				t.Errorf("Expected type 'object', got '%v'", params["type"])
			}

			// Verify properties exist
			if _, ok := params["properties"]; !ok {
				t.Error("Missing properties field")
			}

			// Verify required field exists (can be empty array)
			if _, ok := params["required"]; !ok {
				t.Error("Missing required field")
			}

			// Check that properties is a map
			props, ok := params["properties"].(map[string]interface{})
			if !ok {
				t.Fatal("Properties is not a map")
			}

			// Verify each property has type and description
			for propName, propValue := range props {
				propMap, ok := propValue.(map[string]interface{})
				if !ok {
					t.Errorf("Property %s is not a map", propName)
					continue
				}

				if propMap["type"] == nil {
					t.Errorf("Property %s missing type", propName)
				}

				if propMap["description"] == nil || propMap["description"] == "" {
					t.Errorf("Property %s missing description", propName)
				}
			}
		})
	}
}
