#!/bin/bash

# DDD Microservices Monitoring Script

set -e

echo "📊 DDD Microservices Monitoring Dashboard..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}================================${NC}"
}

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

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed or not in PATH"
    exit 1
fi

# Check if helm is available
if ! command -v helm &> /dev/null; then
    print_error "helm is not installed or not in PATH"
    exit 1
fi

# Function to show pod status
show_pod_status() {
    print_header "POD STATUS"
    kubectl get pods -n ddd-micro -o wide
    echo ""
}

# Function to show service status
show_service_status() {
    print_header "SERVICE STATUS"
    kubectl get services -n ddd-micro
    echo ""
}

# Function to show deployment status
show_deployment_status() {
    print_header "DEPLOYMENT STATUS"
    kubectl get deployments -n ddd-micro
    echo ""
}

# Function to show logs
show_logs() {
    local service=$1
    print_header "LOGS FOR $service"
    kubectl logs -f deployment/$service -n ddd-micro --tail=50
}

# Function to show resource usage
show_resource_usage() {
    print_header "RESOURCE USAGE"
    kubectl top pods -n ddd-micro 2>/dev/null || print_warning "Metrics server not available"
    echo ""
}

# Function to show events
show_events() {
    print_header "RECENT EVENTS"
    kubectl get events -n ddd-micro --sort-by='.lastTimestamp' | tail -10
    echo ""
}

# Function to show helm status
show_helm_status() {
    print_header "HELM STATUS"
    helm status ddd-micro -n ddd-micro
    echo ""
}

# Main monitoring loop
monitor_loop() {
    while true; do
        clear
        print_header "DDD MICROSERVICES MONITORING DASHBOARD"
        echo "Last updated: $(date)"
        echo ""
        
        show_pod_status
        show_service_status
        show_deployment_status
        show_resource_usage
        show_events
        
        echo "Press 'q' to quit, 'l' to view logs, or any other key to refresh..."
        read -t 5 -n 1 key
        
        case $key in
            q|Q)
                print_success "Monitoring stopped!"
                exit 0
                ;;
            l|L)
                echo ""
                echo "Select service to view logs:"
                echo "1) user-service"
                echo "2) product-service"
                echo "3) basket-service"
                echo "4) payment-service"
                echo "5) krakend"
                read -p "Enter choice (1-5): " choice
                
                case $choice in
                    1) show_logs "user-service" ;;
                    2) show_logs "product-service" ;;
                    3) show_logs "basket-service" ;;
                    4) show_logs "payment-service" ;;
                    5) show_logs "krakend" ;;
                    *) print_warning "Invalid choice" ;;
                esac
                ;;
        esac
    done
}

# Check if services are deployed
if ! kubectl get namespace ddd-micro &> /dev/null; then
    print_error "ddd-micro namespace not found. Please deploy the services first."
    exit 1
fi

# Start monitoring
monitor_loop
