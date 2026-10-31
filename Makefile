.PHONY: build run-hub run-agent lint test create-migrate, gen, cover

NAME ?= default_migration

build:
	go build -o bin/hub cmd/hub/main.go
	go build -o bin/agent cmd/agent/main.go

run-hub:
	go run cmd/hub/main.go

run-agent:
	go run cmd/agent/main.go

lint:
	golangci-lint run

test:
	go test -v -race ./...

cover:
	go test -short -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

create-migrate:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

gen:
	sqlc generate
	go generate ./...