// Package storage handles local data persistence
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	path := filepath.Join(s.dataDir, filename)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// LoadJSON loads JSON data from the specified path
func (s *Store) LoadJSON(filename string, target interface{}) error {
	path := filepath.Join(s.dataDir, filename)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

// CreateBackup creates a timestamped backup
func (s *Store) CreateBackup(backupType string, data interface{}) (*BackupMetadata, error) {
	timestamp := time.Now()
	id := timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", backupType, id)
	path := filepath.Join(s.backupDir, filename)

	// Ensure backup directory exists
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode backup: %w", err)
	}

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	return &BackupMetadata{
		ID:        id,
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

	var backups []BackupMetadata
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

		backups = append(backups, BackupMetadata{
			ID:        entry.Name(),
			Timestamp: info.ModTime(),
			Size:      info.Size(),
		})
	}

	return backups, nil
}

// Exists checks if a file exists
func (s *Store) Exists(filename string) bool {
	path := filepath.Join(s.dataDir, filename)
	_, err := os.Stat(path)
	return err == nil
}
