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

	"spotigo/internal/ollama"
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
		doc.Embedding = embedding
	}

	s.mu.Lock()
	s.documents[doc.ID] = doc
	s.mu.Unlock()

	return nil
}

// AddBatch adds multiple documents efficiently
func (s *Store) AddBatch(ctx context.Context, docs []Document) error {
	for i, doc := range docs {
		if len(doc.Embedding) == 0 && s.client != nil {
			embedding, err := s.client.Embed(ctx, s.model, doc.Content)
			if err != nil {
				return fmt.Errorf("failed to generate embedding for doc %s: %w", doc.ID, err)
			}
			docs[i].Embedding = embedding
		}
	}

	s.mu.Lock()
	for _, doc := range docs {
		s.documents[doc.ID] = doc
	}
	s.mu.Unlock()

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

	var results []SearchResult
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
	if err := os.MkdirAll(dir, 0755); err != nil {
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

	if err := os.WriteFile(s.storePath, data, 0644); err != nil {
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
