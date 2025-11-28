# ================================
# Build Stage
# ================================
FROM golang:tip-alpine3.22 AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0: Static binary (no C dependencies)
# -ldflags: Strip debug info and reduce binary size
# -trimpath: Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
  -trimpath \
  -o api \
  cmd/api/main.go

# ================================
# Runtime Stage
# ================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
  ca-certificates \
  tzdata \
  && rm -rf /var/cache/apk/*

# Create non-root user and group
# Use specific UID/GID for consistency
RUN addgroup -g 1000 -S appuser && \
  adduser -u 1000 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /build/api /app/api

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV ENV=production \
  PORT=8080

# Run the application
ENTRYPOINT ["/app/api"]
