BINARY_NAME=streakr
VERSION?=dev
LDFLAGS=-ldflags "-X github.com/Atharva21/streakr/cmd.Version=$(VERSION)"
.DEFAULT_GOAL := install

bootstrap:
	@ echo "Installing development tools..."
	@ go install github.com/spf13/cobra-cli@latest
	@ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@ go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@ echo "✓ cobra-cli installed"
	@ echo "✓ sqlc installed"
	@ echo "✓ migrate installed"
	@ echo ""
	@ echo "Development tools installed successfully!"
	@ echo "Make sure $$GOPATH/bin or $$HOME/go/bin is in your PATH"
	@ echo "Add this to your shell rc file if not already present:"
	@ echo "  export PATH=\$$PATH:\$$HOME/go/bin"

build: generate tidy fmt
	@ GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/linux/${BINARY_NAME}

run: build
	@ ./bin/${BINARY_NAME}

clean:
	@ go clean
	@ rm -rf ./bin/*
	@ rm -rf ~/.config/streakr
	@ rm -rf ./internal/store/generated

test:
	@ echo "Running tests..."
	@ go test ./... -v

test-coverage:
	@ echo "Running tests with coverage..."
	@ go test ./... -cover

test-verbose:
	@ echo "Running tests with verbose output..."
	@ go test ./... -v

test-service:
	@ echo "Running service layer tests..."
	@ go test ./internal/service/... -v

tidy:
	@ go mod tidy

fmt:
	@ go fmt ./... > /dev/null

install: build
	@ go install $(LDFLAGS) .

generate:
	@ sqlc generate

clean-install: clean install

help:
	@ echo "Available commands:"
	@ echo "  make bootstrap         - Install development tools (cobra-cli, sqlc, migrate)"
	@ echo "  make build             - Build the binary"
	@ echo "  make run               - Build and run the binary"
	@ echo "  make clean             - Clean build artifacts"
	@ echo "  make test              - Run all tests with verbose output"
	@ echo "  make test-coverage     - Run tests with coverage report"
	@ echo "  make test-verbose      - Run tests with verbose output (alias for test)"
	@ echo "  make test-service      - Run service layer tests only"
	@ echo "  make tidy              - Tidy go.mod and go.sum files"
	@ echo "  make fmt               - Format the code"
	@ echo "  make install           - Install binary on same machine"
	@ echo "  make generate          - Generate sql code"
	@ echo "  make clean-install     - Clean any temporary files and do a fresh install"
	@ echo "  make help              - Show this help message"

.PHONY: bootstrap build run clean test test-coverage test-verbose test-service tidy fmt install generate clean-install help
