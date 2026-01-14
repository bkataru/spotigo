# Changelog

All notable changes to Spotigo will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CHANGELOG.md for tracking version history
- SECURITY.md with security policy and vulnerability reporting
- CODE_OF_CONDUCT.md with Contributor Covenant
- ARCHITECTURE.md with detailed technical documentation
- GitHub issue templates (bug report, feature request)
- GitHub pull request template
- Dependabot configuration for automated dependency updates
- EditorConfig for consistent code style across IDEs
- Enhanced pre-commit hooks configuration
- Development setup script (`.scripts/setup-dev.sh`)
- Benchmark tests for RAG vector store operations
- Basic integration test framework
- Benchmark CI job for performance regression detection
- Docker Compose for development environment
- FAQ section in README
- Dependency security scanning in CI (govulncheck)
- Additional Makefile targets (bench, security-scan, integration-test)

### Changed
- README.md updated with additional badges (codecov, Go Report Card)
- README.md updated with demo/screenshot placeholders
- Improved pre-commit hooks with additional checks

### Fixed
- Nil pointer dereference in chat.go when config is nil
- Inconsistent backup ID handling between CreateBackup and ListBackups
- Progress display showing stale percentage in model pulling
- Potential data race in parallel embedding generation

### Removed
- Unused `sanitizeInput` function in chat.go

### Documentation
- Updated LICENSE copyright year to 2026 and author to bkataru (Baalateja Kataru)
- Updated Go version requirements from 1.23+ to 1.24+ across all documentation
- Updated Go version badge in README to reflect 1.24+ requirement

## [2.0.0] - 2024-01-14

### Added
- **Complete OAuth2 Authentication System**
  - Spotify OAuth2 flow with browser integration
  - Secure local token storage with AES-256-GCM encryption
  - Authentication status verification with user profile display
  - Graceful logout functionality

- **Full Backup System**
  - Complete Spotify API integration (saved tracks, playlists, followed artists)
  - Local JSON storage with metadata and timestamps
  - Backup management commands (list, restore, status)
  - Incremental and full backup support
  - `--index` flag to build search index after backup

- **AI Chat Interface**
  - Multi-model Ollama integration (chat, fast, reasoning, embeddings)
  - Contextual conversations with conversation memory
  - Fallback model handling for robustness
  - Model configuration via YAML with customizable system prompts
  - Input validation (length, UTF-8, control characters)
  - Conversation history trimming

- **Semantic Search (RAG)**
  - Vector store with Ollama embeddings
  - Cosine similarity search
  - Parallel embedding generation with configurable concurrency
  - Search by type (tracks, artists, playlists)
  - Persistent vector store with JSON serialization

- **Statistics Dashboard**
  - Top artists and albums analysis
  - Genre distribution
  - Playlist analysis
  - Listening timeline

- **Professional Terminal UI**
  - Beautiful BubbleTea interface with navigation
  - Interactive menu system with keyboard controls
  - Themeable design with Lipgloss styling

- **Model Management**
  - List configured models
  - Check Ollama connection and model availability
  - Native model pulling via Ollama API with progress display

- **Production Infrastructure**
  - GitHub Actions CI/CD pipeline
  - Multi-platform builds (Linux, Windows, macOS - amd64 & arm64)
  - Security scanning with gosec
  - Docker multi-platform builds
  - GitHub Container Registry publishing
  - Comprehensive Makefile with development targets
  - Dev container configuration with Ollama

- **Developer Experience**
  - Cobra CLI framework
  - Viper configuration management
  - golangci-lint configuration
  - Pre-commit hooks
  - Comprehensive documentation

### Changed
- Complete rewrite from basic backup tool to AI-powered music intelligence platform
- Architecture redesigned with clean package separation
- Configuration system with YAML files and environment variable support

## [1.0.0] - 2024-01-13

### Added
- Initial commit with basic project structure
- Basic Spotify backup functionality concept

---

## Version History Summary

| Version | Date | Highlights |
|---------|------|------------|
| 2.0.0 | 2024-01-14 | Complete AI-powered platform with OAuth, backup, chat, search, and TUI |
| 1.0.0 | 2024-01-13 | Initial project structure |

[Unreleased]: https://github.com/bkataru/spotigo/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/bkataru/spotigo/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/bkataru/spotigo/releases/tag/v1.0.0