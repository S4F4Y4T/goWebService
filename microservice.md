# Microservice Architecture Documentation

This comprehensively documents the system design and development practices for the newly refactored distributed architecture. The system has migrated from a modular monolith to an API Gateway-driven Microservice pattern while retaining the simplicity of a multi-module monolithic repository (`go.work`).

## 1. Architecture Overview

The system is split into four primary decoupled entities running in Docker containers, sitting behind a central Golang native API Gateway that handles path-rewriting and authentication token guarding.

```mermaid
graph TD
    Client[Client / Browser]
    Gateway[API Gateway :8080]
    Auth[Auth Service :8081]
    User[User Service :8082]
    Product[Product Service :8083]
    DB[(Shared Postgres Container)]
    
    Client -->|HTTP / REST| Gateway
    Gateway -->|/auth/*| Auth
    Gateway -->|/users/* (JWT Protected)| User
    Gateway -->|/products/* (JWT Protected)| Product
    
    Auth -.->|Inter-Service HTTP| User
    
    Auth -->|auth_db| DB
    User -->|user_db| DB
    Product -->|product_db| DB
```

### Components

* **API Gateway (`gateway/`)**: A pure Go reverse proxy (`httputil.NewSingleHostReverseProxy`). Serves as the single public entry point exposing port `8080`.
* **Auth Service (`auth-service/`)**: Handles credential issuance and JWT token signing.
* **User Service (`user-service/`)**: A CRUD service purely managing the `User` aggregate and its persistence. 
* **Product Service (`product-service/`)**: A CRUD service managing the `Product` aggregate and its persistence.

---

## 2. Shared Libraries (The `pkg/` Namespace)

Rather than establishing separate git repositories for every microservice, the system uses the Go `1.18+` **Workspaces** pattern (`go.work`). This drastically reduces friction when sharing Domain configurations, Error definitions, HTTP formats, and JWT parsing logic.

* `pkg/apperror/`: A standardized way to define structured REST errors explicitly mapped to HTTP Status Codes.
* `pkg/response/`: Standard JSON encoders mapping structures explicitly to successes, failures, and app errors.
* `pkg/jwtutil/`: Common logic used by the Auth service to *Sign* the JWT, and by the Gateway to *Validate* the JWT using a shared `HS256` secret.

---

## 3. Database Architecture

While the infrastructure hosts a single `postgres:15-alpine` container instance to conserve memory, **each service connects to its own isolated logical database**. 
* On the first run, `scripts/init-dbs.sql` is invoked by Postgres.
* It initializes three decoupled database instances: `auth_db`, `user_db`, and `product_db`.
* Each Go Microservice independently invokes `github.com/golang-migrate/migrate` on boot to apply SQL upgrades to exclusively its target DB.

---

## 4. Authentication Flow

Authentication relies on stateless JWT tokens.

### A. Registration Flow (Internal Orchestration)
1. `Client` sends `POST /auth/register` to the Gateway.
2. Gateway proxies unmodified payload to `Auth Service`.
3. `Auth Service` executes an internal synchronous `HTTP POST` request to `http://user-service:8082/users/` passing just the `name` and `email` to configure the user profile first.
4. If successful, `Auth Service` captures the `User ID` returned by the User Service. 
5. `Auth Service` hashes the provided plaintext password exclusively via `bcrypt` and stores it into `auth_db`.
6. Returns `201 Created` with a freshly signed JWT.

### B. Route Guarding Example (Fetching Users)
1. `Client` sends `GET /users/` to the Gateway with header `Authorization: Bearer <token>`.
2. Gateway intercepts request at its middleware router. 
3. Gateway invokes `jwtutil.ValidateToken` against the shared `JWT_SECRET`. 
4. If successful, it securely identifies the user context entirely statelessly. 
5. Gateway injects the headers `X-User-ID: <id>` and `X-User-Email: <email>` into the downstream request object, stripping the `/users` prefix.
6. `User Service` receives the request at `/` and handles it purely, unworried about parsing headers for authorization validation.

---

## 5. Deployment & Configuration

The deployment is managed synchronously using Docker Compose (`docker-compose.yml`).

### Core Environment Flags
| Variable | Owner | Purpose |
|----------|-------|---------|
| `JWT_SECRET` | Gateway, Auth | Used to securely sign and guard `HS256` tokens. Must be securely populated in Production. |
| `DB_NAME` | Auth, User, Product | Evaluated explicitly per container mapping to targeted database mapping constraints. |
| `DB_USER` / `DB_PASSWORD` | Postgres | Baseline Database administrative login credentials via container environments. |

### Docker Build Definitions
Every specific microservice manages its own localized `Dockerfile`. Because the services share Go codebase components from `pkg/`, the Dockerfiles are composed by invoking `COPY . .` from the Root directory layout context utilizing multi-stage isolated `/app` workspace containers.

---

## 6. Local Development Operations

### Standing up the Application Environment
To automatically map to port `8080`, simply build:
```bash
docker compose up -d --build
```
> This initiates healthchecks validating Postgres's availability prior to releasing dependent Go Microservices online seamlessly sequentially.

### Rebuilding a Single Service
If iterating heavily explicitly on one domain (e.g. `product-service`):
```bash
docker compose build product-service
docker compose up -d product-service
```

### Formatting Modules
It's important to remember this functions inside Go Workspaces structurally:
```bash
go get -u ./...
go mod tidy
go work sync
```
