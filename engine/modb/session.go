package modb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/acer-red/home/engine/sys"
	"github.com/redis/go-redis/v9"
	log "github.com/tengfei-xy/go-log"
)

const loginCookiePrefix = "session:login:"

var errRedisNotInit = errors.New("redis not initialized")

type loginSession struct {
	UserID   string       `json:"uid"`
	Category sys.CAtegory `json:"category"`
	Expires  int64        `json:"expires"`
}

func saveLoginSession(cookieVal string, userID string, category sys.CAtegory, expiresAt time.Time) error {
	client := redisClient()
	if client == nil {
		return errRedisNotInit
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = 30 * 24 * time.Hour
	}

	data, err := json.Marshal(loginSession{
		UserID:   userID,
		Category: category,
		Expires:  expiresAt.Unix(),
	})
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%s", loginCookiePrefix, cookieVal)
	return client.Set(context.Background(), key, data, ttl).Err()
}

func deleteLoginSession(cookieVal string) error {
	client := redisClient()
	if client == nil {
		return errRedisNotInit
	}
	key := fmt.Sprintf("%s%s", loginCookiePrefix, cookieVal)
	return client.Del(context.Background(), key).Err()
}

// DeleteLoginSession exposes login session deletion for logout flows.
func DeleteLoginSession(cookieVal string) error {
	return deleteLoginSession(cookieVal)
}

func getUserFromLoginSession(cookieVal string) (User, bool, error) {
	client := redisClient()
	if client == nil {
		return User{}, false, errRedisNotInit
	}
	key := fmt.Sprintf("%s%s", loginCookiePrefix, cookieVal)
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return User{}, false, nil
		}
		log.Error(err)
		return User{}, false, err
	}

	var sess loginSession
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		log.Error(err)
		return User{}, false, err
	}

	return GetUserByIDAndCategory(sess.UserID, sess.Category)
}
