.PHONY: run build tidy docker-up docker-down

run:
	go run ./cmd/api

build:
	go build -o bin/server ./cmd/api

tidy:
	go mod tidy

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
