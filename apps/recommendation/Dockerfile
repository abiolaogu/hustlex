# =============================================================================
# HustleX Recommendation Service Dockerfile
# =============================================================================
# Multi-stage build with Go API + Python ML support
# Build: docker build -t hustlex-recommendation .
# Run: docker run -p 8081:8081 --env-file .env hustlex-recommendation

# -----------------------------------------------------------------------------
# Stage 1: Go Builder
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download 2>/dev/null || go mod tidy

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o recommendation-server \
    ./cmd/server

# Compress binary
RUN upx --best --lzma recommendation-server || true

# -----------------------------------------------------------------------------
# Stage 2: Python ML Builder
# -----------------------------------------------------------------------------
FROM python:3.11-slim AS python-builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install
COPY requirements.txt .

# Create virtual environment and install dependencies
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install CPU-only torch first (smaller image)
RUN pip install --no-cache-dir torch --index-url https://download.pytorch.org/whl/cpu

# Install remaining dependencies
RUN pip install --no-cache-dir -r requirements.txt

# -----------------------------------------------------------------------------
# Stage 3: Production Image
# -----------------------------------------------------------------------------
FROM python:3.11-slim AS production

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r hustlex && useradd -r -g hustlex hustlex

# Set timezone
ENV TZ=Africa/Lagos
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

# Copy Python virtual environment
COPY --from=python-builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy Go binary
COPY --from=go-builder /build/recommendation-server .

# Copy Python ML service
COPY ml_service.py .

# Create directories
RUN mkdir -p /app/models /var/log/hustlex && \
    chown -R hustlex:hustlex /app /var/log/hustlex

# Switch to non-root user
USER hustlex

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# Run the Go server (Python ML is called via subprocess/import)
ENTRYPOINT ["./recommendation-server"]

# -----------------------------------------------------------------------------
# Stage 4: Go-only Production (lighter weight)
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS production-go-only

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 -S hustlex && \
    adduser -u 1000 -S hustlex -G hustlex

ENV TZ=Africa/Lagos

WORKDIR /app

# Copy Go binary only
COPY --from=go-builder /build/recommendation-server .

# Create directories
RUN mkdir -p /var/log/hustlex && \
    chown -R hustlex:hustlex /app /var/log/hustlex

USER hustlex

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

ENTRYPOINT ["./recommendation-server"]

# -----------------------------------------------------------------------------
# Stage 5: Development
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS development

# Install development tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    python3 \
    py3-pip

# Install Air for hot reloading
RUN go install github.com/air-verse/air@v1.52.3

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || go mod tidy

# Copy source
COPY . .

EXPOSE 8081

CMD ["air", "-c", ".air.toml"]
