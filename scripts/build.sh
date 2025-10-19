#!/bin/bash

# DDD Microservices Build Script

set -e

echo "🔨 Building DDD Microservices..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')
print_status "Using Go version: $GO_VERSION"

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf bin/
mkdir -p bin/

# Build all services
services=("user" "product" "basket" "payment")

for service in "${services[@]}"; do
    print_status "Building $service-service..."
    
    # Build the service
    if go build -o "bin/${service}-service" "./cmd/${service}/main.go"; then
        print_success "$service-service built successfully"
    else
        print_error "Failed to build $service-service"
        exit 1
    fi
done

# Build client (if needed)
if [ -d "client" ]; then
    print_status "Building client application..."
    cd client
    
    if command -v npm &> /dev/null; then
        if npm run build; then
            print_success "Client built successfully"
        else
            print_warning "Client build failed, but continuing..."
        fi
    else
        print_warning "npm not found, skipping client build"
    fi
    
    cd ..
fi

# Show build results
print_status "Build completed! Binary files:"
ls -la bin/

print_success "All services built successfully! 🎉"
