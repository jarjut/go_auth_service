# Testing Guide

## Prerequisites

Make sure you have:
1. PostgreSQL running
2. Database created: `createdb auth_service`
3. Service running: `go run cmd/main.go`

## Environment Setup

```bash
export API_URL="http://localhost:3000"
```

## Test Scenarios

### Scenario 1: Complete User Journey

```bash
#!/bin/bash
set -e

API_URL="http://localhost:3000"

echo "=== Test 1: Health Check ==="
curl -s "$API_URL/health" | jq .

echo -e "\n=== Test 2: Register New User ==="
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "SecurePass123!",
    "name": "Alice Johnson"
  }')
echo "$REGISTER_RESPONSE" | jq .

# Extract tokens
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refresh_token')

echo -e "\n=== Test 3: Get Profile (Protected) ==="
curl -s "$API_URL/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

echo -e "\n=== Test 4: Login with Same User ==="
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "SecurePass123!"
  }')
echo "$LOGIN_RESPONSE" | jq .

NEW_REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token')

echo -e "\n=== Test 5: Refresh Token ==="
REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$NEW_REFRESH_TOKEN\"}")
echo "$REFRESH_RESPONSE" | jq .

FINAL_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.refresh_token')

echo -e "\n=== Test 6: Logout ==="
curl -s -X POST "$API_URL/auth/logout" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$FINAL_REFRESH_TOKEN\"}" | jq .

echo -e "\n=== Test 7: Try to Use Revoked Token (Should Fail) ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$FINAL_REFRESH_TOKEN\"}" | jq .

echo -e "\n=== Test 8: Get JWKS ==="
curl -s "$API_URL/.well-known/jwks.json" | jq .

echo -e "\n✅ All tests completed!"
```

Save as `test_complete_flow.sh` and run:
```bash
chmod +x test_complete_flow.sh
./test_complete_flow.sh
```

### Scenario 2: Error Handling Tests

```bash
#!/bin/bash

API_URL="http://localhost:3000"

echo "=== Test: Register with Existing Email ==="
curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123",
    "name": "Another Alice"
  }' | jq .
# Expected: 409 Conflict

echo -e "\n=== Test: Register with Short Password ==="
curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "password": "short",
    "name": "Bob"
  }' | jq .
# Expected: 400 Bad Request

echo -e "\n=== Test: Login with Wrong Password ==="
curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "WrongPassword"
  }' | jq .
# Expected: 401 Unauthorized

echo -e "\n=== Test: Login with Non-existent Email ==="
curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "password123"
  }' | jq .
# Expected: 401 Unauthorized

echo -e "\n=== Test: Access Protected Route Without Token ==="
curl -s "$API_URL/auth/profile" | jq .
# Expected: 401 Unauthorized

echo -e "\n=== Test: Access Protected Route With Invalid Token ==="
curl -s "$API_URL/auth/profile" \
  -H "Authorization: Bearer invalid.token.here" | jq .
# Expected: 401 Unauthorized

echo -e "\n=== Test: Refresh with Invalid Token ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "invalid.token.here"}' | jq .
# Expected: 401 Unauthorized
```

### Scenario 3: Multiple Device Logout

```bash
#!/bin/bash

API_URL="http://localhost:3000"

echo "=== Setup: Register User ==="
RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "multidevice@example.com",
    "password": "password123",
    "name": "Multi Device User"
  }')

ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN_1=$(echo "$RESPONSE" | jq -r '.refresh_token')

echo -e "\n=== Simulate Device 2: Login Again ==="
RESPONSE_2=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "multidevice@example.com",
    "password": "password123"
  }')

REFRESH_TOKEN_2=$(echo "$RESPONSE_2" | jq -r '.refresh_token')

echo -e "\n=== Simulate Device 3: Login Again ==="
RESPONSE_3=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "multidevice@example.com",
    "password": "password123"
  }')

REFRESH_TOKEN_3=$(echo "$RESPONSE_3" | jq -r '.refresh_token')

echo -e "\n=== Now we have 3 refresh tokens from 3 'devices' ==="
echo "Token 1: ${REFRESH_TOKEN_1:0:20}..."
echo "Token 2: ${REFRESH_TOKEN_2:0:20}..."
echo "Token 3: ${REFRESH_TOKEN_3:0:20}..."

echo -e "\n=== Logout from ALL devices ==="
curl -s -X POST "$API_URL/auth/logout-all" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

echo -e "\n=== Try to use Token 1 (Should Fail) ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN_1\"}" | jq .

echo -e "\n=== Try to use Token 2 (Should Fail) ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN_2\"}" | jq .

echo -e "\n=== Try to use Token 3 (Should Fail) ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN_3\"}" | jq .

echo -e "\n✅ All tokens revoked successfully!"
```

### Scenario 4: JWKS Validation

```bash
#!/bin/bash

API_URL="http://localhost:3000"

echo "=== Fetch JWKS ==="
JWKS=$(curl -s "$API_URL/.well-known/jwks.json")
echo "$JWKS" | jq .

echo -e "\n=== Verify JWKS Structure ==="
echo "Number of keys: $(echo "$JWKS" | jq '.keys | length')"
echo "Key type: $(echo "$JWKS" | jq -r '.keys[0].kty')"
echo "Algorithm: $(echo "$JWKS" | jq -r '.keys[0].alg')"
echo "Use: $(echo "$JWKS" | jq -r '.keys[0].use')"

echo -e "\n=== Get Token ==="
RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "SecurePass123!"
  }')

TOKEN=$(echo "$RESPONSE" | jq -r '.access_token')

echo -e "\n=== Token Structure ==="
echo "Header:  $(echo $TOKEN | cut -d'.' -f1 | base64 -d 2>/dev/null | jq .)"
echo "Payload: $(echo $TOKEN | cut -d'.' -f2 | base64 -d 2>/dev/null | jq .)"
echo "Signature: $(echo $TOKEN | cut -d'.' -f3)"
```

### Scenario 5: Performance Test

```bash
#!/bin/bash

API_URL="http://localhost:3000"

echo "=== Performance Test: 100 Registrations ==="
time for i in {1..100}; do
  curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"user$i@example.com\",
      \"password\": \"password123\",
      \"name\": \"User $i\"
    }" > /dev/null
done

echo -e "\n=== Performance Test: 100 Logins ==="
time for i in {1..100}; do
  curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"user$i@example.com\",
      \"password\": \"password123\"
    }" > /dev/null
done
```

## Unit Testing with Go

Create `internal/usecase/auth_usecase_test.go`:

```go
package usecase

import (
	"auth-service/internal/domain"
	"context"
	"testing"
)

// Mock repositories would go here
// Example test structure

func TestRegister_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	
	// Mock repositories
	// mockUserRepo := &MockUserRepository{}
	// mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	// mockJWTManager := &MockJWTManager{}
	
	// Create use case
	// useCase := NewAuthUseCase(mockUserRepo, mockRefreshTokenRepo, mockJWTManager)
	
	// Test data
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	
	// Execute
	// resp, err := useCase.Register(ctx, req)
	
	// Assert
	// assert.NoError(t, err)
	// assert.NotNil(t, resp)
	// assert.NotEmpty(t, resp.AccessToken)
	// assert.NotEmpty(t, resp.RefreshToken)
}
```

## Database Testing

```bash
#!/bin/bash

# Connect to database
psql -U postgres -d auth_service

# Check tables
\dt

# Check users
SELECT id, email, name, created_at FROM users;

# Check refresh tokens
SELECT id, user_id, expires_at, is_revoked, created_at FROM refresh_tokens;

# Check expired tokens
SELECT COUNT(*) FROM refresh_tokens WHERE expires_at < NOW();

# Check revoked tokens
SELECT COUNT(*) FROM refresh_tokens WHERE is_revoked = true;

# Exit
\q
```

## Load Testing with Apache Bench

```bash
# Install ab (Apache Bench)
sudo apt-get install apache2-utils  # Ubuntu/Debian
brew install apache2                 # macOS

# Test health endpoint
ab -n 1000 -c 10 http://localhost:3000/health

# Test login endpoint (requires file with POST data)
echo '{"email":"alice@example.com","password":"SecurePass123!"}' > login.json
ab -n 100 -c 10 -T 'application/json' -p login.json http://localhost:3000/auth/login
```

## Expected Results

### Successful Registration
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### Successful Profile Fetch
```json
{
  "id": 1,
  "email": "alice@example.com",
  "name": "Alice Johnson"
}
```

### Error: Invalid Credentials
```json
{
  "error": "invalid credentials"
}
```

### Error: Token Expired
```json
{
  "error": "invalid or expired token"
}
```

## Monitoring

### Check Logs
```bash
# If running with go run
# Logs appear in terminal

# If running with Docker
docker-compose logs -f auth-service

# If running as systemd service
journalctl -u auth-service -f
```

### Database Monitoring
```sql
-- Active tokens by user
SELECT u.email, COUNT(rt.id) as token_count
FROM users u
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.is_revoked = false AND rt.expires_at > NOW()
GROUP BY u.email;

-- Token statistics
SELECT 
  COUNT(*) as total_tokens,
  COUNT(*) FILTER (WHERE is_revoked = true) as revoked_tokens,
  COUNT(*) FILTER (WHERE expires_at < NOW()) as expired_tokens,
  COUNT(*) FILTER (WHERE is_revoked = false AND expires_at > NOW()) as active_tokens
FROM refresh_tokens;
```

## Cleanup

```bash
# Delete test users
psql -U postgres -d auth_service -c "DELETE FROM users WHERE email LIKE '%@example.com';"

# Delete all refresh tokens
psql -U postgres -d auth_service -c "DELETE FROM refresh_tokens;"

# Reset database
make db-reset
```

## Troubleshooting

### Issue: Token validation fails
**Check**: 
- Keys are correctly generated
- Keys path in .env is correct
- Keys have correct permissions

### Issue: Database connection fails
**Check**:
- PostgreSQL is running: `pg_isready`
- Database exists: `psql -l | grep auth_service`
- Credentials in .env are correct

### Issue: Port already in use
**Solution**:
```bash
# Find process
lsof -i :3000

# Kill process
kill -9 <PID>

# Or change port in .env
```

## Continuous Integration

Example GitHub Actions workflow:

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: auth_service
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Generate Keys
        run: ./scripts/generate_keys.sh
      
      - name: Test
        run: go test -v ./...
        env:
          DB_HOST: localhost
          DB_USER: postgres
          DB_PASSWORD: postgres
```

---

Save all test scripts in a `tests/` directory for easy access!
