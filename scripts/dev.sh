#!/bin/bash

# DDD Microservices Development Script

set -e

echo "🚀 Starting DDD Microservices Development Environment..."

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

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed or not in PATH"
    exit 1
fi

# Function to cleanup on exit
cleanup() {
    print_status "Stopping development environment..."
    docker-compose down
    print_success "Development environment stopped!"
}

# Set trap to cleanup on script exit
trap cleanup EXIT

# Start infrastructure services
print_status "Starting infrastructure services (PostgreSQL, Redis, Kafka)..."
docker-compose up -d user-db product-db basket-db payment-db redis zookeeper kafka

# Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Check if services are healthy
print_status "Checking service health..."

# Check PostgreSQL services
for db in user-db product-db basket-db payment-db; do
    if docker-compose exec $db pg_isready -U postgres &> /dev/null; then
        print_success "$db is ready"
    else
        print_warning "$db is not ready yet"
    fi
done

# Check Redis
if docker-compose exec redis redis-cli ping &> /dev/null; then
    print_success "Redis is ready"
else
    print_warning "Redis is not ready yet"
fi

# Check Kafka
if docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list &> /dev/null; then
    print_success "Kafka is ready"
else
    print_warning "Kafka is not ready yet"
fi

# Start microservices
print_status "Starting microservices..."
docker-compose up -d user-service product-service basket-service payment-service krakend

# Show running services
print_status "Development environment is running! Services:"
docker-compose ps

print_success "Development environment started successfully! 🎉"
print_status "API Gateway: http://localhost:8081"
print_status "User Service: http://localhost:8080"
print_status "Product Service: http://localhost:8082"
print_status "Basket Service: http://localhost:8083"
print_status "Payment Service: http://localhost:8084"

print_status "Press Ctrl+C to stop all services"

# Keep script running
wait
