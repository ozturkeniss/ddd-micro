#!/bin/bash

# DDD Microservices Logs Script

set -e

echo "📋 DDD Microservices Logs Viewer..."

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

# Function to show logs for a specific service
show_logs() {
    local service=$1
    local lines=${2:-100}
    local follow=${3:-false}
    
    print_header "LOGS FOR $service"
    
    if [ "$follow" = "true" ]; then
        print_status "Following logs for $service (Press Ctrl+C to stop)..."
        kubectl logs -f deployment/$service -n ddd-micro --tail=$lines
    else
        print_status "Showing last $lines lines for $service..."
        kubectl logs deployment/$service -n ddd-micro --tail=$lines
    fi
}

# Function to show logs for all services
show_all_logs() {
    local lines=${1:-50}
    
    services=("user-service" "product-service" "basket-service" "payment-service" "krakend")
    
    for service in "${services[@]}"; do
        print_header "LOGS FOR $service"
        kubectl logs deployment/$service -n ddd-micro --tail=$lines
        echo ""
    done
}

# Function to show error logs only
show_error_logs() {
    local service=$1
    local lines=${2:-100}
    
    print_header "ERROR LOGS FOR $service"
    kubectl logs deployment/$service -n ddd-micro --tail=$lines | grep -i error || print_warning "No error logs found"
}

# Function to show logs with specific pattern
show_logs_with_pattern() {
    local service=$1
    local pattern=$2
    local lines=${3:-100}
    
    print_header "LOGS FOR $service (Pattern: $pattern)"
    kubectl logs deployment/$service -n ddd-micro --tail=$lines | grep -i "$pattern" || print_warning "No logs found matching pattern: $pattern"
}

# Function to show pod logs
show_pod_logs() {
    local pod_name=$1
    local lines=${2:-100}
    
    print_header "LOGS FOR POD $pod_name"
    kubectl logs $pod_name -n ddd-micro --tail=$lines
}

# Main function
main() {
    local service=""
    local lines=100
    local follow=false
    local pattern=""
    local pod_name=""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--service)
                service="$2"
                shift 2
                ;;
            -l|--lines)
                lines="$2"
                shift 2
                ;;
            -f|--follow)
                follow=true
                shift
                ;;
            -e|--error)
                show_error_logs "$service" "$lines"
                exit 0
                ;;
            -p|--pattern)
                pattern="$2"
                shift 2
                ;;
            --pod)
                pod_name="$2"
                shift 2
                ;;
            -a|--all)
                show_all_logs "$lines"
                exit 0
                ;;
            -h|--help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  -s, --service SERVICE    Service name (user-service, product-service, etc.)"
                echo "  -l, --lines LINES        Number of lines to show (default: 100)"
                echo "  -f, --follow             Follow logs in real-time"
                echo "  -e, --error              Show only error logs"
                echo "  -p, --pattern PATTERN    Filter logs by pattern"
                echo "  --pod POD_NAME           Show logs for specific pod"
                echo "  -a, --all                Show logs for all services"
                echo "  -h, --help               Show this help message"
                echo ""
                echo "Examples:"
                echo "  $0 -s user-service -f"
                echo "  $0 -s product-service -e"
                echo "  $0 -s basket-service -p 'error'"
                echo "  $0 -a -l 50"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if namespace exists
    if ! kubectl get namespace ddd-micro &> /dev/null; then
        print_error "ddd-micro namespace not found. Please deploy the services first."
        exit 1
    fi
    
    # Show logs based on options
    if [ -n "$pod_name" ]; then
        show_pod_logs "$pod_name" "$lines"
    elif [ -n "$pattern" ] && [ -n "$service" ]; then
        show_logs_with_pattern "$service" "$pattern" "$lines"
    elif [ -n "$service" ]; then
        show_logs "$service" "$lines" "$follow"
    else
        print_error "Please specify a service with -s or use -a for all services"
        exit 1
    fi
}

# Run main function with all arguments
main "$@"
