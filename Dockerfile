FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wallets-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wallets-api .

EXPOSE 8080

CMD ["./wallets-api"]