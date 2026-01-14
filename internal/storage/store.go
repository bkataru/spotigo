// Package storage handles local data persistence
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Store handles local file storage for Spotigo data
type Store struct {
	dataDir   string
	backupDir string
}

// NewStore creates a new storage instance
func NewStore(dataDir, backupDir string) *Store {
	return &Store{
		dataDir:   dataDir,
		backupDir: backupDir,
	}
}

// BackupMetadata holds information about a backup
type BackupMetadata struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // full, playlists, tracks, etc.
	Items     int       `json:"items"`
	Size      int64     `json:"size"`
}

// SaveJSON saves data as JSON to the specified path
func (s *Store) SaveJSON(filename string, data interface{}) error {
	// Clean filename to prevent path traversal
	cleanFilename := filepath.Clean(filename)
	path := filepath.Join(s.dataDir, cleanFilename)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: failed to close file %s: %v\n", path, closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// LoadJSON loads JSON data from the specified path
func (s *Store) LoadJSON(filename string, target interface{}) error {
	// Clean filename to prevent path traversal
	cleanFilename := filepath.Clean(filename)
	path := filepath.Join(s.dataDir, cleanFilename)

	file, err := os.Open(path) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: failed to close file %s: %v\n", path, closeErr)
		}
	}()

	if err := json.NewDecoder(file).Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

// CreateBackup creates a timestamped backup
func (s *Store) CreateBackup(backupType string, data interface{}) (*BackupMetadata, error) {
	timestamp := time.Now()
	timestampStr := timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", backupType, timestampStr)
	path := filepath.Join(s.backupDir, filename)

	// Ensure backup directory exists
	if err := os.MkdirAll(s.backupDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	file, err := os.Create(path) // #nosec G304 - path is constructed from controlled backupDir
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: failed to close backup file %s: %v\n", path, closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if encErr := encoder.Encode(data); encErr != nil {
		return nil, fmt.Errorf("failed to encode backup: %w", encErr)
	}

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	// Use filename as ID for consistency with ListBackups
	return &BackupMetadata{
		ID:        filename,
		Timestamp: timestamp,
		Type:      backupType,
		Size:      info.Size(),
	}, nil
}

// ListBackups returns all available backups
func (s *Store) ListBackups() ([]BackupMetadata, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupMetadata{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	backups := make([]BackupMetadata, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Parse backup type from filename (format: type-YYYYMMDD-HHMMSS.json)
		backupType := ""
		name := entry.Name()
		if idx := strings.LastIndex(name, "-"); idx > 0 {
			// Find the second-to-last dash (before timestamp)
			prefix := name[:idx]
			if idx2 := strings.LastIndex(prefix, "-"); idx2 > 0 {
				backupType = prefix[:idx2]
			} else {
				backupType = prefix
			}
		}

		backups = append(backups, BackupMetadata{
			ID:        entry.Name(),
			Timestamp: info.ModTime(),
			Type:      backupType,
			Size:      info.Size(),
		})
	}

	return backups, nil
}

// Exists checks if a file exists
func (s *Store) Exists(filename string) bool {
	// Clean filename to prevent path traversal
	cleanFilename := filepath.Clean(filename)
	path := filepath.Join(s.dataDir, cleanFilename)
	_, err := os.Stat(path)
	return err == nil
}

// GetBackupPath returns the full path to a backup file
func (s *Store) GetBackupPath(backupID string) string {
	return filepath.Join(s.backupDir, backupID)
}

// LoadBackupJSON loads a backup file into the target structure
func (s *Store) LoadBackupJSON(backupID string, target interface{}) error {
	path := s.GetBackupPath(backupID)

	file, err := os.Open(path) // #nosec G304 - path is constructed from controlled backupDir
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: failed to close backup file %s: %v\n", path, closeErr)
		}
	}()

	if err := json.NewDecoder(file).Decode(target); err != nil {
		return fmt.Errorf("failed to decode backup JSON: %w", err)
	}

	return nil
}

// GetDataDir returns the data directory path
func (s *Store) GetDataDir() string {
	return s.dataDir
}

// GetBackupDir returns the backup directory path
func (s *Store) GetBackupDir() string {
	return s.backupDir
}
