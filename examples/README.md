# Spotigo Examples

This directory contains example code demonstrating how to use the Spotigo Go library.

## Runnable Examples

### Basic Usage (`runnable/basic_usage.go`)
A complete runnable example demonstrating the core components of the Spotigo library:
- Spotify client initialization
- Ollama client setup  
- RAG store creation
- Sample document creation

```bash
cd examples/runnable
go run basic_usage.go
```

## Library Reference Examples (`spotigo_examples.go`)
Contains example functions for all major Spotigo components:
- Backup system usage
- RAG (Retrieval-Augmented Generation) functionality
- Chat interface with local LLMs

These examples are meant to be used as reference code rather than runnable programs.

## Library Components

### Spotify Client
The `spotify.Client` provides:
- OAuth2 authentication flow
- Token management with encryption
- Full Spotify Web API coverage
- Rate limiting and error handling

### Ollama Client  
The `ollama.Client` provides:
- Chat completion with local LLMs
- Text embedding generation
- Model management
- Streaming support (planned)

### RAG Store
The `rag.Store` provides:
- In-memory vector storage
- Semantic search capabilities
- Batch embedding generation
- Persistent storage to disk

## Integration Patterns

### Building a Music Assistant
```go
// 1. Authenticate with Spotify
client := spotify.NewClient(config)

// 2. Fetch user library  
tracks, _ := client.GetSavedTracks(ctx)

// 3. Create embeddings
store := rag.NewStore(ollamaClient, "model", "./store.json")
for _, track := range tracks {
    doc := rag.TrackToDocument(track)
    store.Add(ctx, doc)
}

// 4. Enable semantic search
results, _ := store.Search(ctx, "upbeat summer music", 10, "track")
```

### Creating a Music Analysis Tool
```go
// Analyze music preferences
stats := analyzeLibrary(tracks)

// Generate insights with AI
response, _ := ollamaClient.Chat(ctx, ollama.ChatRequest{
    Model: "granite4:3b",
    Messages: []ollama.Message{
        {Role: "system", Content: "You are a music analyst."},
        {Role: "user", Content: fmt.Sprintf("Analyze these stats: %v", stats)},
    },
})
```

## Requirements

For examples to run, you'll need:

1. **Go 1.24+**
2. **Ollama** (for AI features)
   ```bash
   ollama serve
   ollama pull nomic-embed-text
   ollama pull granite4:3b
   ```
3. **Spotify Developer Credentials**
   - Create app at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
   - Set redirect URI to `http://127.0.0.1:8888/callback`

## Next Steps

1. Check the main library documentation in the root README.md
2. Explore the internal package APIs
3. Run tests: `go test ./...`
4. Build the CLI tool: `go build ./cmd/spotigo`