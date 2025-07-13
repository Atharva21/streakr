BINARY_NAME=streakr
.DEFAULT_GOAL := run

build:
	@ GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}

run: build
	@ ./bin/${BINARY_NAME}

clean:
	go clean
	rm -rf ./bin/*

test:
	go test ./...

tidy:
	go mod tidy

help:
	@ echo "Available commands:"
	@ echo "  make build   - Build the binary"
	@ echo "  make run     - Build and run the binary"
	@ echo "  make clean   - Clean build artifacts"
	@ echo "  make test    - Run tests"
	@ echo "  make tidy    - Tidy go.mod and go.sum files"
	@ echo "  make help    - Show this help message"

.PHONY: build run clean test tidy help
