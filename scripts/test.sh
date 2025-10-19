#!/bin/bash

# DDD Microservices Test Script

set -e

echo "🧪 Running DDD Microservices Tests..."

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

# Run Go tests
print_status "Running Go tests..."
if go test ./... -v -race -coverprofile=coverage.out; then
    print_success "All Go tests passed!"
else
    print_error "Some Go tests failed!"
    exit 1
fi

# Generate coverage report
if [ -f "coverage.out" ]; then
    print_status "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
fi

# Run client tests (if exists)
if [ -d "client" ]; then
    print_status "Running client tests..."
    cd client
    
    if command -v npm &> /dev/null; then
        if npm test; then
            print_success "Client tests passed!"
        else
            print_warning "Client tests failed, but continuing..."
        fi
    else
        print_warning "npm not found, skipping client tests"
    fi
    
    cd ..
fi

# Run integration tests (if exists)
if [ -d "tests/integration" ]; then
    print_status "Running integration tests..."
    if go test ./tests/integration/... -v; then
        print_success "Integration tests passed!"
    else
        print_warning "Integration tests failed, but continuing..."
    fi
fi

print_success "Test suite completed! 🎉"
