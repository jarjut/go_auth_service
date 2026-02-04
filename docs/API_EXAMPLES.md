# Auth Service - API Examples

## Environment
```bash
export API_URL="http://localhost:3000"
```

## 1. Health Check
```bash
curl -X GET "$API_URL/health"
```

## 2. Get JWKS (Public Key)
```bash
curl -X GET "$API_URL/.well-known/jwks.json"
```

## 3. Register a New User
```bash
curl -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "name": "John Doe"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

## 4. Login
```bash
curl -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

## 5. Get User Profile (Protected)
```bash
# Save the access token from login/register response
export ACCESS_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET "$API_URL/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Response:
```json
{
  "id": 1,
  "email": "john.doe@example.com",
  "name": "John Doe"
}
```

## 6. Refresh Access Token
```bash
# Save the refresh token from login/register response
export REFRESH_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

Response:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

## 7. Logout (Revoke Refresh Token)
```bash
curl -X POST "$API_URL/auth/logout" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

Response:
```json
{
  "message": "successfully logged out"
}
```

## 8. Logout from All Devices (Protected)
```bash
curl -X POST "$API_URL/auth/logout-all" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Response:
```json
{
  "message": "successfully logged out from all devices"
}
```

## Complete Flow Example

```bash
#!/bin/bash

API_URL="http://localhost:3000"

echo "=== 1. Health Check ==="
curl -X GET "$API_URL/health"
echo -e "\n"

echo "=== 2. Register User ==="
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!",
    "name": "Test User"
  }')
echo $REGISTER_RESPONSE | jq .
echo -e "\n"

# Extract tokens
ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.refresh_token')

echo "=== 3. Get Profile ==="
curl -s -X GET "$API_URL/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo -e "\n"

echo "=== 4. Wait for token to expire (if needed) ==="
# In real scenario, wait for access token to expire
# sleep 901

echo "=== 5. Refresh Token ==="
REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo $REFRESH_RESPONSE | jq .
echo -e "\n"

# Get new tokens
NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.access_token')
NEW_REFRESH_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.refresh_token')

echo "=== 6. Logout ==="
curl -s -X POST "$API_URL/auth/logout" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$NEW_REFRESH_TOKEN\"
  }" | jq .
echo -e "\n"

echo "=== 7. Try to use revoked refresh token (should fail) ==="
curl -s -X POST "$API_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$NEW_REFRESH_TOKEN\"
  }" | jq .
echo -e "\n"
```

## Testing with Different Token Durations

### Short Access Token (5 minutes)
```bash
# Update .env
JWT_ACCESS_TOKEN_DURATION=5m

# Restart service
```

### Very Short Access Token (1 minute for testing)
```bash
# Update .env
JWT_ACCESS_TOKEN_DURATION=1m

# Restart service
```

### Extended Refresh Token (30 days)
```bash
# Update .env
JWT_REFRESH_TOKEN_DURATION=720h

# Restart service
```

## Postman Collection

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Auth Service",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:3000"
    }
  ],
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/health"
      }
    },
    {
      "name": "Register",
      "request": {
        "method": "POST",
        "url": "{{baseUrl}}/auth/register",
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"user@example.com\",\n  \"password\": \"password123\",\n  \"name\": \"Test User\"\n}"
        }
      }
    },
    {
      "name": "Login",
      "request": {
        "method": "POST",
        "url": "{{baseUrl}}/auth/login",
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"user@example.com\",\n  \"password\": \"password123\"\n}"
        }
      }
    }
  ]
}
```
