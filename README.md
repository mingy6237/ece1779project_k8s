# Distributed Inventory Management System - Backend Service

This is a backend service for learning distributed system concepts, focusing on architecture and distributed technologies.

## Project Structure

```
InventoryManagerServer/
├── backend/                 # Go backend service
│   ├── main.go             # Application entry point
│   ├── config/             # Configuration management
│   ├── models/             # Data models
│   │   ├── user.go
│   │   ├── store.go
│   │   ├── store_user.go
│   │   ├── sku.go
│   │   ├── inventory.go
│   │   └── outbox.go
│   ├── handlers/           # HTTP handlers
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── store.go
│   │   ├── sku.go
│   │   ├── inventory.go
│   │   └── report.go
│   ├── middleware/         # Middleware
│   │   └── auth.go
│   ├── services/           # Business logic services
│   │   ├── inventory_service.go
│   │   ├── outbox_service.go
│   │   └── notification_service.go
│   ├── database/           # Database connection
│   ├── cache/              # Redis cache
│   ├── kafka/              # Kafka message queue
│   │   ├── producer.go
│   │   └── consumer.go
│   ├── websocket/          # WebSocket service
│   │   ├── hub.go
│   │   └── client.go
│   ├── routes/             # Route configuration
│   ├── go.mod
│   └── Dockerfile
├── init.sql                # Database initialization script
├── compose.yaml            # Docker Compose configuration
└── README.md
```

## Tech Stack

- **Language**: Go 1.21
- **Web Framework**: Gin
- **ORM**: GORM
- **Real-time Communication**: WebSocket
- **Cache**: Redis
- **Message Queue**: Kafka
- **Database**: PostgreSQL
- **Deployment**: Docker, Docker Compose, Kubernetes

## Feature Modules

### User Management

- [ ] Login/Logout
- [ ] Update personal information/password
- [ ] Create/Delete accounts (Manager)
- [ ] Create/Delete stores (Manager)

### SKU Management

- [ ] Query SKUs (with filtering/sorting)
- [ ] Create/Delete SKUs (Manager)

### Inventory Management

- [ ] Query inventory (overall/store, with filtering/sorting)
- [ ] Add/Update/Delete store inventory
- [ ] Real-time inventory updates (WebSocket)

### System Reports

- [ ] Inventory Dashboard (Manager)
- [ ] Cluster Status Dashboard (Manager)
- [ ] Low Stock Email Notifications (Manager)

## Development Guide

### Local Development

1. Navigate to the backend directory:

```bash
cd backend
```

2. Install dependencies:

```bash
go mod download
```

3. Run the development server:

```bash
go run main.go
```

### Docker Compose Deployment

1. Start all services:

```bash
docker-compose up -d
```

2. Stop all services:

```bash
docker-compose down
```

## Architecture Overview

### Inventory Read Flow

1. Client initiates HTTP request
2. Backend queries Redis cache
3. On cache miss, query database and update cache
4. Return results to client

### Inventory Write Flow

1. Client initiates HTTP request
2. Backend starts database transaction
3. Update inventory table (using version number for consistency)
4. Insert outbox record
5. Delete Redis cache
6. Commit transaction
7. Background thread periodically sends outbox records to Kafka (maintaining consistency via version number)
8. Other servers consume Kafka messages and notify online clients via WebSocket
9. Clients receive updates via WebSocket and refresh view

## API Endpoints

The backend provides RESTful API and WebSocket interfaces for client calls. Main endpoints include:

- **Authentication**: `/api/auth/login`, `/api/auth/logout`
- **User Management**: `/api/users/*` (Manager)
- **Store Management**: `/api/stores/*` (Manager)
- **SKU Management**: `/api/skus/*` (Manager)
- **Inventory Management**: `/api/inventory/*`
- **System Reports**: `/api/reports/*` (Manager)
- **WebSocket**: `/ws` (Real-time inventory updates)

## TODO

All feature modules are marked with TODO and require implementation of specific business logic. Main tasks include:

1. Implement authentication and authorization logic
2. Implement database operations and caching logic
3. Implement Kafka message production and consumption
4. Implement WebSocket real-time communication
5. Implement background tasks (outbox processing, low stock checking, etc.)

## Notes

- All code has TODO comments indicating features to be implemented
- Database models are defined but require implementation of specific business methods
- Routes and middleware framework are set up but require implementation of specific logic
- This is a pure backend service and needs to be used with a frontend client
