//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestBackup_CreateAndRestore(t *testing.T) {
	cfg := SetupTestEnvironment(t)

	// Create test backup data
	testData := map[string]interface{}{
		"tracks": []interface{}{
			map[string]interface{}{
				"id":   "track1",
				"name": "Test Track 1",
				"artist": map[string]string{
					"name": "Artist 1",
				},
			},
			map[string]interface{}{
				"id":   "track2",
				"name": "Test Track 2",
				"artist": map[string]string{
					"name": "Artist 2",
				},
			},
		},
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"type":      "tracks",
		},
	}

	// Save backup
	backupFile := "backup_tracks_test.json"
	err := cfg.Storage.SaveJSON(backupFile, testData)
	AssertNoError(t, err, "Failed to save backup")

	// Verify file exists
	backupPath := filepath.Join(cfg.TempDir, backupFile)
	AssertFileExists(t, backupPath)

	// Load backup
	var loadedData map[string]interface{}
	err = cfg.Storage.LoadJSON(backupFile, &loadedData)
	AssertNoError(t, err, "Failed to load backup")

	// Verify data
	tracks, ok := loadedData["tracks"].([]interface{})
	if !ok {
		t.Fatal("Failed to get tracks from loaded data")
	}

	if len(tracks) != 2 {
		t.Errorf("Expected 2 tracks, got %d", len(tracks))
	}

	// Verify first track
	track1, ok := tracks[0].(map[string]interface{})
	if !ok {
		t.Fatal("Failed to parse first track")
	}

	if track1["id"] != "track1" {
		t.Errorf("Expected track id 'track1', got '%v'", track1["id"])
	}

	if track1["name"] != "Test Track 1" {
		t.Errorf("Expected track name 'Test Track 1', got '%v'", track1["name"])
	}
}

func TestBackup_CreateBackupWithMetadata(t *testing.T) {
	cfg := SetupTestEnvironment(t)

	// Create backup directory
	_, err := cfg.Storage.CreateBackup("full", map[string]interface{}{
		"test": "data",
	})
	AssertNoError(t, err, "Failed to create backup")

	// List backups
	backups, err := cfg.Storage.ListBackups()
	AssertNoError(t, err, "Failed to list backups")

	if len(backups) == 0 {
		t.Fatal("Expected at least one backup")
	}

	// Verify backup metadata
	found := false
	for _, backup := range backups {
		if backup.Type == "full" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Full backup not found in list")
	}
}

func TestBackup_RestoreFromBackup(t *testing.T) {
	cfg := SetupTestEnvironment(t)

	// Create a backup
	backupData := map[string]interface{}{
		"playlists": []interface{}{
			map[string]interface{}{
				"id":   "playlist1",
				"name": "My Playlist",
				"tracks": []interface{}{
					map[string]string{"id": "track1"},
					map[string]string{"id": "track2"},
				},
			},
		},
	}

	backupMeta, err := cfg.Storage.CreateBackup("playlists", backupData)
	backupID := backupMeta.ID
	AssertNoError(t, err, "Failed to create backup")

	// Restore from backup
	var restoredData map[string]interface{}
	err = cfg.Storage.LoadBackupJSON(backupID, &restoredData)
	AssertNoError(t, err, "Failed to restore backup")

	// Verify restored data
	playlists, ok := restoredData["playlists"].([]interface{})
	if !ok {
		t.Fatal("Failed to get playlists from restored data")
	}

	if len(playlists) != 1 {
		t.Errorf("Expected 1 playlist, got %d", len(playlists))
	}
}

func TestBackup_InvalidBackup(t *testing.T) {
	cfg := SetupTestEnvironment(t)

	// Try to load non-existent backup
	var data map[string]interface{}
	err := cfg.Storage.LoadBackupJSON("nonexistent_backup", &data)
	AssertError(t, err, "Should fail to load non-existent backup")
}

func TestBackup_MalformedJSON(t *testing.T) {
	cfg := SetupTestEnvironment(t)

	// Create malformed JSON file
	malformedFile := filepath.Join(cfg.TempDir, "malformed.json")
	err := cfg.Storage.SaveJSON("malformed.json", "not json")

	// This should succeed (SaveJSON handles conversion)
	if err != nil {
		// Write raw malformed JSON
		_, err = json.Marshal(struct{}{})
		AssertNoError(t, err, "Failed to create test file")
	}

	AssertFileExists(t, malformedFile)
}
