#!/bin/bash

# DDD Microservices Health Check Script

set -e

echo "🏥 DDD Microservices Health Check..."

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

# Configuration
BASE_URL=${1:-"http://localhost:8081"}
TIMEOUT=10

# Function to check HTTP endpoint
check_http_endpoint() {
    local url=$1
    local service_name=$2
    
    print_status "Checking $service_name at $url..."
    
    if curl -s --max-time $TIMEOUT "$url/health" > /dev/null; then
        print_success "$service_name is healthy"
        return 0
    else
        print_error "$service_name is unhealthy"
        return 1
    fi
}

# Function to check gRPC endpoint
check_grpc_endpoint() {
    local host=$1
    local port=$2
    local service_name=$3
    
    print_status "Checking $service_name gRPC at $host:$port..."
    
    # Simple TCP connection check for gRPC
    if timeout $TIMEOUT bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        print_success "$service_name gRPC is healthy"
        return 0
    else
        print_error "$service_name gRPC is unhealthy"
        return 1
    fi
}

# Function to check database connection
check_database() {
    local host=$1
    local port=$2
    local db_name=$3
    
    print_status "Checking $db_name database at $host:$port..."
    
    if timeout $TIMEOUT bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        print_success "$db_name database is accessible"
        return 0
    else
        print_error "$db_name database is not accessible"
        return 1
    fi
}

# Function to check Redis
check_redis() {
    local host=$1
    local port=$2
    
    print_status "Checking Redis at $host:$port..."
    
    if timeout $TIMEOUT bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        print_success "Redis is accessible"
        return 0
    else
        print_error "Redis is not accessible"
        return 1
    fi
}

# Function to check Kafka
check_kafka() {
    local host=$1
    local port=$2
    
    print_status "Checking Kafka at $host:$port..."
    
    if timeout $TIMEOUT bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        print_success "Kafka is accessible"
        return 0
    else
        print_error "Kafka is not accessible"
        return 1
    fi
}

# Main health check
main() {
    local failed_checks=0
    
    print_status "Starting health check for DDD Microservices..."
    print_status "Base URL: $BASE_URL"
    echo ""
    
    # Check API Gateway
    if ! check_http_endpoint "$BASE_URL" "API Gateway"; then
        ((failed_checks++))
    fi
    
    # Check microservices through API Gateway
    if ! check_http_endpoint "$BASE_URL/api/v1/users/health" "User Service"; then
        ((failed_checks++))
    fi
    
    if ! check_http_endpoint "$BASE_URL/api/v1/products/health" "Product Service"; then
        ((failed_checks++))
    fi
    
    if ! check_http_endpoint "$BASE_URL/api/v1/basket/health" "Basket Service"; then
        ((failed_checks++))
    fi
    
    if ! check_http_endpoint "$BASE_URL/api/v1/payments/health" "Payment Service"; then
        ((failed_checks++))
    fi
    
    # Check gRPC endpoints (if accessible)
    if ! check_grpc_endpoint "localhost" "9091" "Product Service"; then
        ((failed_checks++))
    fi
    
    if ! check_grpc_endpoint "localhost" "9093" "Basket Service"; then
        ((failed_checks++))
    fi
    
    # Check databases
    if ! check_database "localhost" "5432" "User DB"; then
        ((failed_checks++))
    fi
    
    if ! check_database "localhost" "5433" "Product DB"; then
        ((failed_checks++))
    fi
    
    if ! check_database "localhost" "5434" "Basket DB"; then
        ((failed_checks++))
    fi
    
    if ! check_database "localhost" "5435" "Payment DB"; then
        ((failed_checks++))
    fi
    
    # Check Redis
    if ! check_redis "localhost" "6379"; then
        ((failed_checks++))
    fi
    
    # Check Kafka
    if ! check_kafka "localhost" "9092"; then
        ((failed_checks++))
    fi
    
    echo ""
    if [ $failed_checks -eq 0 ]; then
        print_success "All health checks passed! 🎉"
        exit 0
    else
        print_error "$failed_checks health checks failed!"
        exit 1
    fi
}

# Run main function
main
