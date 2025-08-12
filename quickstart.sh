#!/bin/bash

# Quick Start Script for MinIO Prometheus Service Discovery
# This script helps you get started quickly with the service

set -e

echo "🚀 MinIO Prometheus Service Discovery - Quick Start"
echo "=================================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Check if Docker is installed (optional)
if command -v docker &> /dev/null; then
    echo "✅ Docker detected - you can use Docker Compose for easy testing"
    DOCKER_AVAILABLE=true
else
    echo "⚠️  Docker not detected - you'll need to run MinIO separately"
    DOCKER_AVAILABLE=false
fi

echo ""
echo "📦 Installing dependencies..."
go mod tidy

echo ""
echo "🔨 Building the service..."
go build -o bin/minio-prometheus-sd .

echo ""
echo "✅ Build successful! Binary created at: bin/minio-prometheus-sd"

echo ""
echo "📋 Configuration:"
echo "   The service will use these default settings:"
echo "   - MinIO Endpoint: localhost:9000"
echo "   - MinIO Access Key: minioadmin"
echo "   - MinIO Secret Key: minioadmin"
echo "   - Listen Address: :8080"
echo ""
echo "   To customize, set these environment variables:"
echo "   export MINIO_ENDPOINT='your-minio-server:9000'"
echo "   export MINIO_ACCESS_KEY='your-access-key'"
echo "   export MINIO_SECRET_KEY='your-secret-key'"
echo "   export MINIO_USE_SSL='true'"

echo "   export BUCKET_PATTERN='prod-*'           # Only production buckets"
echo "   export BUCKET_EXCLUDE_PATTERN='temp-*'   # Exclude temporary buckets"

echo ""
echo "🚀 Starting the service..."
echo "   The service will be available at: http://localhost:8080"
echo "   Press Ctrl+C to stop"
echo ""
echo "   # Option 1: Using configuration file (highest priority)"
echo "   ./bin/minio-prometheus-sd -config-file=config.yaml"
echo "   # Option 2: Using environment variables (default)"
echo "   ./bin/minio-prometheus-sd"
echo "   # Option 3: Using command line arguments"
echo "   ./bin/minio-prometheus-sd -minio-endpoint=minio:9000 -minio-access-key=mykey"
echo "   # Option 4: Mix of all methods (config file > command line > environment variables)"
echo "   ./bin/minio-prometheus-sd -config-file=config.yaml -bucket-pattern='prod-*'"
echo ""

echo ""
echo "📖 Available endpoints:"
echo "   - GET /sd?job=minio-server - Service discovery for MinIO server"
echo "   - GET /sd?job=minio-buckets - Service discovery for all MinIO buckets"
echo "   - GET /scrape_configs - All scrape configurations"
echo "   - GET /health - Health check"
echo "   - GET / - Documentation"

echo ""
echo "🔍 Testing the service..."
echo "   In another terminal, you can test:"
echo "   curl http://localhost:8080/health"
echo "   curl http://localhost:8080/scrape_configs"

if [ "$DOCKER_AVAILABLE" = true ]; then
    echo ""
    echo "🐳 Docker Quick Start:"
    echo "   To test with MinIO and Prometheus:"
    echo "   make docker-run"
    echo ""
    echo "   This will start:"
    echo "   - MinIO server at http://localhost:9001 (console)"
    echo "   - Service discovery at http://localhost:8080"
    echo "   - Prometheus at http://localhost:9090"
fi

echo ""
echo "🎯 Next steps:"
echo "   1. Make sure MinIO is running and accessible"
echo "   2. Configure Prometheus to use the service discovery endpoints"
echo "   3. Check the README.md for detailed configuration examples"
echo ""

# Start the service
echo "Starting MinIO Prometheus Service Discovery..."
echo "   Use -help flag to see all available options:"
echo "   ./bin/minio-prometheus-sd -help"
echo ""
echo "   Starting with default configuration..."
./bin/minio-prometheus-sd
