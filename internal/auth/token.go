package auth

import (
	"errors"
	"fmt"
	"time"

	"wallets-api-postgres/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`

	jwt.RegisteredClaims
}

func GenerateToken(user models.User, secret string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:    user.ID,
		UserEmail: user.Email,
		UserRole:  user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}
	return claims, nil
}
