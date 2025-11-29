package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretEnvVar = "ZGATE_JWT_SECRET"
	// Token durations
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

var (
	jSecret     []byte
	jSecretOnce sync.Once
)

func getJWTSecret() []byte {
	jSecretOnce.Do(func() {
		secret := os.Getenv(jwtSecretEnvVar)
		if secret == "" {
			panic("ZGATE_JWT_SECRET must be set")
		}
		jSecret = []byte(secret)
	})
	return jSecret
}

// Claims represents JWT claims (minimal - only username for identification)
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT access token for a user with JTI for revocation tracking
func GenerateToken(user *UserWithPermissions) (string, time.Time, error) {
	expiresAt := time.Now().Add(AccessTokenDuration)

	// Generate unique JWT ID (JTI) for revocation tracking
	jti, err := generateJTI()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate JTI: %w", err)
	}

	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zGate",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// generateJTI generates a unique JWT ID for revocation tracking
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRefreshToken generates a cryptographically secure refresh token and its hash
// Returns: (token, tokenHash, error)
func GenerateRefreshToken() (string, string, error) {
	// Generate 32 bytes of random data
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode as hex string
	token := hex.EncodeToString(tokenBytes)

	// Create SHA-256 hash for storage
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}
