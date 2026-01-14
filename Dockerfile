# Production Dockerfile for Spotigo 2.0
FROM golang:1.23-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o spotigo \
    ./cmd/spotigo

# Runtime stage
FROM alpine:latest

# Install certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 spotigo && \
    adduser -D -s /bin/sh -u 1000 -G spotigo spotigo

# Copy binary from builder
COPY --from=builder /app/spotigo .

# Create data directory
RUN mkdir -p /app/data && \
    chown -R spotigo:spotigo /app

# Switch to non-root user
USER spotigo

# Expose port for OAuth callback (if needed)
EXPOSE 8888

# Default command
ENTRYPOINT ["./spotigo"]
CMD ["--help"]

# Labels for container metadata
LABEL org.opencontainers.image.title="Spotigo 2.0"
LABEL org.opencontainers.image.description="AI-powered local music intelligence platform"
LABEL org.opencontainers.image.version="2.0.0"
LABEL org.opencontainers.image.source="https://github.com/bkataru-workshop/spotigo"

