# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /spotigo ./cmd/spotigo

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' spotigo

WORKDIR /home/spotigo

# Copy binary from builder
COPY --from=builder /spotigo /usr/local/bin/spotigo

# Create data directories
RUN mkdir -p data/backups data/embeddings && \
  chown -R spotigo:spotigo /home/spotigo

# Switch to non-root user
USER spotigo

# Default command
ENTRYPOINT ["spotigo"]
CMD ["--help"]
