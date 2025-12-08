package sys

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtEnvSecretKey    = "HOME_JWT_SECRET"
	jwtDefaultSecret   = "acer-home-jwt-secret"
	jwtDefaultValidity = 30 * 24 * time.Hour
)

type JWTClaims struct {
	UserID   string   `json:"uid"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Category CAtegory `json:"category"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	if v := os.Getenv(jwtEnvSecretKey); v != "" {
		return []byte(v)
	}
	return []byte(jwtDefaultSecret)
}

// CreateJWT builds a signed JWT carrying basic user information.
func CreateJWT(userID, username, email string, category CAtegory, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = jwtDefaultValidity
	}

	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Category: category,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// ParseJWT validates and extracts claims from a token string.
func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// JWTDefaultExpireAt returns the default expiration time from now.
func JWTDefaultExpireAt() time.Time {
	return time.Now().Add(jwtDefaultValidity)
}
