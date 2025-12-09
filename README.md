# Final Report - Inventory Management System

## 0. Deployment Link, Demo Video Link & Repository Memo

Deployment Link:
http://209.38.10.93:3000

Demo Video:
https://drive.google.com/drive/folders/1qoFHWz1hmt2rJ-tTiNDw_6F0OQ8MnOsb?usp=sharing

Repository Memo:  
The frontend was previously developed in separate repositories before being migrated into this project (the commit history was cleaned during the migration). Legacy copies of the frontend code remain available at:

- https://github.com/axiao01/inventory-management-frontend
- https://github.com/konantian/inventory-manager-frontend
- https://github.com/andyzhou23/InventoryManagerServer

## 1. Team Members

| **Team Member** | **Student Number** | **GitHub Account** | **Email**                     |
| --------------- | ------------------ | ------------------ | ----------------------------- |
| Ming Yang       | 1006262223         | mingy6237          | mingy.yang@mail.utoronto.ca   |
| Yuan Wang       | 1002766526         | konantian          | ywang.wang@mail.utoronto.ca   |
| Yanrong Xiao    | 1006026444         | axiao01            | yanrong.xiao@mail.utoronto.ca |
| Yukun Zhou      | 1010122494         | andyzhou23         | yukun.zhou@mail.utornto.ca    |

## 2. Motivation

Many businesses that manage inventory in more than one location face problems like **slow updates, missing data, and manual work**. Existing systems often don’t update in real time, so staff may not see the latest stock information. This can lead to mistakes such as overselling, running out of products, or wasting time checking stock manually.

Our project aims to solve these problems by building a **cloud-based Inventory Management System** that keeps all data **up to date and consistent**. By using **containers**, **Kubernetes**, and **edge deployment**, the system allows updates to appear instantly everywhere.

This project is worth building because it helps businesses:

- Track and manage inventory **in real time**.
- Avoid data loss with **automatic backups**.
- Reduce waiting time through **low-latency, edge-based access**.
- Save time and money on manual work and system maintenance.

The target users include **warehouse managers**, **store employees**, and **administrators** who need fast and reliable access to stock data. Managers can monitor trends and prevent shortages, while staff can process orders or restock items without delays. For smaller companies, the system provides an affordable and scalable solution that can grow with their business.

Compared to existing tools, our approach is more flexible and future-ready. Many off-the-shelf inventory platforms are either too simple, lacking real-time updates, or too complex and expensive for smaller teams. They also tend to be hosted on centralized servers that increase latency and limit customization. Our cloud-native design removes these limitations by offering a lightweight, distributed, and easily deployable solution suitable for modern digital operations.

## 3. Objectives

- **Build a cloud-native, scalable inventory platform** that keeps stock data synchronized in real time between users and locations.
- **Provide a unified inventory view** across all stores and warehouses, with clear per-location breakdowns.
- **Ensure persistent, fault-tolerant data storage** by using PostgreSQL backed by Kubernetes StatefulSets and persistent volumes.
- **Support role-based workflows** so managers can administer users, stores, and SKUs, while staff can only operate on their assigned stores.
- **Deliver a responsive web experience** with instant stock updates, secure authentication, and low-stock/adjustment alerts.
- **Demonstrate distributed systems patterns** including Redis caching, Kafka-based outbox messaging, WebSocket fan-out, and Kubernetes orchestration.
- **Deploy the full stack with Docker and Kubernetes on DigitalOcean** to simplify deployment, scaling, and ongoing maintenance.

## 4. Technical Stack

- **Orchestration & Deployment**
  - Managed Kubernetes cluster on **DigitalOcean Kubernetes Service** with namespace‑scoped resources under `inventory-manager`.
  - Docker images for backend and frontend pushed to **Docker Hub**; `deploy.sh` and manifests in `k8s/` automate deployment to the cluster.
- **Backend**
  - Language: **Go 1.21** with **Gin** framework.
  - Data access: **GORM** ORM, PostgreSQL as the primary database (StatefulSet with persistent volumes).
  - Caching: **Redis** for inventory query caching and performance enhancing.
  - Messaging: **Kafka** for outbox‑driven inventory update events and cross‑instance consistency.
  - Real‑time: WebSocket hub broadcasting inventory events to connected clients (`/api/ws?token=...`).
  - Auth & security: JWT‑based authentication and role‑based authorization (manager vs staff).
- **Frontend**
  - **Next.js (App Router)** with **TypeScript** and **Tailwind CSS**.
  - React context for authentication (`AuthProvider`) and live inventory updates (`InventoryUpdatesProvider`).
  - Fully typed API client and hooks under `src/lib/` and `src/hooks/` consuming the documented REST and WebSocket APIs.
- **Infrastructure & Tooling**
  - Docker Compose for local multi‑service development (`InventoryManagerServer/compose.yaml`).
  - Vitest and React Testing Library for frontend tests.
  - ESLint and TypeScript for static checking.

Chosen approach for this project: **Kubernetes (K8s)** on **DigitalOcean Kubernetes Service** instead of Docker Swarm for higher availability and scalability.

These components together satisfy the core course requirements for containerization and local development (Docker + Docker Compose), state management (PostgreSQL with persistent volumes on Kubernetes), deployment to a cloud provider (DigitalOcean Kubernetes), orchestration with Kubernetes resources (Deployments, Services, PersistentVolumes), and monitoring/observability via DigitalOcean’s Kubernetes dashboard, metrics, and basic CPU/memory/disk alerts; the DigitalOcean dashboard is used to monitor cluster‑level resource usage (CPU, memory, disk), and an alert policy is configured to send email notifications when CPU utilization exceeds a configured threshold so that the team can react to overload conditions.

## 5. Features

### 5.1 Core and Advanced Features

- **Authentication and Authorization**

  - Username/password login via `POST /api/auth/login` returning a JWT.
  - Profile endpoint (`GET /api/profile`) and password change (`PUT /api/profile/password`).
  - Role‑based access: managers can manage users, stores, SKUs, and all inventory; staff are restricted to their assigned stores and limited actions.

- **User and Store Management (Manager‑only)**

  - Manage users via `/api/manager/users` (list, create, update, delete) with validation and duplicate checks.
  - Manage stores via `/api/manager/stores` (list, create, delete) and assign staff to stores via `/api/manager/stores/staff`.
  - UI pages for user and store administration under the protected dashboard provide manager workflows consistent with the API.

- **SKU and Inventory Management**

  - Full SKU lifecycle via `/api/manager/skus*`: list with pagination and filters, categories, detail view, create, update, delete. Only managers can perform these actions.
  - Inventory endpoints (`/api/inventory*` and `/api/manager/inventory*`) for listing, per‑store views, and CRUD on inventory records.
  - Staff and managers can perform stock adjustments via `POST /api/inventory/:id/adjust` with optimistic locking and validation (no negative stock).
  - The frontend exposes:
    - Items list with searching, filtering, sorting, pagination, and CSV export.
    - Item detail page with overview and per‑location quantities.

- **Real‑Time Inventory Updates and Alerts**

  - Kafka topics propagate inventory changes between backend instances. The consistency of the inventory is guaranteed by the outbox pattern, which means each write operation is written to the database within the transaction, and then to the Kafka topic.
  - WebSocket hub (`/api/ws`) broadcasts structured inventory update events (operation type, deltas, new quantity, user, store).
  - Frontend `InventoryUpdatesProvider` listens for WebSocket events and:
    - Highlights updated rows in the items table.
    - Updates counts and badges in near real time.
  - Alerts page in the dashboard shows low‑stock and out‑of‑stock items with severity (Warning/Critical) and status (Open/Acknowledged/Resolved).

- **Dashboard and Monitoring**

  - Dashboard page aggregates key metrics: total SKUs, low‑stock items, active users, total locations.
  - Low‑stock widget links directly to affected items.
  - Backend `GET /api/inventory` is cached in Redis, with cache invalidation on writes to balance performance and freshness.

- **Advanced Feature – Real‑Time Functionality**

  - WebSocket hub and Kafka‑backed outbox pattern deliver **live inventory updates** to all connected sessions.
  - The frontend updates rows, counters, and shows alerts without requiring manual refresh.

- **Advanced Feature – Security Enhancements**

  - JWT‑based authentication with protected routes on both backend and frontend, plus strict role‑based authorization (manager vs staff).
  - Sensitive configuration such as database credentials and JWT secrets is stored in Kubernetes **Secrets**, with CORS configured on the backend to protect cross‑origin access.

- **Advanced Feature – Auto‑scaling and High Availability**
  - Backend and frontend run as Kubernetes **Deployments** with multiple replicas on the managed DigitalOcean Kubernetes cluster.
  - Pods are fronted by Services and an Ingress/load balancer, with liveness/readiness probes and the ability to **scale replicas up or down** (e.g., via `kubectl scale`) to handle higher load while maintaining availability.
  - The backend exposes a dedicated `/health` endpoint, which is used by Kubernetes `livenessProbe` and `readinessProbe` on port `3000`; if the health check fails repeatedly, the Pod is automatically restarted, providing self‑healing behavior.
  - The Next.js frontend container is similarly configured with HTTP liveness/readiness probes on `/` (port `3000`), so unresponsive UI Pods are detected and recreated without manual intervention.

These features collectively satisfy the course requirements around authentication, CRUD operations, consistency, caching, messaging, Kubernetes deployment, real‑time behavior, and self‑healing via liveness probes.

## 6. User Guide

### 6.1 Accessing the Application

- **Cloud deployment (DigitalOcean Kubernetes – live demo)**
  - Open the browser at `http://209.38.10.93:3000`.
- **Default admin login**
  - Username: `admin`
  - Password: `adminadmin`

### 6.2 Main Navigation

Once logged in, the protected layout shows a sidebar/top navigation with:

- **Dashboard** – Overview metrics.
- **Items** – SKU managing for managers, and inventory managing for all roles.
- **Alerts** – Low‑stock and out‑of‑stock alerts.
- **Stores** – Store managing for managers.
- **Users** – User managing for managers.
- **Profile** – Current user info, password change and logout.

A WebSocket status indicator (green/yellow dot) reflects real‑time connection health.

### 6.3 Core User Flows

- **Login**

  - Go to the login page, enter username and password, optionally check “Remember Me,” then submit.
  - On success you are redirected to the dashboard; on failure a clear error message is shown.

- **Browsing and Managing Items**

  - Navigate to **Items** to see the full SKU list with category and location filters, search, sorting, and pagination.
  - Click “Add New Item” to create a SKU with basic info, category, optional image URL, threshold, and initial per‑location quantities.
  - Click on a row to open the **Item Detail** page, where you can view overview, locations; managers can edit or delete.

- **Adjusting Stock**

  - From the items table or item detail page, open the **Adjust Stock** modal.
  - Select a location, choose Add/Remove, enter a positive quantity, then submit.
  - The UI validates input, prevents removing more than existing stock, and updates affected views immediately; other sessions see updates via WebSocket.

- **Working with Alerts**

  - Open the **Alerts** page to view low‑stock and out‑of‑stock alerts
  - Users can resolve alerts by adding stock.

## 7. Development Guide

### 7.1 Prerequisites

- Go 1.21+
- Node.js and npm (for the Next.js frontend)
- Docker and Docker Compose
- `kubectl` (for interacting with the DigitalOcean Kubernetes cluster)

### 7.2 Backend – Local Development (Go)

1. Navigate to the backend:
   - `cd InventoryManagerServer/backend`
2. Install Go dependencies:
   - `go mod download`
3. Ensure PostgreSQL, Redis, and Kafka are available (either via Docker Compose or external services).
4. Run the server locally:
   - `go run main.go`
5. Health check:
   - `curl http://localhost:8080/health` → expects `{"status":"ok"}`.

### 7.3 Backend – Docker Compose

1. From `InventoryManagerServer/`, copy `example.env` to `.env` and adjust settings if needed.
2. Start the stack (Postgres, Redis, Kafka, backend, etc.):
   - `docker compose up --build --wait`
3. Stop and clean up:
   - `docker compose down -v`

Default ports (see `APIDoc.md` for details):

- API server: `http://localhost:8080`
- WebSocket: `ws://localhost:8080/api/ws?token=<jwt>`
- PostgreSQL: `localhost:15432`
- Redis: `localhost:16379`
- Kafka: `localhost:19092`

### 7.4 Frontend – Local Development

1. Navigate to the frontend:
   - `cd inventory-manager-frontend`
2. Create `.env.local`:
   - `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080`
   - `NEXT_PUBLIC_WS_URL=ws://localhost:8080/api/ws` (optional; derived if omitted).
3. Install dependencies and run dev server:
   - `npm install`
   - `npm run dev` → app listens on `http://localhost:3000`.
4. Run tests and lint:
   - `npm run test`
   - `npm run lint`

### 7.5 End‑to‑End Local Testing

1. Start backend stack via Docker Compose as above.
2. Start frontend dev server.
3. Visit `http://localhost:3000`, log in with `admin/adminadmin`.
4. Exercise:
   - Dashboard metrics and low‑stock widget.
   - Items list, item detail, new item creation, and stock adjustments.
   - User and store administration (manager account).
   - WebSocket behavior by opening two browser tabs and adjusting inventory in one tab while observing live updates in the other.

### 7.6 Kubernetes Deployment (DigitalOcean)

For a full Kubernetes deployment to the managed DigitalOcean cluster, follow `DEPLOYMENT.md` and `k8s/README.md`:

1. Configure `kubectl` to point to your DigitalOcean Kubernetes cluster (e.g., via the DigitalOcean UI or CLI).
2. Ensure your backend and frontend images have been built and pushed to Docker Hub.
3. From the repo root, run the deployment script:
   - `chmod +x deploy.sh`
   - `./deploy.sh`
4. The script:
   - Applies manifests in `k8s/` (Postgres, Redis, Kafka, backend, frontend, Ingress) to the current Kubernetes context.
   - Waits for core services to become ready.
5. Verify:
   - `kubectl get pods -n inventory-manager`
   - `kubectl get svc -n inventory-manager`
   - `curl http://209.38.10.93:3000/api/health`

## 8. Deployment Information

- **Deployment model**: Kubernetes, using a managed **DigitalOcean Kubernetes** cluster with the `inventory-manager` namespace.
- **Live demo URL (DigitalOcean)**: `http://209.38.10.93:3000`
- **Backend base path**: `http://209.38.10.93:3000/api`

## 9. Individual Contributions

Please adjust the granularity of the descriptions below to match your Git history.

- **Yuan Wang – Frontend**

  - Co‑designed and implemented the Next.js UI, including dashboard, items list and detail pages, and alerts.
  - Implemented React hooks and components for calling backend APIs, responsive layouts with Tailwind CSS, and integrated WebSocket‑driven updates into the UX.

- **Yanrong Xiao – Frontend**

  - Co‑implemented authentication flow, protected routing, and the `AuthProvider` context managing JWTs and session state.
  - Built UI components for admin workflows (user and store management), form validation, and integrated unit tests.

- **Yukun Zhou – Backend**

  - Implemented core backend modules: handlers, services, and models for users, stores, SKUs, inventory, etc.
  - Added Redis caching, Kafka producer/consumer logic, outbox processing, and the WebSocket hub to support real‑time, distributed inventory updates.

- **Ming Yang – DevOps / Kubernetes**
  - Containerized backend and frontend, authored Dockerfiles, and set up Docker Compose for local development.
  - Designed and implemented Kubernetes manifests in `k8s/`, including Postgres, Redis, Kafka, deployments, services, Ingress, health checks, scaling, and the automated `deploy.sh` script.

## 10. Lessons Learned and Concluding Remarks

- **Distributed consistency requires explicit patterns**: implementing the outbox pattern with Kafka and Redis cache invalidation reinforced the importance of handling concurrency and replication explicitly rather than relying on “best effort” updates.
- **Kubernetes orchestration adds real‑world complexity**: wiring multiple stateful and stateless services together (database, cache, message queue, backend, frontend) highlighted the value of health probes, proper service naming, secrets/config maps, and observability.
- **Full‑stack integration is iterative**: aligning the frontend feature set with the backend API (including roles, validation rules, and error handling) required multiple feedback cycles, but resulted in a consistent and predictable experience.
- **Overall**, the project gave the team hands‑on experience with a realistic cloud‑native architecture, from database and messaging layers up through a modern web UI, and significantly deepened our understanding of scalable inventory systems and distributed application design.
