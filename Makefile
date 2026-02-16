.PHONY: run build test clean dev tidy

# Run the server
run:
	go run ./cmd/server

# Build binary
build:
	go build -o bin/linkbio ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ data/*.db

# Download dependencies
tidy:
	go mod tidy

# Development with hot reload (requires air)
dev:
	air

# Create .env from example
env:
	cp .env.example .env
