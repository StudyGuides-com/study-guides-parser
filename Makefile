.PHONY: fmt test build clean

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