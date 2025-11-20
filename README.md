# Inventory Management System - Kubernetes Deployment

A full-stack inventory management system deployed on Minikube, featuring a Next.js frontend and Go backend with PostgreSQL, Redis, and Kafka.

## Project Structure

```
.
├── inventory-manager-frontend/    # Next.js frontend application
├── InventoryManagerServer/        # Go backend application
├── k8s/                          # Kubernetes deployment manifests
├── deploy.sh                     # Automated deployment script
└── DEPLOYMENT.md                 # Detailed deployment guide
```

## Quick Start

### Prerequisites

- Git
- Docker
- Minikube
- kubectl

See [DEPLOYMENT.md](./DEPLOYMENT.md) for installation instructions.

### Deploy to Minikube

```bash
# 1. Clone the repository
git clone https://github.com/mingy6237/ece1779project_k8s.git
cd ece1779project_k8s

# 2. Start Minikube
minikube start

# 3. Run deployment script
chmod +x deploy.sh
./deploy.sh

# 4. Access the application
# Add to /etc/hosts: $(minikube ip) inventory.local
# Run: sudo minikube tunnel (in a separate terminal)
# Open: http://inventory.local
```

Default Login: `admin` / `adminadmin`

## Documentation

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Complete deployment guide with troubleshooting
- [k8s/README.md](./k8s/README.md) - Kubernetes configuration details
- [inventory-manager-frontend/README.md](./inventory-manager-frontend/README.md) - Frontend documentation
- [InventoryManagerServer/README.md](./InventoryManagerServer/README.md) - Backend documentation

## Architecture

### Components

- Frontend: Next.js application (React, TypeScript)
- Backend: Go REST API with WebSocket support
- Database: PostgreSQL (StatefulSet)
- Cache: Redis (Deployment)
- Message Queue: Kafka (StatefulSet)
- Ingress: Nginx Ingress Controller

### Services

All services are deployed in the `inventory-manager` namespace:
- `frontend` - Frontend application
- `backend` - Backend API (2 replicas)
- `postgres` - PostgreSQL database
- `redis` - Redis cache
- `kafka` - Kafka message broker

## Development

### Local Development (Docker Compose)

See [InventoryManagerServer/README.md](./InventoryManagerServer/README.md) for Docker Compose setup.

### Updating Deployment

```bash
# Rebuild and redeploy backend
eval $(minikube docker-env)
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..
kubectl rollout restart deployment/backend -n inventory-manager

# Rebuild and redeploy frontend
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ..
kubectl rollout restart deployment/frontend -n inventory-manager
```

## Cleanup

```bash
# Delete all resources
kubectl delete namespace inventory-manager
```
