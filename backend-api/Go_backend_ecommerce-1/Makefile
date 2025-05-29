# Makefile

.PHONY: build run test docker-build docker-run migrate

build:
	go build -o bin/server ./cmd/server/main.go

run: build
	./bin/server

test:
	go test ./...

docker-build:
	docker build -t go-backend-ecommerce ./backend-api

docker-run:
	docker run -p 8080:8080 go-backend-ecommerce

migrate:
	./scripts/migrate.sh