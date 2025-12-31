.PHONY: all build frontend backend clean run dev docker-build docker-run test

# Version (can be overridden: make build VERSION=1.0.0)
VERSION ?= dev

# Build directories
FRONTEND_DIR = frontend
BACKEND_DIR = .
BUILD_DIR = bin

# Binary name
BINARY_NAME = days

# Default target
all: build

# Build both frontend and backend
build: frontend backend

# Build frontend
frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

# Build backend (with embedded frontend)
backend:
	@echo "Building backend (version: $(VERSION))..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for Linux AMD64
build-linux-amd64: frontend
	@echo "Building for Linux AMD64 (version: $(VERSION))..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules

# Run the application
run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development mode (run without building)
dev:
	@echo "Starting development server..."
	go run -ldflags "-X main.Version=dev" .

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t days-journal:$(VERSION) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -v $(PWD)/data:/app/data days-journal:$(VERSION)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install frontend dependencies only
frontend-deps:
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

# Download Go dependencies
backend-deps:
	@echo "Downloading Go dependencies..."
	go mod download

# Install all dependencies
deps: frontend-deps backend-deps

# Help
help:
	@echo "Available targets:"
	@echo "  make build              - Build both frontend and backend"
	@echo "  make build-linux-amd64  - Build for Linux AMD64"
	@echo "  make frontend           - Build frontend only"
	@echo "  make backend            - Build backend only"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make run                - Build and run the application"
	@echo "  make dev                - Run in development mode"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make test               - Run tests"
	@echo "  make deps               - Install all dependencies"
	@echo ""
	@echo "Version can be set with: make build VERSION=1.0.0"
