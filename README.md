# AEROS

AEROS (Airline Enterprise Resource Operations System) is a Go-based airline management system.

## Overview


- `auth/` — authentication service, JWT handling, RBAC API, and gRPC auth server
- `flights/` — flight catalog and operational endpoints
- `users/` — user API and auth-client integration
- `pkg/` — shared packages used by Go services (`httpresp`, `httperr`, middleware, RBAC helpers)
- `rbac_config/` — Casbin model and configuration
- `docker-compose.yaml` — local infrastructure stack for Redis, PostgreSQL, Kafka, and admin UIs

### Main modules

| Module | Responsibility |
|---|---|
| `auth/` | Login, refresh/logout flow, access and refresh token validation, RBAC HTTP API, gRPC auth server |
| `flights/` | Flight CRUD, search, status updates, date-based filtering |
| `users/` | User API and gRPC client used for registration/activation workflows |
| `pkg/` | Shared HTTP response helpers, middleware, JWT tooling, RBAC layer |
| `rbac_config/` | Casbin model and config values |
| `aeros-ui/` | Simple UI to create and manage flights |

Each Go module is independent and uses the shared package through `replace pkg => ../pkg`.

## Technology stack

### Backend

| Technology | Purpose |
|---|---|
| Go | Core microservices |
| Gin | HTTP routing |
| Redis | Caching |
| gRPC | Auth service integration |
| Protobuf | Auth contracts |
| GORM | ORM layer for PostgreSQL |
| pgx | PostgreSQL driver |
| YAML | Service config files |
| `validator/v10` | Request validation |
| `google/uuid` | ID generation |
| `log/slog` | Structured logging |
| Kafka | Message broker |

### Auth and authorization

| Technology | Purpose |
|---|---|
| JWT | Access/refresh tokens |
| Redis | Token revocation cache and RBAC update bus |
| Casbin | Role-based access control |
| bcrypt | Password hashing |

### Integration

| Technology | Purpose |
|---|---|
| Kafka | Event bus / async integration support |
| Docker Compose | Local dev infrastructure |
| PostgreSQL | Service persistence |
| Goose | SQL migration management |

## Requirements

- Go 1.26.x (as declared in the module files)
- Docker and Docker Compose
- PostgreSQL 16+
- Redis 7+
- Goose CLI for migrations


## Local environment

The repository includes a Compose stack for data and local tooling:

```bash
docker compose up -d
```

This starts:

- Redis: `localhost:6379` (password `123`)
- RedisInsight: `http://localhost:5540`
- PostgreSQL: `localhost:5432` (user `user`, password `123`, DB `flights` by default)
- Kafka: `localhost:9092`
- Kafka UI: `http://localhost:8080`

## Database setup

The default local configuration expects PostgreSQL on `localhost:5432` with credentials:

- user: `user`
- password: `123`

Create the databases needed by the services:

```sql
CREATE DATABASE auth;
CREATE DATABASE flights;
CREATE DATABASE casbin;
```

Then apply migrations from the service directories:

```bash
goose -dir auth/migrations postgres "host=localhost port=5432 user=user password=123 dbname=auth sslmode=disable" up
goose -dir flights/migrations postgres "host=localhost port=5432 user=user password=123 dbname=flights sslmode=disable" up
goose -dir pkg/rbac/migrations postgres "host=localhost port=5432 user=user password=123 dbname=casbin sslmode=disable" up
```

The config files are in:

- [auth/config/config.yaml](auth/config/config.yaml)
- [flights/config/config.yaml](flights/config/config.yaml)
- [users/config/config.yaml](users/config/config.yaml)
- [rbac_config/config.yaml](rbac_config/config.yaml)

## Running the services

Open separate terminals and run:

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

The users API is defined in [users/api/openapi.yaml](users/api/openapi.yaml).


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
