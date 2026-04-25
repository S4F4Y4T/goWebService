FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service-bin ./user-service/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /user-service-bin ./user-service
COPY ./user-service/db/migrations ./db/migrations

EXPOSE 8082
CMD ["./user-service"]
