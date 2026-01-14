# Spotigo 2.0 Development Container
# Go + Ollama for local AI-powered music intelligence

FROM golang:1.25-bookworm

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    wget \
    curl \
    unzip \
    zsh \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Install Elvish shell (user preference)
RUN wget -q https://dl.elv.sh/linux-amd64/elvish-HEAD.tar.gz \
    && tar -C /usr/local/bin -xzf elvish-HEAD.tar.gz \
    && rm elvish-HEAD.tar.gz

# Install Oh My Zsh for better shell experience
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended

# Install Go tools
RUN go install github.com/air-verse/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/gopls@latest

WORKDIR /workspace

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

# Default command
CMD ["zsh"]
