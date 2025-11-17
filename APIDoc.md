# API Documentation - Inventory Management System

## Backend Setup Guide

### Quick Start

**1. Configuration File (`.env`)**

copy and rename `example.env` to `.env`

**2. Default Ports**

| Service | Port | URL |
|---------|------|-----|
| API Server | 3000 | `http://localhost:3000` |
| WebSocket | 3000 | `ws://localhost:3000/ws` |
| PostgreSQL | 15432 | `localhost:15432` |
| Redis | 16379 | `localhost:16379` |
| Kafka | 19092 | `localhost:19092` |
| Kafka UI | 9094 | `http://localhost:9094` |

**3. Docker Commands**

```bash
docker compose up --build --wait
docker compose down -v
```

**4. Default Admin Account**
(cannot be deleted)

- Username: `admin`
- Password: `adminadmin`
- Role: `manager`

**5. Health Check**

```bash
curl http://localhost:3000/health
# Expected: {"status":"ok"}
```

## Authentication

All endpoints (except `/api/auth/login`) require JWT authentication.

**Header:**

```http
Authorization: Bearer <token>
```

**Token Format:** JWT with payload:

```json
{
  "userID": "uuid",
  "userName": "username",
  "userRole": "manager" | "staff",
  "iat": 1234567890,
  "exp": 1234567890
}
```

---

## All-roles User Endpoints

### POST `/api/auth/login`

Login with username and password.

**Request Body:**

```json
{
  "username": "admin",
  "password": "adminadmin",
  "rememberMe": true
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "admin",
    "email": "admin@admin.com",
    "role": "manager",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

**Errors:**

- 400: Invalid request format
- 401: Invalid credentials

---

### GET `/api/profile`

Get current authenticated user information.

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "admin",
  "email": "admin@admin.com",
  "role": "manager",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**

- 401: Unauthorized
- 404: User not found

---

### PUT `/api/profile/password`

Change the authenticated user's password.

**Request Body:**

```json
{
  "old_password": "currentPass123",
  "new_password": "newStrongPass456"
}
```

**Validation:**

- `new_password`: minimum 8 characters

**Response (200 OK):**

```json
{
  "message": "Password changed successfully"
}
```

**Errors:**

- 400: Invalid old password or validation error
- 401: Unauthorized
- 404: User not found

---

## User Management (Manager Only)

### GET `/api/manager/users`

List all users with pagination.

**Query Parameters:**

- `page` (number, default: 1): Page number
- `limit` (number, default: 20, max: 100): Users per page

**Response (200 OK):**

```json
{
  "users": [
    {
      "id": "a21c1470-9b5d-4f9d-984f-9bd3c8ebc936",
      "username": "employee001",
      "email": "employee@example.com",
      "role": "staff",
      "created_at": "2025-01-03T10:00:00Z",
      "updated_at": "2025-01-03T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20,
  "totalPages": 1
}
```

**Errors:**

- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### POST `/api/manager/users`

Create a new user.

**Request Body:**

```json
{
  "username": "employee001",
  "email": "employee@example.com",
  "password": "securePass123",
  "role": "staff"
}
```

**Validation:**

- `username`: required
- `email`: required, valid email format
- `password`: required
- `role`: required, must be "manager" or "staff"

**Response (201 Created):**

```json
{
  "id": "a21c1470-9b5d-4f9d-984f-9bd3c8ebc936",
  "username": "employee001",
  "email": "employee@example.com",
  "role": "staff",
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

**Errors:**

- 400: Validation error or duplicate username/email
- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### PUT `/api/manager/users`

Update a user's information.

**Request Body:**

```json
{
  "target_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "newUsername",
  "email": "newemail@example.com",
  "role": "staff"
}
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "newUsername",
  "email": "newemail@example.com",
  "role": "staff",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-25T09:00:00Z"
}
```

**Errors:**

- 400: Validation error or duplicate username/email
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: User not found

---

### DELETE `/api/manager/users/:id`

Delete a user.

**Path Parameters:**

- `id` (UUID): User ID to delete

**Response (200 OK):**

```json
{
  "message": "User deleted successfully"
}
```

**Errors:**

- 400: Invalid user ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: User not found

---

## Store Management (Manager Only)

### GET `/api/manager/stores`

List all stores.

**Response (200 OK):**

```json
{
  "items": [
    {
      "id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
      "name": "Main Street Store",
      "address": "123 Main St, New York, NY",
      "created_at": "2025-01-01T09:00:00Z",
      "updated_at": "2025-01-20T09:00:00Z"
    }
  ]
}
```

**Errors:**

- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### POST `/api/manager/stores`

Create a new store.

**Request Body:**

```json
{
  "name": "New Store",
  "address": "456 New Road, Boston, MA"
}
```

**Validation:**

- `name`: required, max 100 characters
- `address`: required, max 255 characters

**Response (201 Created):**

```json
{
  "id": "4a9e7d2d-97de-4287-a2a0-747cad1a7350",
  "name": "New Store",
  "address": "456 New Road, Boston, MA",
  "created_at": "2025-02-01T14:12:00Z",
  "updated_at": "2025-02-01T14:12:00Z"
}
```

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### DELETE `/api/manager/stores/:id`

Delete a store.

**Path Parameters:**

- `id` (UUID): Store ID to delete

**Response (200 OK):**

```json
{
  "message": "Store deleted successfully"
}
```

**Errors:**

- 400: Invalid store ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Store not found
- 409: Cannot delete store with associated users or inventory

---

## Store Staff Management (Manager Only)

### GET `/api/manager/stores/staff`

List all staff members for a specific store.

**Query Parameters:**

- `store_id` (UUID, required): The ID of the store

**Response (200 OK):**

```json
{
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "staff": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "staff",
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

**Errors:**

- 400: Missing or invalid store_id
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Store not found

---

### POST `/api/manager/stores/staff`

Add a staff member to a store.

**Request Body:**

```json
{
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Response (201 Created):**

```json
{
  "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "staff",
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  },
  "store": {
    "id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
    "name": "Main Street Store",
    "address": "123 Main St, New York, NY",
    "created_at": "2025-01-01T09:00:00Z",
    "updated_at": "2025-01-20T09:00:00Z"
  },
  "created_at": "2025-02-01T14:30:00Z",
  "updated_at": "2025-02-01T14:30:00Z"
}
```

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Store or User not found
- 409: User already associated with this store

---

### DELETE `/api/manager/stores/staff/:id`

Remove a staff member from a store.

**Path Parameters:**

- `id` (UUID): The ID of the store-user association

**Response (200 OK):**

```json
{
  "message": "Staff removed from store successfully"
}
```

**Errors:**

- 400: Invalid ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Association not found

---

## SKU Management (Manager Only)

### GET `/api/manager/skus`

List SKUs with pagination, search, and filters.

**Query Parameters:**

- `page` (number, default: 1): Page number
- `page_size` (number, default: 20, max: 100): Items per page
- `search` (string): Search by name, category, or description
- `category` (string): Filter by category
- `sort_by` (string, default: "created_at"): Sort field (name, category, created_at, updated_at, price)
- `order` (string, default: "desc"): Sort order (asc, desc)

**Response (200 OK):**

```json
{
  "items": [
    {
      "id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
      "name": "Wireless Mouse",
      "category": "electronics",
      "description": "Ergonomic wireless mouse",
      "price": 29.99,
      "version": 1,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "total": 74,
  "page": 1,
  "page_size": 20,
  "total_pages": 4
}
```

**Errors:**

- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### GET `/api/manager/skus/categories`

List all available SKU categories.

**Response (200 OK):**

```json
{
  "categories": ["electronics", "furniture", "office-supplies"]
}
```

**Errors:**

- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### GET `/api/manager/skus/:id`

Retrieve SKU details by ID.

**Path Parameters:**

- `id` (UUID): SKU ID

**Response (200 OK):**

```json
{
  "id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "name": "Wireless Mouse",
  "category": "electronics",
  "description": "Ergonomic wireless mouse",
  "price": 29.99,
  "version": 1,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**

- 400: Invalid SKU ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: SKU not found

---

### POST `/api/manager/skus`

Create a new SKU.

**Request Body:**

```json
{
  "name": "Ergonomic Chair",
  "category": "furniture",
  "description": "Comfortable office chair",
  "price": 119.99
}
```

**Validation:**

- `name`: required
- `price`: required, minimum 0

**Response (201 Created):**

```json
{
  "id": "99a095f0-2baa-4f88-81fb-06b89ac8d88b",
  "name": "Ergonomic Chair",
  "category": "furniture",
  "description": "Comfortable office chair",
  "price": 119.99,
  "version": 1,
  "created_at": "2025-02-14T10:15:20Z",
  "updated_at": "2025-02-14T10:15:20Z"
}
```

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)

---

### PUT `/api/manager/skus/:id`

Update a SKU by ID. Only provided fields are updated.

**Path Parameters:**

- `id` (UUID): SKU ID

**Request Body:** (all fields optional)

```json
{
  "name": "Wireless Mouse - Updated",
  "category": "electronics",
  "description": "Updated description",
  "price": 24.99
}
```

**Response (200 OK):**

```json
{
  "id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "name": "Wireless Mouse - Updated",
  "category": "electronics",
  "description": "Updated description",
  "price": 24.99,
  "version": 2,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-02-20T11:23:00Z"
}
```

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: SKU not found

---

### DELETE `/api/manager/skus/:id`

Delete a SKU by ID. Can only delete if SKU has no active inventory.

**Path Parameters:**

- `id` (UUID): SKU ID

**Response (200 OK):**

```json
{
  "message": "SKU deleted successfully"
}
```

**Errors:**

- 400: Invalid SKU ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: SKU not found
- 409: Cannot delete SKU with active inventory

---

## Inventory Management

### GET `/api/inventory`

List inventory items with pagination and filters.

**Access:** All authenticated users

- **Managers**: Can view all stores
- **Staff**: Can only view their assigned stores

**Query Parameters:**

- `store_id` (UUID, optional): Filter by store ID
- `sku_id` (UUID, optional): Filter by SKU ID
- `page` (number, default: 1): Page number
- `page_size` (number, default: 20, max: 100): Items per page
- `sort_by` (string, default: "created_at"): Sort field (quantity, created_at, updated_at)
- `order` (string, default: "desc"): Sort order (asc, desc)

**Response (200 OK):**

```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
      "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
      "quantity": 150,
      "version": 1,
      "sku": {
        "id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
        "name": "Wireless Mouse",
        "category": "electronics",
        "description": "Ergonomic wireless mouse",
        "price": 29.99,
        "version": 1,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      },
      "store": {
        "id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
        "name": "Main Street Store",
        "address": "123 Main St, New York, NY",
        "created_at": "2025-01-01T09:00:00Z",
        "updated_at": "2025-01-20T09:00:00Z"
      },
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

**Notes:**

- Staff users will only see inventory for stores they are assigned to
- Results are cached in Redis for 5 minutes
- Cache is invalidated on inventory updates

**Errors:**

- 401: Unauthorized
- 403: Staff accessing stores they're not assigned to

---

### GET `/api/inventory/:id`

Get single inventory record details.

**Access:** All authenticated users (with store restrictions for staff)

**Path Parameters:**

- `id` (UUID): Inventory ID

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "quantity": 150,
  "version": 1,
  "sku": {
    "id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
    "name": "Wireless Mouse",
    "category": "electronics",
    "description": "Ergonomic wireless mouse",
    "price": 29.99,
    "version": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  "store": {
    "id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
    "name": "Main Street Store",
    "address": "123 Main St, New York, NY",
    "created_at": "2025-01-01T09:00:00Z",
    "updated_at": "2025-01-20T09:00:00Z"
  },
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Errors:**

- 401: Unauthorized
- 403: Staff accessing inventory from non-assigned store
- 404: Inventory not found

---

### POST `/api/manager/inventory`

Create new inventory record.

**Access:** Manager only

**Request Body:**

```json
{
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "quantity": 100
}
```

**Validation:**

- `sku_id`: required, must exist
- `store_id`: required, must exist
- `quantity`: required, minimum 0

**Response (201 Created):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "quantity": 100,
  "version": 1,
  "created_at": "2025-01-20T10:00:00Z",
  "updated_at": "2025-01-20T10:00:00Z"
}
```

**Side Effects:**

- Creates outbox record for cross-instance sync
- Invalidates related cache entries

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: SKU or Store not found
- 409: Inventory already exists for this SKU and store

---

### PUT `/api/manager/inventory/:id`

Update inventory quantity (direct set).

**Access:** Manager only

**Path Parameters:**

- `id` (UUID): Inventory ID

**Request Body:**

```json
{
  "quantity": 150
}
```

**Validation:**

- `quantity`: required, minimum 0

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "quantity": 150,
  "version": 2,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-20T11:00:00Z"
}
```

**Side Effects:**

- Uses optimistic locking (version number)
- Creates outbox record for cross-instance sync
- Invalidates related cache entries

**Errors:**

- 400: Validation error
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Inventory not found

---

### DELETE `/api/manager/inventory/:id`

Delete an inventory record.

**Access:** Manager only

**Path Parameters:**

- `id` (UUID): Inventory ID

**Response (200 OK):**

```json
{
  "message": "Inventory deleted successfully"
}
```

**Side Effects:**

- Creates outbox record for cross-instance sync
- Invalidates related cache entries

**Errors:**

- 400: Invalid inventory ID
- 401: Unauthorized
- 403: Forbidden (not a manager)
- 404: Inventory not found

---

### POST `/api/inventory/:id/adjust`

Adjust inventory quantity by delta (add or subtract).

**Access:** All authenticated users

- **Managers**: Can adjust any inventory
- **Staff**: Can only adjust inventory for their assigned stores

**Path Parameters:**

- `id` (UUID): Inventory ID

**Request Body:**

```json
{
  "delta_quantity": -5
}
```

**Validation:**

- `delta_quantity`: required, non-zero integer (positive = add, negative = subtract)
- Cannot result in negative quantity

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "quantity": 95,
  "version": 2,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-20T14:30:00Z"
}
```

**Side Effects:**

- Uses optimistic locking (version number)
- Creates outbox record for cross-instance sync
- Broadcasts update via Kafka and WebSocket
- Invalidates related cache entries

**Errors:**

- 400: Validation error or insufficient quantity
- 401: Unauthorized
- 403: Staff adjusting inventory from non-assigned store
- 404: Inventory not found

---

## WebSocket Events

### Connection

**URL:** `ws://localhost:3000/ws` or `wss://your-api-domain.com/ws`

**Authentication:** Include JWT token in Authorization header when upgrading connection.

### Server → Client Events

#### Inventory Update Event

Broadcast when inventory is updated from another instance.

**Event Format:**

```json
{
  "operation_type": "adjust",
  "inventory_id": "550e8400-e29b-41d4-a716-446655440000",
  "sku_id": "ed74446a-905b-4ea7-95cf-9e09c92e5c96",
  "sku_name": "Wireless Mouse",
  "store_id": "0f29b0ee-dc5f-4e74-baca-b6eacb56ea89",
  "store_name": "Main Street Store",
  "user_id": "770e8400-e29b-41d4-a716-446655440000",
  "user_name": "John Doe",
  "delta_quantity": -5,
  "new_quantity": 95,
  "version": 2
}
```

**Operation Types:**

- `create`: New inventory record created
- `update`: Inventory quantity updated (direct set)
- `adjust`: Inventory quantity adjusted (delta)
- `delete`: Inventory record deleted

---
