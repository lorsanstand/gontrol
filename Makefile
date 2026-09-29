.PHONY: build run-hub run-agent lint test

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