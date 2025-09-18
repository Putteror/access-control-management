package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/schema"
)

// AuthService handles authentication logic and JWT token management.
type AuthService interface {
	Login(username, password string) (string, error)
}

type authServiceImpl struct {
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() AuthService {
	return &authServiceImpl{}
}

// Login authenticates a user and generates a JWT token.
func (s *authServiceImpl) Login(username, password string) (string, error) {
	// 1. Authenticate user credentials.
	// This is a mock authentication. In a real app, you would query your database.
	if username != "admin" || password != "password123" {
		return "", errors.New("invalid credentials")
	}

	// 2. Define custom claims for the JWT.
	// You can include any data you want to store in the token (e.g., user ID, roles).
	claims := &schema.CustomClaims{
		Username: username,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 3. Create a new token with the defined claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Sign the token with the secret key.
	signedToken, err := token.SignedString(common.JwtSecret)
	if err != nil {
		return "", errors.New("could not sign the token")
	}

	return signedToken, nil
}
