#!/bin/bash

# DDD Microservices Docker Build Script

set -e

echo "🐳 Building DDD Microservices Docker Images..."

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

# Check if Docker is running
if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

# Get version tag (default: latest)
VERSION_TAG=${1:-latest}
REGISTRY=${2:-"ddd-micro"}

print_status "Building images with tag: $VERSION_TAG"
print_status "Registry: $REGISTRY"

# Services to build
services=("user" "product" "basket" "payment")

# Build each service
for service in "${services[@]}"; do
    print_status "Building $service-service image..."
    
    # Build the Docker image
    if docker build -f "dockerfiles/${service}.dockerfile" -t "${REGISTRY}/${service}-service:${VERSION_TAG}" .; then
        print_success "$service-service image built successfully"
    else
        print_error "Failed to build $service-service image"
        exit 1
    fi
done

# Build client image (if exists)
if [ -f "dockerfiles/client.dockerfile" ]; then
    print_status "Building client image..."
    if docker build -f "dockerfiles/client.dockerfile" -t "${REGISTRY}/client:${VERSION_TAG}" ./client; then
        print_success "Client image built successfully"
    else
        print_warning "Failed to build client image, but continuing..."
    fi
fi

# Show built images
print_status "Built images:"
docker images | grep "$REGISTRY"

print_success "All Docker images built successfully! 🎉"

# Optional: Push to registry
if [ "$3" = "--push" ]; then
    print_status "Pushing images to registry..."
    for service in "${services[@]}"; do
        print_status "Pushing ${REGISTRY}/${service}-service:${VERSION_TAG}..."
        if docker push "${REGISTRY}/${service}-service:${VERSION_TAG}"; then
            print_success "$service-service pushed successfully"
        else
            print_error "Failed to push $service-service"
        fi
    done
fi
