# Stage 1: Build the Go services
FROM golang:1.24.3-alpine AS builder

# Create working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build api-service
RUN go build -o bin/api-service ./cmd/services/api-service

# Build channel-service
RUN go build -o bin/channel-service ./cmd/services/channel-service

# Build news-service
RUN go build -o bin/news-service ./cmd/services/news-service

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Create app directory
WORKDIR /app

# Copy built binaries
COPY --from=builder /app/bin/api-service .
COPY --from=builder /app/bin/channel-service .
COPY --from=builder /app/bin/news-service .

# Default command (override with docker run CMD)
#CMD ["./api-service"]
