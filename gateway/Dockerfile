FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire workspace since we're using shared packages
COPY . .

# Build the gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./gateway/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /gateway .

EXPOSE 8080
CMD ["./gateway"]
