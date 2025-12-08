package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/acer-red/home/engine/sys"
	"github.com/redis/go-redis/v9"
	log "github.com/tengfei-xy/go-log"
)

var errRedisNotInit = errors.New("redis not initialized")

type loginSession struct {
	UserID   string       `json:"uid"`
	Category sys.CAtegory `json:"category"`
	Expires  int64        `json:"expires"`
}

// buildSessionKey returns the canonical redis key layout for a JWT token payload.
// Format: <prefix>user:<userID>:jwt:<jti> where prefix already ends with a colon.
func buildSessionKey(prefix, userID, jti string) string {
	return fmt.Sprintf("%suser:%s:jwt:%s", prefix, userID, jti)
}

// buildLookupKey stores the latest session key pointer for a user.
// Format: <prefix>user:<userID>:latest
func buildLookupKey(prefix, userID string) string {
	return fmt.Sprintf("%suser:%s:latest", prefix, userID)
}

// SaveSession stores JWT in Redis using key format <prefix>user:<uid>:jwt:<jti> and tracks the latest key pointer per user.
func SaveSession(jti, userID, token string, category sys.CAtegory, expiresAt time.Time) (string, error) {
	client := redisClient()
	if client == nil {
		return "", errRedisNotInit
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = 30 * 24 * time.Hour
	}

	ctx := context.Background()
	prefix := category.GetAuthCookiePrefix()
	sessionKey := buildSessionKey(prefix, userID, jti)
	lookupKey := buildLookupKey(prefix, userID)
	log.Debugf("sessionKey=%s ", sessionKey)
	log.Debugf("lookupKey=%s", lookupKey)
	pipe := client.TxPipeline()
	pipe.Set(ctx, sessionKey, token, ttl)
	pipe.Set(ctx, lookupKey, sessionKey, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}
	return sessionKey, nil
}

// DeleteSession exposes login session deletion for logout flows.
func DeleteSession(claims sys.JWTClaims) error {
	client := redisClient()
	if client == nil {
		return errRedisNotInit
	}
	ctx := context.Background()
	prefix := claims.GetCategoryPrefix()
	sessionKey := buildSessionKey(string(prefix), claims.GetUID(), claims.GetID())
	lookupKey := buildLookupKey(string(prefix), claims.GetUID())

	pipe := client.TxPipeline()
	pipe.Del(ctx, sessionKey)
	pipe.Del(ctx, lookupKey)
	_, err := pipe.Exec(ctx)
	return err
}

// 根据cookie获取用户信息，用在auth中间件
func GetUserFromCookie(cookie string) (string, sys.CAtegory, *sys.JWTClaims, error) {
	client := redisClient()
	if client == nil {
		return "", "", nil, errRedisNotInit
	}

	claims, err := sys.ParseJWT(cookie)
	if err != nil {
		return "", "", nil, err
	}

	sessionKey := buildSessionKey(string(claims.GetCategoryPrefix()), claims.GetUID(), claims.GetID())

	token, err := client.Get(context.Background(), sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", nil, nil
		}
		log.Error(err)
		return "", "", nil, err
	}

	// token retrieved is the same JWT we issued; parsing validates signature/expiry
	if _, err := sys.ParseJWT(token); err != nil {
		log.Error(err)
		return "", "", nil, err
	}

	return claims.GetUID(), claims.GetCategory(), claims, nil
}

// 根据用户 UID 直接查询是否已有会话（用于无 cookie 登录找回）。
func GetSessionCookieByUID(category sys.CAtegory, uid string) (string, *sys.JWTClaims, bool, error) {
	client := redisClient()
	if client == nil {
		return "", nil, false, errRedisNotInit
	}

	ctx := context.Background()
	prefix := category.GetAuthCookiePrefix()
	lookupKey := buildLookupKey(prefix, uid)

	sessionKey, err := client.Get(ctx, lookupKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil, false, nil
		}
		return "", nil, false, err
	}

	token, err := client.Get(ctx, sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil, false, nil
		}
		return "", nil, false, err
	}

	claims, err := sys.ParseJWT(token)
	if err != nil {
		log.Error(err)
		return "", nil, false, err
	}

	expectedSessionKey := buildSessionKey(prefix, claims.GetUID(), claims.GetID())
	if sessionKey != expectedSessionKey {
		return "", nil, false, fmt.Errorf("session key mismatch")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		_ = DeleteSession(*claims)
		return "", nil, false, nil
	}

	// Return the JWT so it can be reissued to the client as cookie value
	return token, claims, true, nil
}
