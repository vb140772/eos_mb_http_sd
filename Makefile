.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean deps lint

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go binary"
	@echo "  run          - Run the service locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download Go dependencies"
	@echo "  lint         - Run linter"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  docker-clean - Clean Docker resources"

# Build the Go binary
build: deps
	@echo "Building minio-prometheus-sd..."
	go build -o bin/minio-prometheus-sd main.go
	@echo "Build complete: bin/minio-prometheus-sd"

# Run the service locally
run: deps
	@echo "Running minio-prometheus-sd..."

	go run main.go

# Run tests
test: deps
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Download Go dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run linter
lint: deps
	@echo "Running linter..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t minio-prometheus-sd .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "Services started. Access:"
	@echo "  MinIO Console: http://localhost:9001"
	@echo "  Service Discovery: http://localhost:8080"
	@echo "  Prometheus: http://localhost:9090"

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create test buckets in MinIO
create-test-buckets:
	@echo "Creating test buckets..."
	@echo "You can create test buckets manually through the MinIO console at http://localhost:9001"
	@echo "Or use the MinIO client (mc):"
	@echo "  mc alias set myminio http://localhost:9000 minioadmin minioadmin"
	@echo "  mc mb myminio/test-bucket"
	@echo "  mc mb myminio/prod-bucket"

# Show logs
logs:
	@echo "Showing service logs..."
	docker-compose logs -f minio-prometheus-sd

# Health check
health:
	@echo "Checking service health..."
	curl -f http://localhost:8080/health || echo "Service is not healthy"
