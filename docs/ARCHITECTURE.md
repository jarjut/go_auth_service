# Architecture Diagram

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                               │
│                    (Web/Mobile App)                          │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   │ HTTP Requests
                   │ (Authorization: Bearer <token>)
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                    Fiber HTTP Server                         │
│                       (Port 3000)                            │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                  Delivery Layer (HTTP)                       │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Routes & Middleware                                   │ │
│  │  - CORS, Logger, Recovery                             │ │
│  │  - Auth Middleware (JWT Validation)                   │ │
│  └────────────────┬───────────────────────────────────────┘ │
│  ┌────────────────▼───────────────────────────────────────┐ │
│  │  HTTP Handlers                                         │ │
│  │  - Register, Login                                     │ │
│  │  - Refresh Token, Logout                             │ │
│  │  - Get Profile, JWKS Endpoint                        │ │
│  └────────────────┬───────────────────────────────────────┘ │
└───────────────────┼──────────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────────┐
│                  Use Case Layer                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Auth Use Case                                         │ │
│  │  - Business Logic                                      │ │
│  │  - Password Hashing (bcrypt)                          │ │
│  │  - Token Generation & Validation                      │ │
│  │  - Token Revocation Logic                            │ │
│  └────────────┬──────────────────┬───────────────────────┘ │
└───────────────┼──────────────────┼────────────────────────────┘
                │                  │
    ┌───────────▼─────────┐   ┌───▼────────────┐
    │  JWT Manager        │   │  Repositories  │
    │  (RS256)            │   │                │
    └───────────┬─────────┘   └───┬────────────┘
                │                  │
    ┌───────────▼─────────┐        │
    │  RSA Keys           │        │
    │  - private_key.pem  │        │
    │  - public_key.pem   │        │
    └─────────────────────┘        │
                                   │
┌──────────────────────────────────▼───────────────────────────┐
│                  Repository Layer                            │
│  ┌─────────────────────┐  ┌─────────────────────────────┐  │
│  │ User Repository     │  │ RefreshToken Repository     │  │
│  │ - Create User       │  │ - Create Token              │  │
│  │ - Find by Email/ID  │  │ - Find by Token             │  │
│  │ - Update, Delete    │  │ - Revoke Token              │  │
│  └──────────┬──────────┘  └──────────┬──────────────────┘  │
└─────────────┼────────────────────────┼──────────────────────┘
              │                        │
              └────────────┬───────────┘
                          │
┌─────────────────────────▼────────────────────────────────────┐
│                      GORM ORM                                │
└─────────────────────────┬────────────────────────────────────┘
                          │
┌─────────────────────────▼────────────────────────────────────┐
│                  PostgreSQL Database                         │
│  ┌────────────────┐  ┌──────────────────────────────────┐  │
│  │  users         │  │  refresh_tokens                  │  │
│  │  - id          │  │  - id                            │  │
│  │  - email       │  │  - user_id                       │  │
│  │  - password    │  │  - token                         │  │
│  │  - name        │  │  - expires_at                    │  │
│  │  - created_at  │  │  - is_revoked                    │  │
│  │  - updated_at  │  │  - created_at                    │  │
│  └────────────────┘  └──────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

## Request Flow

### 1. Registration Flow
```
Client --> POST /auth/register
    |
    v
Auth Handler --> Validate Request
    |
    v
Auth Use Case --> Check if user exists
    |                 |
    v                 v
    No           User Repo (FindByEmail)
    |
    v
Hash Password (bcrypt)
    |
    v
Create User --> User Repo (Create)
    |
    v
Generate Tokens --> JWT Manager
    |                   |
    v                   v
Store Refresh      Access Token (15m)
Token              Refresh Token (7d)
    |
    v
Refresh Token Repo (Create)
    |
    v
Response {access_token, refresh_token}
```

### 2. Login Flow
```
Client --> POST /auth/login
    |
    v
Auth Handler --> Validate Request
    |
    v
Auth Use Case --> Find User by Email
    |                 |
    v                 v
User Repo (FindByEmail)
    |
    v
Verify Password (bcrypt.Compare)
    |
    v
Valid?
    |
    v
Generate Tokens --> JWT Manager
    |
    v
Store Refresh Token
    |
    v
Response {access_token, refresh_token}
```

### 3. Protected Request Flow
```
Client --> GET /auth/profile
    |      (Authorization: Bearer <token>)
    v
Auth Middleware
    |
    v
Extract Token from Header
    |
    v
JWT Manager --> Validate Token (RS256)
    |              |
    v              v
Valid?        Verify Signature (public key)
    |
    v
Extract Claims (userID, email)
    |
    v
Store in Context
    |
    v
Continue to Handler
    |
    v
Auth Handler --> Get UserID from Context
    |
    v
Auth Use Case --> Get User by ID
    |                 |
    v                 v
User Repo (FindByID)
    |
    v
Response {user data}
```

### 4. Token Refresh Flow
```
Client --> POST /auth/refresh
    |      {refresh_token}
    v
Auth Handler --> Validate Request
    |
    v
Auth Use Case --> Find Refresh Token
    |                 |
    v                 v
Refresh Token Repo (FindByToken)
    |
    v
Check if Valid?
    |
    v
Not Revoked? Not Expired?
    |
    v
Validate JWT Signature
    |
    v
Revoke Old Token
    |
    v
Generate New Tokens
    |
    v
Store New Refresh Token
    |
    v
Response {new_access_token, new_refresh_token}
```

### 5. JWKS Endpoint Flow
```
Client --> GET /.well-known/jwks.json
    |
    v
Auth Handler --> Get JWKS
    |
    v
JWT Manager --> Extract Public Key
    |              |
    v              v
RSA Public Key (public_key.pem)
    |
    v
Convert to JWKS Format
    |
    v
Response {
  "keys": [{
    "kty": "RSA",
    "use": "sig",
    "alg": "RS256",
    "n": "...",
    "e": "..."
  }]
}
```

## Token Structure

### Access Token (Short-lived: 15 minutes)
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "iss": "auth-service",
  "sub": "1",
  "exp": 1234567890,
  "iat": 1234567000,
  "nbf": 1234567000
}
```

### Refresh Token (Long-lived: 7 days)
```json
{
  "user_id": 1,
  "iss": "auth-service",
  "sub": "1",
  "exp": 1234567890,
  "iat": 1234567000,
  "nbf": 1234567000
}
```

## Security Layers

```
┌───────────────────────────────────────────────────┐
│  1. TLS/HTTPS (Transport Layer)                   │
│     - Encrypted communication                     │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│  2. RS256 JWT Signature                           │
│     - Asymmetric encryption                       │
│     - Private key signs, public key verifies      │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│  3. Token Expiration                              │
│     - Access token: 15 minutes                    │
│     - Refresh token: 7 days                       │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│  4. Token Revocation                              │
│     - Refresh tokens stored in DB                 │
│     - Can be revoked at any time                  │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│  5. Password Hashing                              │
│     - bcrypt with default cost                    │
└───────────────────────────────────────────────────┘
```

## Directory Structure

```
auth-service/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/                      # Private application code
│   ├── domain/                    # Domain layer
│   │   ├── user.go               # User entity
│   │   ├── refresh_token.go      # RefreshToken entity
│   │   └── errors.go             # Domain errors
│   ├── repository/                # Data access layer
│   │   ├── repository.go         # Repository interfaces
│   │   ├── user_repository.go    # User repository implementation
│   │   └── refresh_token_repository.go
│   ├── usecase/                   # Business logic layer
│   │   ├── auth_usecase.go       # Auth business logic
│   │   └── dto.go                # Data transfer objects
│   └── delivery/                  # Presentation layer
│       └── http/                  # HTTP delivery
│           ├── auth_handler.go   # HTTP handlers
│           ├── middleware.go     # Auth middleware
│           └── routes.go         # Route definitions
├── pkg/                           # Public packages
│   ├── config/                    # Configuration
│   │   └── config.go
│   ├── database/                  # Database utilities
│   │   └── database.go
│   └── jwt/                       # JWT utilities
│       └── jwt.go
├── keys/                          # RSA keys (gitignored)
│   ├── private_key.pem           # Private key for signing
│   └── public_key.pem            # Public key for verification
├── scripts/                       # Utility scripts
│   └── generate_keys.sh          # Generate RSA keys
├── .env                           # Environment variables
├── .env.example                   # Example environment variables
├── go.mod                         # Go module definition
├── Dockerfile                     # Docker image definition
├── docker-compose.yml             # Docker Compose configuration
├── Makefile                       # Build automation
└── README.md                      # Documentation
```
