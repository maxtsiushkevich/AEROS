# AEROS

**AEROS (Airline Enterprise Resource Operations System)** – a microservices-based system for airline management
## Project Structure

```text
auth/    Authorization service, HTTP and gRPC layers
users/   User Service and Auth gRPC Client
flights/ Flight service, OpenAPI, and migrations
pkg/     Common errors, HTTP responses, and middleware
```

## Flights Service

Flight management service. It supports the creation, search, modification, and deletion of flights, as well as filtering by flight number, airports, status, and date range.

### Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| **Go** | 1.25.0 | Primary development language |
| **net/http** | Standard library | HTTP server and routing |
| **GORM** | 1.31.2 | ORM for PostgreSQL |
| **PostgreSQL** | 16 | Flight data storage |
| **pgx** | 5.10.0 | PostgreSQL driver |
| **validator** | 10.30.3 | DTO validation |
| **google/uuid** | 1.6.0 | UUID identifiers |
| **log/slog** | Standard library | Structured logging |
| **YAML** | yaml.v3 3.0.1 | Configuration |
The REST API is available at `/api/v1/flights`:

| Method | Purpose |
|---|---|
| `GET` | Search flights with filters |
| `POST` | Create a flight |
| `PATCH` | Update a flight |
| `DELETE` | Delete a flight by UUID |

The contract is defined in [flights/api/openapi.yaml](flights/api/openapi.yaml).
The migration creates the flights table, a status enum, an index, and sample data. Applying migrations

```
goose -dir migrations postgres "DB Connection string" up
```

## Auth Service

An authorization service with HTTP and gRPC interfaces. The `Auth.AddUser` gRPC method is used by the `users` service, and tokens are implemented as access and refresh tokens.

### Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| **Go** | 1.25.0 | Primary development language |
| **net/http** | Standard library | HTTP server |
| **gRPC** | 1.83.1 | RPC server |
| **Protocol Buffers** | protobuf 1.36.11 | Auth API contract and Go code generation |
| **JWT** | golang-jwt/jwt v5.3.1 | JWT token handling |
| **GORM** | 1.31.2 | Data layer interaction |
| **validator** | 10.30.3 | Data validation |
| **log/slog** | Standard library | Structured logging |
| **YAML** | yaml.v3 3.0.1 | Configuration |

### API

- `POST /api/v1/auth/refresh` — HTTP interface for refresh tokens;
- `Auth.AddUser` — gRPC method from [auth/api/proto/auth.proto](auth/api/proto/auth.proto).

## Users Service

A user service with an HTTP interface and an `auth` gRPC client. Gin is used for the HTTP layer, while the generated Protocol Buffers client is used to call the Auth service.

### Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| **Go** | 1.26.3 | Primary development language |
| **Gin** | 1.12.0 | HTTP server and routing | | **gRPC** | 1.83.1 | Auth service client |
| **Protocol Buffers** | protobuf 1.36.11 | Inter-service communication contract |
| **MongoDB Driver** | 2.5.0 | User storage driver |
| **GORM** | 1.31.2 | Integrated ORM component |
| **validator** | 10.30.1 | Data validation |
| **google/uuid** | 1.6.0 | UUID handling |
| **YAML** | yaml.v3 3.0.1 | Configuration |