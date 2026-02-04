package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager handles JWT operations with RS256
type JWTManager struct {
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(privateKeyPath, publicKeyPath string, accessTokenDuration, refreshTokenDuration time.Duration) (*JWTManager, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	return &JWTManager{
		privateKey:           privateKey,
		publicKey:            publicKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}, nil
}

// GenerateAccessToken generates a new access token
func (m *JWTManager) GenerateAccessToken(userID string, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

// GenerateRefreshToken generates a cryptographically secure random refresh token
// Returns the token string and expiration time
func (m *JWTManager) GenerateRefreshToken(userID string) (string, time.Time, error) {
	// Generate 32 bytes (256 bits) of random data
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode to hex string for easy storage and transmission
	tokenString := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(m.refreshTokenDuration)

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT access token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetPublicKey returns the public key
func (m *JWTManager) GetPublicKey() *rsa.PublicKey {
	return m.publicKey
}

// GetJWKS returns the JSON Web Key Set
func (m *JWTManager) GetJWKS() (map[string]interface{}, error) {
	// Get the modulus and exponent from the public key
	n := m.publicKey.N
	e := m.publicKey.E

	// Convert to base64 URL encoding
	nBytes := n.Bytes()
	nBase64 := base64.RawURLEncoding.EncodeToString(nBytes)

	eBytes := big.NewInt(int64(e)).Bytes()
	eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)

	// Create JWKS structure
	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"n":   nBase64,
				"e":   eBase64,
			},
		},
	}

	return jwks, nil
}

// GetJWKSJSON returns the JWKS as JSON string
func (m *JWTManager) GetJWKSJSON() (string, error) {
	jwks, err := m.GetJWKS()
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jwks)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// GetAccessTokenDuration returns the access token duration
func (m *JWTManager) GetAccessTokenDuration() time.Duration {
	return m.accessTokenDuration
}

// loadPrivateKey loads RSA private key from PEM file
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
	}

	return privateKey, nil
}

// loadPublicKey loads RSA public key from PEM file
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		// Try PKIX format
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		var ok bool
		publicKey, ok = key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
	}

	return publicKey, nil
}
