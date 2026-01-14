// rag_example.go demonstrates how to use the RAG (Retrieval-Augmented Generation) functionality of Spotigo
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/rag"
)

func main() {
	fmt.Println("=== Spotigo RAG Example ===")
	
	// Create Ollama client (requires Ollama running locally)
	ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)
	
	// Check if Ollama is available
	ctx := context.Background()
	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Println("Warning: Ollama not available - RAG functionality will be limited")
		fmt.Println("To use RAG features, start Ollama: ollama serve")
		fmt.Println("Then pull required models: ollama pull nomic-embed-text")
	} else {
		fmt.Println("Ollama connection successful")
	}

	// Create RAG store
	store := rag.NewStore(ollamaClient, "nomic-embed-text", "./data/embeddings.json")
	fmt.Println("RAG store created successfully")

	// Create sample documents
	documents := []rag.Document{
		{
			ID:      "track_1",
			Type:    "track",
			Content: "Bohemian Rhapsody by Queen",
			Metadata: map[string]string{
				"artist": "Queen",
				"album":  "A Night at the Opera",
				"genre":  "Rock",
				"year":   "1975",
			},
		},
		{
			ID:      "track_2",
			Type:    "track",
			Content: "Stairway to Heaven by Led Zeppelin",
			Metadata: map[string]string{
				"artist": "Led Zeppelin",
				"album":  "Led Zeppelin IV",
				"genre":  "Rock",
				"year":   "1971",
			},
		},
		{
			ID:      "artist_1",
			Type:    "artist",
			Content: "Queen is a British rock band formed in London in 1970",
			Metadata: map[string]string{
				"name":       "Queen",
				"formed":     "1970",
				"origin":     "London, England",
				"genres":     "Rock, Glam rock, Hard rock",
				"members":    "Freddie Mercury, Brian May, Roger Taylor, John Deacon",
			},
		},
	}

	fmt.Printf("\nCreated %d sample documents\n", len(documents))

	// Add documents to store (this would generate embeddings)
	fmt.Println("\n=== Adding Documents to Store ===")
	for _, doc := range documents {
		if err := store.Add(ctx, doc); err != nil {
			log.Printf("Failed to add document %s: %v", doc.ID, err)
		} else {
			fmt.Printf("Added document: %s (%s)\n", doc.ID, doc.Type)
		}
	}

	// Demonstrate search functionality
	fmt.Println("\n=== Semantic Search Example ===")
	query := "classic rock songs"
	fmt.Printf("Searching for: %q\n", query)

	results, err := store.Search(ctx, query, 5, "")
	if err != nil {
		log.Printf("Search failed: %v", err)
	} else {
		fmt.Printf("Found %d results:\n", len(results))
		for i, result := range results {
			fmt.Printf("%d. %s (similarity: %.2f%%)\n", 
				i+1, result.Document.Content, result.Similarity*100)
			fmt.Printf("   Type: %s\n", result.Document.Type)
			if len(result.Document.Metadata) > 0 {
				fmt.Printf("   Metadata: %v\n", result.Document.Metadata)
			}
			fmt.Println()
		}
	}

	// Show store statistics
	fmt.Println("=== Store Statistics ===")
	counts := store.CountByType()
	total := store.Count()
	fmt.Printf("Total documents: %d\n", total)
	for typ, count := range counts {
		fmt.Printf("%s: %d\n", typ, count)
	}

	// Demonstrate persistence
	fmt.Println("\n=== Persistence ===")
	if err := store.Save(); err != nil {
		log.Printf("Failed to save store: %v", err)
	} else {
		fmt.Println("Store saved to disk successfully")
	}

	fmt.Println("\n=== RAG in Spotigo ===")
	fmt.Println("The RAG functionality enables:")
	fmt.Println("- Semantic search across your music library")
	fmt.Println("- AI-powered music recommendations")
	fmt.Println("- Context-aware chat about your music")
	fmt.Println("- Music analysis and insights")
	
	fmt.Println("\nUse the CLI for full RAG features:")
	fmt.Println("$ spotigo search \"upbeat songs for working out\"")
	fmt.Println("$ spotigo search index    # Build search index")
	fmt.Println("$ spotigo chat           # Chat with your music library")
}