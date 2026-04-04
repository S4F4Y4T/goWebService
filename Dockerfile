# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations ./db/migrations
# Create a dummy .env if doesn't exist, we'll rely on environment variables
RUN touch .env

EXPOSE 8080
CMD ["./server"]
