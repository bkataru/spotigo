package rag

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
)

// BenchmarkStore_Add benchmarks adding a single document to the store
func BenchmarkStore_Add(b *testing.B) {
	store := NewStore(nil, "", "")

	doc := Document{
		ID:        "track-1",
		Type:      "track",
		Content:   "Test Track by Test Artist from Test Album",
		Metadata:  map[string]string{"name": "Test Track", "artist": "Test Artist"},
		Embedding: generateRandomEmbedding(768),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc.ID = fmt.Sprintf("track-%d", i)
		_ = store.Add(context.Background(), doc)
	}
}

// BenchmarkStore_AddBatch benchmarks adding documents in batches
func BenchmarkStore_AddBatch(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			docs := generateDocuments(size, 768)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store := NewStore(nil, "", "")
				_ = store.AddBatch(context.Background(), docs)
			}
		})
	}
}

// BenchmarkStore_AddBatchParallel benchmarks parallel batch addition with different concurrency levels
func BenchmarkStore_AddBatchParallel(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8}
	size := 100

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("concurrency=%d", concurrency), func(b *testing.B) {
			docs := generateDocuments(size, 768)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store := NewStore(nil, "", "")
				_ = store.AddBatchParallel(context.Background(), docs, concurrency)
			}
		})
	}
}

// BenchmarkStore_SearchSimilarity benchmarks the similarity calculation part of search
func BenchmarkStore_SearchSimilarity(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	queryEmbedding := generateRandomEmbedding(768)

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			// Setup: populate store
			store := NewStore(nil, "", "")
			docs := generateDocuments(size, 768)
			for _, doc := range docs {
				_ = store.Add(context.Background(), doc)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Simulate search similarity calculation
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

// BenchmarkStore_SearchByTypeSimilarity benchmarks filtered similarity by document type
func BenchmarkStore_SearchByTypeSimilarity(b *testing.B) {
	store := NewStore(nil, "", "")
	queryEmbedding := generateRandomEmbedding(768)

	// Populate with mixed types (50% tracks, 25% albums, 25% artists)
	for i := 0; i < 1000; i++ {
		docType := "track"
		if i%4 == 1 {
			docType = "album"
		} else if i%4 == 2 {
			docType = "artist"
		}

		doc := Document{
			ID:        fmt.Sprintf("doc-%d", i),
			Type:      docType,
			Content:   fmt.Sprintf("Content %d", i),
			Embedding: generateRandomEmbedding(768),
		}
		_ = store.Add(context.Background(), doc)
	}

	types := []string{"", "track", "album", "artist"}
	for _, docType := range types {
		name := "all"
		if docType != "" {
			name = docType
		}
		b.Run(fmt.Sprintf("type=%s", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Simulate filtered search
				store.mu.RLock()
				for _, doc := range store.documents {
					if docType != "" && doc.Type != docType {
						continue
					}
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
	sizes := []int{100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			tmpDir := b.TempDir()
			storePath := filepath.Join(tmpDir, "store.json")

			store := NewStore(nil, "", storePath)
			docs := generateDocuments(size, 768)
			for _, doc := range docs {
				_ = store.Add(context.Background(), doc)
			}

			b.Run("Save", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = store.Save()
				}
			})

			b.Run("Load", func(b *testing.B) {
				// Save once for load tests
				_ = store.Save()

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					newStore := NewStore(nil, "", storePath)
					_ = newStore.Load()
				}
			})
		})
	}
}

// BenchmarkCosineSimilarity_VectorSizes benchmarks the cosine similarity calculation with different vector sizes
func BenchmarkCosineSimilarity_VectorSizes(b *testing.B) {
	sizes := []int{128, 384, 768, 1024}

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

// BenchmarkCosineSimilarity_EdgeCases benchmarks edge cases
func BenchmarkCosineSimilarity_EdgeCases(b *testing.B) {
	b.Run("IdenticalVectors", func(b *testing.B) {
		vec := generateRandomEmbedding(768)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = cosineSimilarity(vec, vec)
		}
	})

	b.Run("ZeroVector", func(b *testing.B) {
		vec := generateRandomEmbedding(768)
		zero := make([]float64, 768)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = cosineSimilarity(vec, zero)
		}
	})

	b.Run("OppositeVectors", func(b *testing.B) {
		vec := generateRandomEmbedding(768)
		opposite := make([]float64, 768)
		for i, v := range vec {
			opposite[i] = -v
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = cosineSimilarity(vec, opposite)
		}
	})
}

// BenchmarkStore_Count benchmarks the count operation
func BenchmarkStore_Count(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			store := NewStore(nil, "", "")
			docs := generateDocuments(size, 768)
			for _, doc := range docs {
				_ = store.Add(context.Background(), doc)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = store.Count()
			}
		})
	}
}

// BenchmarkStore_CountByType benchmarks the count by type operation
func BenchmarkStore_CountByType(b *testing.B) {
	store := NewStore(nil, "", "")

	// Populate with mixed types
	for i := 0; i < 1000; i++ {
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
			Embedding: generateRandomEmbedding(768),
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
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			docs := generateDocuments(size, 768)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				store := NewStore(nil, "", "")
				for _, doc := range docs {
					_ = store.Add(context.Background(), doc)
				}
				b.StartTimer()

				store.Clear()
			}
		})
	}
}

// BenchmarkStore_ConcurrentReads benchmarks concurrent read operations
func BenchmarkStore_ConcurrentReads(b *testing.B) {
	store := NewStore(nil, "", "")
	docs := generateDocuments(1000, 768)
	for _, doc := range docs {
		_ = store.Add(context.Background(), doc)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = store.Count()
			_ = store.CountByType()
		}
	})
}

// Helper functions

// generateRandomEmbedding generates a random embedding vector of given size
func generateRandomEmbedding(size int) []float64 {
	embedding := make([]float64, size)
	for i := range embedding {
		embedding[i] = rand.Float64()*2 - 1 // Random values between -1 and 1
	}
	return embedding
}

// generateDocuments generates n test documents with embeddings of given dimension
func generateDocuments(n, embeddingDim int) []Document {
	docs := make([]Document, n)
	for i := 0; i < n; i++ {
		docs[i] = Document{
			ID:      fmt.Sprintf("track-%d", i),
			Type:    "track",
			Content: fmt.Sprintf("Test Track %d by Artist %d from Album %d", i, i%100, i%50),
			Metadata: map[string]string{
				"name":   fmt.Sprintf("Track %d", i),
				"artist": fmt.Sprintf("Artist %d", i%100),
				"album":  fmt.Sprintf("Album %d", i%50),
			},
			Embedding: generateRandomEmbedding(embeddingDim),
		}
	}
	return docs
}

// Benchmark memory allocations for similarity calculations
func BenchmarkStore_SearchSimilarity_Allocs(b *testing.B) {
	store := NewStore(nil, "", "")
	docs := generateDocuments(1000, 768)
	for _, doc := range docs {
		_ = store.Add(context.Background(), doc)
	}

	queryEmbedding := generateRandomEmbedding(768)

	b.ReportAllocs()
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
}
