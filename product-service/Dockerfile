FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /product-service-bin ./product-service/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /product-service-bin ./product-service
COPY ./product-service/db/migrations ./db/migrations

EXPOSE 8083
CMD ["./product-service"]
