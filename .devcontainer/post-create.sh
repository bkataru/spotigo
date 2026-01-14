#!/bin/bash
set -e

echo ""
echo "===== Spotigo 2.0 Post-Create Setup ====="
echo ""

cd /workspace

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

# Wait for Ollama to be ready
echo "Waiting for Ollama service..."
max_attempts=30
attempt=0
while ! curl -s http://ollama:11434/api/tags > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "Warning: Ollama not ready after ${max_attempts} attempts. Models will be pulled manually."
        break
    fi
    echo "Waiting for Ollama... (attempt $attempt/$max_attempts)"
    sleep 2
done

# Pull required models if Ollama is available
if curl -s http://ollama:11434/api/tags > /dev/null 2>&1; then
    echo ""
    echo "Pulling AI models (this may take a while on first run)..."
    echo ""
    
    # Core models - ordered by priority
    models=(
        "granite4:1b"           # Primary chat (3.3GB)
        "granite4:350m"         # Fast tasks (708MB)
        "nomic-embed-text-v2-moe"  # Embeddings (957MB)
        "qwen3:0.6b"            # Fallback chat (522MB)
    )
    
    for model in "${models[@]}"; do
        echo "Pulling $model..."
        curl -s http://ollama:11434/api/pull -d "{\"name\": \"$model\"}" | while read line; do
            status=$(echo "$line" | jq -r '.status // empty')
            if [ -n "$status" ]; then
                echo "  $status"
            fi
        done
        echo "  Done: $model"
    done
    
    echo ""
    echo "Optional models (pull manually if needed):"
    echo "  ollama pull qwen3:1.7b           # Complex reasoning (1.4GB)"
    echo "  ollama pull functiongemma        # Function calling (300MB)"
    echo "  ollama pull qwen3-embedding:0.6b # Fallback embeddings (639MB)"
fi

echo ""
echo "===== Setup Complete ====="
echo ""
echo "Available commands:"
echo "  go run ./cmd/spotigo           # Run Spotigo CLI"
echo "  go run ./cmd/spotigo backup    # Backup Spotify library"
echo "  go run ./cmd/spotigo chat      # Start AI chat"
echo "  go run ./cmd/spotigo --tui     # Launch TUI mode"
echo ""

exec "$@"
