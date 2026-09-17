# AEROS

**AEROS (Airline Enterprise Resource Operations System)** is a Go-based airline management system.

## Architecture

### Repository layout

| Directory | Responsibility |
|---|---|
| `auth/` | Authentication, JWT access/refresh tokens, RBAC HTTP API, and Auth gRPC server |
| `users/` | User HTTP API and client of the Auth gRPC service |
| `flights/` | Flight CRUD and filtering API |
| `pkg/` | Shared middleware, JWT helpers, HTTP responses, errors, and Casbin RBAC integration |
| `rbac_config/` | Casbin model and runtime configuration |

Each Go service is a separate module. The service modules use the local shared module through `replace pkg => ../pkg`.

## Technology stack

### Backend

| Technology | Version | Usage |
|---|---|---|
| Go | 1.26.5 | Service implementation |
| `net/http` | Standard library | Auth and flights HTTP servers |
| Gin | 1.12.0 | Users HTTP server and routing |
| gRPC | 1.83.1 | Auth/users inter-service communication |
| Protocol Buffers | 1.36.11 | gRPC contract and generated Go clients/servers |
| GORM | 1.31.2 | PostgreSQL data access |
| pgx | 5.10.0 | PostgreSQL driver used by GORM |
| `validator/v10` | 10.30.3 | Request validation |
| `google/uuid` | 1.6.0 | Entity identifiers |
| `log/slog` | Standard library | Structured logging |
| YAML v3 | 3.0.1 | Service configuration |

### Authentication and authorization

| Technology | Version | Usage |
|---|---|---|
| JWT | `golang-jwt/jwt/v5` 5.3.1 | Access and refresh tokens |
| Redis | 7 | Revoked-token cache and Casbin update notifications |
| Casbin | 3.11.0 | Role-based access control |
| bcrypt | `golang.org/x/crypto` 0.56.0 | Password hashing |

### Data and infrastructure

| Technology | Version | Usage |
|---|---|---|
| PostgreSQL | 16 recommended | Auth, flights, and RBAC persistence |
| Goose | CLI | SQL migration management |
| Docker Compose | 3.8 | Local Redis and RedisInsight stack |

## Requirements

- Go 1.26.5 or newer
- PostgreSQL with databases `auth`, `flights`, and `casbin`
- Redis 7 or newer
- [Goose](https://github.com/pressly/goose) for database migrations

The repository currently provides Docker Compose only for Redis and RedisInsight; PostgreSQL must be installed and configured separately.

## Local setup

### 1. Start Redis

From the repository root:

```bash
docker compose -f docker-compose.yaml up -d
```

Redis is available at `localhost:6379` with password `123`. RedisInsight is optional at `http://localhost:5540`.

### 2. Prepare PostgreSQL

The default local configuration expects the PostgreSQL user `user`, password `123`, and databases `auth`, `flights`, and `casbin` on `localhost:5432`. Create those databases or change the service configuration files:

- [auth/config/config.yaml](auth/config/config.yaml)
- [flights/config/config.yaml](flights/config/config.yaml)
- [users/config/config.yaml](users/config/config.yaml)
- [rbac_config/config.yaml](rbac_config/config.yaml)

Apply the migrations:

```bash
goose -dir auth/migrations postgres "postgres://user:123@localhost:5432/auth?sslmode=disable" up
goose -dir flights/migrations postgres "postgres://user:123@localhost:5432/flights?sslmode=disable" up
goose -dir pkg/rbac/migrations postgres "postgres://user:123@localhost:5432/casbin?sslmode=disable" up
```

The flights migration also inserts sample flights.

### 3. Start the services

Run each command from its service directory in a separate terminal:

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

The default addresses are:

| Component | Address |
|---|---|
| Flights HTTP API | `http://localhost:3000` |
| Users HTTP API | `http://localhost:3002` |
| Auth HTTP API | `http://localhost:3001` |
| Auth gRPC API | `localhost:50051` |

## APIs

### Flights

The REST API is defined in [flights/api/openapi.yaml](flights/api/openapi.yaml):

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/flights` | Search and filter flights |
| `POST` | `/api/v1/flights` | Create a flight |
| `PATCH` | `/api/v1/flights` | Update a flight |
| `DELETE` | `/api/v1/flights?id=<uuid>` | Delete a flight |

Search supports flight number, origin, destination, status, RFC3339 date range, page, and limit filters. Supported statuses are `Scheduled`, `CheckIn`, `Boarding`, `Delayed`, `Departed`, `Arrived`, `Cancelled`, and `Redirected`.

### Auth and RBAC

The contract is defined in [auth/api/openapi.yaml](auth/api/openapi.yaml):

| Method | Endpoint | Purpose |
|---|---|---|
| `POST` | `/api/v1/auth/login` | Issue access and refresh tokens |
| `POST` | `/api/v1/auth/refresh` | Refresh tokens using the HttpOnly cookie |
| `POST` | `/api/v1/auth/logout` | Revoke the refresh token |
| `POST` | `/api/v1/auth/change-password` | Change the authenticated user's password |
| `PUT` / `DELETE` | `/api/v1/rbac/...` | Manage roles, permissions, and assignments |

The Auth gRPC contract is [auth/api/proto/auth.proto](auth/api/proto/auth.proto). The `Auth.AddUser` method is consumed by the users service.

### Users

The users service currently exposes:

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/users/` | Read the authenticated account |
| `POST` | `/api/v1/users/registration` | Invoke Auth gRPC user registration flow |

## Configuration and security

Configuration is YAML-based and loaded with the `-config` flag, for example:

```bash
go run ./cmd/auth -config config/config.yaml
```

Auth uses `JWT_SECRET` and `JWT_REFRESH_SECRET`. Development defaults are created automatically when these variables are absent; set both variables explicitly outside local development. Redis and database credentials in the checked-in YAML files are development values and must be replaced for shared or production environments.

Auth and users also read `RBAC_CONFIG_PATH`; by default it points to `../rbac_config/config.yaml` when started from their service directories.

## Development checks

Run tests and compile all Go modules:

```bash
cd auth && go test ./...
cd ../flights && go test ./...
cd ../users && go test ./...
cd ../pkg && go test ./...
```

## Current limitations

- `flights` and `users` both default to `localhost:3000`, so they cannot run simultaneously without separate config files or a routing layer.
- `users/config/config.yaml` currently points to the `flights` database; verify or correct this before enabling persistent user storage.
- The users registration handler currently uses a hard-coded example UUID, email, and password and returns `pong`; it is not a production-ready registration endpoint yet.
- PostgreSQL is not included in the existing Compose file.
