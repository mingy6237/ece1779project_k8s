# 部署指南 - Minikube

本指南将帮助您将库存管理系统前后端部署到 Minikube 上。

## 📋 前置要求

1. **Minikube** - 已安装并运行
2. **kubectl** - Kubernetes 命令行工具
3. **Docker** - 用于构建镜像

## 🚀 快速部署

### 1. 启动 Minikube

```bash
minikube start
```

### 2. 运行部署脚本

```bash
chmod +x deploy.sh
./deploy.sh
```

部署脚本会自动完成所有步骤，包括：
- 启用 Ingress 插件
- 构建 Docker 镜像
- 部署所有 Kubernetes 资源
- 等待服务就绪

## 🌐 访问应用

部署完成后，有两种方式访问应用：

### 方式 1: 通过 Ingress（推荐）

1. 添加 hosts 条目：

```bash
echo "$(minikube ip) inventory.local" | sudo tee -a /etc/hosts
```

2. 访问应用：

```
http://inventory.local
```

### 方式 2: 通过 NodePort

```bash
# 获取 Minikube IP
MINIKUBE_IP=$(minikube ip)

# 获取前端 NodePort
FRONTEND_PORT=$(kubectl get svc frontend -n inventory-manager -o jsonpath='{.spec.ports[0].nodePort}')

# 访问应用
echo "访问地址: http://${MINIKUBE_IP}:${FRONTEND_PORT}"
```

## 🔐 默认登录信息

- **用户名**: `admin`
- **密码**: `adminadmin`

## 📊 架构说明

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
- 前端通过 Ingress 在 `/` 路径暴露
- 后端 API 通过 Ingress 的 `/api` 路径暴露
- WebSocket 连接通过 Ingress 支持

## 🔍 验证部署

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
curl http://inventory.local/api/health

# 或直接测试后端服务
kubectl exec -it deployment/backend -n inventory-manager -- wget -qO- http://localhost:3000/health
```

## 🛠️ 故障排查

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

### Ingress 问题

```bash
# 检查 Ingress 状态
kubectl get ingress -n inventory-manager

# 查看 Ingress 详细信息
kubectl describe ingress inventory-manager-ingress -n inventory-manager

# 检查 Ingress Controller
kubectl get pods -n ingress-nginx
```

## 🧹 清理部署

```bash
# 删除所有资源
kubectl delete namespace inventory-manager

# 或删除单个资源
kubectl delete -f k8s/
```

## 🔄 更新部署

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

## 📈 扩展部署

### 增加后端副本数

```bash
kubectl scale deployment backend --replicas=3 -n inventory-manager
```

### 增加前端副本数

```bash
kubectl scale deployment frontend --replicas=2 -n inventory-manager
```

## 📝 注意事项

1. **镜像拉取策略**: 所有部署使用 `imagePullPolicy: Never`，因为镜像在 Minikube 的 Docker 守护进程中构建
2. **持久化存储**: PostgreSQL 和 Kafka 使用 PersistentVolumeClaim 存储数据
3. **环境变量**: 敏感信息（如密码）存储在 Secrets 中
4. **健康检查**: 所有服务都配置了 liveness 和 readiness 探针
5. **CORS**: 后端已配置 CORS，允许跨域请求
6. **WebSocket**: Ingress 已配置支持 WebSocket 连接

## 🔗 相关文件

- `k8s/` - Kubernetes 配置文件目录
- `deploy.sh` - 自动化部署脚本
- `k8s/README.md` - 详细的 Kubernetes 配置说明

## 💡 提示

- 如果遇到端口冲突，可以修改 Service 的 NodePort
- 如果需要修改配置，编辑 `k8s/configmap.yaml` 和 `k8s/secrets.yaml` 后重新应用
- 生产环境建议使用外部数据库和 Redis 服务
- 建议配置资源限制（requests/limits）以避免资源耗尽

