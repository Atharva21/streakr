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

fmt:
	go fmt ./...

install: build
	@ sudo cp bin/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}
	@ sudo chmod +x /usr/local/bin/${BINARY_NAME}

help:
	@ echo "Available commands:"
	@ echo "  make build   - Build the binary"
	@ echo "  make run     - Build and run the binary"
	@ echo "  make clean   - Clean build artifacts"
	@ echo "  make test    - Run tests"
	@ echo "  make tidy    - Tidy go.mod and go.sum files"
	@ echo "  make fmt     - Format the code"
	@ echo "  make help    - Show this help message"

.PHONY: build run clean test tidy fmt help
