# 项目需求

## 需求概要

我现在要写一个分布式库存管理项目，主要目的是学习分布式的概念本身，所以重点放在架构和分布式上面，而不是具体功能

核心功能很简单，角色分为管理员和用户，前者可以管理用户本身，两者都可以读写库存表，库存的变化要通过websocket实时显示到前端页面，多个用户同时修改库存需要确保一致性

大致架构是多个全功能服务器（websocket连接用户，数据库读写，redis缓存，消息队列实时更新），用docker+k8s部署；数据库+redis+消息队列，同样用docker部署；目前仅考虑拓展服务节点，暂时不考虑其余数据库及中间件的扩容

## 架构

### 对象划分

前端（react）
后端（golang）
数据库（postgresql）
中间件（redis，kafka）

### 部署方式

前端+后端：无状态, 同构分布式，docker部署，k8s扩容
数据库+redis+kafka：docker-compose部署，暂时单实例,kafka不分片

### 技术栈

前端：react
后端：golang+gin+gorm+websocket
数据库：postgresql
中间件：redis+kafka

## 功能设计

### 功能用例

- 用户管理
  - 登录
  - 修改个人信息/密码
  - 创建/删除账号
  - 创建/删除店面
- sku管理
  - 查询sku
  - 创建删除sku
- 库存管理
  - 查询库存（整体/店面，过滤/排序）
  - 增删改本店面库存
- 系统报告
  - 库存dashboard
  - 集群状态dashboard
  - 低库存邮件通知

### 角色划分及权限

-管理员（manager）
    - 基础管理员manager不可删除/修改
    - 增删改用户和管理员账号
    - 增删店面
    - 增删sku
    - 查看系统报告（库存dashboard，集群状态dashboard）
    - 接收低库存邮件通知
-用户（staff）
    - 必须属于某个店面
    修改个人信息/密码
    - 查询全部店面库存（整体/店面，过滤/排序）
    - 增删改本店面库存

- 使用中间件进行鉴权，通过jwt token验证用户身份和权限

### 流程设计

- 库存读：前端发起，后端查询redis缓存，未命中则查询数据库并更新缓存，前端缓存
- 库存写：前端发起，后端发起数据库事务修改库存并新增outbox记录，删除redis缓存，后台线程定期发送outbox记录到mq（通过版本号保持一致性），其余服务器消费后通过websocket通知在线用户，前端更新缓存和视图

### 数据模型设计

- 用户表
  - id
  - email
  - username
  - password_hash
  - role: manager, staff
  - created_at
  - updated_at
- 店面表
  - id
  - name
  - created_at
  - updated_at
- 店面用户表
  - id
  - user_id: 外键，用户表id
  - store_id: 外键，店面表id
  - created_at
  - updated_at
- sku表
  - id
  - name
  - category
  - description
  - price
  - version
  - created_at
  - updated_at
- 库存表
  - id
  - sku_id: 外键，sku表id
  - store_id: 外键，店面表id
  - quantity
  - version
  - created_at
  - updated_at
- outbox记录表
  - id
  - sku_id: 外键，sku表id
  - store_id: 外键，店面表id
  - quantity
  - version
  - created_at
  - updated_at
