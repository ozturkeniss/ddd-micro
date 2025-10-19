#!/bin/bash

# DDD Microservices Kubernetes Cleanup Script

set -e

echo "🧹 Starting DDD Microservices cleanup..."

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

# Uninstall Helm chart
echo "🗑️  Uninstalling DDD Microservices..."
helm uninstall ddd-micro -n ddd-micro || echo "Chart not found, continuing..."

# Delete namespace (this will delete all resources in the namespace)
echo "🗑️  Deleting namespace..."
kubectl delete namespace ddd-micro || echo "Namespace not found, continuing..."

echo "✅ Cleanup completed successfully!"
echo ""
echo "ℹ️  Note: Persistent volumes are not deleted by default."
echo "   To delete them manually, run:"
echo "   kubectl get pv | grep ddd-micro"
echo "   kubectl delete pv <volume-name>"
