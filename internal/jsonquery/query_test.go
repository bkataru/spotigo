package jsonquery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine("./testdata")
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}
	if engine.dataDir != "./testdata" {
		t.Errorf("Expected dataDir './testdata', got %s", engine.dataDir)
	}
}

func TestLoadData(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test data file
	testData := []map[string]interface{}{
		{"id": 1, "name": "Test 1"},
		{"id": 2, "name": "Test 2"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(testFile, data, 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	engine := NewEngine(tmpDir)

	// Test loading data
	loaded, err := engine.loadData("test.json")
	if err != nil {
		t.Fatalf("loadData() error = %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("Expected 2 items, got %d", len(loaded))
	}

	// Test cache
	loaded2, err := engine.loadData("test.json")
	if err != nil {
		t.Fatalf("loadData() cached error = %v", err)
	}

	if len(loaded2) != 2 {
		t.Errorf("Expected 2 cached items, got %d", len(loaded2))
	}
}

func TestClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"id": 1},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "test.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)
	engine.loadData("test.json")

	if len(engine.cache) == 0 {
		t.Error("Cache should not be empty after loading")
	}

	engine.ClearCache()

	if len(engine.cache) != 0 {
		t.Error("Cache should be empty after ClearCache()")
	}
}

func TestCountOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"name": "Track 1"},
		{"name": "Track 2"},
		{"name": "Track 3"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)
	result := engine.Execute(Query{
		Source:    "tracks.json",
		Operation: "count",
	})

	if result.Error != "" {
		t.Errorf("count operation failed: %s", result.Error)
	}

	if result.Count != 3 {
		t.Errorf("Expected count 3, got %d", result.Count)
	}
}

func TestSelectOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"id": 1, "name": "Track 1"},
		{"id": 2, "name": "Track 2"},
		{"id": 3, "name": "Track 3"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	tests := []struct {
		name          string
		query         Query
		expectedCount int
	}{
		{
			name: "select all",
			query: Query{
				Source:    "tracks.json",
				Operation: "select",
			},
			expectedCount: 3,
		},
		{
			name: "select with limit",
			query: Query{
				Source:    "tracks.json",
				Operation: "select",
				Limit:     2,
			},
			expectedCount: 2,
		},
		{
			name: "select with offset",
			query: Query{
				Source:    "tracks.json",
				Operation: "select",
				Offset:    1,
			},
			expectedCount: 2,
		},
		{
			name: "select specific field",
			query: Query{
				Source:    "tracks.json",
				Operation: "select",
				Field:     "name",
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Execute(tt.query)
			if result.Error != "" {
				t.Errorf("select operation failed: %s", result.Error)
			}
			if result.Count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, result.Count)
			}
		})
	}
}

func TestFilterOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"id": 1, "name": "Rock Song", "duration": 180},
		{"id": 2, "name": "Pop Song", "duration": 210},
		{"id": 3, "name": "Rock Anthem", "duration": 240},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	tests := []struct {
		name          string
		filters       []Filter
		expectedCount int
	}{
		{
			name: "filter equals",
			filters: []Filter{
				{Field: "name", Operator: "eq", Value: "Rock Song"},
			},
			expectedCount: 1,
		},
		{
			name: "filter contains",
			filters: []Filter{
				{Field: "name", Operator: "contains", Value: "Rock"},
			},
			expectedCount: 2,
		},
		{
			name: "filter greater than",
			filters: []Filter{
				{Field: "duration", Operator: "gt", Value: 200},
			},
			expectedCount: 2,
		},
		{
			name: "multiple filters",
			filters: []Filter{
				{Field: "name", Operator: "contains", Value: "Rock"},
				{Field: "duration", Operator: "gte", Value: 240},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Execute(Query{
				Source:    "tracks.json",
				Operation: "filter",
				Filters:   tt.filters,
			})
			if result.Error != "" {
				t.Errorf("filter operation failed: %s", result.Error)
			}
			if result.Count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, result.Count)
			}
		})
	}
}

func TestSearchOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"artist": "The Beatles", "track": "Hey Jude"},
		{"artist": "Beatles Cover Band", "track": "Let It Be"},
		{"artist": "Queen", "track": "Bohemian Rhapsody"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	result := engine.Execute(Query{
		Source:     "tracks.json",
		Operation:  "search",
		SearchTerm: "beatles",
	})

	if result.Error != "" {
		t.Errorf("search operation failed: %s", result.Error)
	}

	if result.Count != 2 {
		t.Errorf("Expected 2 results for 'beatles', got %d", result.Count)
	}
}

func TestSortOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"name": "Track C", "plays": 100},
		{"name": "Track A", "plays": 300},
		{"name": "Track B", "plays": 200},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	tests := []struct {
		name          string
		sortBy        string
		sortOrder     string
		expectedFirst interface{}
	}{
		{
			name:          "sort by name ascending",
			sortBy:        "name",
			sortOrder:     "asc",
			expectedFirst: "Track A",
		},
		{
			name:          "sort by name descending",
			sortBy:        "name",
			sortOrder:     "desc",
			expectedFirst: "Track C",
		},
		{
			name:          "sort by plays descending",
			sortBy:        "plays",
			sortOrder:     "desc",
			expectedFirst: float64(300),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Execute(Query{
				Source:    "tracks.json",
				Operation: "sort",
				SortBy:    tt.sortBy,
				SortOrder: tt.sortOrder,
			})

			if result.Error != "" {
				t.Errorf("sort operation failed: %s", result.Error)
			}

			if result.Data != nil {
				items, ok := result.Data.([]interface{})
				if !ok {
					t.Fatal("Expected data to be array")
				}
				if len(items) > 0 {
					first := items[0].(map[string]interface{})
					if first[tt.sortBy] != tt.expectedFirst {
						t.Errorf("Expected first item %s to be %v, got %v",
							tt.sortBy, tt.expectedFirst, first[tt.sortBy])
					}
				}
			}
		})
	}
}

func TestDistinctOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"artist": "Beatles"},
		{"artist": "Queen"},
		{"artist": "Beatles"},
		{"artist": "Led Zeppelin"},
		{"artist": "Queen"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	result := engine.Execute(Query{
		Source:    "tracks.json",
		Operation: "distinct",
		Field:     "artist",
	})

	if result.Error != "" {
		t.Errorf("distinct operation failed: %s", result.Error)
	}

	if result.Count != 3 {
		t.Errorf("Expected 3 distinct artists, got %d", result.Count)
	}
}

func TestAggregateOperations(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"name": "Track 1", "duration": 180},
		{"name": "Track 2", "duration": 200},
		{"name": "Track 3", "duration": 220},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	tests := []struct {
		name     string
		aggFunc  string
		expected float64
	}{
		{
			name:     "sum",
			aggFunc:  "sum",
			expected: 600,
		},
		{
			name:     "avg",
			aggFunc:  "avg",
			expected: 200,
		},
		{
			name:     "min",
			aggFunc:  "min",
			expected: 180,
		},
		{
			name:     "max",
			aggFunc:  "max",
			expected: 220,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Execute(Query{
				Source:    "tracks.json",
				Operation: "aggregate",
				AggFunc:   tt.aggFunc,
				Field:     "duration",
			})

			if result.Error != "" {
				t.Errorf("aggregate operation failed: %s", result.Error)
			}

			if val, ok := result.Data.(float64); ok {
				if val != tt.expected {
					t.Errorf("Expected %s = %v, got %v", tt.aggFunc, tt.expected, val)
				}
			} else {
				t.Errorf("Expected numeric result, got %T", result.Data)
			}
		})
	}
}

func TestGroupOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"artist": "Beatles", "track": "Hey Jude"},
		{"artist": "Queen", "track": "Bohemian Rhapsody"},
		{"artist": "Beatles", "track": "Let It Be"},
		{"artist": "Beatles", "track": "Yesterday"},
		{"artist": "Queen", "track": "We Will Rock You"},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	result := engine.Execute(Query{
		Source:    "tracks.json",
		Operation: "aggregate",
		AggFunc:   "group",
		GroupBy:   "artist",
	})

	if result.Error != "" {
		t.Errorf("group operation failed: %s", result.Error)
	}

	if result.Count != 2 {
		t.Errorf("Expected 2 groups, got %d", result.Count)
	}
}

func TestStatsOperation(t *testing.T) {
	tmpDir := t.TempDir()

	testData := []map[string]interface{}{
		{"duration": 180},
		{"duration": 200},
		{"duration": 220},
	}
	data, _ := json.Marshal(testData)
	testFile := filepath.Join(tmpDir, "tracks.json")
	os.WriteFile(testFile, data, 0600)

	engine := NewEngine(tmpDir)

	result := engine.Execute(Query{
		Source:    "tracks.json",
		Operation: "stats",
		Field:     "duration",
	})

	if result.Error != "" {
		t.Errorf("stats operation failed: %s", result.Error)
	}

	stats, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected stats to be a map")
	}

	expectedFields := []string{"total_count", "numeric_count", "sum", "avg", "min", "max"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Expected stats field '%s' not found", field)
		}
	}
}

func TestGetFieldValue(t *testing.T) {
	data := map[string]interface{}{
		"name": "Test",
		"track": map[string]interface{}{
			"title": "Song Title",
			"artists": []interface{}{
				map[string]interface{}{"name": "Artist 1"},
				map[string]interface{}{"name": "Artist 2"},
			},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{
			name:     "simple field",
			path:     "name",
			expected: "Test",
		},
		{
			name:     "nested field",
			path:     "track.title",
			expected: "Song Title",
		},
		{
			name:     "array index",
			path:     "track.artists.0.name",
			expected: "Artist 1",
		},
		{
			name:     "non-existent field",
			path:     "missing",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFieldValue(data, tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected int
	}{
		{
			name:     "equal numbers",
			a:        10,
			b:        10,
			expected: 0,
		},
		{
			name:     "less than",
			a:        5,
			b:        10,
			expected: -1,
		},
		{
			name:     "greater than",
			a:        10,
			b:        5,
			expected: 1,
		},
		{
			name:     "equal strings",
			a:        "abc",
			b:        "abc",
			expected: 0,
		},
		{
			name:     "string less than",
			a:        "abc",
			b:        "def",
			expected: -1,
		},
		{
			name:     "nil values",
			a:        nil,
			b:        nil,
			expected: 0,
		},
		{
			name:     "nil less than value",
			a:        nil,
			b:        10,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareValues(%v, %v) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMusicQueryHelper(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test music data files
	tracks := []map[string]interface{}{
		{
			"added_at": "2024-01-01T00:00:00Z",
			"track": map[string]interface{}{
				"name": "Song 1",
				"artists": []interface{}{
					map[string]interface{}{"name": "Artist 1"},
				},
			},
		},
		{
			"added_at": "2024-01-02T00:00:00Z",
			"track": map[string]interface{}{
				"name": "Song 2",
				"artists": []interface{}{
					map[string]interface{}{"name": "Artist 1"},
				},
			},
		},
	}
	tracksData, _ := json.Marshal(tracks)
	os.WriteFile(filepath.Join(tmpDir, "saved_tracks.json"), tracksData, 0600)

	playlists := []map[string]interface{}{
		{"name": "My Playlist", "id": "1"},
	}
	playlistsData, _ := json.Marshal(playlists)
	os.WriteFile(filepath.Join(tmpDir, "playlists.json"), playlistsData, 0600)

	artists := []map[string]interface{}{
		{"name": "Artist 1", "id": "1"},
		{"name": "Artist 2", "id": "2"},
	}
	artistsData, _ := json.Marshal(artists)
	os.WriteFile(filepath.Join(tmpDir, "followed_artists.json"), artistsData, 0600)

	helper := NewMusicQueryHelper(tmpDir)

	t.Run("GetLibraryStats", func(t *testing.T) {
		result := helper.GetLibraryStats()
		if result.Error != "" {
			t.Errorf("GetLibraryStats failed: %s", result.Error)
		}

		stats, ok := result.Data.(map[string]interface{})
		if !ok {
			t.Fatal("Expected stats to be a map")
		}

		if stats["saved_tracks"] != 2 {
			t.Errorf("Expected 2 saved tracks, got %v", stats["saved_tracks"])
		}
		if stats["playlists"] != 1 {
			t.Errorf("Expected 1 playlist, got %v", stats["playlists"])
		}
		if stats["followed_artists"] != 2 {
			t.Errorf("Expected 2 followed artists, got %v", stats["followed_artists"])
		}
	})

	t.Run("GetRecentlyAddedTracks", func(t *testing.T) {
		result := helper.GetRecentlyAddedTracks(1)
		if result.Error != "" {
			t.Errorf("GetRecentlyAddedTracks failed: %s", result.Error)
		}
		if result.Count != 1 {
			t.Errorf("Expected 1 track, got %d", result.Count)
		}
	})

	t.Run("SearchAllData", func(t *testing.T) {
		result := helper.SearchAllData("Artist", 10)
		if result.Error != "" {
			t.Errorf("SearchAllData failed: %s", result.Error)
		}
		if result.Count == 0 {
			t.Error("Expected to find results for 'Artist'")
		}
	})
}

func BenchmarkExecuteCount(b *testing.B) {
	tmpDir := b.TempDir()

	testData := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		testData[i] = map[string]interface{}{
			"id":   i,
			"name": "Item " + string(rune(i)),
		}
	}
	data, _ := json.Marshal(testData)
	os.WriteFile(filepath.Join(tmpDir, "test.json"), data, 0600)

	engine := NewEngine(tmpDir)
	query := Query{
		Source:    "test.json",
		Operation: "count",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Execute(query)
	}
}

func BenchmarkExecuteFilter(b *testing.B) {
	tmpDir := b.TempDir()

	testData := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		testData[i] = map[string]interface{}{
			"id":    i,
			"value": i * 10,
		}
	}
	data, _ := json.Marshal(testData)
	os.WriteFile(filepath.Join(tmpDir, "test.json"), data, 0600)

	engine := NewEngine(tmpDir)
	query := Query{
		Source:    "test.json",
		Operation: "filter",
		Filters: []Filter{
			{Field: "value", Operator: "gt", Value: 5000},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Execute(query)
	}
}
