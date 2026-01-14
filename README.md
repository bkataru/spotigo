# Spotigo

[![CI](https://github.com/bkataru/spotigo/actions/workflows/ci.yml/badge.svg)](https://github.com/bkataru/spotigo/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkataru/spotigo)](https://goreportcard.com/report/github.com/bkataru/spotigo)
[![Go Reference](https://pkg.go.dev/badge/github.com/bkataru/spotigo.svg)](https://pkg.go.dev/github.com/bkataru/spotigo)

**A Go library and CLI application for AI-powered Spotify library management with local RAG capabilities.**

Spotigo provides both a comprehensive Go library for developers and a powerful command-line interface for end users to interact with Spotify APIs, manage local music libraries, and implement RAG (Retrieval-Augmented Generation) functionality with Ollama for AI-powered music analysis.

## Features

- **Spotify API Integration** - Comprehensive Go client for Spotify Web API with OAuth2 authentication
- **RAG Vector Store** - In-memory vector store with semantic search capabilities
- **Ollama Integration** - Client for local LLM inference with chat and embedding generation
- **Local Storage** - Encrypted token storage and persistent data management
- **Semantic Search** - Vector-based similarity search across music metadata
- **Batch Processing** - Parallel embedding generation and efficient bulk operations
- **Type-Safe APIs** - Well-documented Go interfaces with comprehensive error handling



## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Configuration](#configuration)
- [Supported Models](#supported-models)
- [Development](#development)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [License](#license)

## Installation

### Prerequisites

- Go 1.24+
- Ollama (optional, for AI functionality)
- Spotify Developer Account (for API access)

### Go Module

```bash
go get github.com/bkataru/spotigo@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/bkataru/spotigo.git
cd spotigo

# Use as a library dependency in your go.mod
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
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
        RedirectURI:  "http://localhost:8888/callback",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Ollama client
    ollamaClient := ollama.NewClient("http://localhost:11434", 30*time.Second)

    // Initialize RAG store
    store := rag.NewStore(ollamaClient, "nomic-embed-text-v2-moe", "./data/store.json")

    // Use the components...
    _ = spotifyClient
    _ = store
}
```

## API Reference

### Spotify Client

```go
// Create Spotify client
client := spotify.NewClient(spotify.Config{
    ClientID:     "your_client_id",
    ClientSecret: "your_client_secret",
    RedirectURI:  "http://localhost:8888/callback",
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

## Configuration

### Spotify Client Configuration

```go
type Config struct {
    ClientID     string
    ClientSecret string
    RedirectURI  string
    TokenFile    string // Optional: encrypted token storage
}
```

### Ollama Client Configuration

```go
// Create with custom timeout
client := ollama.NewClient("http://localhost:11434", 60*time.Second)
```

### RAG Store Configuration

```go
// Initialize with custom embedding model and storage path
store := rag.NewStore(ollamaClient, "qwen3-embedding:0.6b", "/path/to/store.json")
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

### Prerequisites

- Go 1.24+
- Ollama (for testing AI functionality)

### Setup

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
```

### Project Structure

```
spotigo/
├── internal/
│   ├── config/          # Configuration management
│   ├── crypto/          # Token encryption utilities
│   ├── jsonutil/        # JSON utilities
│   ├── ollama/          # Ollama API client
│   ├── rag/             # RAG vector store
│   ├── spotify/         # Spotify API client
│   └── storage/         # Local file storage
├── cmd/spotigo/         # CLI application (example usage)
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
                    │   Local     │     │   Your      │
                    │   Storage   │<───>│   Application │
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

**Q: Does this library send data to external servers?**

A: No. All AI processing happens locally via Ollama. The library does not include telemetry or analytics.

**Q: What Spotify data can I access with this library?**

A: The Spotify client supports read-only access to:
- Saved tracks and albums
- User playlists (including private ones)
- Followed artists
- Recently played tracks
- Top tracks and artists

The library cannot modify your Spotify library or play music.

**Q: Can I use this library offline?**

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
- [zmb3/spotify](https://github.com/zmb3/spotify) - Go Spotify client

## CLI Usage

Spotigo can be used as a standalone CLI application for backing up your Spotify library and performing AI-powered analysis.

### Installation

To install the CLI, you can either:

1. Download a pre-built binary from the [releases page](https://github.com/bkataru/spotigo/releases)
2. Build from source:
   ```bash
   git clone https://github.com/bkataru/spotigo.git
   cd spotigo
   go build -o spotigo ./cmd/spotigo
   ```

### Getting Started

1. First, authenticate with Spotify:
   ```bash
   spotigo auth
   ```

2. Backup your Spotify library:
   ```bash
   spotigo backup
   ```

3. Build search index for semantic search (requires Ollama):
   ```bash
   spotigo search index
   ```

4. Start chatting with your music library:
   ```bash
   spotigo chat
   ```

### Available Commands

- `spotigo backup` - Backup your Spotify library
  - `spotigo backup list` - List available backups
  - `spotigo backup restore [backup-id]` - Restore from a backup
  - `spotigo backup status` - Show backup status and schedule

- `spotigo chat` - Start an AI chat session about your music
  - Uses local LLMs via Ollama for privacy-first conversations

- `spotigo search [query]` - Semantic search across your music library
  - `spotigo search index` - Build or rebuild the search index
  - `spotigo search status` - Show search index status

- `spotigo stats` - View listening statistics and insights
  - `spotigo stats top` - Show top tracks and artists
  - `spotigo stats genres` - Analyze genre distribution
  - `spotigo stats playlists` - Playlist analysis

- `spotigo auth` - Manage Spotify authentication
  - `spotigo auth status` - Check authentication status
  - `spotigo auth logout` - Remove stored credentials

- `spotigo models` - Manage Ollama models
  - `spotigo models list` - List recommended models
  - `spotigo models status` - Show installed models
  - `spotigo models pull` - Pull recommended models

### Configuration

The CLI uses a configuration file (default: `$HOME/.spotigo.yaml`) to manage settings:

```yaml
spotify:
  client_id: "your_spotify_client_id"
  client_secret: "your_spotify_client_secret"
  redirect_uri: "http://localhost:8888/callback"

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

You can specify a custom config file with the `--config` flag:
```bash
spotigo --config /path/to/your/config.yaml backup
```

### TUI Mode

Launch Spotigo in Terminal User Interface mode:
```bash
spotigo --tui
```

This provides an interactive menu-driven interface for all features.