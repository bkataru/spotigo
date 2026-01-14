package rag

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
)

// Lightweight benchmarks with smaller data sizes for faster CI runs

// BenchmarkStore_Add benchmarks adding a single document to the store
func BenchmarkStore_Add(b *testing.B) {
	store := NewStore(nil, "", "")

	doc := Document{
		ID:        "track-1",
		Type:      "track",
		Content:   "Test Track by Test Artist from Test Album",
		Metadata:  map[string]string{"name": "Test Track", "artist": "Test Artist"},
		Embedding: generateRandomEmbedding(128), // Smaller embedding for speed
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc.ID = fmt.Sprintf("track-%d", i)
		_ = store.Add(context.Background(), doc)
	}
}

// BenchmarkStore_AddBatch benchmarks adding documents in batches
func BenchmarkStore_AddBatch(b *testing.B) {
	sizes := []int{10, 50} // Reduced sizes

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			docs := generateDocuments(size, 128) // Smaller embeddings

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store := NewStore(nil, "", "")
				_ = store.AddBatch(context.Background(), docs)
			}
		})
	}
}

// BenchmarkStore_SearchSimilarity benchmarks similarity calculation
func BenchmarkStore_SearchSimilarity(b *testing.B) {
	sizes := []int{50, 100} // Reduced sizes
	queryEmbedding := generateRandomEmbedding(128)

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			store := NewStore(nil, "", "")
			docs := generateDocuments(size, 128)
			_ = store.AddBatch(context.Background(), docs)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store.mu.RLock()
				for _, doc := range store.documents {
					if len(doc.Embedding) > 0 {
						_ = cosineSimilarity(queryEmbedding, doc.Embedding)
					}
				}
				store.mu.RUnlock()
			}
		})
	}
}

// BenchmarkStore_SaveLoad benchmarks persistence operations
func BenchmarkStore_SaveLoad(b *testing.B) {
	tmpDir := b.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store := NewStore(nil, "", storePath)
	docs := generateDocuments(50, 128) // Small dataset
	_ = store.AddBatch(context.Background(), docs)

	b.Run("Save", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.Save()
		}
	})

	b.Run("Load", func(b *testing.B) {
		_ = store.Save()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			newStore := NewStore(nil, "", storePath)
			_ = newStore.Load()
		}
	})
}

// BenchmarkCosineSimilarity benchmarks the core similarity calculation
func BenchmarkCosineSimilarity(b *testing.B) {
	sizes := []int{128, 384}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("dim=%d", size), func(b *testing.B) {
			a := generateRandomEmbedding(size)
			bVec := generateRandomEmbedding(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = cosineSimilarity(a, bVec)
			}
		})
	}
}

// BenchmarkStore_Count benchmarks the count operation
func BenchmarkStore_Count(b *testing.B) {
	store := NewStore(nil, "", "")
	docs := generateDocuments(100, 128)
	_ = store.AddBatch(context.Background(), docs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.Count()
	}
}

// BenchmarkStore_CountByType benchmarks the count by type operation
func BenchmarkStore_CountByType(b *testing.B) {
	store := NewStore(nil, "", "")

	// Populate with mixed types
	for i := 0; i < 100; i++ {
		docType := "track"
		if i%4 == 1 {
			docType = "album"
		} else if i%4 == 2 {
			docType = "artist"
		} else if i%4 == 3 {
			docType = "playlist"
		}

		doc := Document{
			ID:        fmt.Sprintf("doc-%d", i),
			Type:      docType,
			Content:   fmt.Sprintf("Content %d", i),
			Embedding: generateRandomEmbedding(128),
		}
		_ = store.Add(context.Background(), doc)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.CountByType()
	}
}

// BenchmarkStore_Clear benchmarks clearing the store
func BenchmarkStore_Clear(b *testing.B) {
	docs := generateDocuments(50, 128)
	store := NewStore(nil, "", "")
	_ = store.AddBatch(context.Background(), docs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Clear()
		// Re-populate without timing
		b.StopTimer()
		_ = store.AddBatch(context.Background(), docs)
		b.StartTimer()
	}
}

// Helper functions

func generateRandomEmbedding(size int) []float64 {
	embedding := make([]float64, size)
	for i := range embedding {
		embedding[i] = rand.Float64()*2 - 1
	}
	return embedding
}

func generateDocuments(n, embeddingDim int) []Document {
	docs := make([]Document, n)
	for i := 0; i < n; i++ {
		docs[i] = Document{
			ID:      fmt.Sprintf("track-%d", i),
			Type:    "track",
			Content: fmt.Sprintf("Test Track %d", i),
			Metadata: map[string]string{
				"name": fmt.Sprintf("Track %d", i),
			},
			Embedding: generateRandomEmbedding(embeddingDim),
		}
	}
	return docs
}
