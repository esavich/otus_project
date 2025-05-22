FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o resizer ./cmd/resizer/main.go

FROM alpine:latest

COPY --from=builder /app/resizer /app/resizer

COPY .env.dist /app/.env

WORKDIR /app

CMD ["./resizer"]