# Файл: ./Makefile

.PHONY: run test test-all build clean deps setup

# Installing dependencies
deps:
	go mod download
	go install github.com/joho/godotenv/cmd/godotenv@latest

# Setting up a project (first run)
setup: deps
	@if [ ! -f .env ]; then \
		echo "Создаю .env из .env.example"; \
		cp .env.example .env; \
		echo "Теперь отредактируйте .env и добавьте ваш TELEGRAM_BOT_TOKEN"; \
	fi

# Launching the bot
run:
	godotenv -f .env go run cmd/bot/main.go

# Run the bot with hot reloading (requires air: go install github.com/cosmtrek/air@latest)
watch:
	godotenv -f .env air

# Tests
test:
	go test ./...

test-unit:
	go test ./internal/domain/... ./internal/application/... ./internal/infrastructure/...

test-integration:
	go test ./internal/infrastructure/telegram/...

# Build
build:
	godotenv -f .env go build -o bin/bot cmd/bot/main.go

# Clean-up
clean:
	rm -rf bin/
	go clean

# Full rebuild
rebuild: clean build

# Docker build
docker-build:
	docker build -t blood-guess-bot .

docker-run:
	docker run --env-file .env blood-guess-bot