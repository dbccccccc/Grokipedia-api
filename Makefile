.PHONY: help run build test clean docker-build docker-run install

# Default target
help:
	@echo "Grokipedia API - Available commands:"
	@echo "  make install      - Install dependencies"
	@echo "  make run          - Run the server"
	@echo "  make build        - Build the binary"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker containers"

# Install dependencies
install:
	go mod download
	go mod tidy

# Run the server
run:
	go run main.go

# Build the binary
build:
	go build -o grokipedia-api main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o grokipedia-api-linux main.go
	GOOS=windows GOARCH=amd64 go build -o grokipedia-api.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o grokipedia-api-mac main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f grokipedia-api grokipedia-api-* *.exe

# Docker build
docker-build:
	docker build -t grokipedia-api .

# Docker run with compose
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

