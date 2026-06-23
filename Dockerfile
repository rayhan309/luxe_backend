# ─── Build Stage ──────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first for cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/server \
    ./cmd/server

# ─── Run Stage ────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS (MongoDB Atlas TLS)
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /bin/server .

# Create uploads directory
RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./server"]
