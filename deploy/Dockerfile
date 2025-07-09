# этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установим зависимости и подготовим окружение
RUN apk add --no-cache git gcc musl-dev

COPY ../go.mod go.sum ./
RUN go mod download

COPY .. .

# Встроим миграции в бинарь, если используешь embed
RUN CGO_ENABLED=0 go build -o gecko-eats ./cmd/main.go

# финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/gecko-eats .

# (опционально) открываем порт, если нужно
# EXPOSE 8080

ENTRYPOINT ["./gecko-eats"]
