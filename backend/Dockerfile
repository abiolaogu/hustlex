# =============================================================================
# HustleX Backend Dockerfile
# =============================================================================
# Multi-stage build for minimal production image
# Build: docker build -t hustlex-api .
# Run: docker run -p 8080:8080 --env-file .env hustlex-api

# -----------------------------------------------------------------------------
# Stage 1: Build
# -----------------------------------------------------------------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o hustlex-api \
    ./cmd/api

# Compress binary with UPX (reduces size by ~60%)
RUN upx --best --lzma hustlex-api

# -----------------------------------------------------------------------------
# Stage 2: Production
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1000 -S hustlex && \
    adduser -u 1000 -S hustlex -G hustlex

# Set timezone
ENV TZ=Africa/Lagos

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/hustlex-api .

# Copy any required static files or templates
# COPY --from=builder /build/templates ./templates
# COPY --from=builder /build/static ./static

# Create directories for logs and uploads
RUN mkdir -p /var/log/hustlex /app/uploads && \
    chown -R hustlex:hustlex /var/log/hustlex /app/uploads

# Switch to non-root user
USER hustlex

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# Run the application
ENTRYPOINT ["./hustlex-api"]

# -----------------------------------------------------------------------------
# Stage 3: Development (optional - for hot reloading)
# -----------------------------------------------------------------------------
FROM golang:1.22-alpine AS development

# Install development tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (will be overwritten by volume mount)
COPY . .

# Expose port
EXPOSE 8080

# Default command for development
CMD ["air", "-c", ".air.toml"]
