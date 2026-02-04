# Auth Service with Fiber and GORM

A robust authentication service built with Go, Fiber framework, and GORM ORM, implementing JWT authentication with RS256 algorithm and Clean Architecture principles.

## 📚 Documentation

- **[Quick Start Guide](docs/QUICKSTART.md)** - Get started in 3 steps
- **[Architecture Overview](docs/ARCHITECTURE.md)** - System design and diagrams
- **[Development Guide](docs/DEVELOPMENT.md)** - Hot reload setup with Air
- **[Migration Guide](docs/MIGRATIONS.md)** - Database migrations with Atlas
- **[API Examples](docs/API_EXAMPLES.md)** - Usage examples and curl commands
- **[Testing Guide](docs/TESTING.md)** - Comprehensive testing scenarios
- **[Project Summary](docs/PROJECT_SUMMARY.md)** - Complete project overview
- **[File Index](docs/INDEX.md)** - Navigate the project structure

## Features

- 🔐 **JWT Authentication with RS256**: Asymmetric encryption using private/public keys
- 🔑 **JWKS Endpoint**: Public key served in JSON Web Key Set format
- ⏱️ **Token Management**: 
  - Short-lived access tokens (configurable, default 15 minutes)
  - Long-lived refresh tokens (configurable, default 7 days)
  - Refresh token revocation support
- 🗃️ **PostgreSQL Database**: Using GORM for database operations
- 🏗️ **Clean Architecture**: Maintainable and testable code structure
- 🚀 **High Performance**: Built with Fiber framework

## Architecture

The project follows Clean Architecture principles:

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/                 # Domain entities and errors
│   │   ├── user.go
│   │   ├── refresh_token.go
│   │   └── errors.go
│   ├── repository/             # Data access layer
│   │   ├── repository.go
│   │   ├── user_repository.go
│   │   └── refresh_token_repository.go
│   ├── usecase/                # Business logic layer
│   │   ├── auth_usecase.go
│   │   └── dto.go
│   └── delivery/
│       └── http/               # HTTP handlers (presentation layer)
│           ├── auth_handler.go
│           ├── middleware.go
│           └── routes.go
├── pkg/
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── database/               # Database connection
│   │   └── database.go
│   └── jwt/                    # JWT utilities
│       └── jwt.go
└── keys/                       # RSA keys (gitignored)
```

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- OpenSSL (for generating RSA keys)
- Air (optional, for hot reload during development)
- Atlas (optional, for database migrations)

## Setup

### 1. Generate RSA Keys

Run the provided script to generate private and public keys:

```bash
./scripts/generate_keys.sh
```

Or manually with OpenSSL:

```bash
# Generate private key
openssl genrsa -out keys/private_key.pem 4096

# Generate public key
openssl rsa -in keys/private_key.pem -pubout -out keys/public_key.pem
```

### 2. Configure Environment

Copy the example environment file and adjust the values:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_service
DB_SSLMODE=disable

# JWT Configuration
JWT_PRIVATE_KEY_PATH=./keys/private_key.pem
JWT_PUBLIC_KEY_PATH=./keys/public_key.pem
JWT_ACCESS_TOKEN_DURATION=15m      # 15 minutes
JWT_REFRESH_TOKEN_DURATION=168h    # 7 days

# Application Configuration
APP_ENV=development
```

### 3. Setup Database

Create the PostgreSQL database:

```sql
CREATE DATABASE auth_service;
```

The application will automatically run migrations on startup.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```

#### Get JWKS (Public Key)
```
GET /.well-known/jwks.json
```

#### Register
```
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}
```

#### Login
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### Refresh Token
```
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

#### Logout
```
POST /auth/logout
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

### Protected Endpoints

These endpoints require a valid access token in the Authorization header:

#### Get Profile
```
GET /auth/profile
Authorization: Bearer <access_token>
```

#### Logout from All Devices
```
POST /auth/logout-all
Authorization: Bearer <access_token>
```

## Token Configuration

### Access Token
- **Default Duration**: 15 minutes
- **Purpose**: Short-lived token for API authentication
- **Storage**: Client-side (memory or secure storage)
- **Configurable via**: `JWT_ACCESS_TOKEN_DURATION` env variable

### Refresh Token
- **Default Duration**: 7 days (168 hours)
- **Purpose**: Long-lived token for obtaining new access tokens
- **Storage**: Database (can be revoked)
- **Configurable via**: `JWT_REFRESH_TOKEN_DURATION` env variable

### Token Duration Examples

```env
# 5 minutes
JWT_ACCESS_TOKEN_DURATION=5m

# 15 minutes
JWT_ACCESS_TOKEN_DURATION=15m

# 1 hour
JWT_ACCESS_TOKEN_DURATION=1h

# 7 days
JWT_REFRESH_TOKEN_DURATION=168h

# 30 days
JWT_REFRESH_TOKEN_DURATION=720h
```

## Security Features

1. **RS256 Algorithm**: Asymmetric encryption with RSA keys
2. **Password Hashing**: Using bcrypt with default cost
3. **Token Revocation**: Refresh tokens can be revoked
4. **JWKS Support**: Public key available for token validation
5. **Token Expiration**: Both access and refresh tokens expire
6. **Logout Support**: Single device and all devices logout

## Development

### Run with Hot Reload (Recommended)

Install Air for automatic rebuild on file changes:

```bash
# Install Air
make install-tools

# Run with hot reload (Linux/macOS)
make dev

# Run with hot reload (Windows)
make dev-windows
```

See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed Air setup and usage.

### Run without Hot Reload

```bash
go run cmd/main.go
```

### Run Tests

```bash
go test ./...
```

### Build for Production

```bash
go build -o bin/auth-service cmd/main.go
```

## Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/keys ./keys
EXPOSE 3000
CMD ["./main"]
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
