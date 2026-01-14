# Makefile for Spotigo 2.0

.PHONY: help build test clean lint fmt vet install docker-build docker-run release

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
build: ## Build the application
	@echo 'Building Spotigo...'
	@go build -o bin/spotigo ./cmd/spotigo
	@echo '✅ Binary created: bin/spotigo'

build-all: ## Build for all platforms
	@echo 'Building for multiple platforms...'
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/spotigo-linux-amd64 ./cmd/spotigo
	GOOS=windows GOARCH=amd64 go build -o bin/spotigo-windows-amd64.exe ./cmd/spotigo
	GOOS=darwin GOARCH=amd64 go build -o bin/spotigo-darwin-amd64 ./cmd/spotigo
	@echo '✅ Built for Linux, Windows, macOS'

test: ## Run tests
	@echo 'Running tests...'
	@go test -v ./...

clean: ## Clean build artifacts
	@echo 'Cleaning...'
	@rm -rf bin/
	@go clean -cache
	@echo '✅ Cleaned'

lint: ## Run linter
	@echo 'Running linter...'
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

fmt: ## Format code
	@echo 'Formatting code...'
	@go fmt ./...
	@echo '✅ Code formatted'

vet: ## Run go vet
	@echo 'Running go vet...'
	@go vet ./...
	@echo '✅ Vet passed'

install: ## Install the application locally
	@echo 'Installing Spotigo...'
	@go install ./cmd/spotigo
	@echo '✅ Installed to GOPATH'

# Docker targets
docker-build: ## Build Docker image
	@echo 'Building Docker image...'
	@docker build -t spotigo:latest .
	@echo '✅ Docker image built'

docker-run: ## Run Docker container
	@echo 'Running Spotigo in Docker...'
	@docker run --rm -it \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/config:/app/config \
		-p 8888:8888 \
		spotigo:latest

# Development targets
dev-setup: ## Set up development environment
	@echo 'Setting up development environment...'
	@go mod download
	@go mod tidy
	@if command -v ollama >/dev/null 2>&1; then \
		echo "✅ Ollama is installed"; \
	else \
		echo "Installing Ollama..."; \
		curl -fsSL https://ollama.ai/install.sh | sh; \
	fi
	@echo '✅ Dev environment ready'

dev-run: ## Run in development mode
	@echo 'Starting Spotigo in development mode...'
	@go run ./cmd/spotigo

# Production targets
release: ## Create release build
	@echo 'Creating release build...'
	@mkdir -p dist
	@echo 'Building release binaries...'
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/spotigo-linux-amd64 ./cmd/spotigo
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/spotigo-windows-amd64.exe ./cmd/spotigo
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/spotigo-darwin-amd64 ./cmd/spotigo
	@echo 'Creating checksums...'
	@cd dist && sha256sum * > checksums.txt
	@echo '✅ Release build complete'

# Quality checks
quality: ## Run all quality checks
	@echo 'Running all quality checks...'
	@make fmt
	@make vet
	@make test
	@echo '✅ All quality checks passed'

# Status targets
status: ## Show project status
	@echo 'Spotigo 2.0 - Project Status'
	@echo '============================'
	@echo 'Git Branch: $(shell git branch --show-current)'
	@echo 'Git Commit: $(shell git log -1 --oneline)'
	@echo 'Go Version: $(shell go version)'
	@echo 'Dependencies: $(shell go list -m | wc -l)'
	@echo ''
	@echo 'Last build: $(shell ls -la bin/spotigo 2>/dev/null && echo "$(shell date -r bin/spotigo)" || echo "Never")'

# Quick commands
quick-test: ## Quick test (no deps install)
	@go build ./cmd/spotigo && ./spotigo --help > /dev/null && echo '✅ Build and help test passed'

deps: ## Install dependencies
	@echo 'Installing dependencies...'
	@go mod download
	@go mod tidy
	@echo '✅ Dependencies installed'