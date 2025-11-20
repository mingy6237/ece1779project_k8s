# Deployment Guide - Minikube

This guide will help you deploy the inventory management system frontend and backend to Minikube.

## Quick Start

If you already have the prerequisites installed:

```bash
# 1. Clone the repository
git clone https://github.com/mingy6237/ece1779project_k8s.git
cd ece1779project_k8s

# 2. Start Minikube
minikube start

# 3. Deploy everything
chmod +x deploy.sh
./deploy.sh

# 4. Access the application
# Add to /etc/hosts: $(minikube ip) inventory.local
# Run: sudo minikube tunnel (in a separate terminal)
# Open: http://inventory.local
```

Default Login: `admin` / `adminadmin`

## Prerequisites

Before starting, ensure you have the following tools installed:

1. Git - For cloning the repository
2. Docker - For building images
3. Minikube - Kubernetes cluster
4. kubectl - Kubernetes command-line tool

### Installing Prerequisites

#### macOS

```bash
# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Docker Desktop
brew install --cask docker

# Install Minikube
brew install minikube

# Install kubectl
brew install kubectl
```

#### Linux

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

#### Windows

1. Install [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
3. Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/)

### Verify Installation

```bash
docker --version
minikube version
kubectl version --client
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/mingy6237/ece1779project_k8s.git
cd ece1779project_k8s
```

### 2. Start Minikube

```bash
minikube start
minikube status
```

Note: If you encounter issues starting Minikube, you may need to specify a driver:

```bash
# For Docker driver (recommended)
minikube start --driver=docker

# For VirtualBox driver
minikube start --driver=virtualbox

# For HyperKit driver (macOS)
minikube start --driver=hyperkit
```

### 3. Run Deployment Script

```bash
chmod +x deploy.sh
./deploy.sh
```

The deployment script will:
- Check if Minikube is running (starts it if not)
- Enable Ingress addon
- Configure Docker to use Minikube's Docker daemon
- Build backend and frontend Docker images
- Create Kubernetes namespace
- Deploy all services (PostgreSQL, Redis, Kafka, Backend, Frontend)
- Deploy Ingress controller

Expected deployment time: 5-10 minutes

## Accessing the Application

### Method 1: Via Ingress (Recommended)

Step 1: Add hosts entry

macOS/Linux:
```bash
echo "$(minikube ip) inventory.local" | sudo tee -a /etc/hosts
```

Windows:
1. Open `C:\Windows\System32\drivers\etc\hosts` as Administrator
2. Add the line: `<minikube-ip> inventory.local`

Step 2: Start minikube tunnel

```bash
sudo minikube tunnel
```

Keep this terminal window open while using the application.

Step 3: Access the application

```
http://inventory.local
```

### Method 2: Via NodePort

```bash
MINIKUBE_IP=$(minikube ip)
FRONTEND_PORT=$(kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')
echo "Access URL: http://${MINIKUBE_IP}:${FRONTEND_PORT}"
```

## Default Login Credentials

- Username: `admin`
- Password: `adminadmin`

## Architecture Overview

### Service Components

1. PostgreSQL - Primary database (StatefulSet)
2. Redis - Cache service (Deployment)
3. Kafka - Message queue (StatefulSet)
4. Backend - Backend API service (Deployment, 2 replicas, Go, Port 3000)
5. Frontend - Frontend application (Deployment, Next.js, Port 3000)

### Network Configuration

- All services are in the `inventory-manager` namespace
- Services communicate via Service names
- Frontend is exposed via Ingress at `/` path
- Backend API is exposed via Ingress at `/api` path
- WebSocket connections are supported through Ingress

## Verifying Deployment

### Check Pod Status

```bash
kubectl get pods -n inventory-manager
```

All pods should show `Running` status.

### Check Service Status

```bash
kubectl get svc -n inventory-manager
```

### View Logs

```bash
# Backend logs
kubectl logs -f deployment/backend -n inventory-manager

# Frontend logs
kubectl logs -f deployment/frontend -n inventory-manager

# Database logs
kubectl logs -f statefulset/postgres -n inventory-manager
```

### Test Health Check

```bash
curl http://inventory.local/api/health
```

## Troubleshooting

### Pods Not Starting

```bash
kubectl describe pod <pod-name> -n inventory-manager
kubectl logs <pod-name> -n inventory-manager
```

### Service Connection Issues

```bash
kubectl get endpoints -n inventory-manager
```

### Database Connection Issues

```bash
kubectl get pods -l app=postgres -n inventory-manager
kubectl logs -l app=postgres -n inventory-manager
kubectl exec -it statefulset/postgres -n inventory-manager -- psql -U postgres -d inventory_db
```

### Ingress Issues

```bash
kubectl get ingress -n inventory-manager
kubectl describe ingress inventory-manager-ingress -n inventory-manager
kubectl get pods -n ingress-nginx
sudo minikube tunnel
```

## Cleaning Up Deployment

```bash
kubectl delete namespace inventory-manager
```

## Updating Deployment

### Update Backend

```bash
eval $(minikube docker-env)
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..
kubectl rollout restart deployment/backend -n inventory-manager
```

### Update Frontend

```bash
eval $(minikube docker-env)
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ../..
kubectl rollout restart deployment/frontend -n inventory-manager
```

## Scaling Deployment

```bash
# Scale backend replicas
kubectl scale deployment backend --replicas=3 -n inventory-manager

# Scale frontend replicas
kubectl scale deployment frontend --replicas=2 -n inventory-manager
```

## Important Notes

1. Image Pull Policy: All deployments use `imagePullPolicy: Never` because images are built in Minikube's Docker daemon
2. Persistent Storage: PostgreSQL and Kafka use PersistentVolumeClaim for data storage
3. Environment Variables: Sensitive information is stored in Secrets
4. Health Checks: All services are configured with liveness and readiness probes
5. CORS: Backend is configured with CORS to allow cross-origin requests
6. WebSocket: Ingress is configured to support WebSocket connections
7. On macOS with Docker driver, `minikube tunnel` must be running for Ingress to work properly

## FAQ

**Q: How long does the initial deployment take?**

A: 5-10 minutes, as it needs to download base images, build application images, and start all services.

**Q: Can I deploy without Ingress?**

A: Yes, use Method 2 (NodePort) described in the "Accessing the Application" section.

**Q: What if the deployment script fails?**

A: Check the error message. Common issues:
- Minikube not running: Run `minikube start`
- Docker not accessible: Run `eval $(minikube docker-env)` and try again
- Port conflicts: Check if ports are already in use
- Insufficient resources: Increase Minikube's memory/CPU: `minikube start --memory=4096 --cpus=2`

**Q: How do I check if everything is working?**

A: Run:
```bash
kubectl get pods -n inventory-manager
kubectl get svc -n inventory-manager
kubectl exec -it deployment/backend -n inventory-manager -- wget -qO- http://localhost:3000/health
```

**Q: How do I completely remove the deployment?**

A: Run:
```bash
kubectl delete namespace inventory-manager
```

To also remove Docker images:
```bash
eval $(minikube docker-env)
docker rmi inventory-backend:latest inventory-frontend:latest
```

**Q: Can I change the default passwords?**

A: Yes, edit `k8s/secrets.yaml` and update the values, then reapply:
```bash
kubectl apply -f k8s/secrets.yaml
kubectl rollout restart deployment/backend -n inventory-manager
```

**Q: How do I view logs for debugging?**

A:
```bash
kubectl logs -f deployment/backend -n inventory-manager
kubectl logs -f deployment/frontend -n inventory-manager
```

**Q: The application is slow or not responding**

A: Check resource usage:
```bash
kubectl top pods -n inventory-manager
kubectl top node
```

If resources are low, increase Minikube resources:
```bash
minikube stop
minikube start --memory=4096 --cpus=4
```
