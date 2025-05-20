FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o resizer ./cmd/resizer/main.go

FROM alpine:latest

RUN pwd && ls -la
COPY --from=builder /app/resizer /app/resizer

COPY .env.dist /app/.env
RUN pwd && ls -la
WORKDIR /app


CMD ["./resizer"]