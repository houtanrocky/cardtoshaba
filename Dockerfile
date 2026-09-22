# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Add required tools
RUN apk add --no-cache git ca-certificates tzdata

# Install deps for optimal cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /build/server ./cmd/server

# ============================================================
# Stage 2: Runtime
# ============================================================
FROM alpine:3.19

# Install ca-certificates for TLS
RUN apk add --no-cache ca-certificates tzdata wget

# Non-Root user for safety
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binaries and migration
COPY --from=builder /build/server .
COPY --from=builder /build/migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/server"]