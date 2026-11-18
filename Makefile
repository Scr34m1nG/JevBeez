.PHONY: up down dev-redis test-api build

up:
	docker compose up --build -d

down:
	docker compose down

dev-redis:
	docker run -d --name jevbeez-redis -p 6379:6379 -v redis_data:/data redis:7-alpine || docker start jevbeez-redis

test-api:
	docker run --rm -v "$$(pwd)/apps/api:/app" -w /app golang:1.23-alpine go test -v ./...

build:
	docker compose build
