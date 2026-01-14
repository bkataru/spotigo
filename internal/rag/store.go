// Package rag provides RAG (Retrieval-Augmented Generation) functionality
// with vector embeddings for semantic search across music library data.
package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/bkataru/spotigo/internal/ollama"

	"go.uber.org/multierr"
)

// Document represents a searchable item in the vector store
type Document struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // track, artist, album, playlist
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	Embedding []float64         `json:"embedding,omitempty"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Document   Document `json:"document"`
	Similarity float64  `json:"similarity"`
}

// Store is an in-memory vector store with persistence
type Store struct {
	mu        sync.RWMutex
	documents map[string]Document
	client    *ollama.Client
	model     string
	storePath string
}

// NewStore creates a new vector store
func NewStore(client *ollama.Client, model string, storePath string) *Store {
	return &Store{
		documents: make(map[string]Document),
		client:    client,
		model:     model,
		storePath: storePath,
	}
}

// Add adds a document to the store with automatic embedding generation
func (s *Store) Add(ctx context.Context, doc Document) error {
	// Generate embedding if not provided
	if len(doc.Embedding) == 0 && s.client != nil {
		embedding, err := s.client.Embed(ctx, s.model, doc.Content)
		if err != nil {
			return fmt.Errorf("failed to generate embedding: %w", err)
		}
		if len(embedding) == 0 {
			return fmt.Errorf("generated empty embedding for document %s", doc.ID)
		}
		doc.Embedding = embedding
	}

	s.mu.Lock()
	s.documents[doc.ID] = doc
	s.mu.Unlock()

	return nil
}

// AddBatch adds multiple documents efficiently with parallel embedding generation
func (s *Store) AddBatch(ctx context.Context, docs []Document) error {
	return s.AddBatchParallel(ctx, docs, 4) // Default concurrency of 4
}

// AddBatchParallel adds multiple documents with configurable parallelism
func (s *Store) AddBatchParallel(ctx context.Context, docs []Document, concurrency int) error {
	if concurrency < 1 {
		concurrency = 1
	}

	// Make a copy of the documents to avoid modifying the caller's slice
	docsCopy := make([]Document, len(docs))
	copy(docsCopy, docs)

	// Find documents that need embeddings
	var needsEmbedding []int
	for i, doc := range docsCopy {
		if len(doc.Embedding) == 0 && s.client != nil {
			needsEmbedding = append(needsEmbedding, i)
		}
	}

	if len(needsEmbedding) == 0 {
		// No embeddings needed, just add documents
		s.mu.Lock()
		for _, doc := range docsCopy {
			s.documents[doc.ID] = doc
		}
		s.mu.Unlock()
		return nil
	}

	// Create worker pool for parallel embedding generation
	type embeddingResult struct {
		index     int
		embedding []float64
		err       error
	}

	jobs := make(chan int, len(needsEmbedding))
	results := make(chan embeddingResult, len(needsEmbedding))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-ctx.Done():
					results <- embeddingResult{index: idx, err: ctx.Err()}
					return
				default:
					// Make a local copy of the content to avoid race conditions
					content := docsCopy[idx].Content
					embedding, err := s.client.Embed(ctx, s.model, content)
					results <- embeddingResult{index: idx, embedding: embedding, err: err}
				}
			}
		}()
	}

	// Send jobs with context cancellation check
	for _, idx := range needsEmbedding {
		select {
		case <-ctx.Done():
			// Context canceled, stop sending jobs
			close(jobs)
			wg.Wait()
			close(results)
			return ctx.Err()
		case jobs <- idx:
			// Job sent successfully
		}
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results - store embeddings in a separate map to avoid races
	embeddingsMap := make(map[int][]float64)
	var errors []error
	for result := range results {
		if result.err != nil {
			errors = append(errors, fmt.Errorf("doc %s: %w", docsCopy[result.index].ID, result.err))
			continue
		}
		if len(result.embedding) == 0 {
			errors = append(errors, fmt.Errorf("doc %s: generated empty embedding", docsCopy[result.index].ID))
			continue
		}
		embeddingsMap[result.index] = result.embedding
	}

	// Apply embeddings to copies (single-threaded, no race)
	for idx, embedding := range embeddingsMap {
		docsCopy[idx].Embedding = embedding
	}

	// Add all documents to store (even if some embeddings failed)
	s.mu.Lock()
	for _, doc := range docsCopy {
		s.documents[doc.ID] = doc
	}
	s.mu.Unlock()

	// Return combined errors if any
	if len(errors) > 0 {
		return fmt.Errorf("embedding errors (%d/%d failed): %w", len(errors), len(needsEmbedding), multierr.Combine(errors...))
	}

	return nil
}

// Search performs semantic search and returns the most similar documents
func (s *Store) Search(ctx context.Context, query string, limit int, docType string) ([]SearchResult, error) {
	// Generate embedding for query
	queryEmbedding, err := s.client.Embed(ctx, s.model, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]SearchResult, 0, len(s.documents))
	for _, doc := range s.documents {
		// Filter by type if specified
		if docType != "" && docType != "all" && doc.Type != docType {
			continue
		}

		// Skip documents without embeddings
		if len(doc.Embedding) == 0 {
			continue
		}

		similarity := cosineSimilarity(queryEmbedding, doc.Embedding)
		results = append(results, SearchResult{
			Document:   doc,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit results
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// Count returns the number of documents in the store
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.documents)
}

// CountByType returns the count of documents by type
func (s *Store) CountByType() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[string]int)
	for _, doc := range s.documents {
		counts[doc.Type]++
	}
	return counts
}

// Save persists the store to disk
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.storePath == "" {
		return fmt.Errorf("no store path configured")
	}

	// Ensure directory exists
	dir := filepath.Dir(s.storePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Convert to slice for serialization
	docs := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}

	data, err := json.MarshalIndent(docs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal documents: %w", err)
	}

	if err := os.WriteFile(s.storePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write store: %w", err)
	}

	return nil
}

// Load reads the store from disk
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storePath == "" {
		return fmt.Errorf("no store path configured")
	}

	data, err := os.ReadFile(s.storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Empty store is OK
		}
		return fmt.Errorf("failed to read store: %w", err)
	}

	var docs []Document
	if err := json.Unmarshal(data, &docs); err != nil {
		return fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	s.documents = make(map[string]Document, len(docs))
	for _, doc := range docs {
		s.documents[doc.ID] = doc
	}

	return nil
}

// Clear removes all documents from the store
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents = make(map[string]Document)
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
