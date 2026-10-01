# AEROS

**AEROS (Airline Enterprise Resource Operations System)** is a Go-based airline management system.

## Architecture

### Repository layout

| Directory | Responsibility |
|---|---|
| `auth/` | Authentication, JWT access/refresh tokens, RBAC HTTP API, and auth gRPC server |
| `flights/` | Flight catalog, search, and operational actions |
| `users/` | User API and client for the auth gRPC registration flow |
| `pkg/` | Shared middleware, JWT helpers, HTTP responses, errors, and RBAC integration |
| `rbac_config/` | Casbin model and runtime configuration |
| `notifications/`, `passengers/`, `tickets/`, `seats-booking/` | Additional service modules planned 

Each Go module is independent and the shared package is linked through `replace pkg => ../pkg`.

## Technology stack

### Backend

| Technology | Usage |
|---|---|
| Go | Main service implementation |
| `net/http` | Auth and flights HTTP servers |
| Gin | Users HTTP routing |
| gRPC | Auth/users inter-service communication |
| Protocol Buffers | gRPC contracts and generated Go bindings |
| GORM | PostgreSQL persistence |
| pgx | PostgreSQL driver for GORM |
| `validator/v10` | Request validation |
| `google/uuid` | Entity identifiers |
| `log/slog` | Structured logging |
| YAML | Configuration files |

### Authentication and authorization

| Technology | Usage |
|---|---|
| JWT | Access and refresh tokens |
| Redis | Revoked-token cache and RBAC update notifications |
| Casbin | Role-based access control |
| bcrypt | Password hashing |

### Data and infrastructure

| Technology | Usage |
|---|---|
| PostgreSQL | Service persistence |
| Goose | SQL migration management |
| Docker Compose | Local Redis and RedisInsight stack |

## Requirements

- Go 1.26+ recommended
- PostgreSQL databases: `auth`, `flights`, and `casbin`
- Redis 7+
- Goose CLI for migrations

The repository currently ships with Docker Compose for Redis and RedisInsight; PostgreSQL must be prepared manually.

## Local setup

### 1. Start Redis

```bash
docker compose -f docker-compose.yaml up -d
```

Redis is exposed at `localhost:6379` with password `123`; RedisInsight is available at `http://localhost:5540`.

### 2. Prepare PostgreSQL

The default local configuration expects PostgreSQL on `localhost:5432` with the user `user` and password `123`.
Create the required databases:

- `auth`
- `flights`
- `casbin`

Update service configs if needed:

- [auth/config/config.yaml](auth/config/config.yaml)
- [flights/config/config.yaml](flights/config/config.yaml)
- [users/config/config.yaml](users/config/config.yaml)
- [rbac_config/config.yaml](rbac_config/config.yaml)

Apply migrations:

```bash
goose -dir auth/migrations postgres "postgres://user:123@localhost:5432/auth?sslmode=disable" up
goose -dir flights/migrations postgres "postgres://user:123@localhost:5432/flights?sslmode=disable" up
goose -dir pkg/rbac/migrations postgres "postgres://user:123@localhost:5432/casbin?sslmode=disable" up
```

### 3. Start the services

Run each service in its own terminal:

```bash
cd auth
go run ./cmd/auth
```

```bash
cd flights
go run ./cmd/flights
```

```bash
cd users
go run ./cmd/users
```

Default local endpoints:

| Component | Address |
|---|---|
| Auth HTTP API | `http://localhost:3001` |
| Flights HTTP API | `http://localhost:3000` |
| Users HTTP API | `http://localhost:3002` |
| Auth gRPC server | `localhost:50051` |

## APIs

### Flights

The flights REST API is defined in [flights/api/openapi.yaml](flights/api/openapi.yaml).

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/flights` | Search and filter flights |
| `POST` | `/api/v1/flights` | Create a flight |
| `DELETE` | `/api/v1/flights/{id}` | Delete a flight |
| `POST` | `/api/v1/flights/{id}/cancel` | Mark a flight as cancelled |
| `POST` | `/api/v1/flights/{id}/reschedule` | Move a flight to a new date |
| `POST` | `/api/v1/flights/{id}/status` | Change a flight status |

Filters support `flight_number`, `origin`, `destination`, `status`, `date_from`, `date_to`, `page`, and `limit`. Supported statuses are `Scheduled`, `CheckIn`, `Boarding`, `Delayed`, `Departed`, `Arrived`, `Cancelled`, and `Redirected`.

### Auth and RBAC

The auth and RBAC API is defined in [auth/api/openapi.yaml](auth/api/openapi.yaml).

| Method | Endpoint | Purpose |
|---|---|---|
| `POST` | `/api/v1/auth/login` | Authenticate and issue tokens |
| `POST` | `/api/v1/auth/refresh` | Refresh access and refresh tokens using the HttpOnly cookie |
| `POST` | `/api/v1/auth/logout` | Revoke the refresh token |
| `POST` | `/api/v1/auth/change-password` | Change the authenticated user's password |
| `PUT` | `/api/v1/rbac/roles` | Create a role |
| `DELETE` | `/api/v1/rbac/roles/{role}` | Delete a role |
| `PUT` | `/api/v1/rbac/actions` | Create an action |
| `PUT` | `/api/v1/rbac/resources` | Create a resource |
| `PUT` | `/api/v1/rbac/permissions` | Create a permission |
| `PUT` | `/api/v1/rbac/roles/{role_name}/permissions` | Grant a permission to a role |
| `DELETE` | `/api/v1/rbac/roles/{role_name}/permissions` | Revoke a permission from a role |
| `PUT` | `/api/v1/rbac/users/{user_id}/roles` | Assign a role to a user |
| `DELETE` | `/api/v1/rbac/users/{user_id}/roles/{role_name}` | Remove a role from a user |

The gRPC contract is in [auth/api/proto/auth.proto](auth/api/proto/auth.proto) and is consumed by the users service.

### Users

The users service currently exposes:

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/users/` | Read the authenticated account |
| `POST` | `/api/v1/users/registration` | Trigger the auth gRPC registration flow |
| `POST` | `/api/v1/users/:userId/activate` | Activate a user account |

## Configuration and security

Configuration is YAML-based and loaded through the `-config` flag:

```bash
go run ./cmd/auth -config config/config.yaml
```

Auth expects `JWT_SECRET` and `JWT_REFRESH_SECRET`. In local development they can be auto-generated, but they should be set explicitly outside local development. The checked-in YAML files use development credentials and should not be used unchanged in shared or production environments.

## Development checks

```bash
cd auth && go test ./...
cd ../flights && go test ./...
cd ../users && go test ./...
cd ../pkg && go test ./...
```

## Current limitations

- `flights` and `users` both default to `localhost:3000` in their sample configs, so they cannot run simultaneously without separate ports or a reverse proxy.
- `users/config/config.yaml` currently points to the `flights` database and should be verified before enabling persistent user storage.
- The users registration handler still uses example values and should be treated as a prototype rather than production-ready registration logic.
- PostgreSQL is not included in the existing Docker Compose setup.
