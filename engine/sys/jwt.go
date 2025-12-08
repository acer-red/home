package sys

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JwtEnvSecretKey    = "HOME_JWT_SECRET"
	jwtDefaultValidity = 30 * 24 * time.Hour
)

type JWTClaims struct {
	UserID   string   `json:"uid"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Category CAtegory `json:"category"`
	jwt.RegisteredClaims
}

func (j *JWTClaims) GetID() string {
	return j.ID
}
func (j *JWTClaims) GetUID() string {
	return j.UserID
}
func (j *JWTClaims) GetCategory() CAtegory {
	return j.Category
}
func (j *JWTClaims) GetCategoryPrefix() CAtegory {
	return CAtegory(j.Category.GetAuthCookiePrefix())
}

func getJWTSecret() []byte {
	return []byte(os.Getenv(JwtEnvSecretKey))
}

// CreateJWT builds a signed JWT carrying basic user information.
func CreateJWT(uuid, userID, username, email string, category CAtegory, ttl time.Duration) (string, error) {
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
			ID:        uuid,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseJWT validates and extracts claims from a token string.
func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
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
