# TODO

## Testing
- [x] Check test coverage across all packages
- [x] Add tests for Spotify client package (40.9% coverage)
- [x] Add tests for JSON query engine (73.2% coverage)
- [x] Add tests for Windows-specific browser opening (rundll32)
- [x] Add tests for tools package (80.2% coverage)
- [x] Add integration tests for chat tool-calling flow
- [x] Add tests for internal/cmd package (backup, models, search, stats)
- [x] Add integration tests for OAuth flow with mock server
- [x] Add tests for TUI package (98.8% coverage!)
- [x] Increase overall coverage to >60% across all packages (achieved)

## Documentation
- [x] Create comprehensive GitHub Pages landing page
- [x] Add benchmarks page to gh-pages
- [x] Add more usage examples to README (AI chat, tools, query engine)
- [x] Document configuration options in detail
- [x] Add comprehensive tool documentation (docs/TOOLS.md)
- [x] Add query syntax documentation (docs/QUERY_SYNTAX.md)
- [ ] Add API documentation with godoc examples
- [ ] Create CONTRIBUTING.md guide
- [ ] Add architecture documentation

## JSON Query Tool for RAG
- [x] Create powerful JSON query/manipulation tool
- [x] Add support for filtering, sorting, aggregation
- [x] Add music-specific query helpers
- [x] Write comprehensive tests for query engine
- [x] Integrate JSON query tool with AI chat
- [x] Add tool-calling interface for LLM
- [x] Add comprehensive tests for tools package (80.2% coverage)
- [x] Add integration tests for chat tool-calling flow
- [x] Create comprehensive tool usage documentation (docs/TOOLS.md)
- [x] Create documentation for query syntax (docs/QUERY_SYNTAX.md - 818 lines)

## RAG Improvements
- [x] Implement tool-calling for structured JSON queries
- [ ] Add schema-aware chunking for JSON embeddings
- [ ] Create hybrid search (embeddings + structured queries)
- [ ] Optimize context window usage
- [ ] Add query result caching
- [x] Add documentation for tool usage with examples

## Release & CI/CD
- [x] Fix GitHub release workflow
- [x] Successfully publish v2.0.0 release
- [x] Fix Docker build issues
- [x] Create multi-platform binaries
- [ ] Add automated release notes generation
- [ ] Set up automated security scanning
- [ ] Add release checklist documentation

## Future Enhancements
- [ ] Add support for more Spotify data types in backups
  - [ ] Saved albums
  - [ ] Followed podcasts
  - [ ] Listening history (extended)
- [ ] Improve error messages and user feedback
- [ ] Add progress indicators for long-running operations
- [ ] Add configurable concurrency limits via CLI flags
- [ ] Add retry/backoff for Spotify API rate limits
- [ ] Add telemetry/logging (opt-in)
- [ ] Create web UI for music exploration
- [ ] Add playlist generation based on AI recommendations
- [ ] Support for multiple music services (Apple Music, YouTube Music)

## Code Quality
- [x] Fix all linting issues
- [x] Add pre-commit hooks
- [x] Pass CI checks
- [x] Add code coverage reporting to CI (with PR comments)
- [ ] Set up dependabot for dependency updates
- [ ] Add performance regression tests
- [ ] Document all public APIs

## User Experience
- [x] Redesign TUI with cyberpunk nuclear green theme
- [x] Add ASCII art logo to TUI
- [x] Improve navigation with vim-style keys (Ctrl+J/K)
- [x] Add status bar and help text to TUI
- [x] Create compact view for small terminals
- [ ] Add interactive setup wizard
- [ ] Improve CLI help messages
- [ ] Add shell completions (bash, zsh, fish)
- [ ] Create quick-start guide
- [ ] Add troubleshooting guide
- [ ] Improve error messages with actionable suggestions

## Performance
- [x] Optimize backup with concurrency
- [x] Add buffered I/O for JSON operations
- [x] Add query result caching in JSON engine
- [ ] Add incremental backup support
- [ ] Optimize embedding generation
- [ ] Add disk-based vector store option for large libraries
- [ ] Profile and optimize hot paths

## Test Coverage Summary (Latest)
- ✅ internal/tui: 98.8% (EXCELLENT!)
- ✅ internal/jsonutil: 96.8%
- ✅ internal/crypto: 80.5%
- ✅ internal/tools: 80.2%
- ✅ internal/jsonquery: 73.2%
- ✅ internal/storage: 67.9%
- ⚠️ internal/config: 56.8%
- ⚠️ internal/rag: 53.9%
- ⚠️ internal/ollama: 46.3%
- ⚠️ internal/spotify: 40.9%
- ⚠️ internal/cmd: 3.7%
- **Overall**: ~60% average across internal packages