# =========================
# Stage 1: Build
# =========================
FROM golang:1.24-alpine AS builder

# Install CA certificates (important for HTTPS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build binary
# - CGO_ENABLED=0 → static binary
# - -ldflags="-s -w" → smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app

# =========================
# Stage 2: Runtime
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port (adjust as needed)
EXPOSE 8080

# Run binary
CMD ["/app/app"]
