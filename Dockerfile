# ============================================================
# Multi-stage Dockerfile for PMS API
# ============================================================

# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /bin/pms-api ./cmd/api

# --- Stage 2: Runtime ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Non-root user
RUN addgroup -S pms && adduser -S pms -G pms

WORKDIR /app

# Copy binary and config
COPY --from=builder /bin/pms-api .
COPY config.yaml .
COPY config.production.yaml .
COPY migrations/ ./migrations/

# Create log directory
RUN mkdir -p logs && chown -R pms:pms /app

USER pms

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./pms-api"]
