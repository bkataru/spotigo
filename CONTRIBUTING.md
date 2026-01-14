# Contributing to Spotigo Go Library

Thank you for your interest in contributing to the Spotigo Go library! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

Please be respectful and constructive in all interactions. We're building something together!

## Getting Started

### Prerequisites

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Ollama** - [Install](https://ollama.ai/) (for testing AI features)
- **Git** - For version control

### Development Environment

Standard Go development setup is recommended for library development.

## Development Setup

### Local Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/spotigo.git
cd spotigo

# Add upstream remote
git remote add upstream https://github.com/bkataru/spotigo.git

# Install dependencies
go mod download

# Verify setup
go test ./...
```

## Making Changes

### Branching Strategy

- `main` - Stable release branch
- `develop` - Integration branch for features
- `feature/*` - Feature branches
- `fix/*` - Bug fix branches
- `docs/*` - Documentation changes

### Workflow

1. **Fork** the repository
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```
3. **Make changes** following our code style
4. **Test** your changes:
   ```bash
   go test ./...
   golangci-lint run
   ```
5. **Commit** with clear messages:
   ```bash
   git commit -m "feat: add new API feature"
   ```
6. **Push** to your fork:
   ```bash
   git push origin feature/my-feature
   ```
7. **Open a Pull Request** against `main`

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Code style (formatting, semicolons, etc.)
- `refactor` - Code change that neither fixes a bug nor adds a feature
- `perf` - Performance improvement
- `test` - Adding or updating tests
- `chore` - Maintenance tasks

**Examples:**
```
feat(backup): add incremental backup support
fix(auth): handle token refresh edge case
docs(readme): update installation instructions
test(rag): add vector store benchmarks
```

## Code Style

### Go Standards

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project Conventions

1. **Package Comments** - Every package must have a doc comment:
   ```go
   // Package rag provides RAG (Retrieval-Augmented Generation) functionality
   // with vector embeddings for semantic search across music library data.
   package rag
   ```

2. **Error Handling** - Wrap errors with context:
   ```go
   if err != nil {
       return fmt.Errorf("failed to load config: %w", err)
   }
   ```

3. **Interfaces** - Define interfaces where they're used, not where they're implemented

4. **Naming**:
   - Use `MixedCaps` for exported names
   - Use short, descriptive names for local variables
   - Acronyms should be consistent case (`URL`, `ID`, not `Url`, `Id`)

### Linting

We use `golangci-lint` for comprehensive linting:

```bash
golangci-lint run
```

Common linter issues:
- Unused variables or imports
- Missing error handling
- Shadowed variables
- Inefficient code patterns

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/rag/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

1. **Table-Driven Tests** - Prefer table-driven tests for multiple cases:
   ```go
   func TestCosineSimilarity(t *testing.T) {
       tests := []struct {
           name     string
           a, b     []float64
           expected float64
       }{
           {"identical", []float64{1, 0}, []float64{1, 0}, 1.0},
           {"orthogonal", []float64{1, 0}, []float64{0, 1}, 0.0},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result := cosineSimilarity(tt.a, tt.b)
               if result != tt.expected {
                   t.Errorf("got %v, want %v", result, tt.expected)
               }
           })
       }
   }
   ```

2. **Test Files** - Name test files with `_test.go` suffix

3. **Test Helpers** - Use `t.Helper()` for test helper functions

4. **Cleanup** - Use `t.Cleanup()` or `defer` for test cleanup

### Test Coverage Goals

- **Core packages** (`rag`, `config`, `storage`): 80%+ coverage
- **Utility packages** (`jsonutil`, `crypto`): 90%+ coverage
- **Integration tests**: Cover main user flows

## Pull Request Process

### Before Submitting

- [ ] Tests pass locally (`go test ./...`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventions

### PR Template

When opening a PR, include:

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Testing
How were these changes tested?

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Changelog updated (if applicable)
```

### Review Process

1. **Automated Checks** - CI must pass
2. **Code Review** - At least one approval required
3. **Merge** - Squash and merge preferred

## Release Process

Releases are created automatically when a version tag is pushed:

```bash
# Create and push a tag
git tag v2.1.0
git push origin v2.1.0
```

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** - Breaking changes
- **MINOR** - New features (backwards compatible)
- **PATCH** - Bug fixes (backwards compatible)

## Project Structure

```
spotigo/
├── internal/            # Library implementation
│   ├── config/          # Configuration management
│   ├── crypto/          # Encryption utilities
│   ├── jsonutil/        # JSON helpers
│   ├── ollama/          # Ollama API client
│   ├── rag/             # RAG vector store
│   ├── spotify/         # Spotify API client
│   └── storage/         # Local file storage
├── cmd/spotigo/         # CLI example application
└── .github/             # GitHub Actions
```

### Adding New Packages

1. Create package in `internal/`
2. Add package doc comment
3. Add tests (`*_test.go`)
4. Update this documentation if significant

## Getting Help

- **Issues** - Open a GitHub issue for bugs or features
- **Discussions** - Use GitHub Discussions for questions

## Recognition

Contributors are recognized in:
- Release notes
- README acknowledgments (for significant contributions)

Thank you for contributing to Spotigo!
