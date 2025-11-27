# Deployment Guide - Minikube

This guide will help you deploy the inventory management system frontend and backend to Minikube.

## 📋 Prerequisites

1. **Minikube** - Installed and running
2. **kubectl** - Kubernetes command-line tool
3. **Docker** - For building images

## 🚀 Quick Deployment

### 1. Start Minikube

```bash
minikube start
```

### 2. Run Deployment Script

```bash
chmod +x deploy.sh
./deploy.sh
```

The deployment script will automatically complete all steps, including:
- Enable Ingress addon
- Build Docker images
- Deploy all Kubernetes resources
- Wait for services to be ready

## 🌐 Accessing the Application

After deployment, there are two ways to access the application:

### Method 1: Via Ingress (Recommended)

1. Add hosts entry:

```bash
echo '127.0.0.1 inventory.local' | sudo tee -a /etc/hosts
```

2. Start minikube tunnel (required for Ingress on macOS with Docker driver):

```bash
sudo minikube tunnel
```

Keep this terminal window open.

3. Access the application:

```
http://inventory.local
```

### Method 2: Via NodePort

```bash
# Get Minikube IP
MINIKUBE_IP=$(minikube ip)

# Get frontend NodePort
FRONTEND_PORT=$(kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')

# Access the application
echo "Access URL: http://${MINIKUBE_IP}:${FRONTEND_PORT}"
```

## 🔐 Default Login Credentials

- **Username**: `admin`
- **Password**: `adminadmin`

## 📊 Architecture Overview

### Service Components

1. **PostgreSQL** - Primary database
   - Deployed as StatefulSet
   - Persistent storage
   - Automatic database schema initialization

2. **Redis** - Cache service
   - Deployed as Deployment
   - Used for data caching

3. **Kafka** - Message queue
   - Deployed as StatefulSet
   - Used for inventory update events

4. **Backend** - Backend API service
   - Deployed as Deployment (2 replicas)
   - Go application
   - Port: 3000

5. **Frontend** - Frontend application
   - Deployed as Deployment
   - Next.js application
   - Port: 3000

### Network Configuration

- All services are in the `inventory-manager` namespace
- Services communicate via Service names
- Frontend is exposed via Ingress at `/` path
- Backend API is exposed via Ingress at `/api` path
- WebSocket connections are supported through Ingress

## 🔍 Verifying Deployment

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
# Test backend health check
curl http://inventory.local/api/health

# Or test backend service directly
kubectl exec -it deployment/backend -n inventory-manager -- wget -qO- http://localhost:3000/health
```

## 🛠️ Troubleshooting

### Pods Not Starting

```bash
# View pod details
kubectl describe pod <pod-name> -n inventory-manager

# View pod logs
kubectl logs <pod-name> -n inventory-manager
```

### Service Connection Issues

```bash
# Check service endpoints
kubectl get endpoints -n inventory-manager

# Test service connection
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# Inside the container, test connection
# wget -qO- http://backend:3000/health
```

### Database Connection Issues

```bash
# Check database pod
kubectl get pods -l app=postgres -n inventory-manager

# View database logs
kubectl logs -l app=postgres -n inventory-manager

# Access database pod
kubectl exec -it statefulset/postgres -n inventory-manager -- psql -U postgres -d inventory_db
```

### Ingress Issues

```bash
# Check Ingress status
kubectl get ingress -n inventory-manager

# View Ingress details
kubectl describe ingress inventory-manager-ingress -n inventory-manager

# Check Ingress Controller
kubectl get pods -n ingress-nginx

# Ensure minikube tunnel is running (for macOS Docker driver)
sudo minikube tunnel
```

## 🧹 Cleaning Up Deployment

```bash
# Delete all resources
kubectl delete namespace inventory-manager

# Or delete individual resources
kubectl delete -f k8s/
```

## 🔄 Updating Deployment

### Update Backend

```bash
# Rebuild image
eval $(minikube docker-env)
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..

# Restart deployment
kubectl rollout restart deployment/backend -n inventory-manager
```

### Update Frontend

```bash
# Rebuild image
eval $(minikube docker-env)
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ../..

# Restart deployment
kubectl rollout restart deployment/frontend -n inventory-manager
```

## 📈 Scaling Deployment

### Scale Backend Replicas

```bash
kubectl scale deployment backend --replicas=3 -n inventory-manager
```

### Scale Frontend Replicas

```bash
kubectl scale deployment frontend --replicas=2 -n inventory-manager
```

## 📝 Important Notes

1. **Image Pull Policy**: All deployments use `imagePullPolicy: Never` because images are built in Minikube's Docker daemon
2. **Persistent Storage**: PostgreSQL and Kafka use PersistentVolumeClaim for data storage
3. **Environment Variables**: Sensitive information (such as passwords) is stored in Secrets
4. **Health Checks**: All services are configured with liveness and readiness probes
5. **CORS**: Backend is configured with CORS to allow cross-origin requests
6. **WebSocket**: Ingress is configured to support WebSocket connections

## 🔗 Related Files

- `k8s/` - Kubernetes configuration files directory
- `deploy.sh` - Automated deployment script
- `k8s/README.md` - Detailed Kubernetes configuration documentation

## 💡 Tips

- If you encounter port conflicts, you can modify the Service's NodePort
- To modify configuration, edit `k8s/configmap.yaml` and `k8s/secrets.yaml` then reapply
- For production environments, consider using external database and Redis services
- It's recommended to configure resource limits (requests/limits) to avoid resource exhaustion
- On macOS with Docker driver, `minikube tunnel` must be running for Ingress to work properly
