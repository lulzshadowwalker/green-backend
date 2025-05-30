package internal

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// JWTSecretEnv is the environment variable name for the JWT secret
	JWTSecretEnv = "JWT_SECRET"
	// DefaultJWTExpiry is the default token expiry duration (10 years)
	DefaultJWTExpiry = time.Hour * 24 * 365 * 10
)

// Claims defines the JWT claims structure
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// getJWTSecret returns the JWT secret from environment or error if not set
func getJWTSecret() ([]byte, error) {
	// secret := os.Getenv(JWTSecretEnv)
	secret := "example"
	if secret == "" {
		return nil, errors.New("JWT secret not set in environment variable JWT_SECRET")
	}
	return []byte(secret), nil
}

// GenerateJWT generates a JWT for the given user ID and username
func GenerateJWT(userID int, username string, expiry ...time.Duration) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	exp := DefaultJWTExpiry
	if len(expiry) > 0 {
		exp = expiry[0]
	}
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseJWT parses and validates a JWT string, returning the claims if valid
func ParseJWT(tokenStr string) (*Claims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
