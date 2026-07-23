# go-user-service

Go + Fiber + PostgreSQL service, built with a **hexagonal (ports & adapters) architecture**
so new modules can be added without restructuring the codebase.

## Stack
- **Framework:** [Fiber](https://gofiber.io/) (Express-style, fast HTTP framework)
- **ORM:** [GORM](https://gorm.io/) with the Postgres driver
- **Validation:** go-playground/validator
- **Password hashing:** bcrypt (golang.org/x/crypto)
- **Config:** env vars via `.env` (godotenv)

## Architecture

```
cmd/api/main.go              → composition root: wires everything together, starts the server

internal/
  config/                    → env-based configuration loader

  domain/user/                → PORTS (interfaces) + entities — the core, framework-agnostic layer
    entity.go                    - User struct, request/response DTOs
    repository.go                - Repository interface (what persistence must provide)
    service.go                   - Service interface (what business logic must provide)

  usecase/user/               → ADAPTER: implements domain/user.Service (business logic)
    usecase.go                   - create/get/update/delete rules, password hashing, pagination

  repository/postgres/        → ADAPTER: implements domain/user.Repository (persistence)
    db.go                        - GORM Postgres connection
    user_repository.go           - CRUD queries

  delivery/http/              → ADAPTER: HTTP transport (Fiber)
    handler/user_handler.go      - parses requests, calls service, formats responses
    router/router.go             - route registration, grouped by module
    response/response.go         - standard success/error JSON envelope

pkg/
  hash/                       → bcrypt helpers
  validator/                  → struct-tag validation helper

migrations/                   → versioned SQL migrations (recommended for production,
                                 AutoMigrate is used for local dev convenience)
```

**Why this shape:** the `domain` package defines *interfaces only* — it doesn't know
about Fiber or GORM. `usecase` implements the business-logic interface. `repository/postgres`
implements the persistence interface. `delivery/http` only talks to the `Service` interface.
This means you can swap Fiber for Gin, or Postgres for MySQL, by writing a new adapter —
without touching business logic.

## Adding a new module (e.g. "product")

Follow the exact same pattern as `user`:
1. `internal/domain/product/{entity,repository,service}.go` — define the entity + ports
2. `internal/usecase/product/usecase.go` — implement `Service`
3. `internal/repository/postgres/product_repository.go` — implement `Repository`
4. `internal/delivery/http/handler/product_handler.go` — HTTP handlers
5. Register routes in `router.go`, wire it up in `main.go`

## Setup

```bash
cp .env.example .env
# edit .env with your DB credentials

go mod tidy
make run
# or: go run ./cmd/api
```

### With Docker (app + Postgres)
```bash
make docker-up
```

## API

Base path: `/api/v1`

| Method | Path              | Description         |
|--------|-------------------|----------------------|
| GET    | /health            | Health check         |
| POST   | /users             | Create user           |
| GET    | /users?page=&page_size= | List users (paginated) |
| GET    | /users/:id         | Get user by ID        |
| PUT    | /users/:id         | Update user           |
| DELETE | /users/:id         | Delete user           |

### Create user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","password":"secret123","customer":"Acme Inc"}'
```

### Response format

Success:
```json
{
  "success": true,
  "message": "user created successfully",
  "data": { "id": 1, "name": "John Doe", "customer": "Acme Inc", "created_at": "...", "updated_at": "..." }
}
```

Error:
```json
{
  "success": false,
  "message": "validation failed",
  "errors": { "Password": "failed on 'min' validation" }
}
```

List (with pagination meta):
```json
{
  "success": true,
  "message": "users fetched successfully",
  "data": [ ... ],
  "meta": { "page": 1, "page_size": 10, "total_count": 42 }
}
```

## Notes
- Passwords are hashed with bcrypt before storage and never returned in responses (`json:"-"`).
- `go.sum` is not included — run `go mod tidy` after cloning to generate it and download deps.
