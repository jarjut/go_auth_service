# Project File Index

## 📚 Documentation Files (Start Here!)

| File | Description | When to Read |
|------|-------------|--------------|
| **README.md** | Complete project documentation | First read - comprehensive overview |
| **QUICKSTART.md** | Quick start guide (3 steps) | When you want to run it quickly |
| **PROJECT_SUMMARY.md** | Project summary and checklist | High-level overview |
| **ARCHITECTURE.md** | Architecture diagrams and flows | To understand the design |
| **API_EXAMPLES.md** | API usage examples with curl | When testing the API |
| **TESTING.md** | Complete testing guide | When writing tests |
| **THIS FILE** | Navigation guide | You are here! |

## 🔧 Configuration Files

| File | Description |
|------|-------------|
| `.env` | Environment variables (active config) |
| `.env.example` | Example environment variables |
| `go.mod` | Go module dependencies |
| `go.sum` | Go module checksums |
| `.gitignore` | Git ignore rules |

## 🐳 Deployment Files

| File | Description |
|------|-------------|
| `Dockerfile` | Docker image definition |
| `docker-compose.yml` | Docker Compose setup with PostgreSQL |
| `Makefile` | Build automation commands |

## 💻 Source Code Structure

### cmd/ - Application Entry Point
```
cmd/
└── main.go              # Main application file
                         # - Loads configuration
                         # - Connects to database
                         # - Initializes dependencies
                         # - Starts HTTP server
```

### internal/ - Private Application Code

#### internal/domain/ - Domain Layer (Entities & Business Rules)
```
internal/domain/
├── user.go              # User entity with GORM model
├── refresh_token.go     # RefreshToken entity with validation
└── errors.go            # Domain-specific errors
```

**Key Entities:**
- `User`: id, email, password, name, timestamps
- `RefreshToken`: id, user_id, token, expires_at, is_revoked, timestamps

#### internal/repository/ - Data Access Layer
```
internal/repository/
├── repository.go                    # Repository interfaces
├── user_repository.go               # User data access (CRUD)
└── refresh_token_repository.go      # Token data access (CRUD + Revoke)
```

**Key Methods:**
- UserRepository: Create, FindByID, FindByEmail, Update, Delete
- RefreshTokenRepository: Create, FindByToken, Revoke, RevokeAllByUserID

#### internal/usecase/ - Business Logic Layer
```
internal/usecase/
├── auth_usecase.go      # Authentication business logic
└── dto.go               # Data Transfer Objects (Request/Response)
```

**Key Use Cases:**
- Register: Create new user + generate tokens
- Login: Verify credentials + generate tokens
- RefreshToken: Validate + generate new tokens
- Logout: Revoke single token
- LogoutAll: Revoke all user tokens

#### internal/delivery/http/ - HTTP Layer (Presentation)
```
internal/delivery/http/
├── auth_handler.go      # HTTP handlers for auth endpoints
├── middleware.go        # Authentication middleware
└── routes.go            # Route definitions
```

**Endpoints:**
- Public: /health, /auth/register, /auth/login, /auth/refresh, /auth/logout
- Protected: /auth/profile, /auth/logout-all
- JWKS: /.well-known/jwks.json

### pkg/ - Public Packages (Reusable)

#### pkg/config/ - Configuration Management
```
pkg/config/
└── config.go            # Load environment variables
                         # - Server config (port)
                         # - Database config (connection)
                         # - JWT config (keys, durations)
```

#### pkg/database/ - Database Utilities
```
pkg/database/
└── database.go          # Database connection & migrations
                         # - Connect to PostgreSQL
                         # - Run GORM migrations
```

#### pkg/jwt/ - JWT Utilities (Core Security)
```
pkg/jwt/
└── jwt.go               # JWT Manager with RS256
                         # - Generate access tokens
                         # - Generate refresh tokens
                         # - Validate tokens
                         # - Provide JWKS
                         # - Load RSA keys
```

### keys/ - RSA Keys (Security)
```
keys/
├── .gitkeep             # Keep directory in git
├── private_key.pem      # Private key (GITIGNORED)
└── public_key.pem       # Public key (GITIGNORED)
```

⚠️ **Important**: Keys are generated, not committed!

### scripts/ - Utility Scripts
```
scripts/
└── generate_keys.sh     # Generate RSA key pair
                         # - Creates 4096-bit keys
                         # - Sets proper permissions
```

### bin/ - Compiled Binaries (Generated)
```
bin/
└── auth-service         # Compiled binary (18MB)
```

## 🎯 Quick Navigation Guide

### I want to...

**...understand the overall architecture**
→ Read [ARCHITECTURE.md](ARCHITECTURE.md)

**...run the service quickly**
→ Follow [QUICKSTART.md](QUICKSTART.md)

**...test the API**
→ Use [API_EXAMPLES.md](API_EXAMPLES.md)

**...add a new endpoint**
1. Add handler method in [internal/delivery/http/auth_handler.go](internal/delivery/http/auth_handler.go)
2. Add route in [internal/delivery/http/routes.go](internal/delivery/http/routes.go)
3. Add use case method in [internal/usecase/auth_usecase.go](internal/usecase/auth_usecase.go)

**...add a new entity**
1. Create entity in [internal/domain/](internal/domain/)
2. Create repository interface in [internal/repository/repository.go](internal/repository/repository.go)
3. Implement repository in [internal/repository/](internal/repository/)

**...modify JWT logic**
→ Edit [pkg/jwt/jwt.go](pkg/jwt/jwt.go)

**...change database schema**
→ Modify entities in [internal/domain/](internal/domain/) and restart (auto-migration)

**...configure the application**
→ Edit [.env](.env) file

**...deploy with Docker**
→ Run `docker-compose up -d`

**...write tests**
→ Follow [TESTING.md](TESTING.md)

**...change token durations**
→ Update `JWT_ACCESS_TOKEN_DURATION` and `JWT_REFRESH_TOKEN_DURATION` in [.env](.env)

## 📊 File Statistics

| Category | Count |
|----------|-------|
| Go source files | 12 |
| Documentation files | 7 |
| Configuration files | 6 |
| Scripts | 1 |
| **Total project files** | **26** |

## 🔄 Request Flow Through Files

### Registration Flow
```
1. Client Request
   ↓
2. cmd/main.go → Fiber Server
   ↓
3. internal/delivery/http/routes.go → Route to handler
   ↓
4. internal/delivery/http/auth_handler.go → Register handler
   ↓
5. internal/usecase/auth_usecase.go → Register use case
   ↓
6. internal/repository/user_repository.go → Save user
   ↓
7. pkg/jwt/jwt.go → Generate tokens
   ↓
8. internal/repository/refresh_token_repository.go → Save token
   ↓
9. Response with tokens
```

### Protected Request Flow
```
1. Client Request (with Authorization header)
   ↓
2. internal/delivery/http/middleware.go → Auth middleware
   ↓
3. pkg/jwt/jwt.go → Validate token
   ↓
4. internal/delivery/http/auth_handler.go → Handler
   ↓
5. internal/usecase/auth_usecase.go → Use case
   ↓
6. internal/repository/user_repository.go → Get data
   ↓
7. Response
```

## 🎓 Learning Path

### Beginner
1. Read [QUICKSTART.md](QUICKSTART.md)
2. Run the application
3. Test with [API_EXAMPLES.md](API_EXAMPLES.md)
4. Explore [cmd/main.go](cmd/main.go)

### Intermediate
1. Read [ARCHITECTURE.md](ARCHITECTURE.md)
2. Study domain entities in [internal/domain/](internal/domain/)
3. Understand use cases in [internal/usecase/](internal/usecase/)
4. Review JWT implementation in [pkg/jwt/jwt.go](pkg/jwt/jwt.go)

### Advanced
1. Read all documentation
2. Study Clean Architecture layers
3. Implement custom features
4. Write comprehensive tests
5. Deploy to production

## 🔍 Search Hints

| Looking for | File/Directory |
|-------------|----------------|
| User model | internal/domain/user.go |
| Token model | internal/domain/refresh_token.go |
| Register logic | internal/usecase/auth_usecase.go (Register method) |
| Login logic | internal/usecase/auth_usecase.go (Login method) |
| Token generation | pkg/jwt/jwt.go (Generate methods) |
| Token validation | pkg/jwt/jwt.go (ValidateToken method) |
| JWKS | pkg/jwt/jwt.go (GetJWKS method) |
| Database connection | pkg/database/database.go |
| Configuration | pkg/config/config.go |
| Routes | internal/delivery/http/routes.go |
| Middleware | internal/delivery/http/middleware.go |
| Error definitions | internal/domain/errors.go |

## 🚀 Common Commands

```bash
# View this index
cat INDEX.md

# Run the service
make run

# Build the service
make build

# Run tests
make test

# Generate keys
make keys

# View all commands
make help

# Check file locations
find . -name "*.go"
```

## 📞 Quick Reference Card

**Project Name**: Auth Service with Fiber and GORM  
**Language**: Go 1.21+  
**Framework**: Fiber v2  
**ORM**: GORM v1.25  
**Database**: PostgreSQL  
**JWT Algorithm**: RS256 (Asymmetric)  
**Architecture**: Clean Architecture  
**Port**: 3000 (configurable)  
**Access Token**: 15 minutes (configurable)  
**Refresh Token**: 7 days (configurable)  

**Binary Size**: 18MB  
**Total Lines of Code**: ~1,500  
**Documentation Pages**: 7  
**API Endpoints**: 8  

---

**Last Updated**: February 4, 2026  
**Status**: ✅ Production Ready  
**License**: MIT
