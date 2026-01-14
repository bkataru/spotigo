package rag

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore(nil, "test-model", "/tmp/test-store.json")

	if store == nil {
		t.Fatal("NewStore returned nil")
	}

	if store.model != "test-model" {
		t.Errorf("expected model 'test-model', got '%s'", store.model)
	}

	if store.storePath != "/tmp/test-store.json" {
		t.Errorf("expected storePath '/tmp/test-store.json', got '%s'", store.storePath)
	}

	if store.documents == nil {
		t.Error("documents map should be initialized")
	}
}

func TestStore_AddWithEmbedding(t *testing.T) {
	store := NewStore(nil, "", "")

	doc := Document{
		ID:        "test-1",
		Type:      "track",
		Content:   "Test Track by Test Artist",
		Metadata:  map[string]string{"name": "Test Track"},
		Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
	}

	err := store.Add(context.Background(), doc)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if store.Count() != 1 {
		t.Errorf("expected count 1, got %d", store.Count())
	}
}

func TestStore_AddBatch(t *testing.T) {
	store := NewStore(nil, "", "")

	docs := []Document{
		{
			ID:        "track-1",
			Type:      "track",
			Content:   "Track 1",
			Embedding: []float64{0.1, 0.2, 0.3},
		},
		{
			ID:        "track-2",
			Type:      "track",
			Content:   "Track 2",
			Embedding: []float64{0.4, 0.5, 0.6},
		},
		{
			ID:        "artist-1",
			Type:      "artist",
			Content:   "Artist 1",
			Embedding: []float64{0.7, 0.8, 0.9},
		},
	}

	err := store.AddBatch(context.Background(), docs)
	if err != nil {
		t.Fatalf("AddBatch failed: %v", err)
	}

	if store.Count() != 3 {
		t.Errorf("expected count 3, got %d", store.Count())
	}
}

func TestStore_CountByType(t *testing.T) {
	store := NewStore(nil, "", "")

	docs := []Document{
		{ID: "t1", Type: "track", Content: "Track 1", Embedding: []float64{0.1}},
		{ID: "t2", Type: "track", Content: "Track 2", Embedding: []float64{0.2}},
		{ID: "a1", Type: "artist", Content: "Artist 1", Embedding: []float64{0.3}},
		{ID: "p1", Type: "playlist", Content: "Playlist 1", Embedding: []float64{0.4}},
	}

	for _, doc := range docs {
		store.Add(context.Background(), doc)
	}

	counts := store.CountByType()

	if counts["track"] != 2 {
		t.Errorf("expected 2 tracks, got %d", counts["track"])
	}
	if counts["artist"] != 1 {
		t.Errorf("expected 1 artist, got %d", counts["artist"])
	}
	if counts["playlist"] != 1 {
		t.Errorf("expected 1 playlist, got %d", counts["playlist"])
	}
}

func TestStore_Clear(t *testing.T) {
	store := NewStore(nil, "", "")

	docs := []Document{
		{ID: "t1", Type: "track", Content: "Track 1", Embedding: []float64{0.1}},
		{ID: "t2", Type: "track", Content: "Track 2", Embedding: []float64{0.2}},
	}

	for _, doc := range docs {
		store.Add(context.Background(), doc)
	}

	if store.Count() != 2 {
		t.Fatalf("expected count 2 before clear, got %d", store.Count())
	}

	store.Clear()

	if store.Count() != 0 {
		t.Errorf("expected count 0 after clear, got %d", store.Count())
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "spotigo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "vectors.json")

	// Create store and add documents
	store := NewStore(nil, "test-model", storePath)

	docs := []Document{
		{
			ID:        "track-1",
			Type:      "track",
			Content:   "Test Track by Test Artist",
			Metadata:  map[string]string{"name": "Test Track", "artist": "Test Artist"},
			Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		},
		{
			ID:        "artist-1",
			Type:      "artist",
			Content:   "Test Artist. Genres: rock, indie",
			Metadata:  map[string]string{"name": "Test Artist", "genres": "rock, indie"},
			Embedding: []float64{0.5, 0.4, 0.3, 0.2, 0.1},
		},
	}

	for _, doc := range docs {
		store.Add(context.Background(), doc)
	}

	// Save
	err = store.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		t.Fatal("store file was not created")
	}

	// Create new store and load
	store2 := NewStore(nil, "test-model", storePath)
	err = store2.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if store2.Count() != 2 {
		t.Errorf("expected count 2 after load, got %d", store2.Count())
	}

	counts := store2.CountByType()
	if counts["track"] != 1 {
		t.Errorf("expected 1 track after load, got %d", counts["track"])
	}
	if counts["artist"] != 1 {
		t.Errorf("expected 1 artist after load, got %d", counts["artist"])
	}
}

func TestStore_LoadNonexistent(t *testing.T) {
	store := NewStore(nil, "", "/nonexistent/path/store.json")

	// Loading from non-existent path should not error (empty store is OK)
	err := store.Load()
	if err != nil {
		t.Errorf("Load from non-existent path should not error, got: %v", err)
	}

	if store.Count() != 0 {
		t.Errorf("expected count 0, got %d", store.Count())
	}
}

func TestStore_SaveNoPath(t *testing.T) {
	store := NewStore(nil, "", "")

	err := store.Save()
	if err == nil {
		t.Error("Save should fail when no store path is configured")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
		delta    float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{1, 0, 0},
			expected: 1.0,
			delta:    0.0001,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{0, 1, 0},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{-1, 0, 0},
			expected: -1.0,
			delta:    0.0001,
		},
		{
			name:     "similar vectors",
			a:        []float64{1, 2, 3},
			b:        []float64{1, 2, 3.1},
			expected: 0.9998,
			delta:    0.001,
		},
		{
			name:     "empty vectors",
			a:        []float64{},
			b:        []float64{},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "different length vectors",
			a:        []float64{1, 2, 3},
			b:        []float64{1, 2},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "zero vector",
			a:        []float64{0, 0, 0},
			b:        []float64{1, 2, 3},
			expected: 0.0,
			delta:    0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if result < tt.expected-tt.delta || result > tt.expected+tt.delta {
				t.Errorf("expected %.4f (±%.4f), got %.4f", tt.expected, tt.delta, result)
			}
		})
	}
}

func TestStore_SearchWithEmbeddings(t *testing.T) {
	// This test uses pre-computed embeddings to avoid needing Ollama
	store := NewStore(nil, "", "")

	// Add documents with known embeddings
	docs := []Document{
		{
			ID:        "rock-track",
			Type:      "track",
			Content:   "Rock Song by Rock Band",
			Embedding: []float64{0.9, 0.1, 0.0}, // Rock-like embedding
		},
		{
			ID:        "jazz-track",
			Type:      "track",
			Content:   "Jazz Song by Jazz Artist",
			Embedding: []float64{0.0, 0.9, 0.1}, // Jazz-like embedding
		},
		{
			ID:        "pop-track",
			Type:      "track",
			Content:   "Pop Song by Pop Star",
			Embedding: []float64{0.1, 0.0, 0.9}, // Pop-like embedding
		},
	}

	for _, doc := range docs {
		store.Add(context.Background(), doc)
	}

	// Note: Search requires a client for query embedding generation
	// This test verifies the documents are stored correctly
	if store.Count() != 3 {
		t.Errorf("expected 3 documents, got %d", store.Count())
	}
}

// BenchmarkCosineSimilarity benchmarks the cosine similarity calculation
func BenchmarkCosineSimilarity(b *testing.B) {
	// Create vectors typical of embedding models (768-1024 dimensions)
	size := 768
	a := make([]float64, size)
	vecB := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = float64(i) / float64(size)
		vecB[i] = float64(size-i) / float64(size)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cosineSimilarity(a, vecB)
	}
}
