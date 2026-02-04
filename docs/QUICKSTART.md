# Quick Start Guide

## 🚀 Get Started in 3 Steps

### 1. Setup Database

Make sure PostgreSQL is running, then create the database:

```bash
createdb auth_service
```

Or using psql:
```bash
psql -U postgres
CREATE DATABASE auth_service;
\q
```

### 2. Generate Keys (Already Done!)

The RSA keys have been generated in the `keys/` directory.

### 3. Run the Application

**Option A: With Hot Reload (Recommended for Development)**

```bash
# Install Air first
make install-tools

# Run with hot reload (Linux/macOS)
make dev

# Run with hot reload (Windows)
make dev-windows
```

**Option B: Without Hot Reload**

```bash
go run cmd/main.go
```

Or using Make:
```bash
make run
```

The server will start on http://localhost:3000

## 🧪 Test the API

### Register a new user:
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### Login:
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

### Get profile (protected route):
```bash
curl -X GET http://localhost:3000/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Get JWKS (public key):
```bash
curl http://localhost:3000/.well-known/jwks.json
```

## 🐳 Using Docker

### Build and run with Docker Compose:
```bash
docker-compose up -d
```

This will:
- Start PostgreSQL container
- Build and start the auth service
- Automatically run migrations

### View logs:
```bash
docker-compose logs -f
```

### Stop services:
```bash
docker-compose down
```

## 📚 More Information

- Full API documentation: See [API_EXAMPLES.md](API_EXAMPLES.md)
- Configuration options: See [README.md](README.md)
- Available commands: Run `make help`

## 🔧 Useful Commands

```bash
# View all available commands
make help

# Development with hot reload (recommended)
make dev          # Linux/Mac
make dev-windows  # Windows

# Build the application
make build

# Run tests
make test

# Generate new keys
make keys

# Format code
make fmt

# Run without hot reload
make run
```

## 📝 Environment Variables

The application is configured via `.env` file. Key settings:

- `JWT_ACCESS_TOKEN_DURATION`: Access token lifetime (default: 15m)
- `JWT_REFRESH_TOKEN_DURATION`: Refresh token lifetime (default: 168h)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`: Database connection
- `PORT`: Server port (default: 3000)

## 🛡️ Security Notes

1. **Never commit** `keys/private_key.pem` to version control
2. **Change** default database credentials in production
3. **Use HTTPS** in production
4. **Adjust** token durations based on your security requirements
5. **Rotate** keys periodically in production

## 🐛 Troubleshooting

### Database connection failed
- Ensure PostgreSQL is running: `pg_isready`
- Check database exists: `psql -l | grep auth_service`
- Verify credentials in `.env`

### Keys not found
- Run: `make keys` or `./scripts/generate_keys.sh`
- Ensure `keys/` directory exists with .pem files

### Port already in use
- Change `PORT` in `.env`
- Or kill the process: `lsof -ti:3000 | xargs kill -9`

## 🎯 Next Steps

1. ✅ Set up database
2. ✅ Generate keys
3. ✅ Run the application
4. ✅ Test the API
5. 📖 Read the full documentation
6. 🔨 Customize for your needs
7. 🚀 Deploy to production

Happy coding! 🎉
