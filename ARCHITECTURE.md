# Spotigo Architecture

## Overview

Spotigo 2.0 is a CLI/TUI application that combines Spotify library management with AI-powered chat and semantic search capabilities. It uses local AI models through Ollama for privacy-focused, offline-capable features.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          User Interface Layer                        │
├─────────────────────────────────────────────────────────────────────┤
│  CLI (Cobra)                        TUI (Bubbletea)                  │
│  ├─ auth                            └─ Interactive Chat              │
│  ├─ backup/restore                     ├─ Message History            │
│  ├─ models                              ├─ Context Display           │
│  ├─ chat (one-shot)                     └─ RAG Integration           │
│  └─ search                                                            │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Application Core                             │
├──────────────────┬──────────────────┬──────────────────┬────────────┤
│  Spotify Client  │  Ollama Client   │   RAG Engine     │  Storage   │
│  (spotify/)      │  (ollama/)       │   (rag/)         │ (storage/) │
│                  │                  │                  │            │
│  ├─ OAuth Flow   │  ├─ Chat         │  ├─ Embeddings   │ ├─ Backup  │
│  ├─ API Wrapper  │  ├─ Embeddings   │  ├─ Vector Store │ ├─ Restore │
│  ├─ Token Mgmt   │  ├─ Model Mgmt   │  └─ Search       │ └─ Config  │
│  └─ Rate Limit   │  └─ Streaming    │                  │            │
└──────────────────┴──────────────────┴──────────────────┴────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Infrastructure Layer                          │
├──────────────────┬──────────────────┬──────────────────────────────┤
│  Crypto          │  HTTP Clients    │  File System                  │
│  (crypto/)       │                  │                               │
│                  │                  │                               │
│  ├─ AES-256-GCM  │  ├─ Spotify API  │  ├─ Config Dir (~/.spotigo)  │
│  ├─ Key Derive   │  └─ Ollama API   │  ├─ Token Storage (encrypted)│
│  └─ Token Encrypt│                  │  ├─ RAG Store (JSON)          │
│                  │                  │  └─ Backups (JSON)            │
└──────────────────┴──────────────────┴──────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          External Services                           │
├──────────────────────────────────┬──────────────────────────────────┤
│       Spotify Web API            │        Ollama (Local)            │
│       (api.spotify.com)          │        (localhost:11434)         │
│                                  │                                  │
│  ├─ Library Data                 │  ├─ Chat Models (llama3.2, etc) │
│  ├─ User Profile                 │  └─ Embedding Models (nomic)    │
│  └─ OAuth 2.0                    │                                  │
└──────────────────────────────────┴──────────────────────────────────┘
```

## Component Breakdown

### 1. User Interface Layer

#### CLI (Cobra)
- **Purpose**: Command-line interface for all operations
- **Commands**:
  - `auth login/logout/status` - OAuth authentication management
  - `backup create/restore/list` - Library backup/restore operations
  - `models list/pull/set/clear` - AI model management
  - `chat` - One-shot chat with AI assistant
  - `chat interactive` - TUI-based interactive chat
  - `search index/query` - Semantic search operations

#### TUI (Bubbletea)
- **Purpose**: Rich terminal interface for interactive chat
- **Features**:
  - Multi-message conversation flow
  - Streaming AI responses
  - Automatic RAG context integration
  - Message history management
  - Syntax highlighting for code/JSON responses

### 2. Spotify Integration (`internal/spotify`)

#### Client Wrapper
- **OAuth 2.0 Flow**: Authorization code flow with PKCE
- **Token Management**: 
  - Automatic token refresh
  - Encrypted storage (AES-256-GCM)
  - Machine-specific key derivation
- **API Coverage**:
  - Saved tracks, albums, playlists
  - Followed artists
  - Top tracks/artists
  - Recently played tracks

#### Rate Limiting
- Respects Spotify API limits
- Exponential backoff for 429 responses
- Batch operations for efficiency

### 3. Ollama Integration (`internal/ollama`)

#### Chat Client
- **Streaming Support**: Real-time response streaming
- **Context Management**: Conversation history tracking
- **Model Selection**: Dynamic model switching (small/medium/large)
- **Error Handling**: Graceful degradation on model unavailability

#### Embedding Client
- **Vector Generation**: Text-to-embedding conversion via Ollama
- **Model**: Uses `nomic-embed-text` by default
- **Batch Processing**: Parallel embedding generation with configurable concurrency
- **Progress Tracking**: Real-time progress updates for large batches

#### Model Management
- **Pull**: Download models from Ollama registry
- **List**: Display available local models
- **Validation**: Check model availability before use

### 4. RAG Engine (`internal/rag`)

#### Vector Store
- **Storage**: In-memory map with JSON persistence
- **Embeddings**: Generated via Ollama embedding models
- **Similarity**: Cosine similarity for vector comparison
- **Thread Safety**: RWMutex for concurrent read/write

#### Document Schema
```go
type Document struct {
    ID        string            // Unique identifier (track/artist/album ID)
    Type      string            // track, artist, album, playlist
    Content   string            // Searchable text content
    Metadata  map[string]string // Additional context (artist, album, etc)
    Embedding []float64         // Vector representation
}
```

#### Search Algorithm
1. Generate query embedding
2. Filter by document type (optional)
3. Calculate cosine similarity for all documents
4. Sort by similarity (descending)
5. Return top N results

#### Performance Optimizations
- **Parallel Embeddings**: Configurable worker pool (default: 4 workers)
- **Batch Operations**: `AddBatch` for bulk document insertion
- **Lazy Loading**: Documents added without embeddings (generated on-demand)
- **Read Locks**: Multiple concurrent searches allowed

### 5. Storage Layer (`internal/storage`)

#### File Structure
```
~/.spotigo/
├── config.json          # Application configuration
├── token.enc            # Encrypted Spotify OAuth token
├── embeddings.json      # RAG vector store
└── backups/
    ├── backup_full_20240115_120000.json
    ├── backup_tracks_20240115_120000.json
    └── backup_playlists_20240115_120000.json
```

#### Backup System
- **Types**: Full, tracks-only, playlists-only, artists-only
- **Format**: JSON with metadata and timestamps
- **Restore**: Idempotent restoration with duplicate detection
- **Search Indexing**: Optional embedding generation during restore

#### Configuration
- **YAML-based**: `config.yaml` for user preferences
- **Defaults**: Sensible defaults for all settings
- **Validation**: Schema validation on load

### 6. Cryptography (`internal/crypto`)

#### Token Encryption
- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 256 bits (32 bytes)
- **Nonce**: Random 12-byte nonce per encryption
- **Authentication**: Built-in AEAD (Authenticated Encryption with Associated Data)

#### Key Derivation
Machine-specific key derived from:
1. Username (USER/USERNAME env var)
2. Home directory path
3. OS + Architecture (GOOS + GOARCH)
4. Config directory path
5. Application-specific salt

**Hash Function**: SHA-256 to produce 32-byte key

#### Security Properties
- **No Password**: Uses machine-specific entropy (not portable across machines)
- **Deterministic**: Same key on same machine for token persistence
- **File Permissions**: 0600 (owner read/write only)
- **Directory Permissions**: 0700 (owner full access only)

## Data Flow

### 1. Authentication Flow
```
User → CLI auth login
  → Spotify OAuth URL generated
  → User authorizes in browser
  → Callback received at localhost:8080
  → Token exchanged
  → Token encrypted with machine-specific key
  → Token saved to ~/.spotigo/token.enc
  → Success confirmation
```

### 2. Backup Flow
```
User → CLI backup create
  → Spotify Client fetches library data
    → Saved tracks (paginated)
    → Saved albums (paginated)
    → Playlists (paginated)
    → Followed artists (paginated)
  → Data serialized to JSON
  → File saved to ~/.spotigo/backups/
  → Optional: RAG indexing
    → Documents created for each item
    → Embeddings generated in parallel
    → Vectors stored in embeddings.json
```

### 3. Chat Flow (Interactive)
```
User → chat interactive command
  → TUI initialized
  → User types message
  → RAG search triggered (if enabled)
    → Query embedding generated
    → Vector similarity search
    → Top 5 results retrieved
  → Context assembled:
    → System prompt
    → RAG results (if any)
    → Conversation history
    → User message
  → Ollama chat request (streaming)
    → Chunks received in real-time
    → TUI updated incrementally
  → Response complete
  → History saved
  → Await next user input
```

### 4. Search Flow
```
User → search query "upbeat summer songs"
  → Query embedding generated via Ollama
  → Vector store searched
    → Cosine similarity calculated for all docs
    → Results sorted by similarity
  → Top 10 results returned
  → Display:
    → Track name, artist, album
    → Similarity score (%)
    → Metadata (year, genre, etc)
```

## Model Tiering Strategy

### Purpose
Allows users to balance performance vs. quality based on hardware capabilities.

### Tiers

#### Small (Fast)
- **Chat**: `llama3.2:3b`
- **Use Case**: Quick responses, low-end hardware, testing
- **Speed**: ~50 tokens/sec on CPU
- **Quality**: Good for simple queries

#### Medium (Balanced) - Default
- **Chat**: `llama3.2:latest` (7B)
- **Use Case**: General purpose, most users
- **Speed**: ~20 tokens/sec on CPU
- **Quality**: High quality responses

#### Large (Quality)
- **Chat**: `llama3.1:70b` or `qwen2.5:72b`
- **Use Case**: Complex analysis, high-end hardware
- **Speed**: ~2-5 tokens/sec on CPU, faster with GPU
- **Quality**: Best reasoning and context understanding

### Embedding Model
- **Model**: `nomic-embed-text` (fixed)
- **Dimensions**: 768
- **Reason**: Optimized for semantic search, consistent across tiers

## Security Model

### Threat Model

#### In Scope
- **Token Theft**: Prevent plaintext OAuth token exposure
- **Local Attacks**: Protect against other users on same machine
- **Accidental Exposure**: Prevent tokens in logs/errors

#### Out of Scope
- **Remote Attacks**: Assumes attacker doesn't have machine access
- **Root/Admin Access**: Cannot protect against privileged users
- **Memory Dumps**: Tokens decrypted in-memory during use

### Security Measures

#### Token Storage
1. **Encryption**: AES-256-GCM for token file
2. **Key Derivation**: Machine-specific, no user password needed
3. **File Permissions**: Restrictive (0600) to prevent other users
4. **No Network**: Tokens never transmitted except to Spotify API

#### API Communication
- **HTTPS Only**: All Spotify/Ollama communication encrypted in transit
- **Token Refresh**: Short-lived access tokens, refresh tokens stored encrypted
- **Scope Limiting**: Minimal required scopes (read-only library access)

#### RAG Data
- **No Encryption**: Library data in embeddings.json is plaintext
- **Justification**: Public data, no sensitive info beyond music preferences
- **Privacy**: All processing local, no cloud uploads

#### Ollama Communication
- **Local Only**: Assumes Ollama on localhost (default)
- **No Auth**: Ollama typically runs without authentication on localhost
- **Data Exposure**: Music data sent to Ollama for embeddings/chat

### Privacy Considerations

#### Data Collection
- **None**: No telemetry, analytics, or cloud services
- **Spotify API**: Only data explicitly requested by user
- **Ollama**: All AI processing local

#### Data Retention
- **Backups**: User-controlled, stored locally
- **Embeddings**: Persisted until user clears
- **Tokens**: Valid until user logs out or token expires
- **Chat History**: Stored in-memory only (cleared on exit)

## Performance Characteristics

### Benchmarks (Approximate)

#### Embedding Generation
- **Single Document**: ~100-200ms (Ollama overhead)
- **Batch (100 docs, 4 workers)**: ~10-15 seconds
- **Batch (1000 docs, 4 workers)**: ~90-120 seconds

#### Vector Search
- **1K documents**: <10ms
- **10K documents**: ~50ms
- **100K documents**: ~500ms

#### Backup Operations
- **Full Library (1000 tracks)**: ~30 seconds (Spotify API)
- **Restore (1000 tracks)**: ~5 seconds (no API calls)
- **Index (1000 tracks)**: ~2 minutes (with embeddings)

### Scalability Limits

#### Memory
- **Vector Store**: ~1MB per 1000 documents (with 768-dim embeddings)
- **Example**: 10K tracks ≈ 10MB RAM

#### Disk
- **Embeddings**: ~500KB per 1000 documents (JSON serialized)
- **Backups**: ~1MB per 1000 tracks (depends on metadata)

#### Network
- **Spotify API**: Rate limited to ~100 requests/sec
- **Ollama**: No rate limit (local), bottleneck is model inference speed

### Optimization Opportunities

1. **Batch Embeddings**: Already implemented with configurable parallelism
2. **Incremental Indexing**: Only embed new/changed documents
3. **Approximate Search**: Use HNSW or FAISS for >100K documents
4. **Caching**: Cache embeddings for common queries
5. **Compression**: Use msgpack or protobuf instead of JSON

## Technology Stack

### Core
- **Language**: Go 1.23+
- **Build Tool**: Make, Go modules

### Dependencies
- **CLI Framework**: [spf13/cobra](https://github.com/spf13/cobra)
- **TUI Framework**: [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- **Spotify SDK**: [zmb3/spotify](https://github.com/zmb3/spotify)
- **OAuth2**: [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)

### External Services
- **Spotify Web API**: REST API for music library access
- **Ollama**: Local AI inference server (HTTP API)

### Development Tools
- **Testing**: Go standard library testing
- **Linting**: golangci-lint
- **Formatting**: gofmt, goimports
- **Pre-commit**: pre-commit hooks for quality checks
- **CI/CD**: GitHub Actions

## Deployment

### Installation
```bash
# From source
git clone https://github.com/bkataru/spotigo
cd spotigo
make build

# Or use pre-built binary
make install  # Installs to /usr/local/bin
```

### Prerequisites
1. **Ollama**: Must be running on localhost:11434
   ```bash
   ollama serve
   ```

2. **Models**: Pull required models
   ```bash
   ollama pull llama3.2
   ollama pull nomic-embed-text
   ```

3. **Spotify App**: Register app at https://developer.spotify.com
   - Note Client ID and Secret
   - Set redirect URI to `http://localhost:8080/callback`

### Configuration
```yaml
# ~/.spotigo/config.yaml
spotify:
  client_id: "your-client-id"
  client_secret: "your-client-secret"
  redirect_uri: "http://localhost:8080/callback"

ollama:
  base_url: "http://localhost:11434"
  model_tier: "medium"  # small, medium, large

rag:
  embedding_model: "nomic-embed-text"
  search_limit: 10
  enable_auto_index: true
```

## Future Enhancements

### Planned Features
- [ ] Playlist generation from natural language
- [ ] Duplicate detection in library
- [ ] Advanced filters (genre, year, mood)
- [ ] Export to other formats (CSV, XML)
- [ ] Multi-user support (separate configs)

### Performance Improvements
- [ ] Vector index optimization (HNSW)
- [ ] Parallel Spotify API requests
- [ ] Streaming backup/restore
- [ ] Delta backups (only changed items)

### Security Enhancements
- [ ] Optional password-based encryption
- [ ] Keychain integration (macOS/Linux)
- [ ] Token rotation policies
- [ ] Audit logging

## References

- [Spotify Web API Documentation](https://developer.spotify.com/documentation/web-api)
- [Ollama API Documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)
- [RAG Overview](https://arxiv.org/abs/2005.11401)
- [AES-GCM Specification](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf)
