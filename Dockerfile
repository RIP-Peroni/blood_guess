# Файл: ./Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /blood-guess-bot cmd/bot/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /blood-guess-bot .

# Копируем .env.example (опционально)
COPY .env.example .

# Запускаем приложение
CMD ["./blood-guess-bot"]