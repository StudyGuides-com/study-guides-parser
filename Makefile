.PHONY: fmt test build clean server

# Format Go code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Build the project
build:
	go build ./...

# Clean build artifacts
clean:
	go clean ./...

# Run development server (recommended for local testing)
server: build
	go run ./cmd/server 