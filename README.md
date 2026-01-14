# Spotigo

[![CI](https://github.com/bkataru/spotigo/actions/workflows/ci.yml/badge.svg)](https://github.com/bkataru/spotigo/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkataru/spotigo)](https://goreportcard.com/report/github.com/bkataru/spotigo)
[![Go Reference](https://pkg.go.dev/badge/github.com/bkataru/spotigo.svg)](https://pkg.go.dev/github.com/bkataru/spotigo)

**A powerful command-line application and Go library for AI-powered Spotify library management with local RAG capabilities.**

Spotigo provides both a comprehensive CLI application for end users and a powerful Go library for developers. Interact with Spotify APIs, manage local music libraries, and implement RAG (Retrieval-Augmented Generation) functionality with Ollama for AI-powered music analysis.

## Quick Start

### As a CLI Application

1. **Install Spotigo CLI:**
   ```bash
   # Download pre-built binary from releases
   # OR build from source:
   git clone https://github.com/bkataru/spotigo.git
   cd spotigo
   go build -o spotigo ./cmd/spotigo
   ```

2. **Authenticate with Spotify:**
   ```bash
   spotigo auth
   ```

3. **Backup your Spotify library:**
   ```bash
   spotigo backup
   ```

4. **Chat with your music library:**
   ```bash
   spotigo chat
   ```

### As a Go Library

```bash
go get github.com/bkataru/spotigo@latest
```

```go
package main

import (
    "context"
    "log"

    "github.com/bkataru/spotigo/internal/spotify"
    "github.com/bkataru/spotigo/internal/ollama"
    "github.com/bkataru/spotigo/internal/rag"
)

func main() {
    ctx := context.Background()

    // Initialize Spotify client
    spotifyClient, err := spotify.NewClient(spotify.Config{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
        RedirectURI:  "http://127.0.0.1:8888/callback",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Ollama client
    ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)

    // Initialize RAG store
    store := rag.NewStore(ollamaClient, "nomic-embed-text-v2-moe", "./data/store.json")
}
```

## CLI Usage

Spotigo comes with a comprehensive command-line interface for managing your Spotify library and performing AI-powered analysis.

### Available Commands

```bash
# Backup and restore Spotify library
spotigo backup                    # Backup your library
spotigo backup list              # List available backups
spotigo backup restore <id>      # Restore from backup
spotigo backup status            # Show backup status

# AI-powered chat about your music
spotigo chat                     # Start interactive chat session
spotigo chat --tools=true        # Enable function calling (default)
spotigo chat --data-dir ./data   # Specify data directory

# Semantic search across music library
spotigo search "rock music"      # Search with natural language
spotigo search index             # Build/rebuild search index
spotigo search status            # Show search index status

# Statistics and insights
spotigo stats                    # Overall listening statistics
spotigo stats top                # Top tracks and artists
spotigo stats genres             # Genre distribution analysis
spotigo stats playlists          # Playlist analysis

# Authentication management
spotigo auth                     # Authenticate with Spotify
spotigo auth status              # Check authentication status
spotigo auth logout              # Remove credentials

# Ollama model management
spotigo models list              # List recommended models
spotigo models status            # Show installed models
spotigo models pull              # Pull recommended models

# Interactive TUI mode
spotigo --tui                    # Launch terminal UI interface
```

### Configuration

Create `~/.spotigo.yaml`:

```yaml
spotify:
  client_id: "your_spotify_client_id"
  client_secret: "your_spotify_client_secret"
  redirect_uri: "http://127.0.0.1:8888/callback"

ollama:
  host: "http://localhost:11434"
  timeout: 30

storage:
  data_dir: "./data"
  backup_dir: "./data/backups"
  embeddings_dir: "./data/embeddings"

backup:
  schedule: "daily"
  retain_days: 30
```

Use a custom configuration file:
```bash
spotigo --config /path/to/config.yaml backup
```

## Features

- **Spotify API Integration** - Comprehensive Go client for Spotify Web API with OAuth2 authentication
- **CLI Application** - Full-featured command-line interface for end users
- **AI Chat with Function Calling** - Natural language queries with structured JSON tool execution
- **RAG Vector Store** - In-memory vector store with semantic search capabilities
- **Ollama Integration** - Client for local LLM inference with chat and embedding generation
- **JSON Query Engine** - Powerful structured queries for music data (filtering, sorting, aggregation)
- **Local Storage** - Encrypted token storage and persistent data management
- **Semantic Search** - Vector-based similarity search across music metadata
- **Batch Processing** - Parallel embedding generation and efficient bulk operations
- **Type-Safe APIs** - Well-documented Go interfaces with comprehensive error handling
- **Terminal UI** - Interactive menu-driven interface for all features

## Installation

### Prerequisites

- Go 1.24+ (for building from source)
- Ollama (optional, for AI functionality)
- Spotify Developer Account (for API access)

### Install as CLI Application

**Option 1: Download pre-built binary**
- Visit the [releases page](https://github.com/bkataru/spotigo/releases)
- Download the binary for your platform
- Add to your PATH

**Option 2: Build from source**
```bash
git clone https://github.com/bkataru/spotigo.git
cd spotigo
go build -o spotigo ./cmd/spotigo
sudo mv spotigo /usr/local/bin/  # or add to your PATH
```

### Install as Go Library

```bash
go get github.com/bkataru/spotigo@latest
```

## API Reference

### Spotify Client

```go
// Create Spotify client
client := spotify.NewClient(spotify.Config{
    ClientID:     "your_client_id",
    ClientSecret: "your_client_secret",
    RedirectURI:  "http://127.0.0.1:8888/callback",
})

// Get authentication URL
authURL := client.GetAuthURL("state")

// Handle OAuth callback
err := client.HandleCallback(ctx, "state", request)

// Fetch user data
tracks, err := client.GetSavedTracks(ctx)
playlists, err := client.GetPlaylists(ctx)
artists, err := client.GetFollowedArtists(ctx)
```

### Ollama Client

```go
// Create Ollama client
client := ollama.NewClient("http://localhost:11434", 30*time.Second)

// Generate embeddings
embedding, err := client.Embed(ctx, "nomic-embed-text-v2-moe", "text to embed")

// Chat completion
response, err := client.Chat(ctx, ollama.ChatRequest{
    Model: "granite4:1b",
    Messages: []ollama.Message{
        {Role: "user", Content: "Hello!"},
    },
})
```

### RAG Store

```go
// Create vector store
store := rag.NewStore(ollamaClient, "nomic-embed-text-v2-moe", "./data/store.json")

// Add documents
err := store.Add(ctx, rag.Document{
    ID:      "track_123",
    Type:    "track",
    Content: "Artist - Song Name",
    Metadata: map[string]string{"genre": "rock"},
})

// Semantic search
results, err := store.Search(ctx, "upbeat rock music", 10, "track")
```

### AI Chat with Tool Calling

The AI chat uses function calling to efficiently query your music library without loading all data into context.

```bash
# Start chat with tool calling enabled
spotigo chat

# Example conversation:
You: How many tracks do I have?
🔧 Calling tool: get_library_stats
Spotigo: You have 1,234 saved tracks, 25 playlists, and 42 followed artists.

You: Find my most popular Queen songs
🔧 Calling tool: get_tracks_by_artist
🔧 Calling tool: query_music_data
Spotigo: Here are your top Queen tracks:
1. "Bohemian Rhapsody" (popularity: 95)
2. "We Will Rock You" (popularity: 90)
...
```

**Available Tools:**
- `get_library_stats` - Overall library statistics
- `search_tracks` - Search by artist, song, album
- `get_tracks_by_artist` - All tracks by specific artist
- `get_recently_added_tracks` - Recently saved tracks
- `get_all_artists` - List all unique artists
- `get_playlist_by_name` - Find playlists
- `query_music_data` - Custom queries with filters, sorting, aggregation

**Why Function Calling?**
- ✅ **Efficient** - Only retrieves relevant data, minimal context usage
- ✅ **Accurate** - Structured queries are more precise than text embeddings
- ✅ **Fast** - Direct JSON queries with caching
- ✅ **Structured** - Preserves data relationships and schema

For detailed documentation on AI chat and tool calling, see [docs/TOOLS.md](docs/TOOLS.md).

### JSON Query Engine

```go
import "github.com/bkataru/spotigo/internal/jsonquery"

// Create query engine
engine := jsonquery.NewEngine("./data")

// Execute structured query
result := engine.Execute(jsonquery.Query{
    Source:    "saved_tracks.json",
    Operation: "filter",
    Filters: []jsonquery.Filter{
        {
            Field:    "track.popularity",
            Operator: "gte",
            Value:    90,
        },
    },
    SortBy:    "track.popularity",
    SortOrder: "desc",
    Limit:     10,
})

// Use music query helpers
helper := jsonquery.NewMusicQueryHelper("./data")
stats := helper.GetLibraryStats()
tracks := helper.GetTracksByArtist("Queen")
recent := helper.GetRecentlyAddedTracks(20)
```

## Supported Models

### Embedding Models
- `nomic-embed-text-v2-moe` - Recommended for embeddings
- `qwen3-embedding:0.6b` - Smaller alternative
- Any Ollama-compatible embedding model

### Chat Models
- `granite4:1b` - Balanced performance
- `qwen3:0.6b` - Lightweight option
- `granite4:350m` - Fast inference
- Any Ollama-compatible chat model

## Development

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/bkataru/spotigo.git
cd spotigo

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run

# Build CLI
go build -o spotigo ./cmd/spotigo
```

### Project Structure

```
spotigo/
├── cmd/spotigo/         # CLI application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── crypto/          # Token encryption utilities
│   ├── jsonutil/        # JSON utilities
│   ├── ollama/          # Ollama API client
│   ├── rag/             # RAG vector store
│   ├── spotify/         # Spotify API client
│   ├── storage/         # Local file storage
│   └── cmd/             # CLI command implementations
├── examples/            # Example usage (library and CLI)
├── .github/workflows/   # CI/CD pipelines
└── go.mod              # Module definition
```

## Architecture

### Component Relationships

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Spotify   │────>│   Spotify   │────>│   RAG       │
│     API     │     │   Client    │     │   Store     │
└─────────────┘     └─────────────┘     └─────────────┘
                                             │
                   ┌─────────────┐           │
                   │   Ollama    │<──────────┤
                   │   Client    │           │
                   └─────────────┘           │
                         │                   │
                   ┌─────────────┐     ┌─────────────┐
                   │   Local     │     │   CLI/API   │
                   │   Storage   │<───>│   Interface │
                   └─────────────┘     └─────────────┘
```

### Security Features

- **Token Encryption**: OAuth tokens encrypted using AES-256-GCM
- **Local Processing**: All AI inference happens locally via Ollama
- **No External Dependencies**: Minimal reliance on external services
- **Type Safety**: Comprehensive Go interfaces with proper error handling

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run`
6. Commit your changes
7. Open a Pull Request

## Troubleshooting

### Ollama Connection Issues

Ensure Ollama is running:
```bash
curl http://localhost:11434/api/tags
```

### Spotify Authentication Issues

Check your OAuth configuration and ensure your Spotify app has the correct redirect URI.

### Build Issues

```bash
# Clean and rebuild
go mod tidy
go test ./...
```

## FAQ

### General Questions

**Q: Does this application send data to external servers?**

A: No. All AI processing happens locally via Ollama. The application does not include telemetry or analytics.

**Q: What Spotify data can I access?**

A: The application supports read-only access to:
- Saved tracks and albums
- User playlists (including private ones)
- Followed artists
- Recently played tracks
- Top tracks and artists

The application cannot modify your Spotify library or play music.

**Q: Can I use this application offline?**

A: After initial Spotify authentication, most features work offline:
- Semantic search (with pre-generated embeddings)
- Local data processing
- RAG functionality (requires Ollama running locally)

Internet access is required for:
- Spotify OAuth authentication
- Spotify API calls

### Technical Questions

**Q: How do I handle OAuth authentication?**

A: The Spotify client provides methods for OAuth flow:
```go
// Get authentication URL
authURL := client.GetAuthURL("state")

// Handle callback
err := client.HandleCallback(ctx, "state", request)
```

**Q: How do I persist authentication tokens?**

A: The Spotify client includes encrypted token storage:
```go
err := client.SaveToken(".spotify_token")
```

**Q: What models are recommended for embeddings?**

A: `nomic-embed-text-v2-moe` is recommended for embeddings, but any Ollama-compatible embedding model will work.

**Q: How do I handle errors?**

A: All library functions return proper Go errors with descriptive messages for easy debugging.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Ollama](https://ollama.ai/) - Local LLM inference
- [Spotify Web API](https://developer.spotify.com/) - Music data
- [zmb3's Spotify library](https://github.com/zmb3/spotify) - Go client library