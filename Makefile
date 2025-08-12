.PHONY: build run dev test clean deps

# Variables
BINARY_NAME=samskipnad
MAIN_PATH=./cmd/server

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# Development mode with auto-reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Test the application
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f *.db

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	mkdir -p bin
	@echo "Development environment ready!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Generate admin user (for production setup)
admin:
	@echo "This would create an admin user. Not implemented yet."

# Database operations
db-reset:
	@echo "Resetting database..."
	rm -f *.db
	@echo "Database reset. It will be recreated on next run."

# Run with specific port
run-port:
	@echo "Running on port $(PORT)..."
	PORT=$(PORT) ./bin/$(BINARY_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Build and run the application"
	@echo "  dev       - Run in development mode with hot reload"
	@echo "  test      - Run tests"
	@echo "  deps      - Download dependencies"
	@echo "  clean     - Clean build artifacts"
	@echo "  setup     - Setup development environment"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  db-reset  - Reset database"
	@echo "  help      - Show this help"