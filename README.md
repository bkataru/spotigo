# Spotigo 2.0

[![CI](https://github.com/bkataru-workshop/spotigo/actions/workflows/ci.yml/badge.svg)](https://github.com/bkataru-workshop/spotigo/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkataru-workshop/spotigo)](https://goreportcard.com/report/github.com/bkataru-workshop/spotigo)
[![codecov](https://codecov.io/gh/bkataru-workshop/spotigo/branch/main/graph/badge.svg)](https://codecov.io/gh/bkataru-workshop/spotigo)
[![Release](https://img.shields.io/github/v/release/bkataru-workshop/spotigo?include_prereleases)](https://github.com/bkataru-workshop/spotigo/releases)

**AI-powered local music intelligence platform for Spotify library management and analysis.**

Spotigo backs up your complete Spotify library locally and provides AI-powered semantic search, statistics, and natural language conversations about your music - all running 100% offline via Ollama.

## Features

- **Complete Library Backup** - Save all your saved tracks, playlists, and followed artists locally
- **AI Chat** - Natural language conversations about your music library
- **Semantic Search** - Find songs by mood, vibe, or description using vector embeddings
- **Deep Statistics** - Insights into your listening habits, top artists, genres, and more
- **Secure Token Storage** - OAuth tokens are encrypted using AES-256-GCM
- **100% Offline AI** - All AI processing runs locally via Ollama, your data never leaves your machine
- **Beautiful TUI** - Interactive terminal user interface

## Screenshots & Demos

> 🎬 **Screencasts coming soon!**
> 
> Demos will be created using [powersession](https://github.com/Hanaasagi/powersession-rs).  
> Install with: `scoop install powersession-rs`

### Planned Demos

- [ ] **Initial Setup** - Authentication and first-time configuration
- [ ] **Library Backup** - Creating and restoring full library backups
- [ ] **AI Chat** - Interactive conversations about your music library
- [ ] **Semantic Search** - Finding songs by mood, vibe, or description
- [ ] **TUI Walkthrough** - Complete tour of the interactive interface

## Table of Contents

- [Screenshots & Demos](#screenshots--demos)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [AI Models](#ai-models)
- [Development](#development)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [License](#license)

## Installation

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Ollama](https://ollama.ai/) for local AI inference
- Spotify Developer Account ([Create one here](https://developer.spotify.com/dashboard))

### From Source

```bash
# Clone the repository
git clone https://github.com/bkataru-workshop/spotigo.git
cd spotigo

# Build
go build -o spotigo ./cmd/spotigo

# Or install to GOPATH
go install ./cmd/spotigo
```

### Using Make

```bash
make build       # Build for current platform
make build-all   # Build for Linux, Windows, macOS
make install     # Install to GOPATH
```

### Docker

```bash
docker build -t spotigo:latest .
docker run --rm -it -v $(pwd)/data:/app/data spotigo:latest --help
```

## Quick Start

### 1. Set Up Spotify API Credentials

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new application
3. Add `http://localhost:8888/callback` as a Redirect URI
4. Copy your Client ID and Client Secret

### 2. Configure Spotigo

```bash
# Copy the example configuration
cp spotigo.example.yaml spotigo.yaml

# Edit with your credentials
# Or use environment variables:
export SPOTIFY_CLIENT_ID=your_client_id
export SPOTIFY_CLIENT_SECRET=your_client_secret
```

### 3. Start Ollama and Pull Models

```bash
# Start Ollama service
ollama serve

# Pull required models (in another terminal)
ollama pull granite4:1b              # Chat model (~3.3GB)
ollama pull nomic-embed-text-v2-moe  # Embeddings (~957MB)

# Optional: additional models
ollama pull qwen3:0.6b               # Fallback chat model
ollama pull granite4:350m            # Fast model for quick tasks
```

### 4. Authenticate and Use

```bash
# Authenticate with Spotify (opens browser)
./spotigo auth

# Backup your library
./spotigo backup

# Start chatting about your music!
./spotigo chat
```

## Commands

### Authentication

```bash
spotigo auth              # Start OAuth flow (opens browser)
spotigo auth status       # Check authentication status
spotigo auth logout       # Clear stored tokens
```

### Backup & Restore

```bash
spotigo backup                    # Backup entire library
spotigo backup --type tracks      # Backup only saved tracks
spotigo backup --type playlists   # Backup only playlists
spotigo backup --type artists     # Backup only followed artists
spotigo backup --index            # Backup and build search index
spotigo backup list               # List available backups
spotigo backup restore            # Restore from latest backup
spotigo backup restore <id>       # Restore from specific backup
spotigo backup status             # Show backup status
```

### AI Chat

```bash
spotigo chat                      # Start interactive chat
spotigo chat --model qwen3:1.7b   # Use specific model
spotigo chat --context 8192       # Set context window size
```

**Example conversations:**
- "What are my most played genres?"
- "Find upbeat songs for working out"
- "When did I start listening to The Beatles?"
- "Recommend something based on my recent listening"

### Search

```bash
spotigo search "chill electronic music"   # Semantic search
spotigo search --type track "summer"      # Search only tracks
spotigo search --type artist "jazz"       # Search only artists
spotigo search index                      # Build/rebuild search index
```

### Statistics

```bash
spotigo stats                     # Overview statistics
spotigo stats top artists         # Top artists
spotigo stats top tracks          # Top tracks
spotigo stats genres              # Genre distribution
spotigo stats timeline            # Listening timeline
```

### Model Management

```bash
spotigo models list               # List configured models
spotigo models status             # Check Ollama connection & available models
spotigo models pull granite4:1b   # Download a model via Ollama API
```

### TUI Mode

```bash
spotigo --tui                     # Launch interactive terminal UI
spotigo                           # Also launches TUI by default
```

## Configuration

### Configuration File

Create `spotigo.yaml` in the project root:

```yaml
spotify:
  client_id: "your_client_id"
  client_secret: "your_client_secret"
  redirect_uri: "http://localhost:8888/callback"
  token_file: ".spotify_token"

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
  format: "json"

app:
  verbose: false
  theme: "dark"
```

### Environment Variables

All configuration can be overridden with environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SPOTIFY_CLIENT_ID` | Spotify API Client ID | - |
| `SPOTIFY_CLIENT_SECRET` | Spotify API Client Secret | - |
| `SPOTIFY_REDIRECT_URI` | OAuth callback URL | `http://localhost:8888/callback` |
| `OLLAMA_HOST` | Ollama API endpoint | `http://localhost:11434` |
| `SPOTIGO_DATA_DIR` | Data storage directory | `./data` |

## AI Models

Spotigo uses a tiered model architecture for optimal performance:

| Role | Primary Model | Fallback | Size | Use Case |
|------|---------------|----------|------|----------|
| Chat | `granite4:1b` | `qwen3:0.6b` | 3.3GB | Main conversations |
| Fast | `granite4:350m` | `qwen3:0.6b` | 708MB | Quick classifications |
| Reasoning | `qwen3:1.7b` | `granite4:1b` | 1.4GB | Deep analysis |
| Embeddings | `nomic-embed-text-v2-moe` | `qwen3-embedding:0.6b` | 957MB | Semantic search |

### Minimum Requirements

- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: ~6GB for all models
- **GPU**: Optional but significantly improves performance

### Custom Models

Edit `config/models.yaml` to customize model selection:

```yaml
models:
  chat:
    primary: llama3.2:3b
    fallback: qwen3:0.6b
    temperature: 0.7
```

## Development

### Prerequisites

- Go 1.23+
- Make (optional)
- Docker (optional)

### Setup

```bash
# Clone and enter
git clone https://github.com/bkataru-workshop/spotigo.git
cd spotigo

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o spotigo ./cmd/spotigo
```

### Using Dev Container

The project includes a full dev container configuration with Ollama:

1. Open in VS Code with the Dev Containers extension
2. Select "Reopen in Container"
3. Ollama is pre-configured at `http://ollama:11434`

### Make Targets

```bash
make help         # Show all targets
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make lint         # Run linter
make fmt          # Format code
make quality      # Run all quality checks
make docker-build # Build Docker image
make dev-run      # Run in development mode
make release      # Create release builds
```

### Project Structure

```
spotigo/
├── cmd/spotigo/         # Main entry point
├── internal/
│   ├── cmd/             # CLI commands (Cobra)
│   ├── config/          # Configuration loading
│   ├── crypto/          # Token encryption
│   ├── jsonutil/        # JSON utilities
│   ├── ollama/          # Ollama API client
│   ├── rag/             # Vector store & embeddings
│   ├── spotify/         # Spotify API client
│   ├── storage/         # Local file storage
│   └── tui/             # Terminal UI (Bubbletea)
├── config/              # Configuration files
├── data/                # Local data storage
│   ├── backups/         # Library backups
│   └── embeddings/      # Vector embeddings
├── .devcontainer/       # Dev container config
└── .github/workflows/   # CI/CD pipelines
```

## Architecture

### Data Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Spotify   │────>│   Backup    │────>│   Local     │
│     API     │     │   Engine    │     │   Storage   │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                    ┌─────────────┐           │
                    │   Vector    │<──────────┤
                    │   Store     │           │
                    └─────────────┘           │
                          │                   │
                    ┌─────────────┐     ┌─────────────┐
                    │   Ollama    │     │    CLI/     │
                    │   (LLMs)    │<───>│    TUI      │
                    └─────────────┘     └─────────────┘
```

### Security

- **Token Encryption**: OAuth tokens are encrypted using AES-256-GCM with machine-specific key derivation
- **Local Processing**: All AI inference happens locally via Ollama
- **No Telemetry**: Spotigo does not send any data to external servers
- **Minimal Permissions**: Only requests necessary Spotify scopes

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linter: `make lint`
6. Commit: `git commit -am 'Add my feature'`
7. Push: `git push origin feature/my-feature`
8. Open a Pull Request

## Troubleshooting

### Ollama Connection Issues

```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama
ollama serve

# Check model availability
spotigo models status
```

### Authentication Issues

```bash
# Check token status
spotigo auth status

# Re-authenticate
spotigo auth logout
spotigo auth
```

### Build Issues

```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## FAQ

### General Questions

**Q: Does Spotigo send my data to external servers?**

A: No. All AI processing happens locally via Ollama. Your music library data never leaves your machine. Spotigo does not include any telemetry or analytics.

**Q: What Spotify data does Spotigo access?**

A: Spotigo requests read-only access to:
- Your saved tracks and albums
- Your playlists (including private ones)
- Artists you follow
- Recently played tracks
- Your top tracks and artists

Spotigo cannot modify your Spotify library or play music.

**Q: Can I use Spotigo without an internet connection?**

A: After the initial authentication with Spotify and backing up your library, most features work offline:
- AI chat (requires Ollama running locally)
- Semantic search
- Statistics
- Browsing backed-up data

You need internet access to:
- Authenticate with Spotify
- Sync/update your library backup

**Q: How much disk space does Spotigo need?**

A: Space requirements depend on your library size:
- **Application**: ~20MB
- **AI Models**: 1-6GB (depending on which models you install)
- **Library Backup**: Varies (typically 1-50MB for most users)
- **Embeddings**: Roughly 1-5MB per 1000 tracks

### Model Questions

**Q: Which AI model should I use?**

A: It depends on your hardware:
- **8GB RAM**: Use `granite4:1b` or `qwen3:0.6b`
- **16GB RAM**: Can run larger models like `qwen3:1.7b`
- **GPU available**: Any model will run faster with GPU acceleration

**Q: Can I use other Ollama models?**

A: Yes! Edit `config/models.yaml` to use any Ollama-compatible model. Just ensure:
- Chat models support conversation format
- Embedding models output vector embeddings

**Q: Why is the first query slow?**

A: Ollama loads models into memory on first use. Subsequent queries are much faster. You can pre-load a model with:
```bash
ollama run granite4:1b "hello"
```

### Technical Questions

**Q: Where is my data stored?**

A: By default:
- **Config**: `./spotigo.yaml`
- **Backups**: `./data/backups/`
- **Embeddings**: `./data/embeddings/`
- **Token**: `./.spotify_token` (encrypted)

**Q: How do I reset everything?**

A: To start fresh:
```bash
rm -rf ./data/
rm .spotify_token
spotigo auth logout
```

**Q: Can I use Spotigo with multiple Spotify accounts?**

A: Currently, Spotigo supports one account at a time. To switch accounts:
1. Run `spotigo auth logout`
2. Run `spotigo auth` with a different account

**Q: Does Spotigo work with Spotify Free accounts?**

A: Yes, Spotigo works with both Free and Premium Spotify accounts. All features are available regardless of your subscription type.

### Troubleshooting Questions

**Q: I get "token expired" errors, what do I do?**

A: Re-authenticate:
```bash
spotigo auth logout
spotigo auth
```

**Q: Ollama is running but Spotigo can't connect?**

A: Check the Ollama host in your config:
```yaml
ollama:
  host: "http://localhost:11434"  # Default
```

For Docker or remote Ollama, update the host accordingly.

**Q: Search returns no results?**

A: Ensure you've built the search index:
```bash
spotigo search index
# Or backup with indexing:
spotigo backup --index
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Ollama](https://ollama.ai/) - Local LLM inference
- [Spotify Web API](https://developer.spotify.com/) - Music data
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI
- [zmb3/spotify](https://github.com/zmb3/spotify) - Go Spotify client
