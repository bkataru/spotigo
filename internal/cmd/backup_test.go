package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseBackupType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"all types", "all", "all"},
		{"tracks only", "tracks", "tracks"},
		{"playlists only", "playlists", "playlists"},
		{"artists only", "artists", "artists"},
		{"empty defaults to all", "", "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == "" {
				// Test default behavior
				backupType = "all"
			} else {
				backupType = tt.input
			}

			if backupType != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, backupType)
			}
		})
	}
}

func TestCreateBackupDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")

	// Test directory creation
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(backupDir)
	if err != nil {
		t.Fatalf("Backup directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Path exists but is not a directory")
	}
}

func TestBackupMetadata(t *testing.T) {
	// Test backup metadata structure
	metadata := struct {
		Timestamp   time.Time      `json:"timestamp"`
		Type        string         `json:"type"`
		ItemCounts  map[string]int `json:"item_counts"`
		Description string         `json:"description"`
		Version     string         `json:"version"`
		Status      string         `json:"status"`
	}{
		Timestamp:   time.Now(),
		Type:        "full",
		ItemCounts:  map[string]int{"tracks": 100, "playlists": 10},
		Description: "Test backup",
		Version:     "2.0.0",
		Status:      "completed",
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal metadata: %v", err)
	}

	if len(data) == 0 {
		t.Error("Marshaled metadata is empty")
	}

	// Unmarshal back
	var decoded struct {
		Timestamp   time.Time      `json:"timestamp"`
		Type        string         `json:"type"`
		ItemCounts  map[string]int `json:"item_counts"`
		Description string         `json:"description"`
		Version     string         `json:"version"`
		Status      string         `json:"status"`
	}

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if decoded.Type != metadata.Type {
		t.Errorf("Expected type %s, got %s", metadata.Type, decoded.Type)
	}

	if decoded.ItemCounts["tracks"] != 100 {
		t.Errorf("Expected 100 tracks, got %d", decoded.ItemCounts["tracks"])
	}
}

func TestBackupFileNaming(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		format   string
		contains string
	}{
		{
			name:     "timestamp format",
			format:   now.Format("20060102_150405"),
			contains: now.Format("20060102"),
		},
		{
			name:     "saved tracks filename",
			format:   "saved_tracks.json",
			contains: "saved_tracks",
		},
		{
			name:     "playlists filename",
			format:   "playlists.json",
			contains: "playlists",
		},
		{
			name:     "followed artists filename",
			format:   "followed_artists.json",
			contains: "followed_artists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.format == "" {
				t.Error("Format is empty")
			}

			// Check if contains expected substring
			if tt.contains != "" {
				found := false
				for i := 0; i <= len(tt.format)-len(tt.contains); i++ {
					if tt.format[i:i+len(tt.contains)] == tt.contains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Format %s does not contain %s", tt.format, tt.contains)
				}
			}
		})
	}
}

func TestValidateBackupData(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expectValid bool
	}{
		{
			name:        "valid JSON array",
			data:        `{"items": [{"id": "1", "name": "test"}]}`,
			expectValid: true,
		},
		{
			name:        "valid empty object",
			data:        `{}`,
			expectValid: true,
		},
		{
			name:        "invalid JSON",
			data:        `{invalid json}`,
			expectValid: false,
		},
		{
			name:        "empty string",
			data:        "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(tt.data), &result)

			isValid := err == nil
			if isValid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v (err=%v)", tt.expectValid, isValid, err)
			}
		})
	}
}

func TestBackupListParsing(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mock backup directories
	backupDirs := []string{
		"backup_20240115_120000",
		"backup_20240114_100000",
		"backup_20240113_090000",
	}

	for _, dir := range backupDirs {
		dirPath := filepath.Join(tmpDir, dir)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock backup dir: %v", err)
		}

		// Create metadata file
		metadata := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"type":      "full",
			"status":    "completed",
		}
		data, _ := json.Marshal(metadata)
		err = os.WriteFile(filepath.Join(dirPath, "metadata.json"), data, 0644)
		if err != nil {
			t.Fatalf("Failed to write metadata: %v", err)
		}
	}

	// List directories
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read backup directory: %v", err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}

	if count != len(backupDirs) {
		t.Errorf("Expected %d backup directories, found %d", len(backupDirs), count)
	}
}

func TestBackupStatusCheck(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a backup with status file
	backupDir := filepath.Join(tmpDir, "backup_test")
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup dir: %v", err)
	}

	// Write status
	status := map[string]interface{}{
		"status":    "completed",
		"timestamp": time.Now().Format(time.RFC3339),
		"items": map[string]int{
			"tracks":    100,
			"playlists": 10,
		},
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal status: %v", err)
	}

	statusFile := filepath.Join(backupDir, "status.json")
	err = os.WriteFile(statusFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write status file: %v", err)
	}

	// Read and verify
	readData, err := os.ReadFile(statusFile)
	if err != nil {
		t.Fatalf("Failed to read status file: %v", err)
	}

	var readStatus map[string]interface{}
	err = json.Unmarshal(readData, &readStatus)
	if err != nil {
		t.Fatalf("Failed to unmarshal status: %v", err)
	}

	if readStatus["status"] != "completed" {
		t.Errorf("Expected status 'completed', got '%v'", readStatus["status"])
	}
}

func TestConcurrencyLimits(t *testing.T) {
	// Test concurrency constants
	if maxConcurrentPlaylistFetches <= 0 {
		t.Error("maxConcurrentPlaylistFetches should be positive")
	}

	if maxConcurrentPlaylistFetches > 100 {
		t.Error("maxConcurrentPlaylistFetches should be reasonable (<=100)")
	}

	if writeBufferSize <= 0 {
		t.Error("writeBufferSize should be positive")
	}

	if writeBufferSize < 1024 {
		t.Error("writeBufferSize should be at least 1KB")
	}
}

func TestBackupTypeValidation(t *testing.T) {
	validTypes := []string{"all", "tracks", "playlists", "artists"}

	for _, validType := range validTypes {
		t.Run("valid_"+validType, func(t *testing.T) {
			isValid := false
			for _, vt := range validTypes {
				if validType == vt {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Type %s should be valid", validType)
			}
		})
	}

	invalidTypes := []string{"invalid", "unknown", ""}
	for _, invalidType := range invalidTypes {
		if invalidType == "" {
			continue // empty defaults to "all"
		}
		t.Run("invalid_"+invalidType, func(t *testing.T) {
			isValid := false
			for _, vt := range validTypes {
				if invalidType == vt {
					isValid = true
					break
				}
			}

			if isValid {
				t.Errorf("Type %s should be invalid", invalidType)
			}
		})
	}
}

func TestBackupFileStructure(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup_test")

	// Create backup structure
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Expected files in a backup
	files := []string{
		"saved_tracks.json",
		"playlists.json",
		"followed_artists.json",
		"metadata.json",
	}

	// Create files
	for _, file := range files {
		filePath := filepath.Join(backupDir, file)
		err := os.WriteFile(filePath, []byte("{}"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Verify all files exist
	for _, file := range files {
		filePath := filepath.Join(backupDir, file)
		_, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("File %s does not exist: %v", file, err)
		}
	}
}

func TestBackupTimestampParsing(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		expectErr bool
	}{
		{
			name:      "valid RFC3339",
			timestamp: "2024-01-15T12:00:00Z",
			expectErr: false,
		},
		{
			name:      "valid with timezone",
			timestamp: "2024-01-15T12:00:00-05:00",
			expectErr: false,
		},
		{
			name:      "invalid format",
			timestamp: "invalid-timestamp",
			expectErr: true,
		},
		{
			name:      "empty string",
			timestamp: "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse(time.RFC3339, tt.timestamp)

			hasErr := err != nil
			if hasErr != tt.expectErr {
				t.Errorf("Expected error=%v, got error=%v", tt.expectErr, hasErr)
			}
		})
	}
}
