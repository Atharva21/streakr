BINARY_NAME=streakr
.DEFAULT_GOAL := run

build: generate tidy fmt
	@ GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}

run: build
	@ ./bin/${BINARY_NAME}

clean:
	@ go clean
	@ rm -rf ./bin/*
	@ rm -rf ~/.config/streakr
	@ rm -rf ./internal/store/generated

test:
	go test ./...

tidy:
	@ go mod tidy

fmt:
	@ go fmt ./...

install: build
	@ go install .

generate:
	@ sqlc generate

help:
	@ echo "Available commands:"
	@ echo "  make build   - Build the binary"
	@ echo "  make run     - Build and run the binary"
	@ echo "  make clean   - Clean build artifacts"
	@ echo "  make test    - Run tests"
	@ echo "  make tidy    - Tidy go.mod and go.sum files"
	@ echo "  make fmt     - Format the code"
	@ echo "  make help    - Show this help message"

.PHONY: build run clean test tidy fmt install generate help
