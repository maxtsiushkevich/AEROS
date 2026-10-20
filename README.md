# AEROS

AEROS (Airline Enterprise Resource Operations System) is a Go-based airline management system built as a set of independently runnable services.

## Overview

- `auth/` — authentication, JWT issuance/revocation, RBAC HTTP API, and gRPC server
- `users/` — user registration, profile, OTP-based activation, and Kafka event publishing
- `flights/` — flight operations and per-flight cargo manifests/items
- `pkg/` — shared HTTP helpers, JWT/middleware, and Casbin RBAC implementation
- `rbac_config/` — Casbin model and local RBAC configuration
- `docker-compose.yaml` — local infrastructure
- `configs/` — configuration files mounted into the Compose auth container

### Main modules

| Module | Responsibility |
|---|---|
| `auth/` | Login, refresh/logout, password changes, RBAC administration, and gRPC auth server |
| `flights/` | Flight CRUD/search/status operations and cargo manifest/item management |
| `users/` | User registration/profile, OTP activation through Redis and Kafka, and auth gRPC client |
| `pkg/` | Shared HTTP helpers, JWT tooling, auth middleware, and PostgreSQL-backed RBAC |
| `rbac_config/` | Casbin model and RBAC configuration |

Each Go module is independent and uses the shared package through `replace pkg => ../pkg`.

## Technology stack

### Backend

| Technology | Purpose |
|---|---|
| Go | Core microservices|
| Gin | HTTP routing|
| `net/http` | HTTP routing |
| Redis | Refresh-token revocation, OTP cache, and RBAC policy-change notifications |
| gRPC | Auth service integration |
| Protobuf | Auth contracts |
| GORM | ORM layer for PostgreSQL |
| pgx | PostgreSQL driver |
| YAML | Service config files |
| `validator/v10` | Request validation |
| `google/uuid` | ID generation |
| `log/slog` | Structured logging |
| Kafka | User OTP event transport |

### Auth and authorization

| Technology | Purpose |
|---|---|
| JWT | Access/refresh tokens |
| Redis | Refresh-token revocation, user OTP cache, and RBAC update notifications |
| Casbin | Role-based access control |
| bcrypt | Password hashing |

### Integration

| Technology | Purpose |
|---|---|
| Kafka | Asynchronous OTP event publishing by the users service |
| Docker Compose | Local dev infrastructure |
| PostgreSQL | Service persistence |
| Goose | SQL migration management |

## Requirements

- Go 1.26.5
- Docker and Docker Compose
- PostgreSQL 16+
- Redis 7+
- Goose CLI for migrations


## Local environment

The Compose stack starts Redis, RedisInsight, PostgreSQL, Kafka, Kafka UI, and the auth service.

```bash
docker compose up -d
```

This starts:

- Redis: `localhost:6379` (password `123`)
- RedisInsight: `http://localhost:5540`
- PostgreSQL: `localhost:5432` (user `user`, password `123`, DB `flights` by default)
- Kafka: `localhost:9092`
- Kafka UI: `http://localhost:8080`
- Auth HTTP API: `http://localhost:3001`
- Auth gRPC API: `localhost:50051`

## Database setup

The default local configuration expects PostgreSQL on `localhost:5432` with credentials:

- user: `user`
- password: `123`

Create the databases needed by the services:

```sql
CREATE DATABASE auth;
CREATE DATABASE users;
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
- [configs/auth.yaml](configs/auth.yaml) and [configs/rbac.yaml](configs/rbac.yaml) (Compose-mounted auth/RBAC settings)

## Running the services

Start the infrastructure first with `docker compose up -d`. Then start services in separate terminals. The auth service can be run in Docker Compose; to run it locally instead:

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

Each flight gets a cargo manifest. Cargo endpoints are implemented in the service but are not yet listed in the Flights OpenAPI document:

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/cargo/{id}` | Get a flight's cargo manifest |
| `GET` | `/api/v1/cargo/{id}/items` | List manifest items |
| `GET` | `/api/v1/cargo/{id}/items/{item_id}` | Get one cargo item |
| `POST` | `/api/v1/cargo/{id}/add` | Add an item to a manifest |
| `DELETE` | `/api/v1/cargo/{id}/delete/{item_id}` | Delete an item |
| `POST` | `/api/v1/cargo/{id}/move` | Move an item to another manifest |

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
| `GET` | `/api/v1/rbac/roles` | List roles |
| `GET` | `/api/v1/rbac/actions` | List actions |
| `GET` | `/api/v1/rbac/resources` | List resources |
| `GET` | `/api/v1/rbac/permissions` | List permissions |
| `PUT` | `/api/v1/rbac/roles/{role_name}/permissions` | Grant a permission to a role |
| `DELETE` | `/api/v1/rbac/roles/{role_name}/permissions` | Revoke a permission from a role |
| `GET` | `/api/v1/rbac/users/{user_id}/roles` | List a user's roles |
| `PUT` | `/api/v1/rbac/users/{user_id}/roles` | Assign a role to a user |
| `DELETE` | `/api/v1/rbac/users/{user_id}/roles/{role_name}` | Remove a role from a user |

The gRPC contract is in [auth/api/proto/auth.proto](auth/api/proto/auth.proto) and is consumed by the users service.
The service also serves a browser-based RBAC management page at `/admin/rbac`. This page currently has no authentication middleware; use it only in a trusted local environment.

### Users

The users API is defined in [users/api/openapi.yaml](users/api/openapi.yaml).


| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/v1/users` | Read the authenticated account (Bearer access token) |
| `POST` | `/api/v1/users/registration` | Register and create the auth account; returns an access token and refresh-token cookie, publish activation code |
| `POST` | `/api/v1/users/request-activation` | Generate and publish a new activation code |
| `POST` | `/api/v1/users/activate` | Activate the authenticated account using the submitted OTP code |
| `POST` | `/api/v1/users/{userId}/activate` | Admin activation of the specified account |

Activation codes are cached in Redis for 10 minutes. The users service publishes OTP events to the configured Kafka topic (default `users`). The request and self-activation routes require a valid Bearer token. See [users/api/openapi.yaml](users/api/openapi.yaml) for the API contract; it currently needs to be brought in sync with the OTP routes.

## Configuration and security

Configuration is YAML-based and loaded through the `-config` flag:

```bash
go run ./cmd/auth -config config/config.yaml
```

Auth expects `JWT_SECRET` and `JWT_REFRESH_SECRET`. The current auth entry point supplies development fallback values when they are unset; set strong secrets explicitly outside local development. Both auth and users use `RBAC_CONFIG_PATH` for the shared RBAC configuration (the entry points default to `../rbac_config/config.yaml`). Checked-in YAML and Compose credentials are for local development only and must not be used unchanged in shared or production environments.

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
