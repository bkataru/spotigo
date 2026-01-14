#!/bin/bash
# Development setup script for Spotigo
# Usage: ./.scripts/setup-dev.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}                 Spotigo Development Setup                      ${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print status
print_status() {
    if [ "$2" = "ok" ]; then
        echo -e "  ${GREEN}✓${NC} $1"
    elif [ "$2" = "warn" ]; then
        echo -e "  ${YELLOW}⚠${NC} $1"
    elif [ "$2" = "error" ]; then
        echo -e "  ${RED}✗${NC} $1"
    else
        echo -e "  ${BLUE}→${NC} $1"
    fi
}

# Check Go installation
echo -e "${YELLOW}Checking prerequisites...${NC}"
echo ""

if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Go installed: $GO_VERSION" "ok"
else
    print_status "Go not installed" "error"
    echo "    Install Go from: https://go.dev/dl/"
    exit 1
fi

# Check Go version (require 1.23+)
GO_MAJOR=$(go version | awk '{print $3}' | sed 's/go//' | cut -d. -f1)
GO_MINOR=$(go version | awk '{print $3}' | sed 's/go//' | cut -d. -f2)
if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 23 ]); then
    print_status "Go version 1.23+ required (you have go$GO_MAJOR.$GO_MINOR)" "warn"
fi

# Check Git
if command_exists git; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    print_status "Git installed: $GIT_VERSION" "ok"
else
    print_status "Git not installed" "error"
    exit 1
fi

# Check Ollama
if command_exists ollama; then
    print_status "Ollama installed" "ok"
else
    print_status "Ollama not installed" "warn"
    echo "    Install from: https://ollama.ai/"
    echo "    Ollama is required for AI features"
fi

# Check Make
if command_exists make; then
    print_status "Make installed" "ok"
else
    print_status "Make not installed (optional)" "warn"
fi

# Check Docker (optional)
if command_exists docker; then
    print_status "Docker installed" "ok"
else
    print_status "Docker not installed (optional)" "warn"
fi

echo ""
echo -e "${YELLOW}Installing Go dependencies...${NC}"
go mod download
print_status "Go modules downloaded" "ok"

go mod tidy
print_status "Go modules tidied" "ok"

echo ""
echo -e "${YELLOW}Installing development tools...${NC}"

# Install golangci-lint
if command_exists golangci-lint; then
    print_status "golangci-lint already installed" "ok"
else
    print_status "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    print_status "golangci-lint installed" "ok"
fi

# Install govulncheck
if command_exists govulncheck; then
    print_status "govulncheck already installed" "ok"
else
    print_status "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
    print_status "govulncheck installed" "ok"
fi

# Install goimports
if command_exists goimports; then
    print_status "goimports already installed" "ok"
else
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    print_status "goimports installed" "ok"
fi

# Check/install pre-commit (optional)
echo ""
if command_exists pre-commit; then
    print_status "pre-commit installed" "ok"
    echo ""
    echo -e "${YELLOW}Setting up pre-commit hooks...${NC}"
    pre-commit install
    pre-commit install --hook-type commit-msg
    print_status "Pre-commit hooks installed" "ok"
else
    print_status "pre-commit not installed (optional)" "warn"
    echo "    Install with: pip install pre-commit"
    echo "    Then run: pre-commit install"
fi

# Create config file if it doesn't exist
echo ""
echo -e "${YELLOW}Checking configuration...${NC}"
if [ -f "spotigo.yaml" ]; then
    print_status "spotigo.yaml exists" "ok"
else
    print_status "Creating spotigo.yaml from template..."
    cp spotigo.example.yaml spotigo.yaml
    print_status "spotigo.yaml created" "ok"
    echo ""
    echo -e "    ${YELLOW}⚠ Please edit spotigo.yaml with your Spotify API credentials${NC}"
    echo "    Get credentials from: https://developer.spotify.com/dashboard"
fi

# Create data directories
echo ""
echo -e "${YELLOW}Creating data directories...${NC}"
mkdir -p data/backups data/embeddings
print_status "data/backups created" "ok"
print_status "data/embeddings created" "ok"

# Build the project
echo ""
echo -e "${YELLOW}Building Spotigo...${NC}"
go build -o bin/spotigo ./cmd/spotigo
print_status "Binary built: bin/spotigo" "ok"

# Run tests
echo ""
echo -e "${YELLOW}Running tests...${NC}"
if go test -race ./... > /dev/null 2>&1; then
    print_status "All tests passed" "ok"
else
    print_status "Some tests failed" "warn"
    echo "    Run 'go test -v ./...' for details"
fi

# Summary
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}                    Setup Complete!                             ${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo "Next steps:"
echo ""
echo "  1. Edit spotigo.yaml with your Spotify API credentials"
echo "     Get them from: https://developer.spotify.com/dashboard"
echo ""
echo "  2. Start Ollama (in a separate terminal):"
echo "     $ ollama serve"
echo ""
echo "  3. Pull AI models:"
echo "     $ ollama pull granite4:1b"
echo "     $ ollama pull nomic-embed-text-v2-moe"
echo ""
echo "  4. Authenticate with Spotify:"
echo "     $ ./bin/spotigo auth"
echo ""
echo "  5. Back up your library:"
echo "     $ ./bin/spotigo backup --index"
echo ""
echo "  6. Start chatting!"
echo "     $ ./bin/spotigo chat"
echo ""
echo "Development commands:"
echo "  make build       - Build the binary"
echo "  make test        - Run tests"
echo "  make lint        - Run linter"
echo "  make quality     - Run all quality checks"
echo "  make help        - Show all available commands"
echo ""
