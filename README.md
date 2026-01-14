# Spotigo 2.0

AI-powered local music intelligence platform.

## Features

- **Complete Spotify Backup** - Save your entire library locally
- **AI Chat** - Natural language conversations about your music
- **Semantic Search** - Find songs by mood, vibe, or description
- **Statistics** - Deep insights into your listening habits
- **100% Offline** - All AI runs locally via Ollama

## Quick Start

```bash
# 1. Clone and enter the project
cd spotigo

# 2. Copy and configure
cp spotigo.example.yaml spotigo.yaml
# Edit spotigo.yaml with your Spotify API credentials

# 3. Install dependencies
go mod tidy

# 4. Start Ollama and pull models
ollama serve
ollama pull granite4:1b
ollama pull nomic-embed-text-v2-moe

# 5. Run Spotigo
go run ./cmd/spotigo

# Or build and run
go build -o spotigo ./cmd/spotigo
./spotigo
```

## Commands

```bash
# Authenticate with Spotify
spotigo auth

# Backup your library
spotigo backup

# Start AI chat
spotigo chat

# Search your library
spotigo search "upbeat songs for working out"

# View statistics
spotigo stats
spotigo stats top artists
spotigo stats genres

# Check model status
spotigo models status

# Launch TUI mode
spotigo --tui
```

## Configuration

Create `spotigo.yaml` from `spotigo.example.yaml`:

```yaml
spotify:
  client_id: "your_client_id"
  client_secret: "your_client_secret"

ollama:
  host: "http://localhost:11434"
```

Or use environment variables:
```bash
export SPOTIFY_CLIENT_ID=your_client_id
export SPOTIFY_CLIENT_SECRET=your_client_secret
export OLLAMA_HOST=http://localhost:11434
```

## AI Models

Spotigo uses tiered local models via Ollama:

| Role | Model | Size | Use Case |
|------|-------|------|----------|
| Chat | granite4:1b | 3.3GB | Main conversations |
| Fast | granite4:350m | 708MB | Quick classifications |
| Reasoning | qwen3:1.7b | 1.4GB | Deep analysis |
| Embeddings | nomic-embed-text-v2-moe | 957MB | Semantic search |

## Development

```bash
# Run in dev container (includes Ollama)
# Open in VS Code with Dev Containers extension

# Run tests
go test ./...

# Build
go build -o spotigo ./cmd/spotigo
```

## License

MIT
