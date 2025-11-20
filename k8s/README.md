# Kubernetes Deployment Guide

本指南将帮助您将库存管理系统部署到 Minikube。

## 前置要求

1. **Minikube** - 已安装并运行
2. **kubectl** - Kubernetes 命令行工具
3. **Docker** - 用于构建镜像

## 快速开始

### 1. 启动 Minikube

```bash
minikube start
```

### 2. 运行部署脚本

```bash
chmod +x deploy.sh
./deploy.sh
```

部署脚本会自动：
- 启用 Ingress 插件
- 构建 Docker 镜像
- 部署所有 Kubernetes 资源
- 等待服务就绪

## 手动部署步骤

如果您想手动部署，可以按照以下步骤：

### 1. 配置 Docker 环境

```bash
eval $(minikube docker-env)
```

### 2. 构建 Docker 镜像

```bash
# 构建后端镜像
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..

# 构建前端镜像
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ..
```

### 3. 启用 Ingress

```bash
minikube addons enable ingress
```

### 4. 部署 Kubernetes 资源

按顺序部署：

```bash
# 1. 创建命名空间
kubectl apply -f k8s/namespace.yaml

# 2. 创建配置和密钥
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres-init-configmap.yaml

# 3. 部署数据库和缓存
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml

# 4. 部署 Kafka
kubectl apply -f k8s/kafka.yaml

# 5. 等待 Kafka 就绪
kubectl wait --for=condition=ready pod -l app=kafka -n inventory-manager --timeout=300s

# 6. 等待 Kafka topic 初始化完成
kubectl wait --for=condition=complete job/kafka-topic-init -n inventory-manager --timeout=300s

# 7. 等待 PostgreSQL 就绪
kubectl wait --for=condition=ready pod -l app=postgres -n inventory-manager --timeout=300s

# 8. 部署后端
kubectl apply -f k8s/backend.yaml

# 9. 部署前端
kubectl apply -f k8s/frontend.yaml

# 10. 部署 Ingress
kubectl apply -f k8s/ingress.yaml
```

## 访问应用

### 方式 1: 通过 NodePort

```bash
# 获取 Minikube IP
MINIKUBE_IP=$(minikube ip)

# 获取前端 NodePort
FRONTEND_PORT=$(kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')

# 访问应用
echo "访问地址: http://${MINIKUBE_IP}:${FRONTEND_PORT}"
```

### 方式 2: 通过 Ingress

1. 添加 hosts 条目：

```bash
echo "$(minikube ip) inventory.local" | sudo tee -a /etc/hosts
```

2. 访问应用：

```
http://inventory.local
```

## 验证部署

### 检查 Pod 状态

```bash
kubectl get pods -n inventory-manager
```

所有 Pod 应该显示 `Running` 状态。

### 检查服务状态

```bash
kubectl get svc -n inventory-manager
```

### 查看日志

```bash
# 后端日志
kubectl logs -f deployment/backend -n inventory-manager

# 前端日志
kubectl logs -f deployment/frontend -n inventory-manager

# 数据库日志
kubectl logs -f statefulset/postgres -n inventory-manager
```

### 测试健康检查

```bash
# 测试后端健康检查
kubectl exec -it deployment/backend -n inventory-manager -- wget -qO- http://localhost:3000/health

# 或从外部测试
curl http://$(minikube ip):$(kubectl get svc backend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')/health
```

## 默认登录信息

- **用户名**: `admin`
- **密码**: `adminadmin`

## 架构说明

### 服务组件

1. **PostgreSQL** - 主数据库
   - StatefulSet 部署
   - 持久化存储
   - 自动初始化数据库结构

2. **Redis** - 缓存服务
   - Deployment 部署
   - 用于缓存数据

3. **Kafka** - 消息队列
   - StatefulSet 部署
   - 用于库存更新事件

4. **Backend** - 后端 API 服务
   - Deployment 部署（2 个副本）
   - Go 应用
   - 端口: 3000

5. **Frontend** - 前端应用
   - Deployment 部署
   - Next.js 应用
   - 端口: 3000

### 网络配置

- 所有服务在 `inventory-manager` 命名空间中
- 服务间通过 Service 名称通信
- 前端通过 Ingress 暴露给外部访问
- 后端 API 通过 Ingress 的 `/api` 路径暴露

## 故障排查

### Pod 无法启动

```bash
# 查看 Pod 详细信息
kubectl describe pod <pod-name> -n inventory-manager

# 查看 Pod 日志
kubectl logs <pod-name> -n inventory-manager
```

### 服务无法连接

```bash
# 检查服务端点
kubectl get endpoints -n inventory-manager

# 测试服务连接
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# 在容器内测试连接
# wget -qO- http://backend:3000/health
```

### 数据库连接问题

```bash
# 检查数据库 Pod
kubectl get pods -l app=postgres -n inventory-manager

# 查看数据库日志
kubectl logs -l app=postgres -n inventory-manager

# 进入数据库 Pod
kubectl exec -it statefulset/postgres -n inventory-manager -- psql -U postgres -d inventory_db
```

## 清理部署

```bash
# 删除所有资源
kubectl delete namespace inventory-manager

# 或删除单个资源
kubectl delete -f k8s/
```

## 更新部署

### 更新后端

```bash
# 重新构建镜像
cd InventoryManagerServer/backend
docker build -t inventory-backend:latest .
cd ../..

# 重启部署
kubectl rollout restart deployment/backend -n inventory-manager
```

### 更新前端

```bash
# 重新构建镜像
cd inventory-manager-frontend
docker build -t inventory-frontend:latest .
cd ..

# 重启部署
kubectl rollout restart deployment/frontend -n inventory-manager
```

## 扩展部署

### 增加后端副本数

```bash
kubectl scale deployment backend --replicas=3 -n inventory-manager
```

### 增加前端副本数

```bash
kubectl scale deployment frontend --replicas=2 -n inventory-manager
```

## 注意事项

1. **镜像拉取策略**: 所有部署使用 `imagePullPolicy: Never`，因为镜像在 Minikube 的 Docker 守护进程中构建
2. **持久化存储**: PostgreSQL 和 Kafka 使用 PersistentVolumeClaim 存储数据
3. **环境变量**: 敏感信息（如密码）存储在 Secrets 中
4. **健康检查**: 所有服务都配置了 liveness 和 readiness 探针

## 生产环境建议

在生产环境中，您应该：

1. 使用外部镜像仓库（如 Docker Hub、GCR 等）
2. 更改 `imagePullPolicy` 为 `Always` 或 `IfNotPresent`
3. 使用更强的密码和 JWT Secret
4. 配置 TLS/SSL 证书
5. 设置资源限制（requests/limits）
6. 配置 Horizontal Pod Autoscaler
7. 使用外部数据库和 Redis（如云服务）
8. 配置监控和日志收集

