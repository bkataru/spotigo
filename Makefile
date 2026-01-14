# Makefile for Spotigo 2.0
# AI-powered local music intelligence platform

.PHONY: help build build-all test clean lint fmt vet install docker-build docker-run release dev-setup dev-run quality status deps quick-test coverage bench bench-rag bench-compare security-scan outdated integration-test test-short pre-commit hooks-install hooks-run ci-bench

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Default target
help: ## Show this help message
	@echo 'Spotigo 2.0 - Development Commands'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# Build Targets
# =============================================================================

build: ## Build for current platform
	@echo 'Building Spotigo $(VERSION)...'
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/spotigo ./cmd/spotigo
	@echo '✅ Binary created: bin/spotigo'

build-all: ## Build for all platforms (including ARM64)
	@echo 'Building for multiple platforms...'
	@mkdir -p bin
	@echo '  Building linux/amd64...'
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/spotigo-linux-amd64 ./cmd/spotigo
	@echo '  Building linux/arm64...'
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/spotigo-linux-arm64 ./cmd/spotigo
	@echo '  Building windows/amd64...'
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/spotigo-windows-amd64.exe ./cmd/spotigo
	@echo '  Building darwin/amd64...'
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/spotigo-darwin-amd64 ./cmd/spotigo
	@echo '  Building darwin/arm64 (Apple Silicon)...'
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/spotigo-darwin-arm64 ./cmd/spotigo
	@echo '✅ Built for Linux, Windows, macOS (amd64 + arm64)'

install: ## Install to GOPATH/bin
	@echo 'Installing Spotigo...'
	@go install $(LDFLAGS) ./cmd/spotigo
	@echo '✅ Installed to $(shell go env GOPATH)/bin/spotigo'

# =============================================================================
# Testing & Quality
# =============================================================================

test: ## Run tests
	@echo 'Running tests...'
	@go test -race ./...
	@echo '✅ All tests passed'

test-v: ## Run tests with verbose output
	@echo 'Running tests (verbose)...'
	@go test -v -race ./...

coverage: ## Generate coverage report
	@echo 'Generating coverage report...'
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo ''
	@echo 'To view HTML report: go tool cover -html=coverage.out'

coverage-html: coverage ## Generate and open HTML coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo '✅ Coverage report: coverage.html'

lint: ## Run golangci-lint
	@echo 'Running linter...'
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi
	@echo '✅ Lint passed'

fmt: ## Format code
	@echo 'Formatting code...'
	@go fmt ./...
	@echo '✅ Code formatted'

vet: ## Run go vet
	@echo 'Running go vet...'
	@go vet ./...
	@echo '✅ Vet passed'

quality: fmt vet test ## Run all quality checks
	@echo ''
	@echo '✅ All quality checks passed'

# =============================================================================
# Docker Targets
# =============================================================================

docker-build: ## Build Docker image
	@echo 'Building Docker image...'
	@docker build -t spotigo:$(VERSION) -t spotigo:latest .
	@echo '✅ Docker image built: spotigo:$(VERSION)'

docker-run: ## Run Docker container
	@echo 'Running Spotigo in Docker...'
	@docker run --rm -it \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/config:/app/config \
		-p 8888:8888 \
		spotigo:latest

docker-compose-up: ## Start development environment with docker-compose
	@echo 'Starting development environment...'
	@docker-compose -f .devcontainer/docker-compose.yml up -d
	@echo '✅ Development environment started'

docker-compose-down: ## Stop development environment
	@docker-compose -f .devcontainer/docker-compose.yml down

# =============================================================================
# Development Targets
# =============================================================================

deps: ## Download and tidy dependencies
	@echo 'Installing dependencies...'
	@go mod download
	@go mod tidy
	@echo '✅ Dependencies installed'

dev-setup: deps ## Set up development environment
	@echo 'Setting up development environment...'
	@echo 'Checking for Ollama...'
	@if command -v ollama >/dev/null 2>&1; then \
		echo "✅ Ollama is installed"; \
	else \
		echo "⚠️  Ollama not found. Install from: https://ollama.ai/"; \
	fi
	@echo 'Checking for golangci-lint...'
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "✅ golangci-lint is installed"; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo ''
	@echo '✅ Development environment ready'
	@echo ''
	@echo 'Next steps:'
	@echo '  1. Copy spotigo.example.yaml to spotigo.yaml'
	@echo '  2. Add your Spotify API credentials'
	@echo '  3. Start Ollama: ollama serve'
	@echo '  4. Run: make dev-run'

dev-run: ## Run in development mode
	@echo 'Starting Spotigo in development mode...'
	@go run ./cmd/spotigo

dev-run-tui: ## Run TUI in development mode
	@go run ./cmd/spotigo --tui

# =============================================================================
# Release Targets
# =============================================================================

release: clean ## Create release build
	@echo 'Creating release build $(VERSION)...'
	@mkdir -p dist
	@echo 'Building release binaries...'
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/spotigo-linux-amd64 ./cmd/spotigo
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/spotigo-linux-arm64 ./cmd/spotigo
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/spotigo-windows-amd64.exe ./cmd/spotigo
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/spotigo-darwin-amd64 ./cmd/spotigo
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/spotigo-darwin-arm64 ./cmd/spotigo
	@echo 'Creating checksums...'
	@cd dist && sha256sum * > checksums.txt
	@echo ''
	@echo '✅ Release build complete'
	@echo 'Files in dist/:'
	@ls -la dist/

# =============================================================================
# Utility Targets
# =============================================================================

clean: ## Clean build artifacts
	@echo 'Cleaning...'
	@rm -rf bin/ dist/ coverage.out coverage.html
	@go clean -cache
	@echo '✅ Cleaned'

status: ## Show project status
	@echo 'Spotigo 2.0 - Project Status'
	@echo '============================'
	@echo ''
	@echo 'Version:      $(VERSION)'
	@echo 'Commit:       $(COMMIT)'
	@echo 'Go Version:   $(shell go version | cut -d" " -f3)'
	@echo 'Git Branch:   $(shell git branch --show-current 2>/dev/null || echo "N/A")'
	@echo ''
	@echo 'Test Packages:'
	@go list ./... | grep -v cmd | wc -l | xargs echo "  Total:"
	@echo ''
	@echo 'Binary Status:'
	@if [ -f bin/spotigo ]; then \
		echo "  bin/spotigo: $(shell ls -lh bin/spotigo | awk '{print $$5}')"; \
	else \
		echo "  bin/spotigo: Not built (run 'make build')"; \
	fi

quick-test: build ## Quick build and smoke test
	@./bin/spotigo --help > /dev/null
	@echo '✅ Build and help test passed'

# =============================================================================
# Benchmark & Performance Targets
# =============================================================================

bench: ## Run benchmarks
	@echo 'Running benchmarks...'
	@go test -bench=. -benchmem ./...
	@echo '✅ Benchmarks complete'

bench-rag: ## Run RAG-specific benchmarks
	@echo 'Running RAG benchmarks...'
	@go test -bench=. -benchmem -run=^$$ ./internal/rag/...
	@echo '✅ RAG benchmarks complete'

bench-compare: ## Run benchmarks and save for comparison
	@echo 'Running benchmarks for comparison...'
	@go test -bench=. -benchmem -count=5 ./... > bench-current.txt
	@echo '✅ Benchmark results saved to bench-current.txt'
	@echo 'Compare with: benchstat bench-baseline.txt bench-current.txt'

# =============================================================================
# Security & Analysis Targets
# =============================================================================

security-scan: ## Run security scans
	@echo 'Running security scans...'
	@echo 'Installing/updating gosec...'
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo 'Running gosec...'
	@gosec -quiet ./...
	@echo ''
	@echo 'Installing/updating govulncheck...'
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo 'Running govulncheck...'
	@govulncheck ./...
	@echo '✅ Security scans passed'

outdated: ## Check for outdated dependencies
	@echo 'Checking for outdated dependencies...'
	@go install github.com/psampaz/go-mod-outdated@latest
	@go list -u -m -json all | go-mod-outdated -update -direct
	@echo ''
	@echo 'To update all: go get -u ./... && go mod tidy'

# =============================================================================
# Testing Targets (Extended)
# =============================================================================

integration-test: build ## Run integration tests
	@echo 'Running integration tests...'
	@go test -v -tags=integration ./test/integration/...
	@echo '✅ Integration tests passed'

test-short: ## Run short tests only
	@echo 'Running short tests...'
	@go test -short ./...
	@echo '✅ Short tests passed'

# =============================================================================
# Pre-commit & Hooks
# =============================================================================

pre-commit: fmt vet lint test ## Run pre-commit checks (fmt, vet, lint, test)
	@echo ''
	@echo '✅ All pre-commit checks passed'

hooks-install: ## Install git hooks via pre-commit
	@echo 'Installing pre-commit hooks...'
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		echo "✅ Pre-commit hooks installed"; \
	else \
		echo "⚠️  pre-commit not found. Install with: pip install pre-commit"; \
		exit 1; \
	fi

hooks-run: ## Run pre-commit on all files
	@pre-commit run --all-files

# =============================================================================
# CI Targets (used by GitHub Actions)
# =============================================================================

ci-test: ## Run CI tests
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

ci-lint: ## Run CI lint
	@golangci-lint run --timeout=5m

ci-build: ## CI build for current platform
	@go build -v ./cmd/spotigo

ci-bench: ## Run CI benchmarks
	@go test -bench=. -benchmem -run=^$$ -count=3 ./... | tee bench-results.txt
