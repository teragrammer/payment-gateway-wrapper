package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"github.com/teragrammer/payment-gateway-wrapper/internal/models"
)

type Claims struct {
	Id   string      `json:"id"`
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const UserContextKey = contextKey("user")

func GenerateAccessToken(Id string, role models.Role) (string, error) {
	cfg := config.Load()
	accessTokenTTL := time.Duration(cfg.JWTExpiredDays) * (time.Hour * 24)

	claims := Claims{
		Id:   Id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateAccessToken(token string) (*Claims, error) {
	cfg := config.Load()
	parsed, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func ExtractClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
