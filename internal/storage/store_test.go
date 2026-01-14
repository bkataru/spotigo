package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	store := NewStore("/data", "/backups")

	if store == nil {
		t.Fatal("NewStore returned nil")
	}

	if store.dataDir != "/data" {
		t.Errorf("expected dataDir '/data', got '%s'", store.dataDir)
	}

	if store.backupDir != "/backups" {
		t.Errorf("expected backupDir '/backups', got '%s'", store.backupDir)
	}
}

func TestStore_SaveAndLoadJSON(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	// Test data
	type TestData struct {
		Name    string   `json:"name"`
		Count   int      `json:"count"`
		Tags    []string `json:"tags"`
		Enabled bool     `json:"enabled"`
	}

	original := TestData{
		Name:    "Test Item",
		Count:   42,
		Tags:    []string{"tag1", "tag2", "tag3"},
		Enabled: true,
	}

	// Save
	err = store.SaveJSON("test-data.json", original)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	// Verify file exists
	expectedPath := filepath.Join(tmpDir, "test-data.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatal("SaveJSON did not create file")
	}

	// Load
	var loaded TestData
	err = store.LoadJSON("test-data.json", &loaded)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	// Compare
	if loaded.Name != original.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", original.Name, loaded.Name)
	}
	if loaded.Count != original.Count {
		t.Errorf("Count mismatch: expected %d, got %d", original.Count, loaded.Count)
	}
	if len(loaded.Tags) != len(original.Tags) {
		t.Errorf("Tags length mismatch: expected %d, got %d", len(original.Tags), len(loaded.Tags))
	}
	if loaded.Enabled != original.Enabled {
		t.Errorf("Enabled mismatch: expected %v, got %v", original.Enabled, loaded.Enabled)
	}
}

func TestStore_SaveJSON_CreateSubdirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	data := map[string]string{"key": "value"}

	// Save to nested path
	err = store.SaveJSON("subdir/nested/file.json", data)
	if err != nil {
		t.Fatalf("SaveJSON with subdirectory failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "subdir", "nested", "file.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatal("SaveJSON did not create nested directories")
	}
}

func TestStore_LoadJSON_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	var data map[string]string
	err = store.LoadJSON("nonexistent.json", &data)
	if err == nil {
		t.Error("LoadJSON should fail for nonexistent file")
	}
}

func TestStore_CreateBackup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	backupDir := filepath.Join(tmpDir, "backups")
	store := NewStore(tmpDir, backupDir)

	data := map[string]interface{}{
		"tracks":    []string{"track1", "track2"},
		"playlists": []string{"playlist1"},
	}

	metadata, err := store.CreateBackup("full", data)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Check metadata
	if metadata.ID == "" {
		t.Error("metadata ID should not be empty")
	}
	if metadata.Type != "full" {
		t.Errorf("expected type 'full', got '%s'", metadata.Type)
	}
	if metadata.Size <= 0 {
		t.Error("metadata size should be positive")
	}
	if metadata.Timestamp.IsZero() {
		t.Error("metadata timestamp should be set")
	}
	if time.Since(metadata.Timestamp) > time.Minute {
		t.Error("metadata timestamp should be recent")
	}

	// Verify backup file exists
	files, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup dir: %v", err)
	}
	if len(files) == 0 {
		t.Error("backup file was not created")
	}
}

func TestStore_ListBackups(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	backupDir := filepath.Join(tmpDir, "backups")
	store := NewStore(tmpDir, backupDir)

	// Initially empty
	backups, err := store.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected 0 backups, got %d", len(backups))
	}

	// Create some backups
	store.CreateBackup("tracks", map[string]interface{}{"data": "tracks"})
	time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	store.CreateBackup("playlists", map[string]interface{}{"data": "playlists"})
	time.Sleep(10 * time.Millisecond)
	store.CreateBackup("all", map[string]interface{}{"data": "all"})

	// List backups
	backups, err = store.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 3 {
		t.Errorf("expected 3 backups, got %d", len(backups))
	}

	// Check that each backup has valid metadata
	for _, backup := range backups {
		if backup.ID == "" {
			t.Error("backup ID should not be empty")
		}
		if backup.Size <= 0 {
			t.Error("backup size should be positive")
		}
	}
}

func TestStore_ListBackups_IgnoresNonJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	store := NewStore(tmpDir, backupDir)

	// Create a JSON backup
	store.CreateBackup("test", map[string]interface{}{"data": "test"})

	// Create non-JSON files that should be ignored
	os.WriteFile(filepath.Join(backupDir, "readme.txt"), []byte("readme"), 0644)
	os.WriteFile(filepath.Join(backupDir, "notes.md"), []byte("notes"), 0644)

	// Create a subdirectory that should be ignored
	os.MkdirAll(filepath.Join(backupDir, "subdir"), 0755)

	backups, err := store.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("expected 1 backup (ignoring non-JSON), got %d", len(backups))
	}
}

func TestStore_Exists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	// Create a file
	testFile := "test-exists.json"
	store.SaveJSON(testFile, map[string]string{"test": "data"})

	// Test exists
	if !store.Exists(testFile) {
		t.Error("Exists should return true for existing file")
	}

	// Test not exists
	if store.Exists("nonexistent.json") {
		t.Error("Exists should return false for nonexistent file")
	}
}

func TestBackupMetadata_Structure(t *testing.T) {
	meta := BackupMetadata{
		ID:        "20240115-120000",
		Timestamp: time.Now(),
		Type:      "full",
		Items:     100,
		Size:      1024,
	}

	if meta.ID == "" {
		t.Error("ID should be set")
	}
	if meta.Type == "" {
		t.Error("Type should be set")
	}
	if meta.Items != 100 {
		t.Errorf("expected Items 100, got %d", meta.Items)
	}
	if meta.Size != 1024 {
		t.Errorf("expected Size 1024, got %d", meta.Size)
	}
}

func TestStore_SaveJSON_ComplexData(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spotigo-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	// Complex nested structure
	data := map[string]interface{}{
		"tracks": []map[string]interface{}{
			{
				"id":   "track1",
				"name": "Song One",
				"artists": []map[string]string{
					{"name": "Artist A"},
					{"name": "Artist B"},
				},
			},
			{
				"id":   "track2",
				"name": "Song Two",
				"artists": []map[string]string{
					{"name": "Artist C"},
				},
			},
		},
		"metadata": map[string]interface{}{
			"count":     2,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	err = store.SaveJSON("complex.json", data)
	if err != nil {
		t.Fatalf("SaveJSON with complex data failed: %v", err)
	}

	// Load it back
	var loaded map[string]interface{}
	err = store.LoadJSON("complex.json", &loaded)
	if err != nil {
		t.Fatalf("LoadJSON with complex data failed: %v", err)
	}

	// Basic verification
	if loaded["tracks"] == nil {
		t.Error("tracks should be present")
	}
	if loaded["metadata"] == nil {
		t.Error("metadata should be present")
	}
}

// BenchmarkSaveJSON benchmarks JSON saving
func BenchmarkSaveJSON(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "spotigo-bench-*")
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir, filepath.Join(tmpDir, "backups"))

	data := map[string]interface{}{
		"tracks": make([]map[string]string, 100),
	}
	for i := 0; i < 100; i++ {
		data["tracks"].([]map[string]string)[i] = map[string]string{
			"id":   "track",
			"name": "Song",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.SaveJSON("bench.json", data)
	}
}
