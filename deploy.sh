#!/bin/bash

set -e

echo "🚀 Starting deployment to Minikube..."

# Check if minikube is running
if ! minikube status > /dev/null 2>&1; then
    echo "❌ Minikube is not running. Starting Minikube..."
    minikube start
fi

# Enable ingress addon
echo "📦 Enabling ingress addon..."
minikube addons enable ingress

# Set docker environment to use minikube's docker daemon
echo "🐳 Configuring Docker to use Minikube's daemon..."
eval $(minikube docker-env)

# Build backend image
echo "🔨 Building backend Docker image..."
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..

# Build frontend image
echo "🔨 Building frontend Docker image..."
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ..

# Apply Kubernetes manifests
echo "📝 Applying Kubernetes manifests..."
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres-init-configmap.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/kafka.yaml

# Wait for Kafka to be ready before initializing topics
echo "⏳ Waiting for Kafka to be ready..."
kubectl wait --for=condition=ready pod -l app=kafka -n inventory-manager --timeout=80s

# Wait for Kafka topic init job
echo "⏳ Waiting for Kafka topic initialization..."
kubectl wait --for=condition=complete job/kafka-topic-init -n inventory-manager --timeout=80s || true

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n inventory-manager --timeout=80s

# Wait for Redis to be ready
echo "⏳ Waiting for Redis to be ready..."
kubectl wait --for=condition=ready pod -l app=redis -n inventory-manager --timeout=80s

# Deploy backend
echo "🚀 Deploying backend..."
kubectl apply -f k8s/backend.yaml

# Wait for backend to be ready
echo "⏳ Waiting for backend to be ready..."
kubectl wait --for=condition=ready pod -l app=backend -n inventory-manager --timeout=80s || true

# Deploy frontend
echo "🚀 Deploying frontend..."
kubectl apply -f k8s/frontend.yaml

# Deploy ingress
echo "🚀 Deploying ingress..."
kubectl apply -f k8s/ingress.yaml

# Get service URLs
echo ""
echo "✅ Deployment complete!"
echo ""
echo "📋 Service Information:"
echo "======================"
echo ""
echo "Frontend NodePort:"
kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}' && echo "" || echo "N/A"
echo ""
echo "To access the application:"
echo "1. Get minikube IP: minikube ip"
echo "2. Access frontend via NodePort: http://\$(minikube ip):\$(kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')"
echo ""
echo "Or add to /etc/hosts:"
echo "\$(minikube ip) inventory.local"
echo "Then access: http://inventory.local"
echo ""
echo "To view logs:"
echo "  Backend: kubectl logs -f deployment/backend -n inventory-manager"
echo "  Frontend: kubectl logs -f deployment/frontend -n inventory-manager"
echo ""
echo "To check pod status:"
echo "  kubectl get pods -n inventory-manager"

