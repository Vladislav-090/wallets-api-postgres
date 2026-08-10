# Wallets API

English | [Русский](README_RU.md)

Wallets API is a learning REST API written in Go for managing users, wallets, and money transfers.

The project is built with clear separation of responsibilities between HTTP handlers, the service layer, and the database access layer.

The API supports user registration and authentication, multiple wallets per user, transfers between wallets, role-based access control, administrative operations, SQL migrations, structured logging, unit tests, and Swagger/OpenAPI documentation.

## Features

### Users

- User registration
- User authentication
- Password hashing with bcrypt
- JWT authentication
- `user` and `admin` roles
- Admin access to the list of all users
- Admin access to a user by ID

### Wallets

- Create a wallet
- Create multiple wallets for one user
- Get all wallets of the authenticated user
- Get a wallet by ID
- Update wallet name
- Delete a wallet
- Associate wallets with a specific user
- Restrict users to their own wallets
- Get all wallets as admin
- Get wallets of a specific user as admin
- Get any wallet by ID as admin

### Transfers

- Create transfers between wallets
- Get transfer history of the authenticated user
- Get a transfer by ID
- Check sufficient balance
- Prevent transfers between the same wallet
- Check wallet currency compatibility
- Execute transfers inside an SQL transaction
- Get all transfers as admin
- Get any transfer by ID as admin

## Architecture

The project is divided into several layers:

- `handlers` — HTTP request handling and response generation
- `service` — application business logic
- `repository` — interaction with PostgreSQL
- `models` — application data structures
- `middleware` — authentication, authorization, logging, and panic recovery
- `router` — HTTP route registration
- `config` — application configuration loading
- `database` — PostgreSQL connection setup
- `response` — unified HTTP response formatting
- `auth` — JWT-related logic
- `migrations` — SQL database migrations

Main HTTP request flow:

```text
Client
  ↓
Router
  ↓
Middleware
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
PostgreSQL
```

This separation keeps HTTP logic, business logic, and database access independent from each other.

## Technologies

- Go
- `net/http`
- PostgreSQL
- SQL
- Docker
- Docker Compose
- JWT
- bcrypt
- `shopspring/decimal`
- `context.Context`
- `log/slog`
- golang-migrate
- Swagger / OpenAPI
- Swaggo
- Git
- GitHub

## Authentication

After successful login, the user receives a JWT token.

Protected endpoints require the token in the HTTP header:

```text
Authorization: Bearer <token>
```

`AuthMiddleware` validates the JWT, extracts the claims, and stores user information in the current request `context.Context`.

Handlers can then retrieve the authenticated user data from the request context.

Administrative routes additionally use `AdminMiddleware`, which checks the user role.

## Roles

The application uses two roles.

### user

A regular user can:

- create wallets
- get their wallets
- get their wallet by ID
- update their wallets
- delete their wallets
- create transfers
- view their transfers

### admin

An administrator can additionally:

- get all users
- get a user by ID
- get all wallets
- get wallets of a specific user
- get any wallet by ID
- get all transfers
- get any transfer by ID

## API Endpoints

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register a user |
| POST | `/login` | Authenticate and receive JWT |
| GET | `/health` | Check API health |

### Wallets

JWT authentication is required.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/wallets` | Create a wallet |
| GET | `/wallets` | Get authenticated user's wallets |
| GET | `/wallets/{id}` | Get authenticated user's wallet by ID |
| PATCH | `/wallets/{id}` | Update wallet name |
| DELETE | `/wallets/{id}` | Delete a wallet |

### Transfers

JWT authentication is required.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/transfers` | Create a transfer |
| GET | `/transfers` | Get authenticated user's transfers |
| GET | `/transfers/{id}` | Get a transfer by ID |

### Admin

JWT authentication and the `admin` role are required.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/users` | Get all users |
| GET | `/admin/users/{id}` | Get a user by ID |
| GET | `/admin/wallets` | Get all wallets |
| GET | `/admin/wallets/{id}` | Get a wallet by ID |
| GET | `/admin/users/{id}/wallets` | Get wallets of a specific user |
| GET | `/admin/transfers` | Get all transfers |
| GET | `/admin/transfers/{id}` | Get a transfer by ID |

## Swagger / OpenAPI

The project includes Swagger API documentation.

After starting the application, Swagger UI is available at:

```text
http://localhost:8080/swagger/index.html
```

Swagger displays available endpoints, request parameters, data models, and possible HTTP responses.

For protected endpoints, use the `Authorize` button and provide the JWT through the `Authorization` header.

Generate Swagger documentation with:

```bash
swag init -g ./cmd/api/main.go
```

Generated files:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

## Database

The project uses PostgreSQL.

Main tables:

### users

Stores application users:

- ID
- email
- password hash
- role
- created timestamp
- updated timestamp

### wallets

Stores user wallets.

Each wallet belongs to a specific user.

Main fields:

- ID
- user ID
- name
- currency
- balance
- created timestamp
- updated timestamp

### transfers

Stores information about transfers between wallets:

- ID
- sender wallet ID
- receiver wallet ID
- amount
- status
- created timestamp

## SQL Transactions

Money transfers are executed inside an SQL transaction.

The sender wallet balance update, receiver wallet balance update, and transfer record creation are executed as a single operation.

If an error occurs, the transaction is rolled back.

This prevents situations where money is deducted from one wallet but not credited to another.

## Money Values

The project uses:

```text
github.com/shopspring/decimal
```

for money values.

Using `decimal.Decimal` avoids precision issues commonly associated with `float32` and `float64`.

## Migrations

Database schema changes are managed with SQL migrations.

Migration files are stored in:

```text
migrations/
```

Each migration has two files:

```text
*.up.sql
*.down.sql
```

`up` applies the database change.

`down` rolls it back.

Apply all migrations:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" up
```

Roll back the latest migration:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" down 1
```

Check the current migration version:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" version
```

## Middleware

The application uses several middleware components.

### AuthMiddleware

Responsible for JWT authentication:

- reads JWT from the `Authorization` header
- validates the token
- extracts claims
- stores claims in `context.Context`

### AdminMiddleware

Responsible for administrative access checks:

- reads claims from the context
- checks the user role
- allows access only to users with the `admin` role

### LoggingMiddleware

Logs information about each HTTP request:

- HTTP method
- URL path
- HTTP status code
- request processing duration

Example:

```text
level=INFO msg="http request" method=GET path=/wallets status=200 duration=1.5ms
```

### RecoveryMiddleware

Recovers from `panic` during HTTP request processing:

- catches panic using `recover`
- logs the error
- returns `500 Internal Server Error` to the client

## context.Context

The HTTP request context is passed through the application layers:

```text
HTTP Request
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
Database
```

The handler passes:

```go
r.Context()
```

to the service layer, and the service passes the same context to the repository.

The repository uses the context when executing SQL queries.

This ties database operations to the lifecycle of the current HTTP request.

## Logging

Structured logging is implemented using the standard Go package:

```text
log/slog
```

The application logs:

- server startup
- configuration loading errors
- database connection errors
- HTTP requests
- HTTP status codes
- request duration
- recovered panics

## Testing

The project includes unit tests for the service layer.

Repositories are replaced with fake implementations during tests.

This allows business logic to be tested independently from a real PostgreSQL database.

Run all tests:

```bash
go test ./...
```

Run tests without cached results:

```bash
go test -count=1 ./...
```

## Docker

The project includes a `Dockerfile` for the Go API and `docker-compose.yml` for running the API together with PostgreSQL.

Docker Compose starts two services:

```text
api
postgres
```

Inside the Docker network, the API connects to PostgreSQL using the service name:

```text
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

From the host machine, PostgreSQL is available through the port configured in `.env`, for example:

```text
5433:5432
```

The API is exposed as:

```text
8080:8080
```

The Go application is built using a multi-stage Docker build. The first stage compiles the application binary, while the final image contains only the compiled binary and the minimal runtime environment.

## Running with Docker Compose

### 1. Clone the repository

```bash
git clone <repository-url>
cd wallets-api-postgres
```

### 2. Create `.env`

Use `.env.example` as a template.

Example:

```env
SERVER_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=wallets_db
POSTGRES_SSLMODE=disable

JWT_SECRET=your_secret
```

The `.env` file contains sensitive data and must not be committed to Git or included in the Docker image.

### 3. Build and start containers

```bash
docker compose up --build
```

After startup, the following containers will be running:

```text
wallets-api
wallets-postgres
```

### 4. Apply migrations

After PostgreSQL starts, apply migrations:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" up
```

### 5. Check the API

Health check:

```text
http://localhost:8080/health
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

### Stop the project

```bash
docker compose down
```

To remove containers together with the PostgreSQL volume:

```bash
docker compose down -v
```

> `docker compose down -v` removes PostgreSQL data stored in the Docker volume.

## Local API Run

PostgreSQL can stay running in Docker while the Go API runs locally:

```bash
go run ./cmd/api
```

In this mode, the application uses `.env` values such as:

```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_SSLMODE=disable
```

When the API runs inside Docker Compose, database connection values are overridden:

```text
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

This allows the same Go code to run both locally and inside Docker.

## Project Structure

```text
wallets-api-postgres/
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── auth/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── response/
│   ├── router/
│   └── service/
├── migrations/
├── tests/
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── README_RU.md
```

## Main Topics Practiced

The project was built to practice and apply:

- REST APIs in Go
- `net/http`
- PostgreSQL and SQL
- layered architecture
- Repository pattern
- Service layer
- dependency injection through constructors
- JWT authentication
- role-based access control
- bcrypt
- middleware
- `context.Context`
- SQL transactions
- money handling with `decimal.Decimal`
- SQL migrations
- structured logging with `slog`
- `defer`, `panic`, and `recover`
- unit testing with fake repositories
- Swagger / OpenAPI
- Docker
- Docker Compose
- multi-stage Docker builds

## Project Goal

This project was created to practice backend development in Go and move from a simple CRUD API toward an application with separated responsibilities, authentication, authorization, business logic, transactions, tests, migrations, documentation, and containerization.