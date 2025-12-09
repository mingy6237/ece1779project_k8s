package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret []byte

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// InitJWT initializes JWT secret from environment variable
func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development (should be changed in production)
		secret = "your-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uuid.UUID, username, email, role string, rememberMe bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Default 24 hours
	if rememberMe {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 days if remember me
	}

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
