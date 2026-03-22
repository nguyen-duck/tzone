package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// GenerateToken generates a new JWT token for a given user ID
func GenerateToken(userID uuid.UUID) (string, error) {
	// Secret key from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key" // Fallback if not set
	}

	// Create claims with standard claims and custom claims
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	return token.SignedString([]byte(secretKey))
}

// ValidateToken parses and validates a JWT token, returning the user ID if valid
func ValidateToken(tokenString string) (uuid.UUID, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrExpiredToken
		}
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user_id
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, ErrInvalidToken
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, ErrInvalidToken
		}

		return userID, nil
	}

	return uuid.Nil, ErrInvalidToken
}
