# TODO

## Testing
- [x] Check test coverage across all packages
- [x] Add tests for Spotify client package (40.9% coverage)
- [x] Add tests for JSON query engine (73.2% coverage)
- [x] Add tests for Windows-specific browser opening (rundll32)
- [ ] Improve test coverage for internal/cmd package (currently 2.9%)
- [ ] Add integration tests for OAuth flow with mock server
- [ ] Add tests for TUI package (currently 0.0%)
- [ ] Increase overall coverage to >60% across all packages

## Documentation
- [x] Create comprehensive GitHub Pages landing page
- [x] Add benchmarks page to gh-pages
- [ ] Add more usage examples to README
- [ ] Document configuration options in detail
- [ ] Add API documentation with godoc examples
- [ ] Create CONTRIBUTING.md guide
- [ ] Add architecture documentation

## JSON Query Tool for RAG
- [x] Create powerful JSON query/manipulation tool
- [x] Add support for filtering, sorting, aggregation
- [x] Add music-specific query helpers
- [x] Write comprehensive tests for query engine
- [ ] Integrate JSON query tool with AI chat
- [ ] Add tool-calling interface for LLM
- [ ] Create documentation for query syntax

## RAG Improvements
- [ ] Implement tool-calling for structured JSON queries
- [ ] Add schema-aware chunking for JSON embeddings
- [ ] Create hybrid search (embeddings + structured queries)
- [ ] Optimize context window usage
- [ ] Add query result caching

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
- [ ] Add code coverage reporting to CI
- [ ] Set up dependabot for dependency updates
- [ ] Add performance regression tests
- [ ] Document all public APIs

## User Experience
- [ ] Add interactive setup wizard
- [ ] Improve CLI help messages
- [ ] Add shell completions (bash, zsh, fish)
- [ ] Create quick-start guide
- [ ] Add troubleshooting guide
- [ ] Improve error messages with actionable suggestions

## Performance
- [x] Optimize backup with concurrency
- [x] Add buffered I/O for JSON operations
- [ ] Add incremental backup support
- [ ] Optimize embedding generation
- [ ] Add disk-based vector store option for large libraries
- [ ] Profile and optimize hot paths