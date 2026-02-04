# AI Agent Instructions for Auth Service

## Architecture Overview

This is a **Clean Architecture** Go auth service with 4 distinct layers:
- **Domain** (`internal/domain/`): Entities (User, RefreshToken) with GORM tags and business validation
- **Repository** (`internal/repository/`): Data access with interface/implementation pattern
- **Use Case** (`internal/usecase/`): Business logic orchestration
- **Delivery** (`internal/delivery/http/`): Fiber HTTP handlers and routes

**Critical**: Never cross layer boundaries incorrectly. Domain doesn't import anything. Repository only imports domain. Use case imports domain + repository. Delivery imports everything.

## Dependency Injection Pattern

This project uses a **Container pattern** for DI (see `internal/delivery/http/container.go`):

```go
type Container struct {
    AuthUseCase usecase.AuthUseCase
    JWTManager  *jwt.JWTManager
    // Add new dependencies here
}
```

**When adding new features**:
1. Add use case to Container struct
2. Add parameter to `NewContainer()`
3. Routes signature stays `SetupRoutes(app *fiber.App, container *Container)`

This scales to 10+ use cases without parameter bloat.

## Database Migrations: Dual System

**GORM AutoMigrate** (fallback): Runs automatically on startup if Atlas not detected  
**Atlas Migrations** (production): Versioned, reviewable SQL migrations

Check in `pkg/database/atlas.go` - looks for `atlas_schema_revisions` table to decide which system to use.

### Migration Workflow
```bash
# 1. Modify domain entities (internal/domain/*.go)
# 2. Generate migration from GORM models
make atlas-generate name=add_user_role
# 3. Review generated SQL in migrations/
# 4. Apply migration
make atlas-apply
```

**Important**: Atlas reads GORM structs via `atlas-provider-gorm`. The `atlas.hcl` uses `data "external_schema" "gorm"` to introspect `internal/domain/` directory.

## Key Workflows

### Development
```bash
make dev              # Hot reload with Air (Linux/macOS)
make dev-windows      # Hot reload with Air (Windows)
```

Air watches `.go` files and rebuilds on change. Config in `.air.toml` (Unix) and `.air.windows.toml` (Windows).

### Building
```bash
make build            # Produces bin/auth-service
go build -o tmp/main cmd/main.go  # Manual build
```

### Database Setup
```bash
# First time setup
make keys             # Generate RSA keys for JWT
make db-create        # Create PostgreSQL database
make atlas-generate name=initial  # Create first migration
make atlas-apply      # Apply migrations
```

## JWT Architecture

**RS256 asymmetric encryption** with 4096-bit RSA keys in `keys/`:
- `private_key.pem`: Signs tokens (server only)
- `public_key.pem`: Verifies tokens (can be shared)

**JWKS endpoint**: `GET /.well-known/jwks.json` serves public key in JSON Web Key Set format for external verification.

**Token types**:
- Access token: Short-lived (15m default), in-memory validation
- Refresh token: Long-lived (7d default), stored in DB for revocation

Configured via `JWT_ACCESS_TOKEN_DURATION` and `JWT_REFRESH_TOKEN_DURATION` env vars.

## Project-Specific Conventions

### Error Handling
Domain errors in `internal/domain/errors.go`:
```go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidPassword  = errors.New("invalid password")
)
```

Repository layer translates `gorm.ErrRecordNotFound` → domain errors.  
HTTP layer translates domain errors → HTTP status codes.

### Repository Pattern
Interface in `internal/repository/repository.go`, implementation in `*_repository.go`:
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
```

Always pass `context.Context` as first parameter for request-scoped operations.

### Testing Strategy
No tests implemented yet. When adding:
- Unit tests: Test use cases with mocked repositories
- Integration tests: Test repositories with real DB (use testcontainers)
- E2E tests: Test HTTP handlers with Fiber test harness

## Configuration

Environment variables in `.env` (see `.env.example`):
```
PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=auth_service
DB_SSLMODE=disable
JWT_PRIVATE_KEY_PATH=./keys/private_key.pem
JWT_PUBLIC_KEY_PATH=./keys/public_key.pem
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
```

Loaded via `godotenv` in `pkg/config/config.go`.

## Common Tasks

### Adding a New Endpoint
1. Add method to use case interface + implementation (`internal/usecase/`)
2. Add DTO structs if needed (`internal/usecase/dto.go`)
3. Add handler method (`internal/delivery/http/auth_handler.go`)
4. Register route in `SetupRoutes()` (`internal/delivery/http/routes.go`)

### Adding a New Entity
1. Create struct in `internal/domain/` with GORM tags
2. Create repository interface + implementation in `internal/repository/`
3. Update `pkg/database/database.go` AutoMigrate to include new entity
4. Generate Atlas migration: `make atlas-generate name=add_entity_name`

### Debugging
- Server logs to stdout (structured with Fiber's logger middleware)
- GORM logs SQL queries when enabled (see verbose output in terminal)
- Air shows build errors and restarts automatically

## File Organization

```
cmd/main.go                          # Wire up all dependencies, start server
internal/
  domain/                            # Pure business entities (no external deps)
  repository/repository.go           # All repo interfaces
  repository/*_repository.go         # Implementations
  usecase/auth_usecase.go           # Business logic
  usecase/dto.go                    # Request/response structures
  delivery/http/container.go        # DI container
  delivery/http/routes.go           # Route registration
  delivery/http/auth_handler.go     # HTTP handlers
  delivery/http/middleware.go       # JWT auth middleware
pkg/
  config/                           # Env var loading
  database/                         # DB connection & migration detection
  jwt/                              # JWT operations, JWKS generation
docs/                               # Comprehensive documentation
```

## External Dependencies

- **Fiber v2**: Fast HTTP framework (Express-like API)
- **GORM**: ORM with auto-migrations and query building
- **golang-jwt/jwt v5**: JWT operations
- **Air**: Hot reload in development
- **Atlas**: Production-grade schema migrations
- **godotenv**: .env file loading

## When in Doubt

1. Check `Makefile` for available commands
2. Review `docs/` directory for detailed guides
3. Look at existing patterns in `internal/delivery/http/auth_handler.go` for HTTP examples
4. Check `internal/usecase/auth_usecase.go` for business logic patterns
