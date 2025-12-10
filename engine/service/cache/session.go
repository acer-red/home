package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/acer-red/official/engine/util"
	"github.com/redis/go-redis/v9"
	log "github.com/tengfei-xy/go-log"
)

var errRedisNotInit = errors.New("redis not initialized")

// / 关键步骤: 设置Key的名称
func buildSessionKey(prefix, accountID, uid string) string {
	return fmt.Sprintf("%s:account:%s:%s", prefix, accountID, uid)
}

func buildLookupKey(prefix, accountID string) string {
	return fmt.Sprintf("%s:account:%s:latest", prefix, accountID)
}

// SaveSession stores JWT in Redis using key format <prefix>user:<uid>:jwt:<jti> and tracks the latest key pointer per user.
func SaveSession(accountID, userID, token string, category util.CAtegory, expiresAt time.Time) (string, error) {
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
	sessionKey := buildSessionKey(prefix, accountID, userID)
	lookupKey := buildLookupKey(prefix, accountID)
	log.Infof("保存cookie %s %s", sessionKey, token)
	pipe := client.TxPipeline()
	pipe.Set(ctx, sessionKey, token, ttl)
	pipe.Set(ctx, lookupKey, sessionKey, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}
	return sessionKey, nil
}

// DeleteSession exposes login session deletion for logout flows.
func DeleteSession(claims util.JWTClaims) error {
	client := redisClient()
	if client == nil {
		return errRedisNotInit
	}
	ctx := context.Background()
	prefix := claims.GetCategoryPrefix()
	sessionKey := claims.GetKeyName()
	lookupKey := buildLookupKey(string(prefix), claims.AccountID)

	pipe := client.TxPipeline()
	pipe.Del(ctx, sessionKey)
	pipe.Del(ctx, lookupKey)
	_, err := pipe.Exec(ctx)
	return err
}

// 根据cookie获取用户信息，用在auth中间件
func GetUserFromCookie(cookie string) (string, util.CAtegory, *util.JWTClaims, error) {
	log.Debug3f("从redis获取cookie获取用户信息")
	client := redisClient()
	if client == nil {
		log.Error(errRedisNotInit)
		return "", "", nil, errRedisNotInit
	}

	claims, err := util.ParseJWT(cookie)
	if err != nil {
		log.Error(err)
		return "", "", nil, err
	}

	sessionKey := claims.GetKeyName()
	log.Debug3f("session key: %s", sessionKey)
	token, err := client.Get(context.Background(), sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", nil, nil
		}
		log.Error(err)
		return "", "", nil, err
	}

	// token retrieved is the same JWT we issued; parsing validates signature/expiry
	if _, err := util.ParseJWT(token); err != nil {
		log.Error(err)
		return "", "", nil, err
	}

	return claims.GetUID(), claims.GetCategory(), claims, nil
}

// 根据用户 UID 直接查询是否已有会话（用于无 cookie 登录找回）。
func GetSessionCookieByUID(category util.CAtegory, uid string) (string, *util.JWTClaims, bool, error) {
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

	claims, err := util.ParseJWT(token)
	if err != nil {
		log.Error(err)
		return "", nil, false, err
	}

	expectedSessionKey := claims.GetKeyName()
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
