FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка бинарника
#RUN go build -o video-manage ./cmd/video-manage
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet ./cmd/wallet
FROM alpine:3.23

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates

# Копируем из builder-а бинарь, конфиги и миграции
COPY --from=builder /app/wallet .
COPY config.env .
COPY internal/migrations /app/internal/migrations

# Открываем порт для gRPC (опционально, для понимания)
EXPOSE 9000

CMD ["./wallet"]