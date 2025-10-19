#!/bin/bash

# DDD Microservices Kubernetes Deployment Script

set -e

echo "🚀 Starting DDD Microservices deployment..."

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed or not in PATH"
    exit 1
fi

# Check if helm is available
if ! command -v helm &> /dev/null; then
    echo "❌ helm is not installed or not in PATH"
    exit 1
fi

# Add Bitnami Helm repository
echo "📦 Adding Bitnami Helm repository..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create namespace if it doesn't exist
echo "🏗️  Creating namespace..."
kubectl create namespace ddd-micro --dry-run=client -o yaml | kubectl apply -f -

# Deploy the Helm chart
echo "🚀 Deploying DDD Microservices..."
helm upgrade --install ddd-micro ../helm/ddd-micro \
    --namespace ddd-micro \
    --create-namespace \
    --wait \
    --timeout=10m

echo "✅ Deployment completed successfully!"
echo ""
echo "📋 Useful commands:"
echo "  kubectl get pods -n ddd-micro"
echo "  kubectl get services -n ddd-micro"
echo "  kubectl logs -f deployment/user-service -n ddd-micro"
echo "  helm status ddd-micro -n ddd-micro"
echo ""
echo "🌐 To access the API Gateway:"
echo "  kubectl port-forward service/krakend 8080:8080 -n ddd-micro"
